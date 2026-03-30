// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"errors"

	api "github.com/go-vela/server/api/types"
)

var (
	// ErrEmptyInstallID defines the error type when a
	// Installation type has an empty InstallID field provided.
	ErrEmptyInstallID = errors.New("empty install id provided")

	// ErrEmptyTarget defines the error type when a
	// Installation type has an empty Target field provided.
	ErrEmptyTarget = errors.New("empty target provided")
)

// Installation is the database representation of an installation.
type Installation struct {
	InstallID sql.NullInt64  `sql:"install_id"`
	Target    sql.NullString `sql:"target"`
}

// ToAPI converts the Installation type
// to an API Installation type.
func (i *Installation) ToAPI() *api.Installation {
	installation := new(api.Installation)

	installation.SetInstallID(i.InstallID.Int64)
	installation.SetTarget(i.Target.String)

	return installation
}

// Validate verifies the necessary fields for
// the Installation type are populated correctly.
func (i *Installation) Validate() error {
	// verify the InstallID field is populated
	if i.InstallID.Int64 == 0 {
		return ErrEmptyInstallID
	}

	// verify the Target field is populated
	if i.Target.String == "" {
		return ErrEmptyTarget
	}

	return nil
}

// InstallationFromAPI converts the API Installation type
// to a database Installation type.
func InstallationFromAPI(u *api.Installation) *Installation {
	installation := &Installation{
		InstallID: sql.NullInt64{Int64: u.GetInstallID(), Valid: true},
		Target:    sql.NullString{String: u.GetTarget(), Valid: true},
	}

	return installation
}
