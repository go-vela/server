// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"

	"github.com/google/go-github/v76/github"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// newOAuthTokenClient returns the GitHub OAuth client.
func (c *Client) newOAuthTokenClient(ctx context.Context, token string) *github.Client {
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
func (c *Client) newGithubAppClient() (*github.Client, error) {
	if c.AppsTransport == nil {
		return nil, errors.New("unable to create github app client: no AppsTransport configured")
	}

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
func (c *Client) newGithubAppInstallationRepoToken(ctx context.Context, r *api.Repo, repos []string, permissions *github.InstallationPermissions) (*github.InstallationToken, int64, error) {
	// create a github client based off the existing GitHub App configuration
	client, err := c.newGithubAppClient()
	if err != nil {
		return nil, 0, err
	}

	opts := &github.InstallationTokenOptions{
		Repositories: repos,
		Permissions:  permissions,
	}

	id := r.GetInstallID()

	// if the source scm repo has an install ID but the Vela db record does not
	// then use the source repo to create an installation token
	if id == 0 {
		// list all installations (a.k.a. orgs) where the GitHub App is installed
		installations, _, err := client.Apps.ListInstallations(ctx, &github.ListOptions{})
		if err != nil {
			return nil, 0, err
		}

		// iterate through the list of installations
		for _, install := range installations {
			// find the installation that matches the org for the repo
			if strings.EqualFold(install.GetAccount().GetLogin(), r.GetOrg()) {
				if install.GetRepositorySelection() == constants.AppInstallRepositoriesSelectionSelected {
					installationCanReadRepo, err := c.installationCanReadRepo(ctx, r, install)
					if err != nil {
						return nil, 0, fmt.Errorf("installation for org %s exists but unable to check if it can read repo %s: %w", install.GetAccount().GetLogin(), r.GetFullName(), err)
					}

					if !installationCanReadRepo {
						return nil, 0, fmt.Errorf("installation for org %s exists but does not have access to repo %s", install.GetAccount().GetLogin(), r.GetFullName())
					}
				}

				id = install.GetID()
			}
		}
	}

	// failsafe in case the repo does not belong to an org where the GitHub App is installed
	if id == 0 {
		return nil, 0, errors.New("unable to find installation ID for repo")
	}

	// create installation token for the repo
	t, _, err := client.Apps.CreateInstallationToken(ctx, id, opts)
	if err != nil {
		return nil, 0, err
	}

	return t, id, nil
}

// installationCanReadRepo checks if the installation can read the repo.
func (c *Client) installationCanReadRepo(ctx context.Context, r *api.Repo, installation *github.Installation) (bool, error) {
	installationCanReadRepo := false

	if installation.GetRepositorySelection() == constants.AppInstallRepositoriesSelectionSelected {
		client, err := c.newGithubAppClient()
		if err != nil {
			return false, err
		}

		t, _, err := client.Apps.CreateInstallationToken(ctx, installation.GetID(), &github.InstallationTokenOptions{})
		if err != nil {
			return false, err
		}

		client = c.newOAuthTokenClient(ctx, t.GetToken())

		repos, _, err := client.Apps.ListRepos(ctx, &github.ListOptions{})
		if err != nil {
			return false, err
		}

		for _, repo := range repos.Repositories {
			if strings.EqualFold(repo.GetFullName(), r.GetFullName()) {
				installationCanReadRepo = true
				break
			}
		}
	}

	return installationCanReadRepo, nil
}
