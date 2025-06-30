// SPDX-License-Identifier: Apache-2.0

package vault

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awssign "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/pkg/errors"
)

// initialize obtains the vault token from the given auth method
//
// docs: https://www.vaultproject.io/docs/auth
func (c *Client) initialize() error {
	c.Logger.Trace("initializing token for vault")

	ctx := context.Background()

	// declare variables to be utilized within the switch
	var (
		token string
		ttl   time.Duration
	)

	switch c.config.AuthMethod {
	case "aws":
		// load AWS config using SDK v2
		cfg, err := awsconfig.LoadDefaultConfig(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to load AWS config for vault")
		}

		// generate sts client for future API calls
		c.AWS.StsClient = sts.NewFromConfig(cfg)

		// obtain token from vault, passing the loaded config to avoid reloading
		token, ttl, err = c.getAwsTokenWithConfig(ctx, &cfg)
		if err != nil {
			return errors.Wrap(err, "failed to get AWS token from vault")
		}
	}

	c.Vault.SetToken(token)
	c.config.TokenTTL = ttl

	return nil
}

// getAwsTokenWithConfig allows passing a pre-loaded AWS config
func (c *Client) getAwsTokenWithConfig(ctx context.Context, cfg *aws.Config) (string, time.Duration, error) {
	headers, err := c.generateAwsAuthHeaders(ctx, *cfg)
	if err != nil {
		return "", 0, err
	}

	c.Logger.Trace("getting AWS token from vault")

	secret, err := c.Vault.Logical().Write("auth/aws/login", headers)
	if err != nil {
		return "", 0, err
	}

	if secret.Auth.ClientToken == "" {
		return "", 0, fmt.Errorf("vault failed to return a token")
	}

	return secret.Auth.ClientToken, time.Duration(secret.Auth.LeaseDuration) * time.Second, nil
}

// generateAwsAuthHeaders gets AWS auth headers for Vault authentication (requires config)
func (c *Client) generateAwsAuthHeaders(ctx context.Context, cfg aws.Config) (map[string]interface{}, error) {
	c.Logger.Trace("generating AWS auth headers for vault")

	// create credentials from the provided config
	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, err
	}

	stsURL := "https://sts.amazonaws.com/"
	requestBody := "Action=GetCallerIdentity&Version=2011-06-15"

	reqCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(reqCtx, "POST", stsURL, strings.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	signer := awssign.NewSigner()
	err = signer.SignHTTP(ctx, creds, httpReq, requestBody, "sts", cfg.Region, time.Now())
	if err != nil {
		return nil, err
	}

	headersJSON, err := json.Marshal(httpReq.Header)
	if err != nil {
		return nil, err
	}

	// see https://developer.hashicorp.com/vault/docs/auth/aws#perform-the-login-operation
	// or https://developer.hashicorp.com/vault/api-docs/auth/aws#iam_request_headers
	loginData := map[string]any{
		"role":                    c.AWS.Role,
		"iam_http_request_method": httpReq.Method,
		"iam_request_url":         base64.StdEncoding.EncodeToString([]byte(httpReq.URL.String())),
		"iam_request_headers":     base64.StdEncoding.EncodeToString(headersJSON),
		"iam_request_body":        base64.StdEncoding.EncodeToString([]byte(requestBody)),
	}

	return loginData, nil
}

// refreshToken will refresh the token used for Vault.
func (c *Client) refreshToken() {
	for {
		c.Logger.Tracef("sleeping for configured vault token duration %v", c.config.TokenDuration)
		// sleep for the configured token duration before refreshing the token
		time.Sleep(c.config.TokenDuration)

		// reinitialize the client to refresh the token
		err := c.initialize()
		if err != nil {
			c.Logger.Errorf("failed to refresh vault token: %s", err)
		} else {
			c.Logger.Trace("successfully refreshed vault token")
		}
	}
}
