// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"errors"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyStepBuildID defines the error type when a
	// Step type has an empty BuildID field provided.
	ErrEmptyStepBuildID = errors.New("empty step build_id provided")

	// ErrEmptyStepName defines the error type when a
	// Step type has an empty Name field provided.
	ErrEmptyStepName = errors.New("empty step name provided")

	// ErrEmptyStepImage defines the error type when a
	// Step type has an empty Image field provided.
	ErrEmptyStepImage = errors.New("empty step image provided")

	// ErrEmptyStepNumber defines the error type when a
	// Step type has an empty Number field provided.
	ErrEmptyStepNumber = errors.New("empty step number provided")

	// ErrEmptyStepRepoID defines the error type when a
	// Step type has an empty RepoID field provided.
	ErrEmptyStepRepoID = errors.New("empty step repo_id provided")
)

// Step is the database representation of a step in a build.
type Step struct {
	ID           sql.NullInt64  `sql:"id"`
	BuildID      sql.NullInt64  `sql:"build_id"`
	RepoID       sql.NullInt64  `sql:"repo_id"`
	Number       sql.NullInt64  `sql:"number"`
	Name         sql.NullString `sql:"name"`
	Image        sql.NullString `sql:"image"`
	Stage        sql.NullString `sql:"stage"`
	Status       sql.NullString `sql:"status"`
	Error        sql.NullString `sql:"error"`
	ExitCode     sql.NullInt64  `sql:"exit_code"`
	Created      sql.NullInt64  `sql:"created"`
	Started      sql.NullInt64  `sql:"started"`
	Finished     sql.NullInt64  `sql:"finished"`
	Host         sql.NullString `sql:"host"`
	Runtime      sql.NullString `sql:"runtime"`
	Distribution sql.NullString `sql:"distribution"`
	ReportAs     sql.NullString `sql:"report_as"`
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the Step type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (s *Step) Nullify() *Step {
	if s == nil {
		return nil
	}

	// check if the ID field should be false
	if s.ID.Int64 == 0 {
		s.ID.Valid = false
	}

	// check if the BuildID field should be false
	if s.BuildID.Int64 == 0 {
		s.BuildID.Valid = false
	}

	// check if the RepoID field should be false
	if s.RepoID.Int64 == 0 {
		s.RepoID.Valid = false
	}

	// check if the Number field should be false
	if s.Number.Int64 == 0 {
		s.Number.Valid = false
	}

	// check if the Name field should be false
	if len(s.Name.String) == 0 {
		s.Name.Valid = false
	}

	// check if the Image field should be false
	if len(s.Image.String) == 0 {
		s.Image.Valid = false
	}

	// check if the Stage field should be false
	if len(s.Stage.String) == 0 {
		s.Stage.Valid = false
	}

	// check if the Status field should be false
	if len(s.Status.String) == 0 {
		s.Status.Valid = false
	}

	// check if the Error field should be false
	if len(s.Error.String) == 0 {
		s.Error.Valid = false
	}

	// check if the ExitCode field should be false
	if s.ExitCode.Int64 == 0 {
		s.ExitCode.Valid = false
	}

	// check if Created field should be false
	if s.Created.Int64 == 0 {
		s.Created.Valid = false
	}

	// check if Started field should be false
	if s.Started.Int64 == 0 {
		s.Started.Valid = false
	}

	// check if Finished field should be false
	if s.Finished.Int64 == 0 {
		s.Finished.Valid = false
	}

	// check if the Host field should be false
	if len(s.Host.String) == 0 {
		s.Host.Valid = false
	}

	// check if the Runtime field should be false
	if len(s.Runtime.String) == 0 {
		s.Runtime.Valid = false
	}

	// check if the Distribution field should be false
	if len(s.Distribution.String) == 0 {
		s.Distribution.Valid = false
	}

	// check if the ReportAs field should be false
	if len(s.ReportAs.String) == 0 {
		s.ReportAs.Valid = false
	}

	return s
}

// ToAPI converts the Step type
// to a API Step type.
func (s *Step) ToAPI() *api.Step {
	step := new(api.Step)

	step.SetID(s.ID.Int64)
	step.SetBuildID(s.BuildID.Int64)
	step.SetRepoID(s.RepoID.Int64)
	step.SetNumber(int(s.Number.Int64))
	step.SetName(s.Name.String)
	step.SetImage(s.Image.String)
	step.SetStage(s.Stage.String)
	step.SetStatus(s.Status.String)
	step.SetError(s.Error.String)
	step.SetExitCode(int(s.ExitCode.Int64))
	step.SetCreated(s.Created.Int64)
	step.SetStarted(s.Started.Int64)
	step.SetFinished(s.Finished.Int64)
	step.SetHost(s.Host.String)
	step.SetRuntime(s.Runtime.String)
	step.SetDistribution(s.Distribution.String)
	step.SetReportAs(s.ReportAs.String)

	return step
}

// Validate verifies the necessary fields for
// the Step type are populated correctly.
func (s *Step) Validate() error {
	// verify the BuildID field is populated
	if s.BuildID.Int64 <= 0 {
		return ErrEmptyStepBuildID
	}

	// verify the RepoID field is populated
	if s.RepoID.Int64 <= 0 {
		return ErrEmptyStepRepoID
	}

	// verify the Number field is populated
	if s.Number.Int64 <= 0 {
		return ErrEmptyStepNumber
	}

	// verify the Name field is populated
	if len(s.Name.String) == 0 {
		return ErrEmptyStepName
	}

	// verify the Image field is populated
	if len(s.Image.String) == 0 {
		return ErrEmptyStepImage
	}

	// ensure that all Step string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	s.Name = sql.NullString{String: util.Sanitize(s.Name.String), Valid: s.Name.Valid}
	s.Image = sql.NullString{String: util.Sanitize(s.Image.String), Valid: s.Image.Valid}
	s.Stage = sql.NullString{String: util.Sanitize(s.Stage.String), Valid: s.Stage.Valid}
	s.Status = sql.NullString{String: util.Sanitize(s.Status.String), Valid: s.Status.Valid}
	s.Error = sql.NullString{String: util.Sanitize(s.Error.String), Valid: s.Error.Valid}
	s.Host = sql.NullString{String: util.Sanitize(s.Host.String), Valid: s.Host.Valid}
	s.Runtime = sql.NullString{String: util.Sanitize(s.Runtime.String), Valid: s.Runtime.Valid}
	s.Distribution = sql.NullString{String: util.Sanitize(s.Distribution.String), Valid: s.Distribution.Valid}
	s.ReportAs = sql.NullString{String: util.Sanitize(s.ReportAs.String), Valid: s.ReportAs.Valid}

	return nil
}

// StepFromAPI converts the API Step type
// to a database Step type.
func StepFromAPI(s *api.Step) *Step {
	step := &Step{
		ID:           sql.NullInt64{Int64: s.GetID(), Valid: true},
		BuildID:      sql.NullInt64{Int64: s.GetBuildID(), Valid: true},
		RepoID:       sql.NullInt64{Int64: s.GetRepoID(), Valid: true},
		Number:       sql.NullInt64{Int64: int64(s.GetNumber()), Valid: true},
		Name:         sql.NullString{String: s.GetName(), Valid: true},
		Image:        sql.NullString{String: s.GetImage(), Valid: true},
		Stage:        sql.NullString{String: s.GetStage(), Valid: true},
		Status:       sql.NullString{String: s.GetStatus(), Valid: true},
		Error:        sql.NullString{String: s.GetError(), Valid: true},
		ExitCode:     sql.NullInt64{Int64: int64(s.GetExitCode()), Valid: true},
		Created:      sql.NullInt64{Int64: s.GetCreated(), Valid: true},
		Started:      sql.NullInt64{Int64: s.GetStarted(), Valid: true},
		Finished:     sql.NullInt64{Int64: s.GetFinished(), Valid: true},
		Host:         sql.NullString{String: s.GetHost(), Valid: true},
		Runtime:      sql.NullString{String: s.GetRuntime(), Valid: true},
		Distribution: sql.NullString{String: s.GetDistribution(), Valid: true},
		ReportAs:     sql.NullString{String: s.GetReportAs(), Valid: true},
	}

	return step.Nullify()
}
