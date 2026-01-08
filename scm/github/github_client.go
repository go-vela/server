// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"net/http/httptrace"
	"net/url"
	"strings"

	"github.com/google/go-github/v81/github"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/cache/models"
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

// NewAppInstallationToken returns the GitHub App installation token for a particular repo with granular permissions.
func (c *Client) NewAppInstallationToken(ctx context.Context, installID int64, repos []string, permissions map[string]string) (*models.InstallToken, error) {
	var err error

	ghPermissions := new(github.InstallationPermissions)

	for resource, perm := range permissions {
		ghPermissions, err = ApplyInstallationPermissions(resource, perm, ghPermissions)
		if err != nil {
			return nil, err
		}
	}

	opts := &github.InstallationTokenOptions{
		Repositories: repos,
		Permissions:  ghPermissions,
	}

	// create installation token for the repo
	t, _, err := c.AppClient.Apps.CreateInstallationToken(ctx, installID, opts)
	if err != nil {
		return nil, err
	}

	return &models.InstallToken{
		Token:        t.GetToken(),
		InstallID:    installID,
		Repositories: repos,
		Permissions:  permissions,
		Expiration:   t.GetExpiresAt().Unix(),
	}, nil
}

func (c *Client) IsInstallationToken(ctx context.Context, token string) bool {
	return strings.HasPrefix(token, "ghs_")
}

// installationCanReadRepo checks if the installation can read the repo.
func (c *Client) installationCanReadRepo(ctx context.Context, r *api.Repo, installation *github.Installation) (bool, error) {
	installationCanReadRepo := false

	if installation.GetRepositorySelection() == constants.AppInstallRepositoriesSelectionSelected {
		t, _, err := c.AppClient.Apps.CreateInstallationToken(ctx, installation.GetID(), &github.InstallationTokenOptions{})
		if err != nil {
			return false, err
		}

		client := c.newOAuthTokenClient(ctx, t.GetToken())

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
