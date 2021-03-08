// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"context"
	"fmt"
	"net/url"

	"github.com/google/go-github/v29/github"
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
)

var ctx = context.Background()

type config struct {
	// specifies the GitHub address to use
	Address string
	// specifies the GitHub API path to use
	API string
	// specifies the OAuth client ID to use from GitHub
	ClientID string
	// specifies the OAuth client secret to use from GitHub
	ClientSecret string
	// specifies the Vela server address to use
	ServerAddress string
	// specifies the context for the commit status for GitHub
	StatusContext string
	// specifies the Vela web UI address to use
	WebUIAddress string
}

type client struct {
	config  *config
	OAuth   *oauth2.Config
	AuthReq *github.AuthorizationRequest
}

// New returns a Source implementation that integrates with
// a GitHub or a GitHub Enterprise instance.
//
// nolint: golint // ignore returning unexported client
func New(opts ...ClientOpt) (*client, error) {
	// create new GitHub client
	c := new(client)

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
		Scopes:       []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/login/oauth/authorize", c.config.Address),
			TokenURL: fmt.Sprintf("%s/login/oauth/access_token", c.config.Address),
		},
	}

	// create the GitHub authorization object
	c.AuthReq = &github.AuthorizationRequest{
		ClientID:     &c.config.ClientID,
		ClientSecret: &c.config.ClientSecret,
		Scopes:       []github.Scope{"repo", "repo:status", "user:email", "read:user", "read:org"},
	}

	return c, nil
}

// NewTest returns a Source implementation that integrates with the provided
// mock server. Only the url from the mock server is required.
//
// This function is intended for running tests only.
//
// nolint: golint // ignore returning unexported client
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

	// ensure the proper URL is set
	github.BaseURL, _ = url.Parse(c.config.API)

	return github
}
