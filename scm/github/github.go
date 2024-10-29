// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v65/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"

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
	// specifies the Vela server address to use for the GitHub client
	ServerAddress string
	// specifies the Vela server address that the scm provider should use to send Vela webhooks
	ServerWebhookAddress string
	// specifies the context for the commit status to use for the GitHub client
	StatusContext string
	// specifies the Vela web UI address to use for the GitHub client
	WebUIAddress string
	// specifies the OAuth scopes to use for the GitHub client
	Scopes []string
}

type client struct {
	config        *config
	OAuth         *oauth2.Config
	AuthReq       *github.AuthorizationRequest
	Tracing       *tracing.Client
	AppsTransport *AppsTransport
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	Logger *logrus.Entry
}

// New returns a SCM implementation that integrates with
// a GitHub or a GitHub Enterprise instance.
//
//nolint:revive // ignore returning unexported client
func New(ctx context.Context, opts ...ClientOpt) (*client, error) {
	// create new GitHub client
	c := new(client)

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
		Scopes:       c.config.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/login/oauth/authorize", c.config.Address),
			TokenURL: fmt.Sprintf("%s/login/oauth/access_token", c.config.Address),
		},
	}

	var githubScopes []github.Scope
	for _, scope := range c.config.Scopes {
		githubScopes = append(githubScopes, github.Scope(scope))
	}

	// create the GitHub authorization object
	c.AuthReq = &github.AuthorizationRequest{
		ClientID:     &c.config.ClientID,
		ClientSecret: &c.config.ClientSecret,
		Scopes:       githubScopes,
	}

	if c.config.AppID != 0 && len(c.config.AppPrivateKey) > 0 {
		c.Logger.Infof("setting up GitHub App integration for App ID %d", c.config.AppID)

		transport, err := c.newGitHubAppTransport(c.config.AppID, c.config.AppPrivateKey, c.config.API)
		if err != nil {
			return nil, err
		}

		c.AppsTransport = transport

		err = c.ValidateGitHubApp(ctx)
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

// ValidateGitHubApp ensures the GitHub App configuration is valid.
func (c *client) ValidateGitHubApp(ctx context.Context) error {
	client, err := c.newGithubAppClient()
	if err != nil {
		return fmt.Errorf("error creating github app client: %w", err)
	}

	app, _, err := client.Apps.Get(ctx, "")
	if err != nil {
		return fmt.Errorf("error getting github app: %w", err)
	}

	perms := app.GetPermissions()

	type perm struct {
		resource           string
		requiredPermission string
		actualPermission   string
	}

	// GitHub App installation requires the following permissions
	// - contents:read
	// - checks:write
	requiredPermissions := []perm{
		{
			resource:           constants.AppInstallResourceContents,
			requiredPermission: constants.AppInstallPermissionRead,
			actualPermission:   perms.GetContents(),
		},
		{
			resource:           constants.AppInstallResourceChecks,
			requiredPermission: constants.AppInstallPermissionWrite,
			actualPermission:   perms.GetChecks(),
		},
	}

	for _, p := range requiredPermissions {
		err := hasPermission(p.resource, p.requiredPermission, p.actualPermission)
		if err != nil {
			return err
		}
	}

	return nil
}

// hasPermission takes a resource:perm pair and checks if the actual permission matches the expected permission or is supersceded by a higher permission.
func hasPermission(resource, requiredPerm, actualPerm string) error {
	if len(actualPerm) == 0 {
		return fmt.Errorf("github app missing permission %s:%s", resource, requiredPerm)
	}

	permitted := false

	switch requiredPerm {
	case constants.AppInstallPermissionNone:
		permitted = true
	case constants.AppInstallPermissionRead:
		if actualPerm == constants.AppInstallPermissionRead ||
			actualPerm == constants.AppInstallPermissionWrite {
			permitted = true
		}
	case constants.AppInstallPermissionWrite:
		if actualPerm == constants.AppInstallPermissionWrite {
			permitted = true
		}
	default:
		return fmt.Errorf("invalid required permission type: %s", requiredPerm)
	}

	if !permitted {
		return fmt.Errorf("github app requires permission %s:%s, found: %s", constants.AppInstallResourceContents, constants.AppInstallPermissionRead, actualPerm)
	}

	return nil
}

// NewTest returns a SCM implementation that integrates with the provided
// mock server. Only the url from the mock server is required.
//
// This function is intended for running tests only.
//
//nolint:revive // ignore returning unexported client
func NewTest(urls ...string) (*client, error) {
	address := urls[0]
	server := address

	// check if multiple URLs were provided
	if len(urls) > 1 {
		server = urls[1]
	}

	return New(
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
}
