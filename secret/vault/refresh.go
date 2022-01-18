package vault

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
)

// initialize obtains the vault token from the given auth method
//
// docs: https://www.vaultproject.io/docs/auth
func (c *client) initialize() error {
	c.Logger.Trace("initializing AWS auth headers for vault")

	// declare variables to be utilized within the switch
	var (
		token string
		ttl   time.Duration
	)

	switch c.config.AuthMethod {
	case "aws":
		// create session for aws
		sess, err := session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				CredentialsChainVerboseErrors: aws.Bool(true),
			},
			SharedConfigState: session.SharedConfigEnable,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create aws session for vault")
		}

		// generate sts client for later api calls
		c.AWS.StsClient = sts.New(sess)

		// obtain token from vault
		token, ttl, err = c.getAwsToken()
		if err != nil {
			return err
		}
	}

	c.Vault.SetToken(token)
	c.config.TokenTTL = ttl

	return nil
}

// getAwsToken will retrieve a Vault token for the given IAM principal
//
// docs: https://www.vaultproject.io/docs/auth/aws
func (c *client) getAwsToken() (string, time.Duration, error) {
	headers, err := c.generateAwsAuthHeader()
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

// generateAwsAuthHeader will generate the necessary data
// to send to the Vault server for generating a token.
func (c *client) generateAwsAuthHeader() (map[string]interface{}, error) {
	c.Logger.Trace("generating AWS auth headers for vault")

	req, _ := c.AWS.StsClient.GetCallerIdentityRequest(&sts.GetCallerIdentityInput{})

	// sign the request
	err := req.Sign()
	// will return error if credentials are invalid or expired
	if err != nil {
		return nil, err
	}

	// extract headers from the STS Request
	headersJSON, err := json.Marshal(req.HTTPRequest.Header)
	if err != nil {
		return nil, err
	}

	// read the STS request body
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// construct the vault STS auth header
	//
	// nolint: lll // ignore long line length due to variable names
	loginData := map[string]interface{}{
		"role":                    c.AWS.Role,
		"iam_http_request_method": req.HTTPRequest.Method,
		"iam_request_url":         base64.StdEncoding.EncodeToString([]byte(req.HTTPRequest.URL.String())),
		"iam_request_headers":     base64.StdEncoding.EncodeToString(headersJSON),
		"iam_request_body":        base64.StdEncoding.EncodeToString(requestBody),
	}

	return loginData, nil
}

// refreshToken will refresh the given token if possible or generate a new one entirely.
func (c *client) refreshToken() {
	for {
		// create a channel to signal the need to refresh the token
		refresh := make(chan bool)

		go func() {
			// check if the token TTL is less than the configured duration
			if c.config.TokenTTL < c.config.TokenDuration {
				c.Logger.Debugf("token TTL %v is less than duration %v",
					c.config.TokenTTL,
					c.config.TokenDuration,
				)

				// push a boolean to the channel to signal refresh of the token
				refresh <- true
			}
		}()

		select {
		// refresh the token since the TTL is below the configured duration
		case <-refresh:
			err := c.initialize()
			if err != nil {
				c.Logger.Errorf("failed to refresh vault token: %s", err)
			} else {
				c.Logger.Trace("successfully refreshed vault token")
			}
		// renew the token after sleeping for the configured token duration
		case <-time.After(c.config.TokenDuration):
			// token renewal may fail since the refresh timeframe varies depending on the auth method
			_, err := c.Vault.Auth().Token().RenewSelf(int(c.config.TokenTTL / time.Second))
			if err != nil {
				c.Logger.Errorf("failed to renew vault token: %s", err)

				// fall back to refreshing the token if the renewal fails
				err = c.initialize()
				if err != nil {
					c.Logger.Errorf("failed to refresh vault token: %s", err)
				} else {
					c.Logger.Trace("successfully refreshed vault token")
				}
			}
		}
	}
}
