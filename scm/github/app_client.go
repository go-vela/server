// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"

	"github.com/google/go-github/v65/github"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"

	api "github.com/go-vela/server/api/types"
)

// newOAuthTokenClient returns the GitHub OAuth client.
func (c *client) newOAuthTokenClient(ctx context.Context, token string) *github.Client {
	// create the token object for the client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	// create the OAuth client
	tc := oauth2.NewClient(ctx, ts)

	if c.Tracing.Config.EnableTracing {
		tc.Transport = otelhttp.NewTransport(
			tc.Transport,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx, otelhttptrace.WithoutSubSpans())
			}),
		)
	}

	// create the GitHub client from the OAuth client
	github := github.NewClient(tc)

	// ensure the proper URL is set in the GitHub client
	github.BaseURL, _ = url.Parse(c.config.API)

	return github
}

// newGithubAppClient returns the GitHub App client for authenticating as the GitHub App itself using the RoundTripper.
func (c *client) newGithubAppClient() (*github.Client, error) {
	// create a github client based off the existing GitHub App configuration
	client, err := github.NewClient(
		&http.Client{
			Transport: c.AppsTransport,
		}).
		WithEnterpriseURLs(c.config.API, c.config.API)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// newGithubAppInstallationRepoToken returns the GitHub App installation token for a particular repo with granular permissions.
func (c *client) newGithubAppInstallationRepoToken(ctx context.Context, r *api.Repo, repos []string, permissions *github.InstallationPermissions) (string, error) {
	// create a github client based off the existing GitHub App configuration
	client, err := c.newGithubAppClient()
	if err != nil {
		return "", err
	}

	opts := &github.InstallationTokenOptions{
		Repositories: repos,
		Permissions:  permissions,
	}

	// if repo has an install ID, use it to create an installation token
	if r.GetInstallID() != 0 {
		// create installation token for the repo
		t, _, err := client.Apps.CreateInstallationToken(ctx, r.GetInstallID(), opts)
		if err != nil {
			return "", err
		}

		return t.GetToken(), nil
	}

	// list all installations (a.k.a. orgs) where the GitHub App is installed
	installations, _, err := client.Apps.ListInstallations(ctx, &github.ListOptions{})
	if err != nil {
		return "", err
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
		return "", errors.New("unable to find installation ID for repo")
	}

	// create installation token for the repo
	t, _, err := client.Apps.CreateInstallationToken(ctx, id, opts)
	if err != nil {
		return "", err
	}

	return t.GetToken(), nil
}
