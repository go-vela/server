// SPDX-License-Identifier: Apache-2.0

package schedule

import (
	"context"
	"errors"
	"fmt"

	"database/sql"

	"github.com/adhocore/gronx"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
)

var (
	// ErrEmptyScheduleEntry defines the error type when a Schedule type has an empty Entry field provided.
	ErrEmptyScheduleEntry = errors.New("empty schedule entry provided")

	// ErrEmptyScheduleName defines the error type when a Schedule type has an empty Name field provided.
	ErrEmptyScheduleName = errors.New("empty schedule name provided")

	// ErrEmptyScheduleRepoID defines the error type when a Schedule type has an empty RepoID field provided.
	ErrEmptyScheduleRepoID = errors.New("empty schedule repo_id provided")

	// ErrInvalidScheduleEntry defines the error type when a Schedule type has an invalid Entry field provided.
	ErrInvalidScheduleEntry = errors.New("invalid schedule entry provided")
)

type (
	// config represents the settings required to create the engine that implements the ScheduleInterface interface.
	config struct {
		// specifies to skip creating tables and indexes for the Schedule engine
		SkipCreation bool
	}

	// engine represents the schedule functionality that implements the ScheduleInterface interface.
	engine struct {
		// engine configuration settings used in schedule functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in schedule functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in schedule functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}

	Schedule struct {
		ID          sql.NullInt64  `sql:"id"`
		RepoID      sql.NullInt64  `sql:"repo_id"`
		Active      sql.NullBool   `sql:"active"`
		Name        sql.NullString `sql:"name"`
		Entry       sql.NullString `sql:"entry"`
		CreatedAt   sql.NullInt64  `sql:"created_at"`
		CreatedBy   sql.NullString `sql:"created_by"`
		UpdatedAt   sql.NullInt64  `sql:"updated_at"`
		UpdatedBy   sql.NullString `sql:"updated_by"`
		ScheduledAt sql.NullInt64  `sql:"scheduled_at"`
		Branch      sql.NullString `sql:"branch"`
		Error       sql.NullString `sql:"error"`
	}
)

// New creates and returns a Vela service for integrating with schedules in the database.
//
//nolint:revive // ignore returning unexported engine
func New(opts ...EngineOpt) (*engine, error) {
	// create new Schedule engine
	e := new(engine)

	// create new fields
	e.client = new(gorm.DB)
	e.config = new(config)
	e.logger = new(logrus.Entry)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(e)
		if err != nil {
			return nil, err
		}
	}

	// check if we should skip creating schedule database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of schedules table and indexes in the database")

		return e, nil
	}

	// create the schedules table
	err := e.CreateScheduleTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableSchedule, err)
	}

	// create the indexes for the schedules table
	err = e.CreateScheduleIndexes(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableSchedule, err)
	}

	return e, nil
}

// ScheduleFromLibrary converts the api.Schedule type to a database Schedule type.
func FromAPI(s *api.Schedule) *Schedule {
	schedule := &Schedule{
		ID:          sql.NullInt64{Int64: s.GetID(), Valid: true},
		RepoID:      sql.NullInt64{Int64: s.GetRepoID(), Valid: true},
		Active:      sql.NullBool{Bool: s.GetActive(), Valid: true},
		Name:        sql.NullString{String: s.GetName(), Valid: true},
		Entry:       sql.NullString{String: s.GetEntry(), Valid: true},
		CreatedAt:   sql.NullInt64{Int64: s.GetCreatedAt(), Valid: true},
		CreatedBy:   sql.NullString{String: s.GetCreatedBy(), Valid: true},
		UpdatedAt:   sql.NullInt64{Int64: s.GetUpdatedAt(), Valid: true},
		UpdatedBy:   sql.NullString{String: s.GetUpdatedBy(), Valid: true},
		ScheduledAt: sql.NullInt64{Int64: s.GetScheduledAt(), Valid: true},
		Branch:      sql.NullString{String: s.GetBranch(), Valid: true},
		Error:       sql.NullString{String: s.GetError(), Valid: true},
	}

	return schedule.Nullify()
}

// Nullify ensures the valid flag for the sql.Null types are properly set.
//
// When a field within the Schedule type is the zero value for the field, the
// valid flag is set to false causing it to be NULL in the database.
func (s *Schedule) Nullify() *Schedule {
	if s == nil {
		return nil
	}

	// check if the ID field should be valid
	s.ID.Valid = s.ID.Int64 != 0
	// check if the RepoID field should be valid
	s.RepoID.Valid = s.RepoID.Int64 != 0
	// check if the ID field should be valid
	s.Active.Valid = s.RepoID.Int64 != 0
	// check if the Name field should be valid
	s.Name.Valid = len(s.Name.String) != 0
	// check if the Entry field should be valid
	s.Entry.Valid = len(s.Entry.String) != 0
	// check if the CreatedAt field should be valid
	s.CreatedAt.Valid = s.CreatedAt.Int64 != 0
	// check if the CreatedBy field should be valid
	s.CreatedBy.Valid = len(s.CreatedBy.String) != 0
	// check if the UpdatedAt field should be valid
	s.UpdatedAt.Valid = s.UpdatedAt.Int64 != 0
	// check if the UpdatedBy field should be valid
	s.UpdatedBy.Valid = len(s.UpdatedBy.String) != 0
	// check if the ScheduledAt field should be valid
	s.ScheduledAt.Valid = s.ScheduledAt.Int64 != 0
	// check if the Branch field should be valid
	s.Branch.Valid = len(s.Branch.String) != 0
	// check if the Error field should be valid
	s.Error.Valid = len(s.Error.String) != 0

	return s
}

// ToLibrary converts the Schedule type to a api.Schedule type.
func (s *Schedule) ToAPI() *api.Schedule {
	return &api.Schedule{
		ID:          &s.ID.Int64,
		RepoID:      &s.RepoID.Int64,
		Active:      &s.Active.Bool,
		Name:        &s.Name.String,
		Entry:       &s.Entry.String,
		CreatedAt:   &s.CreatedAt.Int64,
		CreatedBy:   &s.CreatedBy.String,
		UpdatedAt:   &s.UpdatedAt.Int64,
		UpdatedBy:   &s.UpdatedBy.String,
		ScheduledAt: &s.ScheduledAt.Int64,
		Branch:      &s.Branch.String,
		Error:       &s.Error.String,
	}
}

// Validate verifies the necessary fields for the Schedule type are populated correctly.
func (s *Schedule) Validate() error {
	// verify the RepoID field is populated
	if s.RepoID.Int64 <= 0 {
		return ErrEmptyScheduleRepoID
	}

	// verify the Name field is populated
	if len(s.Name.String) <= 0 {
		return ErrEmptyScheduleName
	}

	// verify the Entry field is populated
	if len(s.Entry.String) <= 0 {
		return ErrEmptyScheduleEntry
	}

	gron := gronx.New()
	if !gron.IsValid(s.Entry.String) {
		return ErrInvalidScheduleEntry
	}

	// ensure that all Schedule string fields that can be returned as JSON are sanitized to avoid unsafe HTML content
	s.Name = sql.NullString{String: util.Sanitize(s.Name.String), Valid: s.Name.Valid}
	s.Entry = sql.NullString{String: util.Sanitize(s.Entry.String), Valid: s.Entry.Valid}

	return nil
}
