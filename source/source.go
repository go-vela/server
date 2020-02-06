// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package source

import (
	"net/http"

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
	// LoginCLI defines a function that begins
	// the OAuth workflow for the session.
	LoginCLI(string, string, string) (*library.User, error)

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

	// Repo Source Interface Functions

	// Config defines a function that captures
	// the pipeline configuration from a repo.
	Config(user *library.User, org, name, ref string) ([]byte, error)
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

	// Webhook Source Interface Functions

	// ProcessWebhook defines a function that
	// parses the webhook from a repo.
	ProcessWebhook(*http.Request) (*library.Hook, *library.Repo, *library.Build, error)
	// VerifyWebhook defines a function that
	// verifies the webhook from a repo.
	VerifyWebhook(*http.Request, *library.Repo) error

	// TODO: Add convert functions to interface?
}
