// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package source

import (
	"net/http"

	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
)

// Service represents the interface for Vela integrating
// with the different supported source providers.
type Service interface {
	// Authentication Source Interface Functions

	// Authorize defines a function that uses the
	// given access token to authorize the user.
	Authorize(string) (string, error)
	// Authenticate defines a function that completes
	// the OAuth workflow for the session.
	Authenticate(http.ResponseWriter, *http.Request, string) (*library.User, error)
	// Login defines a function that begins
	// the OAuth workflow for the session.
	Login(http.ResponseWriter, *http.Request) (string, error)

	// Access Source Interface Functions

	// OrgAccess defines a function that captures
	// the user's access level for an org.
	OrgAccess(*library.User, string) (string, error)
	// RepoAccess defines a function that captures
	// the user's access level for a repo.
	RepoAccess(*library.User, string, string) (string, error)
	// TeamAccess defines a function that captures
	// the user's access level for a team.
	TeamAccess(*library.User, string, string) (string, error)

	// Changeset Source Interface Functions

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

	// Deployment Source Interface Functions

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

	// Repo Source Interface Functions

	// Config defines a function that captures
	// the pipeline configuration from a repo.
	Config(*library.User, string, string, string) ([]byte, error)
	// ConfigBackoff is a truncated constant backoff wrapper for Config.
	// Retry again in five seconds if Config fails to retrieve yaml/yml file.
	// Will return an error after five failed attempts.
	ConfigBackoff(*library.User, string, string, string) ([]byte, error)
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

	// Webhook Source Interface Functions

	// ProcessWebhook defines a function that
	// parses the webhook from a repo.
	ProcessWebhook(*http.Request) (*types.Webhook, error)
	// VerifyWebhook defines a function that
	// verifies the webhook from a repo.
	VerifyWebhook(*http.Request, *library.Repo) error

	// TODO: Add convert functions to interface?
}
