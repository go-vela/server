// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyDashName defines the error type when a
	// User type has an empty Name field provided.
	ErrEmptyDashName = errors.New("empty dashboard name provided")

	// ErrExceededAdminLimit defines the error type when a
	// User type has Admins field provided that exceeds the database limit.
	ErrExceededAdminLimit = errors.New("exceeded admins limit")
)

type (
	// Dashboard is the database representation of a dashboard.
	Dashboard struct {
		ID        uuid.UUID      `sql:"id"          gorm:"type:uuid;default:uuid_generate_v7()"`
		Name      sql.NullString `sql:"name"`
		CreatedAt sql.NullInt64  `sql:"created_at"`
		CreatedBy sql.NullString `sql:"created_by"`
		UpdatedAt sql.NullInt64  `sql:"updated_at"`
		UpdatedBy sql.NullString `sql:"updated_by"`
		Admins    AdminsJSON     `sql:"admins"`
		Repos     DashReposJSON  `sql:"repos"`
	}

	DashReposJSON []*api.DashboardRepo
	AdminsJSON    []*api.User
)

// Value - Implementation of valuer for database/sql for DashReposJSON.
func (r DashReposJSON) Value() (driver.Value, error) {
	valueString, err := json.Marshal(r)
	return string(valueString), err
}

// Scan - Implement the database/sql scanner interface for DashReposJSON.
func (r *DashReposJSON) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &r)
	case string:
		return json.Unmarshal([]byte(v), &r)
	default:
		return fmt.Errorf("wrong type for repos: %T", v)
	}
}

// Value - Implementation of valuer for database/sql for AdminsJSON.
func (a AdminsJSON) Value() (driver.Value, error) {
	valueString, err := json.Marshal(a)
	return string(valueString), err
}

// Scan - Implement the database/sql scanner interface for AdminsJSON.
func (a *AdminsJSON) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &a)
	case string:
		return json.Unmarshal([]byte(v), &a)
	default:
		return fmt.Errorf("wrong type for admins: %T", v)
	}
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the Dashboard type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (d *Dashboard) Nullify() *Dashboard {
	if d == nil {
		return nil
	}

	// check if the Name field should be false
	if len(d.Name.String) == 0 {
		d.Name.Valid = false
	}

	// check if the CreatedAt field should be false
	if d.CreatedAt.Int64 == 0 {
		d.CreatedAt.Valid = false
	}

	// check if the CreatedBy field should be false
	if len(d.CreatedBy.String) == 0 {
		d.CreatedBy.Valid = false
	}

	// check if the UpdatedAt field should be false
	if d.UpdatedAt.Int64 == 0 {
		d.UpdatedAt.Valid = false
	}

	// check if the UpdatedBy field should be false
	if len(d.UpdatedBy.String) == 0 {
		d.UpdatedBy.Valid = false
	}

	return d
}

// ToAPI converts the Dashboard type
// to an API Dashboard type.
func (d *Dashboard) ToAPI() *api.Dashboard {
	dashboard := new(api.Dashboard)

	dashboard.SetID(d.ID.String())
	dashboard.SetName(d.Name.String)
	dashboard.SetAdmins(d.Admins)
	dashboard.SetCreatedAt(d.CreatedAt.Int64)
	dashboard.SetCreatedBy(d.CreatedBy.String)
	dashboard.SetUpdatedAt(d.UpdatedAt.Int64)
	dashboard.SetUpdatedBy(d.UpdatedBy.String)
	dashboard.SetRepos(d.Repos)

	return dashboard
}

// Validate verifies the necessary fields for
// the Dashboard type are populated correctly.
func (d *Dashboard) Validate() error {
	// verify the Name field is populated
	if len(d.Name.String) == 0 {
		return ErrEmptyDashName
	}

	// verify the number of repos
	if len(d.Repos) > constants.DashboardRepoLimit {
		return fmt.Errorf("exceeded repos limit of %d", constants.DashboardRepoLimit)
	}

	// ensure that all Dashboard string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	d.Name = sql.NullString{String: util.Sanitize(d.Name.String), Valid: d.Name.Valid}

	return nil
}

// DashboardFromAPI converts the API Dashboard type
// to a database Dashboard type.
func DashboardFromAPI(d *api.Dashboard) *Dashboard {
	var (
		id  uuid.UUID
		err error
	)

	if d.GetID() == "" {
		id = uuid.New()
	} else {
		id, err = uuid.Parse(d.GetID())
		if err != nil {
			return nil
		}
	}

	dashboard := &Dashboard{
		ID:        id,
		Name:      sql.NullString{String: d.GetName(), Valid: true},
		CreatedAt: sql.NullInt64{Int64: d.GetCreatedAt(), Valid: true},
		CreatedBy: sql.NullString{String: d.GetCreatedBy(), Valid: true},
		UpdatedAt: sql.NullInt64{Int64: d.GetUpdatedAt(), Valid: true},
		UpdatedBy: sql.NullString{String: d.GetUpdatedBy(), Valid: true},
		Admins:    d.GetAdmins(),
		Repos:     d.GetRepos(),
	}

	return dashboard.Nullify()
}
