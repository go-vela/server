// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"errors"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyServiceBuildID defines the error type when a
	// Service type has an empty BuildID field provided.
	ErrEmptyServiceBuildID = errors.New("empty service build_id provided")

	// ErrEmptyServiceName defines the error type when a
	// Service type has an empty Name field provided.
	ErrEmptyServiceName = errors.New("empty service name provided")

	// ErrEmptyServiceImage defines the error type when a
	// Service type has an empty Image field provided.
	ErrEmptyServiceImage = errors.New("empty service image provided")

	// ErrEmptyServiceNumber defines the error type when a
	// Service type has an empty Number field provided.
	ErrEmptyServiceNumber = errors.New("empty service number provided")

	// ErrEmptyServiceRepoID defines the error type when a
	// Service type has an empty RepoID field provided.
	ErrEmptyServiceRepoID = errors.New("empty service repo_id provided")
)

// Service is the database representation of a service in a build.
type Service struct {
	ID           sql.NullInt64  `sql:"id"`
	BuildID      sql.NullInt64  `sql:"build_id"`
	RepoID       sql.NullInt64  `sql:"repo_id"`
	Number       sql.NullInt32  `sql:"number"`
	Name         sql.NullString `sql:"name"`
	Image        sql.NullString `sql:"image"`
	Status       sql.NullString `sql:"status"`
	Error        sql.NullString `sql:"error"`
	ExitCode     sql.NullInt32  `sql:"exit_code"`
	Created      sql.NullInt64  `sql:"created"`
	Started      sql.NullInt64  `sql:"started"`
	Finished     sql.NullInt64  `sql:"finished"`
	Host         sql.NullString `sql:"host"`
	Runtime      sql.NullString `sql:"runtime"`
	Distribution sql.NullString `sql:"distribution"`
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the Service type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (s *Service) Nullify() *Service {
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
	if s.Number.Int32 == 0 {
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

	// check if the Status field should be false
	if len(s.Status.String) == 0 {
		s.Status.Valid = false
	}

	// check if the Error field should be false
	if len(s.Error.String) == 0 {
		s.Error.Valid = false
	}

	// check if the ExitCode field should be false
	if s.ExitCode.Int32 == 0 {
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

	return s
}

// ToAPI converts the Service type
// to a API Service type.
func (s *Service) ToAPI() *api.Service {
	service := new(api.Service)

	service.SetID(s.ID.Int64)
	service.SetBuildID(s.BuildID.Int64)
	service.SetRepoID(s.RepoID.Int64)
	service.SetNumber(s.Number.Int32)
	service.SetName(s.Name.String)
	service.SetImage(s.Image.String)
	service.SetStatus(s.Status.String)
	service.SetError(s.Error.String)
	service.SetExitCode(s.ExitCode.Int32)
	service.SetCreated(s.Created.Int64)
	service.SetStarted(s.Started.Int64)
	service.SetFinished(s.Finished.Int64)
	service.SetHost(s.Host.String)
	service.SetRuntime(s.Runtime.String)
	service.SetDistribution(s.Distribution.String)

	return service
}

// Validate verifies the necessary fields for
// the Service type are populated correctly.
func (s *Service) Validate() error {
	// verify the BuildID field is populated
	if s.BuildID.Int64 <= 0 {
		return ErrEmptyServiceBuildID
	}

	// verify the RepoID field is populated
	if s.RepoID.Int64 <= 0 {
		return ErrEmptyServiceRepoID
	}

	// verify the Number field is populated
	if s.Number.Int32 <= 0 {
		return ErrEmptyServiceNumber
	}

	// verify the Name field is populated
	if len(s.Name.String) == 0 {
		return ErrEmptyServiceName
	}

	// verify the Image field is populated
	if len(s.Image.String) == 0 {
		return ErrEmptyServiceImage
	}

	// ensure that all Service string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	s.Name = sql.NullString{String: util.Sanitize(s.Name.String), Valid: s.Name.Valid}
	s.Image = sql.NullString{String: util.Sanitize(s.Image.String), Valid: s.Image.Valid}
	s.Status = sql.NullString{String: util.Sanitize(s.Status.String), Valid: s.Status.Valid}
	s.Error = sql.NullString{String: util.Sanitize(s.Error.String), Valid: s.Error.Valid}
	s.Host = sql.NullString{String: util.Sanitize(s.Host.String), Valid: s.Host.Valid}
	s.Runtime = sql.NullString{String: util.Sanitize(s.Runtime.String), Valid: s.Runtime.Valid}
	s.Distribution = sql.NullString{String: util.Sanitize(s.Distribution.String), Valid: s.Distribution.Valid}

	return nil
}

// ServiceFromAPI converts the API Service type
// to a database Service type.
func ServiceFromAPI(s *api.Service) *Service {
	service := &Service{
		ID:           sql.NullInt64{Int64: s.GetID(), Valid: true},
		BuildID:      sql.NullInt64{Int64: s.GetBuildID(), Valid: true},
		RepoID:       sql.NullInt64{Int64: s.GetRepoID(), Valid: true},
		Number:       sql.NullInt32{Int32: s.GetNumber(), Valid: true},
		Name:         sql.NullString{String: s.GetName(), Valid: true},
		Image:        sql.NullString{String: s.GetImage(), Valid: true},
		Status:       sql.NullString{String: s.GetStatus(), Valid: true},
		Error:        sql.NullString{String: s.GetError(), Valid: true},
		ExitCode:     sql.NullInt32{Int32: s.GetExitCode(), Valid: true},
		Created:      sql.NullInt64{Int64: s.GetCreated(), Valid: true},
		Started:      sql.NullInt64{Int64: s.GetStarted(), Valid: true},
		Finished:     sql.NullInt64{Int64: s.GetFinished(), Valid: true},
		Host:         sql.NullString{String: s.GetHost(), Valid: true},
		Runtime:      sql.NullString{String: s.GetRuntime(), Valid: true},
		Distribution: sql.NullString{String: s.GetDistribution(), Valid: true},
	}

	return service.Nullify()
}
