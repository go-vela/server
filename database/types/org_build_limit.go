// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"errors"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyOrgBuildLimitOrg defines the error type when an
	// OrgBuildLimit type has an empty Org field provided.
	ErrEmptyOrgBuildLimitOrg = errors.New("empty org build limit org provided")

	// ErrInvalidOrgBuildLimit defines the error type when an
	// OrgBuildLimit type has a BuildLimit field that is not positive.
	ErrInvalidOrgBuildLimit = errors.New("org build limit must be greater than zero")
)

// OrgBuildLimit is the database representation of an
// organization concurrent build limit.
type OrgBuildLimit struct {
	ID         sql.NullInt64  `sql:"id"`
	Org        sql.NullString `sql:"org"`
	BuildLimit sql.NullInt32  `sql:"build_limit"`
	CreatedAt  sql.NullInt64  `sql:"created_at"`
	UpdatedAt  sql.NullInt64  `sql:"updated_at"`
	UpdatedBy  sql.NullString `sql:"updated_by"`
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the OrgBuildLimit type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (o *OrgBuildLimit) Nullify() *OrgBuildLimit {
	if o == nil {
		return nil
	}

	// check if the ID field should be false
	if o.ID.Int64 == 0 {
		o.ID.Valid = false
	}

	// check if the Org field should be false
	if len(o.Org.String) == 0 {
		o.Org.Valid = false
	}

	// check if the CreatedAt field should be false
	if o.CreatedAt.Int64 < 0 {
		o.CreatedAt.Valid = false
	}

	// check if the UpdatedAt field should be false
	if o.UpdatedAt.Int64 < 0 {
		o.UpdatedAt.Valid = false
	}

	// check if the UpdatedBy field should be false
	if len(o.UpdatedBy.String) == 0 {
		o.UpdatedBy.Valid = false
	}

	return o
}

// ToAPI converts the OrgBuildLimit type
// to an API OrgBuildLimit type.
func (o *OrgBuildLimit) ToAPI() *api.OrgBuildLimit {
	oAPI := new(api.OrgBuildLimit)

	oAPI.SetID(o.ID.Int64)
	oAPI.SetOrg(o.Org.String)
	oAPI.SetBuildLimit(o.BuildLimit.Int32)
	oAPI.SetCreatedAt(o.CreatedAt.Int64)
	oAPI.SetUpdatedAt(o.UpdatedAt.Int64)
	oAPI.SetUpdatedBy(o.UpdatedBy.String)

	return oAPI
}

// Validate verifies the necessary fields for
// the OrgBuildLimit type are populated correctly.
func (o *OrgBuildLimit) Validate() error {
	// verify the Org field is populated
	if len(o.Org.String) == 0 {
		return ErrEmptyOrgBuildLimitOrg
	}

	// verify the BuildLimit field is positive
	if o.BuildLimit.Int32 <= 0 {
		return ErrInvalidOrgBuildLimit
	}

	// ensure that the Org field is sanitized
	// to avoid unsafe HTML content
	o.Org = sql.NullString{String: util.Sanitize(o.Org.String), Valid: o.Org.Valid}

	return nil
}

// OrgBuildLimitFromAPI converts the API OrgBuildLimit type
// to a database OrgBuildLimit type.
func OrgBuildLimitFromAPI(o *api.OrgBuildLimit) *OrgBuildLimit {
	orgBuildLimit := &OrgBuildLimit{
		ID:         sql.NullInt64{Int64: o.GetID(), Valid: true},
		Org:        sql.NullString{String: o.GetOrg(), Valid: true},
		BuildLimit: sql.NullInt32{Int32: o.GetBuildLimit(), Valid: true},
		CreatedAt:  sql.NullInt64{Int64: o.GetCreatedAt(), Valid: true},
		UpdatedAt:  sql.NullInt64{Int64: o.GetUpdatedAt(), Valid: true},
		UpdatedBy:  sql.NullString{String: o.GetUpdatedBy(), Valid: true},
	}

	return orgBuildLimit.Nullify()
}
