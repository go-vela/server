// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-github/v70/github"
	"golang.org/x/oauth2"
)

const (
	defaultURL = "https://github.com/"     // Default GitHub URL
	defaultAPI = "https://api.github.com/" // Default GitHub API URL
)

type Client struct {
	githubClient *github.Client
	URL          string
	API          string
}

func (c *Client) Equal(other *Client) bool {
	return (reflect.DeepEqual(c.githubClient.Client(), other.githubClient.Client())) && c.URL == other.URL && c.API == other.API
}

// New returns a Registry implementation that integrates
// with GitHub or a GitHub Enterprise instance.
func New(ctx context.Context, address, token string) (*Client, error) {
	// create the client object
	c := &Client{
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
	c.githubClient = gitClient

	return c, nil
}

// newOAuthTokenClient is a helper function to return the GitHub oauth2 client.
func (c *Client) newOAuthTokenClient(ctx context.Context, token string) *github.Client {
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
