// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package scm

import (
	"context"
	"net/http"

	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
)

// Service represents the interface for Vela integrating
// with the different supported scm providers.
type Service interface {
	// Service Interface Functions

	// Driver defines a function that outputs
	// the configured scm driver.
	Driver() string

	// Authentication SCM Interface Functions

	// Authorize defines a function that uses the
	// given access token to authorize the user.
	Authorize(string) (string, error)
	// Authenticate defines a function that completes
	// the OAuth workflow for the session.
	Authenticate(http.ResponseWriter, *http.Request, string) (*library.User, error)

	// AuthenticateToken defines a function that completes
	// the OAuth workflow for the session using PAT Token
	AuthenticateToken(*http.Request) (*library.User, error)

	// Login defines a function that begins
	// the OAuth workflow for the session.
	Login(http.ResponseWriter, *http.Request) (string, error)

	// Access SCM Interface Functions

	// OrgAccess defines a function that captures
	// the user's access level for an org.
	OrgAccess(*library.User, string) (string, error)
	// RepoAccess defines a function that captures
	// the user's access level for a repo.
	RepoAccess(*library.User, string, string, string) (string, error)
	// TeamAccess defines a function that captures
	// the user's access level for a team.
	TeamAccess(*library.User, string, string) (string, error)

	// Teams SCM Interface Functions

	// ListUsersTeamsForOrg defines a function that captures
	// the user's teams for an org
	ListUsersTeamsForOrg(*library.User, string) ([]string, error)

	// Changeset SCM Interface Functions

	// Changeset defines a function that captures the list
	// of files changed for a commit.
	//
	// https://en.wikipedia.org/wiki/Changeset.
	Changeset(*library.User, *library.Repo, string) ([]string, error)
	// ChangesetPR defines a function that captures the list
	// of files changed for a pull request.
	//
	// https://en.wikipedia.org/wiki/Changeset.
	ChangesetPR(*library.User, *library.Repo, int) ([]string, error)
	// ChangesetRelease defines a function that captures the list
	// of files changed for a release.
	//
	// https://en.wikipedia.org/wiki/Changeset.
	ChangesetRelease(*library.User, *library.Repo, string) ([]string, string, error)

	// Deployment SCM Interface Functions

	// GetDeployment defines a function that
	// gets a deployment by number and repo.
	GetDeployment(*library.User, *library.Repo, int64) (*library.Deployment, error)
	// GetDeploymentCount defines a function that
	// counts a list of all deployment for a repo.
	GetDeploymentCount(*library.User, *library.Repo) (int64, error)
	// GetDeploymentList defines a function that gets
	// a list of all deployments for a repo.
	GetDeploymentList(*library.User, *library.Repo, int, int) ([]*library.Deployment, error)
	// CreateDeployment defines a function that
	// creates a new deployment.
	CreateDeployment(*library.User, *library.Repo, *library.Deployment) error

	// Repo SCM Interface Functions

	// Config defines a function that captures
	// the pipeline configuration from a repo.
	Config(*library.User, *library.Repo, string) ([]byte, error)
	// ConfigBackoff is a truncated constant backoff wrapper for Config.
	// Retry again in five seconds if Config fails to retrieve yaml/yml file.
	// Will return an error after five failed attempts.
	ConfigBackoff(*library.User, *library.Repo, string) ([]byte, error)
	// Disable defines a function that deactivates
	// a repo by destroying the webhook.
	Disable(*library.User, string, string) error
	// Enable defines a function that activates
	// a repo by creating the webhook.
	Enable(*library.User, string, string, string) (string, error)
	// Status defines a function that sends the
	// commit status for the given SHA from a repo.
	Status(*library.User, *library.Build, string, string) error
	// ListUserRepos defines a function that retrieves
	// all repos with admin rights for the user.
	ListUserRepos(*library.User) ([]*library.Repo, error)
	// GetPullRequest defines a function that retrieves
	// a pull request for a repo.
	GetPullRequest(*library.User, *library.Repo, int) (string, string, string, string, error)
	// GetRepo defines a function that retrieves
	// details for a repo.
	GetRepo(*library.User, *library.Repo) (*library.Repo, error)
	// GetOrgAndRepoName defines a function that retrieves
	// the name of the org and repo in the SCM.
	GetOrgAndRepoName(*library.User, string, string) (string, string, error)
	// GetOrg defines a function that retrieves
	// the name for an org in the SCM.
	GetOrgName(*library.User, string) (string, error)
	// GetHTMLURL defines a function that retrieves
	// a repository file's html_url.
	GetHTMLURL(*library.User, string, string, string, string) (string, error)

	// Webhook SCM Interface Functions

	// ProcessWebhook defines a function that
	// parses the webhook from a repo.
	ProcessWebhook(*http.Request) (*types.Webhook, error)
	// VerifyWebhook defines a function that
	// verifies the webhook from a repo.
	VerifyWebhook(*http.Request, *library.Repo) error
	// RedeliverWebhook defines a function that
	// redelivers the webhook from the SCM.
	RedeliverWebhook(context.Context, *library.User, *library.Repo, *library.Hook) error

	// TODO: Add convert functions to interface?
}
