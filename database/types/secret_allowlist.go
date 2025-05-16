// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"errors"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/util"
)

// SecretAllowlist is the database representation of a secret allowlist.
type SecretAllowlist struct {
	ID       sql.NullInt64  `sql:"id"`
	SecretID sql.NullInt64  `sql:"secret_id"`
	Repo     sql.NullString `sql:"repo"`
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the SecretAllowlist type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (s *SecretAllowlist) Nullify() *SecretAllowlist {
	if s == nil {
		return nil
	}

	// check if the ID field should be false
	if s.ID.Int64 == 0 {
		s.ID.Valid = false
	}

	if s.SecretID.Int64 == 0 {
		s.SecretID.Valid = false
	}

	// check if the Repo field should be false
	if len(s.Repo.String) == 0 {
		s.Repo.Valid = false
	}

	return s
}

// Validate verifies the necessary fields for
// the Secret type are populated correctly.
func (s *SecretAllowlist) Validate() error {
	// verify the SecretID field is populated
	if s.SecretID.Int64 == 0 {
		return errors.New("secret id cannot be empty")
	}

	// verify the Repo field is populated
	if len(s.Repo.String) == 0 {
		return errors.New("repo cannot be empty")
	}

	return nil
}

// SecretAllowlistFromAPI converts the API Secret and a repo full name
// to a database SecretAllowlist type.
func SecretAllowlistFromAPI(s *api.Secret, repo string) *SecretAllowlist {
	secret := &SecretAllowlist{
		SecretID: sql.NullInt64{Int64: s.GetID(), Valid: true},
		Repo:     sql.NullString{String: util.Sanitize(repo), Valid: true},
	}

	return secret.Nullify()
}
