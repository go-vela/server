package vault

import (
	"encoding/base64"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
)

// initialize obtains the vault token from the given auth method
//
// docs: https://www.vaultproject.io/docs/auth
func (c *client) initialize() error {
	// declare variables to be utilize within the switch
	var token string
	var ttl time.Duration

	switch c.AuthMethod {
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
		c.Aws.StsClient = sts.New(sess)

		// obtain token from vault
		token, ttl, err = c.getAwsToken()
		if err != nil {
			return err
		}
	}

	c.Vault.SetToken(token)
	c.TTL = ttl

	return nil
}

// getAwsToken retrieve a Vault token for the given IAM principal
//
// docs: https://www.vaultproject.io/docs/auth/aws
func (c *client) getAwsToken() (string, time.Duration, error) {
	headers, err := c.generateAwsAuthHeader()
	if err != nil {
		return "", 0, err
	}

	logrus.Trace("getting token from vault")
	secret, err := c.Vault.Logical().Write("auth/aws/login", headers)
	if err != nil {
		return "", 0, err
	}

	return secret.Auth.ClientToken, time.Duration(secret.Auth.LeaseDuration) * time.Second, nil
}

// generateAwsAuthHeader generate the necessary data to send to the Vault server for generating a token
func (c *client) generateAwsAuthHeader() (map[string]interface{}, error) {
	logrus.Trace("generating auth headers for vault")
	req, _ := c.Aws.StsClient.GetCallerIdentityRequest(&sts.GetCallerIdentityInput{})

	// sign the request
	// will return error is credentials are invalid or expired
	err := req.Sign()
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
	loginData := make(map[string]interface{})
	loginData["role"] = c.Aws.Role
	loginData["iam_http_request_method"] = req.HTTPRequest.Method
	loginData["iam_request_url"] = base64.StdEncoding.EncodeToString([]byte(req.HTTPRequest.URL.String()))
	loginData["iam_request_headers"] = base64.StdEncoding.EncodeToString(headersJSON)
	loginData["iam_request_body"] = base64.StdEncoding.EncodeToString(requestBody)

	return loginData, nil
}

// refreshToken will refresh the given token if possible or generate a new one entirely
func (c *client) refreshToken() {
	for {
		select {
		case <-time.After(c.Renewal):
			// token refresh ttl varies depending on the auth method
			// aws method provides a variety of config options https://www.vaultproject.io/docs/auth/aws#expiration-times-and-tidying-of-blacklist-and-whitelist-entries
			_, err := c.Vault.Auth().Token().RenewSelf(int(c.TTL / time.Second))
			// fall back to obtaining a new token if the refresh fails
			if err != nil {
				err = c.initialize()
			}

			if err != nil {
				logrus.Errorf("failed to refresh vault token: %s", err)
			} else {
				logrus.Trace("successfully refreshed vault token")
			}
			//case <-c.Finished:
			//	return
		}
	}
}
