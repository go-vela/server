// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/bradleyfalzon/ghinstallation/v2"
	api "github.com/go-vela/server/api/types"
	"github.com/google/go-github/v62/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
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
	// optional and experimental
	GithubAppID         int64
	GithubAppPrivateKey string
}

type client struct {
	config        *config
	OAuth         *oauth2.Config
	AuthReq       *github.AuthorizationRequest
	AppsTransport *ghinstallation.AppsTransport
	// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
	Logger *logrus.Entry
}

// New returns a SCM implementation that integrates with
// a GitHub or a GitHub Enterprise instance.
//
//nolint:revive // ignore returning unexported client
func New(opts ...ClientOpt) (*client, error) {
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

	if c.config.GithubAppID != 0 && len(c.config.GithubAppPrivateKey) > 0 {
		c.Logger.Infof("sourcing private key from path: %s", c.config.GithubAppPrivateKey)

		decodedPEM, err := base64.StdEncoding.DecodeString(c.config.GithubAppPrivateKey)
		if err != nil {
			return nil, fmt.Errorf("error decoding base64: %w", err)
		}

		block, _ := pem.Decode(decodedPEM)
		if block == nil {
			return nil, fmt.Errorf("failed to parse PEM block containing the key")
		}

		privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
		}

		transport := ghinstallation.NewAppsTransportFromPrivateKey(http.DefaultTransport, c.config.GithubAppID, privateKey)

		transport.BaseURL = c.config.API
		c.AppsTransport = transport
	}

	return c, nil
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
		WithAddress(address),
		WithClientID("foo"),
		WithClientSecret("bar"),
		WithServerAddress(server),
		WithServerWebhookAddress(""),
		WithStatusContext("continuous-integration/vela"),
		WithWebUIAddress(address),
	)
}

// helper function to return the GitHub OAuth client.
func (c *client) newClientToken(token string) *github.Client {
	// create the token object for the client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	// create the OAuth client
	tc := oauth2.NewClient(context.Background(), ts)
	// if c.SkipVerify {
	// 	tc.Transport.(*oauth2.Transport).Base = &http.Transport{
	// 		Proxy: http.ProxyFromEnvironment,
	// 		TLSClientConfig: &tls.Config{
	// 			InsecureSkipVerify: true,
	// 		},
	// 	}
	// }

	// create the GitHub client from the OAuth client
	github := github.NewClient(tc)

	// ensure the proper URL is set in the GitHub client
	github.BaseURL, _ = url.Parse(c.config.API)

	return github
}

// helper function to return the GitHub App token.
func (c *client) newGithubAppToken(r *api.Repo) (*github.Client, error) {
	// create a github client based off the existing GitHub App configuration
	client, err := github.NewClient(&http.Client{Transport: c.AppsTransport}).WithEnterpriseURLs(c.config.API, c.config.API)
	if err != nil {
		return nil, err
	}

	// if repo has an install ID, use it to create an installation token
	if r.GetInstallID() != 0 {
		// create installation token for the repo
		t, _, err := client.Apps.CreateInstallationToken(context.Background(), r.GetInstallID(), &github.InstallationTokenOptions{})
		if err != nil {
			panic(err)
		}

		return c.newClientToken(t.GetToken()), nil
	}

	// list all installations (a.k.a. orgs) where the GitHub App is installed
	installations, _, err := client.Apps.ListInstallations(context.Background(), &github.ListOptions{})
	if err != nil {
		return nil, err
	}

	var id int64
	// iterate through the list of installations
	for _, install := range installations {
		// find the installation that matches the org for the repo
		if strings.EqualFold(install.GetAccount().GetLogin(), r.GetOrg()) {
			id = install.GetID()
		}
	}

	// failsafe in case the repo does not belong to an org where the GitHub App is installed
	if id == 0 {
		return nil, err
	}

	// create installation token for the repo
	t, _, err := client.Apps.CreateInstallationToken(context.Background(), id, &github.InstallationTokenOptions{})
	if err != nil {
		panic(err)
	}

	return c.newClientToken(t.GetToken()), nil
}
