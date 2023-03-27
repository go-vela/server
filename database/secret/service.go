// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package secret

import (
	"github.com/go-vela/types/library"
)

// SecretService represents the Vela interface for secret
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type SecretService interface {
	// Secret Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateSecretIndexes defines a function that creates the indexes for the secrets table.
	CreateSecretIndexes() error
	// CreateSecretTable defines a function that creates the secrets table.
	CreateSecretTable(string) error

	// Secret Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountSecrets defines a function that gets the count of all secrets.
	CountSecrets() (int64, error)
	// CountSecretsForOrg defines a function that gets the count of secrets by org name.
	CountSecretsForOrg(string, map[string]interface{}) (int64, error)
	// CountSecretsForRepo defines a function that gets the count of secrets by org and repo name.
	CountSecretsForRepo(*library.Repo, map[string]interface{}) (int64, error)
	// CountSecretsForTeam defines a function that gets the count of secrets by org and team name.
	CountSecretsForTeam(string, string, map[string]interface{}) (int64, error)
	// CountSecretsForTeams defines a function that gets the count of secrets by teams within an org.
	CountSecretsForTeams(string, []string, map[string]interface{}) (int64, error)
	// CreateSecret defines a function that creates a new secret.
	CreateSecret(*library.Secret) error
	// DeleteSecret defines a function that deletes an existing secret.
	DeleteSecret(*library.Secret) error
	// GetSecret defines a function that gets a secret by ID.
	GetSecret(int64) (*library.Secret, error)
	// GetSecretForOrg defines a function that gets a secret by org name.
	GetSecretForOrg(string, string) (*library.Secret, error)
	// GetSecretForRepo defines a function that gets a secret by org and repo name.
	GetSecretForRepo(string, *library.Repo) (*library.Secret, error)
	// GetSecretForTeam defines a function that gets a secret by org and team name.
	GetSecretForTeam(string, string, string) (*library.Secret, error)
	// ListSecrets defines a function that gets a list of all secrets.
	ListSecrets() ([]*library.Secret, error)
	// ListSecretsForOrg defines a function that gets a list of secrets by org name.
	ListSecretsForOrg(string, map[string]interface{}, int, int) ([]*library.Secret, int64, error)
	// ListSecretsForRepo defines a function that gets a list of secrets by org and repo name.
	ListSecretsForRepo(*library.Repo, map[string]interface{}, int, int) ([]*library.Secret, int64, error)
	// ListSecretsForTeam defines a function that gets a list of secrets by org and team name.
	ListSecretsForTeam(string, string, map[string]interface{}, int, int) ([]*library.Secret, int64, error)
	// ListSecretsForTeams defines a function that gets a list of secrets by teams within an org.
	ListSecretsForTeams(string, []string, map[string]interface{}, int, int) ([]*library.Secret, int64, error)
	// UpdateSecret defines a function that updates an existing secret.
	UpdateSecret(*library.Secret) error
}
