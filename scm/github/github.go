// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/go-github/v81/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/tracing"
)

const (
	defaultURL = "https://github.com/"     // Default GitHub URL
	defaultAPI = "https://api.github.com/" // Default GitHub API URL

	// events for repo webhooks.
	eventPush         = "push"
	eventPullRequest  = "pull_request"
	eventDeployment   = "deployment"
	eventIssueComment = "issue_comment"
	eventRepository   = "repository"
	eventInitialize   = "initialize"

	// merge queue prefix.
	mergeQueueBranchPrefix = "gh-readonly-queue/"
)

type config struct {
	// specifies the address to use for the GitHub client
	Address string
	// specifies the API endpoint to use for the GitHub client
	API string
	// specifies the OAuth client ID from GitHub to use for the GitHub client
	ClientID string
	// specifies the OAuth client secret from GitHub to use for the GitHub client
	ClientSecret string
	// specifies the ID for the Vela GitHub App
	AppID int64
	// specifies the App private key to use for the GitHub client when interacting with App resources
	AppPrivateKey string
	// specifies the App private key to use for the GitHub client when interacting with App resources
	AppPrivateKeyPath string
	// specifics the App permissions set
	AppPermissions []string
	// specifies the Vela server address to use for the GitHub client
	ServerAddress string
	// specifies the Vela server address that the scm provider should use to send Vela webhooks
	ServerWebhookAddress string
	// specifies the context for the commit status to use for the GitHub client
	StatusContext string
	// specifies the Vela web UI address to use for the GitHub client
	WebUIAddress string
	// specifies the OAuth scopes to use for the GitHub client
	OAuthScopes []string
}

type Client struct {
	config    *config
	OAuth     *oauth2.Config
	AuthReq   *github.AuthorizationRequest
	Tracing   *tracing.Client
	AppClient *github.Client

	settings.SCM

	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	Logger *logrus.Entry
}

// New returns a SCM implementation that integrates with
// a GitHub or a GitHub Enterprise instance.
func New(ctx context.Context, opts ...ClientOpt) (*Client, error) {
	// create new GitHub client
	c := new(Client)

	// create new fields
	c.config = new(config)
	c.OAuth = new(oauth2.Config)
	c.AuthReq = new(github.AuthorizationRequest)

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#StandardLogger
	logger := logrus.StandardLogger()

	// create new logger for the client
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#NewEntry
	c.Logger = logrus.NewEntry(logger).WithField("scm", c.Driver())

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(c)
		if err != nil {
			return nil, err
		}
	}

	// create the GitHub OAuth config object
	c.OAuth = &oauth2.Config{
		ClientID:     c.config.ClientID,
		ClientSecret: c.config.ClientSecret,
		Scopes:       c.config.OAuthScopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/login/oauth/authorize", c.config.Address),
			TokenURL: fmt.Sprintf("%s/login/oauth/access_token", c.config.Address),
		},
	}

	var oauthScopes []github.Scope
	for _, scope := range c.config.OAuthScopes {
		oauthScopes = append(oauthScopes, github.Scope(scope))
	}

	// create the GitHub authorization object
	c.AuthReq = &github.AuthorizationRequest{
		ClientID:     &c.config.ClientID,
		ClientSecret: &c.config.ClientSecret,
		Scopes:       oauthScopes,
	}

	var err error

	if c.config.AppID != 0 {
		c.Logger.Infof("configurating github app integration for app_id %d", c.config.AppID)

		var privateKeyPEM []byte

		if len(c.config.AppPrivateKey) == 0 && len(c.config.AppPrivateKeyPath) == 0 {
			return nil, errors.New("GitHub App ID provided but no valid private key was provided in either VELA_SCM_APP_PRIVATE_KEY or VELA_SCM_APP_PRIVATE_KEY_PATH")
		}

		if len(c.config.AppPrivateKey) > 0 {
			privateKeyPEM, err = base64.StdEncoding.DecodeString(c.config.AppPrivateKey)
			if err != nil {
				return nil, fmt.Errorf("error decoding base64: %w", err)
			}
		} else {
			// try reading from path if necessary
			c.Logger.Infof("no VELA_SCM_APP_PRIVATE_KEY provided, reading github app private key from path %s", c.config.AppPrivateKeyPath)

			privateKeyPEM, err = os.ReadFile(c.config.AppPrivateKeyPath)
			if err != nil {
				return nil, err
			}
		}

		if len(privateKeyPEM) == 0 {
			return nil, errors.New("GitHub App ID provided but no valid private key was provided in either VELA_SCM_APP_PRIVATE_KEY or VELA_SCM_APP_PRIVATE_KEY_PATH")
		}

		block, _ := pem.Decode(privateKeyPEM)
		if block == nil {
			return nil, fmt.Errorf("failed to parse GitHub App private key PEM block containing the key")
		}

		parsedPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse GitHub App RSA private key: %w", err)
		}

		// create a github client based off the existing GitHub App configuration
		c.AppClient = github.NewClient(
			&http.Client{
				Transport: c.newGitHubAppTransport(c.config.AppID, c.config.API, parsedPrivateKey),
			})

		if c.config.API != defaultAPI {
			c.AppClient, err = c.AppClient.WithEnterpriseURLs(c.config.API, c.config.API)
			if err != nil {
				return nil, err
			}
		}

		err = c.ValidateGitHubApp(ctx)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// ValidateGitHubApp ensures the GitHub App configuration is valid.
func (c *Client) ValidateGitHubApp(ctx context.Context) error {
	app, _, err := c.AppClient.Apps.Get(ctx, "")
	if err != nil {
		return fmt.Errorf("error getting github app: %w", err)
	}

	appPermissions := app.GetPermissions()

	type perm struct {
		resource           string
		requiredPermission string
		actualPermission   string
	}

	// the GitHub App installation requires the same permissions as provided at runtime
	requiredPermissions := []perm{}

	// retrieve the required permissions for checking
	for _, permission := range c.config.AppPermissions {
		splitPerm := strings.Split(permission, ":")
		if len(splitPerm) != 2 {
			return fmt.Errorf("invalid app permission format %s, expected resource:permission", permission)
		}

		resource := splitPerm[0]
		requiredPermission := splitPerm[1]

		actual, err := GetInstallationPermission(resource, appPermissions)
		if err != nil {
			return err
		}

		perm := perm{
			resource:           resource,
			requiredPermission: requiredPermission,
			actualPermission:   actual,
		}
		requiredPermissions = append(requiredPermissions, perm)
	}

	// verify the app permissions
	for _, p := range requiredPermissions {
		err := InstallationHasPermission(p.resource, p.requiredPermission, p.actualPermission)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewTest returns a SCM implementation that integrates with the provided
// mock server. Only the url from the mock server is required.
//
// This function is intended for running tests only.
func NewTest(urls ...string) (*Client, error) {
	var (
		repoRoleMap = map[string]string{
			"admin":    constants.PermissionAdmin,
			"write":    constants.PermissionWrite,
			"maintain": constants.PermissionWrite,
			"triage":   constants.PermissionRead,
			"read":     constants.PermissionRead,
		}

		orgRoleMap = map[string]string{
			"admin":  constants.PermissionAdmin,
			"member": constants.PermissionRead,
		}

		teamRoleMap = map[string]string{
			"maintainer": constants.PermissionAdmin,
			"member":     constants.PermissionRead,
		}
	)

	address := urls[0]
	server := address

	// check if multiple URLs were provided
	if len(urls) > 1 {
		server = urls[1]
	}

	c, err := New(
		context.Background(),
		WithAddress(address),
		WithClientID("foo"),
		WithClientSecret("bar"),
		WithServerAddress(server),
		WithServerWebhookAddress(""),
		WithStatusContext("continuous-integration/vela"),
		WithWebUIAddress(address),
		WithTracing(&tracing.Client{Config: tracing.Config{EnableTracing: false}}),
	)
	if err != nil {
		return nil, err
	}

	c.SetRepoRoleMap(repoRoleMap)
	c.SetOrgRoleMap(orgRoleMap)
	c.SetTeamRoleMap(teamRoleMap)

	return c, nil
}
