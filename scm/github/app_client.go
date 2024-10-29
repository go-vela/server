// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/go-github/v65/github"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// NewGitHubAppTransport creates a new GitHub App transport for authenticating as the GitHub App.
func NewGitHubAppTransport(appID int64, privateKey, baseURL string) (*AppsTransport, error) {
	decodedPEM, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("error decoding base64: %w", err)
	}

	block, _ := pem.Decode(decodedPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	_privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
	}

	transport := NewAppsTransportFromPrivateKey(http.DefaultTransport, appID, _privateKey)
	transport.BaseURL = baseURL

	return transport, nil
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
	if len(perms.GetContents()) == 0 ||
		(perms.GetContents() != constants.AppInstallPermissionRead && perms.GetContents() != constants.AppInstallPermissionWrite) {
		return fmt.Errorf("github app requires contents:read permissions, found: %s", perms.GetContents())
	}

	if len(perms.GetChecks()) == 0 ||
		perms.GetChecks() != constants.AppInstallPermissionWrite {
		return fmt.Errorf("github app requires checks:write permissions, found: %s", perms.GetChecks())
	}

	return nil
}

// newGithubAppClient returns the GitHub App client for authenticating as the GitHub App itself using the RoundTripper.
func (c *client) newGithubAppClient() (*github.Client, error) {
	// todo: create transport using context to apply tracing
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
