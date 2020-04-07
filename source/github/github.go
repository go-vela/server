// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"strings"

	"github.com/google/go-github/v29/github"
	"golang.org/x/oauth2"

	"github.com/urfave/cli"
)

const (
	defaultURL = "https://github.com"     // Default GitHub URL
	defaultAPI = "https://api.github.com" // Default GitHub API URL

	// events for repo webhooks
	eventPush         = "push"
	eventPullRequest  = "pull_request"
	eventTag          = "tag"
	eventDeployment   = "deployment"
	eventIssueComment = "issue_comment"
)

var ctx = context.Background()

type client struct {
	URL           string
	API           string
	LocalHost     string
	WebUIHost     string
	StatusContext string
	OConfig       *oauth2.Config
	AuthReq       *github.AuthorizationRequest
}

// New returns a Source implementation that integrates with GitHub or a
// GitHub Enterprise instance.
func New(c *cli.Context) (*client, error) {
	// create the client object
	client := &client{
		URL:           defaultURL,
		API:           defaultAPI,
		LocalHost:     c.String("server-addr"),
		WebUIHost:     c.String("webui-addr"),
		StatusContext: c.String("source-context"),
	}

	// ensure we have the URL and API set
	if defaultURL != c.String("source-url") {
		client.URL = strings.TrimSuffix(c.String("source-url"), "/")
		client.API = client.URL + "/api/v3/"
	}

	sourceClient := c.String("source-client")
	sourceClientSecret := c.String("source-secret")

	// create the OAuth config object
	client.OConfig = &oauth2.Config{
		ClientID:     sourceClient,
		ClientSecret: sourceClientSecret,
		Scopes:       []string{"repo", "repo:status", "user:email", "read:user", "read:org"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  fmt.Sprintf("%s/login/oauth/authorize", client.URL),
			TokenURL: fmt.Sprintf("%s/login/oauth/access_token", client.URL),
		},
	}

	// create the GitHub authorization object
	client.AuthReq = &github.AuthorizationRequest{
		ClientID:     &sourceClient,
		ClientSecret: &sourceClientSecret,
		Scopes:       []github.Scope{"repo", "repo:status", "user:email", "read:user", "read:org"},
	}

	return client, nil
}

// NewTest returns a Source implementation that integrates with the provided
// mock server. Only the url from the mock server is required.
//
// This function is intended for running tests only.
func NewTest(urls ...string) (*client, error) {
	return New(createTestContext(urls...))
}

// helper function to create a cli.Context for the NewTest function.
//
// This function is intended for running tests only.
func createTestContext(urls ...string) *cli.Context {
	set := flag.NewFlagSet("test", 0)
	set.String("server-addr", urls[0], "doc")

	if len(urls) > 1 {
		set.Set("server-addr", urls[1])
	}

	set.String("webui-addr", urls[0], "doc")
	set.String("source-url", urls[0], "doc")
	set.String("source-client", "foo", "doc")
	set.String("source-secret", "bar", "doc")
	set.String("source-context", "continuous-integration/vela", "doc")

	return cli.NewContext(nil, set, nil)
}

// helper function to return the GitHub OAuth client
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
	github.BaseURL, _ = url.Parse(c.API)

	return github
}

// helper function to return the GitHub Basic auth client
func (c *client) newClientBasicAuth(username, password, otp string) *github.Client {
	// create the transport object for the client
	auth := github.BasicAuthTransport{
		Username: username,
		Password: password,
		OTP:      otp,
	}

	// create the GitHub client from the OAuth client
	github := github.NewClient(auth.Client())

	// ensure the proper URL is set
	github.BaseURL, _ = url.Parse(c.API)

	return github
}
