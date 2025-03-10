// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/lib/pq"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptySecretName defines the error type when a
	// Secret type has an empty Name field provided.
	ErrEmptySecretName = errors.New("empty secret name provided")

	// ErrEmptySecretOrg defines the error type when a
	// Secret type has an empty Org field provided.
	ErrEmptySecretOrg = errors.New("empty secret org provided")

	// ErrEmptySecretOrgSCMID defines the error type when a
	// Secret type has an empty OrgSCMID field provided.
	ErrEmptySecretOrgSCMID = errors.New("empty secret org scm id provided")

	// ErrEmptySecretRepo defines the error type when a
	// Secret type has an empty Repo field provided.
	ErrEmptySecretRepo = errors.New("empty secret repo provided")

	// ErrEmptySecretRepoSCMID defines the error type when a
	// Secret type has an empty RepoSCMID field provided.
	ErrEmptySecretRepoSCMID = errors.New("empty secret repo scm id provided")

	// ErrEmptySecretTeam defines the error type when a
	// Secret type has an empty Team field provided.
	ErrEmptySecretTeam = errors.New("empty secret team provided")

	// ErrEmptySecretTeamSCMID defines the error type when a
	// Secret type has an empty SCMID field provided.
	ErrEmptySecretTeamSCMID = errors.New("empty secret scm id provided")

	// ErrEmptySecretType defines the error type when a
	// Secret type has an empty Type field provided.
	ErrEmptySecretType = errors.New("empty secret type provided")

	// ErrEmptySecretValue defines the error type when a
	// Secret type has an empty Value field provided.
	ErrEmptySecretValue = errors.New("empty secret value provided")
)

// Secret is the database representation of a secret.
type Secret struct {
	ID                sql.NullInt64  `sql:"id"`
	Org               sql.NullString `sql:"org"`
	OrgSCMID          sql.NullInt64  `sql:"org_scm_id"`
	Repo              sql.NullString `sql:"repo"`
	RepoSCMID         sql.NullInt64  `sql:"repo_scm_id"`
	Team              sql.NullString `sql:"team"`
	TeamSCMID         sql.NullInt64  `sql:"team_scm_id"`
	Name              sql.NullString `sql:"name"`
	Value             sql.NullString `sql:"value"`
	Type              sql.NullString `sql:"type"`
	Images            pq.StringArray `sql:"images"             gorm:"type:varchar(1000)"`
	AllowEvents       sql.NullInt64  `sql:"allow_events"`
	AllowCommand      sql.NullBool   `sql:"allow_command"`
	AllowSubstitution sql.NullBool   `sql:"allow_substitution"`
	CreatedAt         sql.NullInt64  `sql:"created_at"`
	CreatedBy         sql.NullString `sql:"created_by"`
	UpdatedAt         sql.NullInt64  `sql:"updated_at"`
	UpdatedBy         sql.NullString `sql:"updated_by"`
}

// Decrypt will manipulate the existing secret value by
// base64 decoding that value. Then, a AES-256 cipher
// block is created from the encryption key in order to
// decrypt the base64 decoded secret value.
func (s *Secret) Decrypt(key string) error {
	// base64 decode the encrypted secret value
	decoded, err := base64.StdEncoding.DecodeString(s.Value.String)
	if err != nil {
		return err
	}

	// decrypt the base64 decoded secret value
	decrypted, err := util.Decrypt(key, decoded)
	if err != nil {
		return err
	}

	// set the decrypted secret value
	s.Value = sql.NullString{
		String: string(decrypted),
		Valid:  true,
	}

	return nil
}

// Encrypt will manipulate the existing secret value by
// creating a AES-256 cipher block from the encryption
// key in order to encrypt the secret value. Then, the
// secret value is base64 encoded for transport across
// network boundaries.
func (s *Secret) Encrypt(key string) error {
	// encrypt the secret value
	encrypted, err := util.Encrypt(key, []byte(s.Value.String))
	if err != nil {
		return err
	}

	// base64 encode the encrypted secret data to make it network safe
	s.Value = sql.NullString{
		String: base64.StdEncoding.EncodeToString(encrypted),
		Valid:  true,
	}

	return nil
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the Secret type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (s *Secret) Nullify() *Secret {
	if s == nil {
		return nil
	}

	// check if the ID field should be false
	if s.ID.Int64 == 0 {
		s.ID.Valid = false
	}

	// check if the Org field should be false
	if len(s.Org.String) == 0 {
		s.Org.Valid = false
	}

	// check if the OrgSCMID field should be false
	if s.OrgSCMID.Int64 == 0 {
		s.OrgSCMID.Valid = false
	}

	// check if the Repo field should be false
	if len(s.Repo.String) == 0 {
		s.Repo.Valid = false
	}

	// check if the RepoSCMID field should be false
	if s.RepoSCMID.Int64 == 0 {
		s.RepoSCMID.Valid = false
	}

	// check if the Team field should be false
	if len(s.Team.String) == 0 {
		s.Team.Valid = false
	}

	// check if the TeamSCMID field should be false
	if s.TeamSCMID.Int64 == 0 {
		s.TeamSCMID.Valid = false
	}

	// check if the Name field should be false
	if len(s.Name.String) == 0 {
		s.Name.Valid = false
	}

	// check if the Value field should be false
	if len(s.Value.String) == 0 {
		s.Value.Valid = false
	}

	// check if the Value should be false
	if len(s.Type.String) == 0 {
		s.Type.Valid = false
	}

	// check if the AllowEvents field should be false
	if s.AllowEvents.Int64 == 0 {
		s.AllowEvents.Valid = false
	}

	// check if the CreatedAt field should be false
	if s.CreatedAt.Int64 == 0 {
		s.CreatedAt.Valid = false
	}

	// check if the CreatedBy field should be false
	if len(s.CreatedBy.String) == 0 {
		s.CreatedBy.Valid = false
	}

	// check if the UpdatedAt field should be false
	if s.UpdatedAt.Int64 == 0 {
		s.UpdatedAt.Valid = false
	}

	// check if the UpdatedBy field should be false
	if len(s.UpdatedBy.String) == 0 {
		s.UpdatedBy.Valid = false
	}

	return s
}

// ToAPI converts the Secret type
// to a API Secret type.
func (s *Secret) ToAPI() *api.Secret {
	secret := new(api.Secret)

	secret.SetID(s.ID.Int64)
	secret.SetOrg(s.Org.String)
	secret.SetOrgSCMID(s.OrgSCMID.Int64)
	secret.SetRepo(s.Repo.String)
	secret.SetRepoSCMID(s.RepoSCMID.Int64)
	secret.SetTeam(s.Team.String)
	secret.SetTeamSCMID(s.TeamSCMID.Int64)
	secret.SetName(s.Name.String)
	secret.SetValue(s.Value.String)
	secret.SetType(s.Type.String)
	secret.SetImages(s.Images)
	secret.SetAllowEvents(api.NewEventsFromMask(s.AllowEvents.Int64))
	secret.SetAllowCommand(s.AllowCommand.Bool)
	secret.SetAllowSubstitution(s.AllowSubstitution.Bool)
	secret.SetCreatedAt(s.CreatedAt.Int64)
	secret.SetCreatedBy(s.CreatedBy.String)
	secret.SetUpdatedAt(s.UpdatedAt.Int64)
	secret.SetUpdatedBy(s.UpdatedBy.String)

	return secret
}

// Validate verifies the necessary fields for
// the Secret type are populated correctly.
func (s *Secret) Validate() error {
	// verify the Type field is populated
	if len(s.Type.String) == 0 {
		return ErrEmptySecretType
	}

	// verify the Org field is populated
	if len(s.Org.String) == 0 {
		return ErrEmptySecretOrg
	}

	if s.OrgSCMID.Int64 <= 0 {
		return ErrEmptySecretOrgSCMID
	}

	// check if an org or repo secret
	if strings.EqualFold(s.Type.String, constants.SecretRepo) ||
		strings.EqualFold(s.Type.String, constants.SecretOrg) {
		// verify the Repo field is populated
		if len(s.Repo.String) == 0 {
			return ErrEmptySecretRepo
		}

		if s.Type.String == constants.SecretRepo && s.RepoSCMID.Int64 <= 0 {
			return ErrEmptySecretRepoSCMID
		}
	}

	// check if a shared secret
	if strings.EqualFold(s.Type.String, constants.SecretShared) {
		// verify the Team field is populated
		if len(s.Team.String) == 0 {
			return ErrEmptySecretTeam
		}

		if s.TeamSCMID.Int64 <= 0 {
			return ErrEmptySecretTeamSCMID
		}
	}

	// verify the Name field is populated
	if len(s.Name.String) == 0 {
		return ErrEmptySecretName
	}

	// verify the Value field is populated
	if len(s.Value.String) == 0 {
		return ErrEmptySecretValue
	}

	// ensure that all Secret string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	s.Org = sql.NullString{String: util.Sanitize(s.Org.String), Valid: s.Org.Valid}
	s.Repo = sql.NullString{String: util.Sanitize(s.Repo.String), Valid: s.Repo.Valid}
	s.Team = sql.NullString{String: util.Sanitize(s.Team.String), Valid: s.Team.Valid}
	s.Name = sql.NullString{String: util.Sanitize(s.Name.String), Valid: s.Name.Valid}
	s.Type = sql.NullString{String: util.Sanitize(s.Type.String), Valid: s.Type.Valid}

	// ensure that all Images are sanitized
	// to avoid unsafe HTML content
	for i, v := range s.Images {
		s.Images[i] = util.Sanitize(v)
	}

	return nil
}

// SecretFromAPI converts the API Secret type
// to a database Secret type.
func SecretFromAPI(s *api.Secret) *Secret {
	secret := &Secret{
		ID:                sql.NullInt64{Int64: s.GetID(), Valid: true},
		Org:               sql.NullString{String: s.GetOrg(), Valid: true},
		OrgSCMID:          sql.NullInt64{Int64: s.GetOrgSCMID(), Valid: true},
		Repo:              sql.NullString{String: s.GetRepo(), Valid: true},
		RepoSCMID:         sql.NullInt64{Int64: s.GetRepoSCMID(), Valid: true},
		Team:              sql.NullString{String: s.GetTeam(), Valid: true},
		TeamSCMID:         sql.NullInt64{Int64: s.GetTeamSCMID(), Valid: true},
		Name:              sql.NullString{String: s.GetName(), Valid: true},
		Value:             sql.NullString{String: s.GetValue(), Valid: true},
		Type:              sql.NullString{String: s.GetType(), Valid: true},
		Images:            pq.StringArray(s.GetImages()),
		AllowEvents:       sql.NullInt64{Int64: s.GetAllowEvents().ToDatabase(), Valid: true},
		AllowCommand:      sql.NullBool{Bool: s.GetAllowCommand(), Valid: true},
		AllowSubstitution: sql.NullBool{Bool: s.GetAllowSubstitution(), Valid: true},
		CreatedAt:         sql.NullInt64{Int64: s.GetCreatedAt(), Valid: true},
		CreatedBy:         sql.NullString{String: s.GetCreatedBy(), Valid: true},
		UpdatedAt:         sql.NullInt64{Int64: s.GetUpdatedAt(), Valid: true},
		UpdatedBy:         sql.NullString{String: s.GetUpdatedBy(), Valid: true},
	}

	return secret.Nullify()
}
