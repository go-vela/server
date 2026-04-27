// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"

	"github.com/google/go-github/v84/github"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// newOAuthTokenClient creates a new GitHub client using the provided OAuth token.
func (c *Client) newUserOAuthTokenClient(ctx context.Context, user *api.User) *github.Client {
	var oauthToken *oauth2.Token

	if user.GetTokenExp() > 0 && time.Now().Unix() >= user.GetTokenExp() {
		oauthToken = &oauth2.Token{
			RefreshToken: user.GetOAuthRefreshToken(),
		}

		ts := c.OAuth.TokenSource(ctx, oauthToken)
		newToken, err := ts.Token()
		if err == nil {
			oauthToken = newToken

			user.SetToken(newToken.AccessToken)
			user.SetTokenExp(newToken.Expiry.Unix())

			_, err = c.Database.UpdateUser(ctx, user)
			if err != nil {
				c.Logger.Errorf("unable to update user token for user %s: %v", user.GetName(), err)
			}
		}
	} else {
		oauthToken = &oauth2.Token{
			AccessToken: user.GetToken(),
		}
	}

	ts := oauth2.StaticTokenSource(oauthToken)

	tc := oauth2.NewClient(ctx, ts)

	if c.Tracing.Config.EnableTracing {
		tc.Transport = otelhttp.NewTransport(
			tc.Transport,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx, otelhttptrace.WithoutSubSpans())
			}),
		)
	}

	githubClient := github.NewClient(tc)

	githubClient.BaseURL, _ = url.Parse(c.config.API)

	return githubClient
}

// newTokenClient creates a new GitHub client using the provided token.
func (c *Client) newTokenClient(ctx context.Context, token string) *github.Client {
	c.Logger.Debugf("DEBUG: github token client")

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})

	tc := oauth2.NewClient(ctx, ts)

	if c.Tracing.Config.EnableTracing {
		tc.Transport = otelhttp.NewTransport(
			tc.Transport,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx, otelhttptrace.WithoutSubSpans())
			}),
		)
	}

	githubClient := github.NewClient(tc)

	githubClient.BaseURL, _ = url.Parse(c.config.API)

	return githubClient
}

// scopedAccessTokenResponse is the response from creating a scoped access token.
type scopedAccessTokenResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

// basicAuthRoundTripper is an http.RoundTripper that applies HTTP Basic
// Authentication to each request.
type basicAuthRoundTripper struct {
	username string
	password string
	base     http.RoundTripper
}

func (t *basicAuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.SetBasicAuth(t.username, t.password)

	if t.base == nil {
		return http.DefaultTransport.RoundTrip(req)
	}

	return t.base.RoundTrip(req)
}

// applyUserTokenPermissions converts a map of resource:level pairs to a
// UserAccessTokenPermissions struct for use with the go-github SDK.
func applyUserTokenPermissions(perms map[string]string) *github.UserAccessTokenPermissions {
	p := &github.UserAccessTokenPermissions{}

	for resource, level := range perms {
		level := level

		switch strings.ToLower(resource) {
		case AppInstallResourceContents:
			p.Contents = &level
		case AppInstallResourceChecks:
			p.Checks = &level
		case AppInstallResourcePackages:
			p.Packages = &level
		case AppInstallResourceStatuses:
			p.Statuses = &level
		case AppInstallResourceDeployments:
			p.Deployments = &level
		case AppInstallResourcePullRequests:
			p.PullRequests = &level
		case AppInstallResourceIssues:
			p.Issues = &level
		}
	}

	return p
}

// CreateScopedAccessToken creates a scoped access token from a user's existing token
// using the GitHub App's client credentials.
//
// This calls POST /applications/{client_id}/token/scoped with BasicAuth.
// See: https://docs.github.com/en/rest/apps/apps#create-a-scoped-access-token
func (c *Client) CreateScopedAccessToken(ctx context.Context, repo *api.Repo, repos []string, permissions map[string]string) (*scopedAccessTokenResponse, error) {
	owner := repo.GetOwner()

	if owner.GetTokenExp() > 0 && time.Now().Unix() >= owner.GetTokenExp() {
		oauthToken := &oauth2.Token{
			RefreshToken: owner.GetOAuthRefreshToken(),
		}

		ts := c.OAuth.TokenSource(ctx, oauthToken)
		newToken, err := ts.Token()
		if err != nil {
			return nil, fmt.Errorf("unable to refresh token for repository owner %s: %w", owner.GetName(), err)
		}

		owner.SetToken(newToken.AccessToken)
		owner.SetTokenExp(newToken.Expiry.Unix())

		_, err = c.Database.UpdateUser(ctx, owner)
		if err != nil {
			c.Logger.Errorf("unable to update user token for repository owner %s: %v", owner.GetName(), err)
		}
	}

	opts := &github.ScopedUserAccessTokenOptions{
		AccessToken:  owner.GetToken(),
		Target:       repo.GetOrg(),
		Repositories: repos,
		Permissions:  applyUserTokenPermissions(permissions),
	}

	var transport http.RoundTripper = http.DefaultTransport

	if c.Tracing != nil && c.Tracing.Config.EnableTracing {
		transport = otelhttp.NewTransport(
			transport,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx, otelhttptrace.WithoutSubSpans())
			}),
		)
	}

	client := github.NewClient(&http.Client{
		Transport: &basicAuthRoundTripper{
			username: c.config.ClientID,
			password: c.config.ClientSecret,
			base:     transport,
		},
	})
	client.BaseURL, _ = url.Parse(c.config.API)

	token, _, err := client.Apps.CreateScopedUserAccessToken(ctx, c.config.ClientID, opts)
	if err != nil {
		return nil, fmt.Errorf("unable to create scoped access token: %w", err)
	}

	return &scopedAccessTokenResponse{
		Token:     token.GetToken(),
		ExpiresAt: token.ExpiresAt.Time,
	}, nil
}

// NewAppInstallationToken returns the GitHub App installation token for a particular repo with granular permissions.
func (c *Client) NewAppInstallationToken(ctx context.Context, installID int64, repos []string, permissions map[string]string) (string, error) {
	c.Logger.Debugf("DEBUG: app installation token")

	var err error

	ghPermissions := new(github.InstallationPermissions)

	for resource, perm := range permissions {
		ghPermissions, err = ApplyInstallationPermissions(resource, perm, ghPermissions)
		if err != nil {
			return "", err
		}
	}

	opts := &github.InstallationTokenOptions{
		Repositories: repos,
		Permissions:  ghPermissions,
	}

	// create installation token for the repo
	t, _, err := c.AppClient.Apps.CreateInstallationToken(ctx, installID, opts)
	if err != nil {
		return "", err
	}

	return t.GetToken(), nil
}

func (c *Client) IsGitToken(ctx context.Context, token string) bool {
	return strings.HasPrefix(token, "ghu_")
}

// installationCanReadRepo checks if the installation can read the repo.
func (c *Client) installationCanReadRepo(ctx context.Context, org, repo string, installation *github.Installation) (bool, error) {
	installationCanReadRepo := false

	if installation.GetRepositorySelection() == constants.AppInstallRepositoriesSelectionSelected {
		t, _, err := c.AppClient.Apps.CreateInstallationToken(ctx, installation.GetID(), &github.InstallationTokenOptions{Repositories: []string{repo}})
		if err != nil {
			return false, err
		}

		client := c.newTokenClient(ctx, t.GetToken())

		_, _, err = client.Repositories.Get(ctx, org, repo)
		if err == nil {
			installationCanReadRepo = true
		}
	}

	if installation.GetRepositorySelection() == constants.AppInstallRepositoriesSelectionAll {
		installationCanReadRepo = true
	}

	return installationCanReadRepo, nil
}
