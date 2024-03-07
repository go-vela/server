// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"context"
	"database/sql"
)

// PlatformSettings_DB is the database representation of platform settings.
type PlatformSettings_DB struct {
	ID     sql.NullInt64  `sql:"id"`
	FooNum sql.NullInt64  `sql:"foo_num"`
	BarStr sql.NullString `sql:"bar_str"`
}

// CreateSettings updates a platform settings in the database.
func (e *engine) CreateSettings(ctx context.Context, s *string) (*string, error) {
	e.logger.Tracef("creating platform settings in the database with %s", *s)

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#SettingsFromLibrary
	// s := database.SettingsFromLibrary(r)

	// todo: settings.validate()

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Settings.Validate
	// err := settings.Validate()
	// if err != nil {
	// 	return nil, err
	// }

	s2 := PlatformSettings_DB{
		ID:     sql.NullInt64{Int64: 1, Valid: true},
		FooNum: sql.NullInt64{Int64: 420, Valid: true},
		BarStr: sql.NullString{String: *s, Valid: true},
	}

	// send query to the database
	err := e.client.Table(constantsTableSettings).Create(&s2).Error
	if err != nil {
		return nil, err
	}

	return s, nil
}

// GetSettings gets platform settings from the database.
func (e *engine) GetSettings(ctx context.Context) (*string, error) {
	e.logger.Trace("getting platform settings from the database")

	// variable to store query results
	s := new(PlatformSettings_DB)

	// send query to the database and store result in variable
	err := e.client.
		Table(constantsTableSettings).
		// todo: how to ensure this is always a singleton at the first row
		Where("id = ?", 1).
		Take(s).
		Error
	if err != nil {
		return nil, err
	}

	// return the settings
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Settings.ToLibrary
	return &s.BarStr.String, nil
}

// UpdateSettings updates a platform settings in the database.
func (e *engine) UpdateSettings(ctx context.Context, s *string) (*string, error) {
	e.logger.Trace("updating platform settings in the database")

	// cast the library type to database type
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#SettingsFromLibrary
	// s := database.SettingsFromLibrary(r)

	// todo: settings.validate()

	// validate the necessary fields are populated
	//
	// https://pkg.go.dev/github.com/go-vela/types/database#Settings.Validate
	// err := settings.Validate()
	// if err != nil {
	// 	return nil, err
	// }

	s2 := PlatformSettings_DB{
		ID:     sql.NullInt64{Int64: 1, Valid: true},
		FooNum: sql.NullInt64{Int64: 421, Valid: true},
		BarStr: sql.NullString{String: *s, Valid: true},
	}

	// send query to the database
	err := e.client.Table(constantsTableSettings).Save(s2).Error
	if err != nil {
		return nil, err
	}

	return s, nil
}
