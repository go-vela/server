// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-github/v69/github"
	"golang.org/x/oauth2"
)

const (
	defaultURL = "https://github.com/"     // Default GitHub URL
	defaultAPI = "https://api.github.com/" // Default GitHub API URL
)

type client struct {
	Github *github.Client
	URL    string
	API    string
}

func (c *client) Equal(other *client) bool {
	return (reflect.DeepEqual(c.Github.Client(), other.Github.Client())) && c.URL == other.URL && c.API == other.API
}

// New returns a Registry implementation that integrates
// with GitHub or a GitHub Enterprise instance.
//
//nolint:revive,unused // ignore returning unexported client
func New(ctx context.Context, address, token string) (*client, error) {
	// create the client object
	c := &client{
		URL: defaultURL,
		API: defaultAPI,
	}

	// ensure we have the URL and API set
	if len(address) > 0 {
		if !strings.EqualFold(c.URL, address) {
			c.URL = strings.Trim(address, "/")
			if !strings.Contains(c.URL, "https://github.com") {
				c.API = c.URL + "/api/v3/"
			}
		}
	}

	// create the GitHub client
	gitClient := github.NewClient(nil)
	// ensure the proper URL is set
	gitClient.BaseURL, _ = url.Parse(c.API)

	if len(token) > 0 {
		// create GitHub OAuth client with user's token
		gitClient = c.newOAuthTokenClient(ctx, token)
	}

	// overwrite the github client
	c.Github = gitClient

	//nolint:revive // ignore returning unexported engine
	return c, nil
}

// newOAuthTokenClient is a helper function to return the GitHub oauth2 client.
func (c *client) newOAuthTokenClient(ctx context.Context, token string) *github.Client {
	// create the token object for the client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	// create the OAuth client
	tc := oauth2.NewClient(ctx, ts)

	// create the GitHub client from the OAuth client
	github := github.NewClient(tc)

	// ensure the proper URL is set
	github.BaseURL, _ = url.Parse(c.API)

	return github
}
