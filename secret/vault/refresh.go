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
	"github.com/sirupsen/logrus"
)

// initialize obtains the vault token from the given auth method
//
// docs: https://www.vaultproject.io/docs/auth
func (c *client) initialize() error {
	logrus.Trace("initializing token for vault")

	// declare variables to be utilized within the switch
	var (
		token string
		ttl   time.Duration
	)

	switch c.config.AuthMethod {
	case "aws":
		// create session for AWS
		sess, err := session.NewSessionWithOptions(session.Options{
			Config: aws.Config{
				CredentialsChainVerboseErrors: aws.Bool(true),
			},
			SharedConfigState: session.SharedConfigEnable,
		})
		if err != nil {
			return errors.Wrap(err, "failed to create aws session for vault")
		}

		// generate sts client for future API calls
		c.AWS.StsClient = sts.New(sess)

		// obtain token from vault
		token, ttl, err = c.getAwsToken()
		if err != nil {
			return errors.Wrap(err, "failed to get AWS token from vault")
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

	logrus.Trace("getting AWS token from vault")
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
	logrus.Trace("generating AWS auth headers for vault")
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

// refreshToken will refresh the token used for Vault.
func (c *client) refreshToken() {
	for {
		logrus.Tracef("sleeping for configured vault token duration %v", c.config.TokenDuration)
		// sleep for the configured token duration before refreshing the token
		time.Sleep(c.config.TokenDuration)

		// reinitialize the client to refresh the token
		err := c.initialize()
		if err != nil {
			logrus.Errorf("failed to refresh vault token: %s", err)
		} else {
			logrus.Trace("successfully refreshed vault token")
		}
	}
}
