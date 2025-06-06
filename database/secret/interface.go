// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// SecretInterface represents the Vela interface for secret
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type SecretInterface interface {
	// Secret Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateSecretIndexes defines a function that creates the indexes for the secrets table.
	CreateSecretIndexes(context.Context) error
	// CreateSecretTables defines a function that creates the secrets table.
	CreateSecretTables(context.Context, string) error

	// Secret Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountSecrets defines a function that gets the count of all secrets.
	CountSecrets(context.Context) (int64, error)
	// CountSecretsForOrg defines a function that gets the count of secrets by org name.
	CountSecretsForOrg(context.Context, string, map[string]interface{}) (int64, error)
	// CountSecretsForRepo defines a function that gets the count of secrets by org and repo name.
	CountSecretsForRepo(context.Context, *api.Repo, map[string]interface{}) (int64, error)
	// CountSecretsForTeam defines a function that gets the count of secrets by org and team name.
	CountSecretsForTeam(context.Context, string, string, map[string]interface{}) (int64, error)
	// CountSecretsForTeams defines a function that gets the count of secrets by teams within an org.
	CountSecretsForTeams(context.Context, string, []string, map[string]interface{}) (int64, error)
	// CreateSecret defines a function that creates a new secret.
	CreateSecret(context.Context, *api.Secret) (*api.Secret, error)
	// DeleteSecret defines a function that deletes an existing secret.
	DeleteSecret(context.Context, *api.Secret) error
	// FillSecretAllowlist defines a function that fills a secret with its allowlist.
	FillSecretAllowlist(context.Context, *api.Secret) (*api.Secret, error)
	// FillSecretsAllowlists defines a function that takes a list of secrets and fills the allowlists on each.
	FillSecretsAllowlists(context.Context, []*api.Secret) ([]*api.Secret, error)
	// GetSecret defines a function that gets a secret by ID.
	GetSecret(context.Context, int64) (*api.Secret, error)
	// GetSecretForOrg defines a function that gets a secret by org name.
	GetSecretForOrg(context.Context, string, string) (*api.Secret, error)
	// GetSecretForRepo defines a function that gets a secret by org and repo name.
	GetSecretForRepo(context.Context, string, *api.Repo) (*api.Secret, error)
	// GetSecretForTeam defines a function that gets a secret by org and team name.
	GetSecretForTeam(context.Context, string, string, string) (*api.Secret, error)
	// ListSecrets defines a function that gets a list of all secrets.
	ListSecrets(context.Context) ([]*api.Secret, error)
	// ListSecretsForOrg defines a function that gets a list of secrets by org name.
	ListSecretsForOrg(context.Context, string, map[string]interface{}, int, int) ([]*api.Secret, error)
	// ListSecretsForRepo defines a function that gets a list of secrets by org and repo name.
	ListSecretsForRepo(context.Context, *api.Repo, map[string]interface{}, int, int) ([]*api.Secret, error)
	// ListSecretsForTeam defines a function that gets a list of secrets by org and team name.
	ListSecretsForTeam(context.Context, string, string, map[string]interface{}, int, int) ([]*api.Secret, error)
	// ListSecretsForTeams defines a function that gets a list of secrets by teams within an org.
	ListSecretsForTeams(context.Context, string, []string, map[string]interface{}, int, int) ([]*api.Secret, error)
	// MigrateSecrets defines a function that updates the org and name of all repo secrets when there is a name change.
	MigrateSecrets(context.Context, string, string, string, string) error
	// UpdateSecret defines a function that updates an existing secret.
	UpdateSecret(context.Context, *api.Secret) (*api.Secret, error)
}
