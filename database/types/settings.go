// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyCloneImage defines the error type when a
	// Settings type has an empty CloneImage field provided.
	ErrEmptyCloneImage = errors.New("empty settings clone image provided")
)

type (
	// Platform is the database representation of platform settings.
	Platform struct {
		ID sql.NullInt64 `sql:"id"`
		Compiler
		Queue

		RepoAllowlist     pq.StringArray `json:"repo_allowlist" sql:"repo_allowlist" gorm:"type:varchar(1000)"`
		ScheduleAllowlist pq.StringArray `json:"schedule_allowlist" sql:"schedule_allowlist" gorm:"type:varchar(1000)"`

		CreatedAt sql.NullInt64  `sql:"created_at"`
		UpdatedAt sql.NullInt64  `sql:"updated_at"`
		UpdatedBy sql.NullString `sql:"updated_by"`
	}

	// Compiler is the database representation of compiler settings.
	Compiler struct {
		CloneImage        sql.NullString `json:"clone_image" sql:"clone_image"`
		TemplateDepth     sql.NullInt64  `json:"template_depth" sql:"template_depth"`
		StarlarkExecLimit sql.NullInt64  `json:"starlark_exec_limit" sql:"starlark_exec_limit"`
	}

	// Queue is the database representation of queue settings.
	Queue struct {
		Routes pq.StringArray `json:"routes" sql:"routes" gorm:"type:varchar(1000)"`
	}
)

// Value - Implementation of valuer for database/sql for Compiler.
func (r Compiler) Value() (driver.Value, error) {
	valueString, err := json.Marshal(r)
	return string(valueString), err
}

// Scan - Implement the database/sql scanner interface for Compiler.
func (r *Compiler) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &r)
	case string:
		return json.Unmarshal([]byte(v), &r)
	default:
		return fmt.Errorf("wrong type for compiler: %T", v)
	}
}

// Value - Implementation of valuer for database/sql for Queue.
func (r Queue) Value() (driver.Value, error) {
	valueString, err := json.Marshal(r)
	return string(valueString), err
}

// Scan - Implement the database/sql scanner interface for Queue.
func (r *Queue) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &r)
	case string:
		return json.Unmarshal([]byte(v), &r)
	default:
		return fmt.Errorf("wrong type for queue: %T", v)
	}
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the Settings type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (s *Platform) Nullify() *Platform {
	if s == nil {
		return nil
	}

	// check if the ID field should be false
	if s.ID.Int64 == 0 {
		s.ID.Valid = false
	}

	// check if the CloneImage field should be false
	if len(s.CloneImage.String) == 0 {
		s.CloneImage.Valid = false
	}

	// check if the CreatedAt field should be false
	if s.CreatedAt.Int64 < 0 {
		s.CreatedAt.Valid = false
	}

	// check if the UpdatedAt field should be false
	if s.UpdatedAt.Int64 < 0 {
		s.UpdatedAt.Valid = false
	}

	return s
}

// ToAPI converts the Settings type
// to an API Settings type.
func (s *Platform) ToAPI() *settings.Platform {
	ss := new(settings.Platform)
	ss.SetID(s.ID.Int64)

	ss.SetRepoAllowlist(s.RepoAllowlist)
	ss.SetScheduleAllowlist(s.ScheduleAllowlist)

	ss.Compiler = &settings.Compiler{}
	ss.SetCloneImage(s.CloneImage.String)
	ss.SetTemplateDepth(int(s.TemplateDepth.Int64))
	ss.SetStarlarkExecLimit(uint64(s.StarlarkExecLimit.Int64))

	ss.Queue = &settings.Queue{}
	ss.SetRoutes(s.Routes)

	ss.SetCreatedAt(s.CreatedAt.Int64)
	ss.SetUpdatedAt(s.UpdatedAt.Int64)
	ss.SetUpdatedBy(s.UpdatedBy.String)

	return ss
}

// Validate verifies the necessary fields for
// the Settings type are populated correctly.
func (s *Platform) Validate() error {
	// verify the CloneImage field is populated
	if len(s.CloneImage.String) == 0 {
		return ErrEmptyCloneImage
	}

	// verify compiler settings are within limits
	if s.TemplateDepth.Int64 <= 0 {
		return fmt.Errorf("template depth must be greater than zero, got: %d", s.TemplateDepth.Int64)
	}

	if s.StarlarkExecLimit.Int64 <= 0 {
		return fmt.Errorf("starlark exec limit must be greater than zero, got: %d", s.StarlarkExecLimit.Int64)
	}

	// ensure that all Settings string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	s.CloneImage = sql.NullString{String: util.Sanitize(s.CloneImage.String), Valid: s.CloneImage.Valid}

	// ensure that all Queue.Routes are sanitized
	// to avoid unsafe HTML content
	for i, v := range s.Routes {
		s.Routes[i] = util.Sanitize(v)
	}

	// ensure that all RepoAllowlist are sanitized
	// to avoid unsafe HTML content
	for i, v := range s.RepoAllowlist {
		s.RepoAllowlist[i] = util.Sanitize(v)
	}

	// ensure that all ScheduleAllowlist are sanitized
	// to avoid unsafe HTML content
	for i, v := range s.ScheduleAllowlist {
		s.ScheduleAllowlist[i] = util.Sanitize(v)
	}

	if s.CreatedAt.Int64 < 0 {
		return fmt.Errorf("created_at must be greater than zero, got: %d", s.CreatedAt.Int64)
	}

	if s.UpdatedAt.Int64 < 0 {
		return fmt.Errorf("updated_at must be greater than zero, got: %d", s.UpdatedAt.Int64)
	}

	return nil
}

// FromAPI converts the API Settings type
// to a database Settings type.
func FromAPI(s *settings.Platform) *Platform {
	settings := &Platform{
		ID: sql.NullInt64{Int64: s.GetID(), Valid: true},
		Compiler: Compiler{
			CloneImage:        sql.NullString{String: s.GetCloneImage(), Valid: true},
			TemplateDepth:     sql.NullInt64{Int64: int64(s.GetTemplateDepth()), Valid: true},
			StarlarkExecLimit: sql.NullInt64{Int64: int64(s.GetStarlarkExecLimit()), Valid: true},
		},
		Queue: Queue{
			Routes: pq.StringArray(s.GetRoutes()),
		},
		RepoAllowlist:     pq.StringArray(s.GetRepoAllowlist()),
		ScheduleAllowlist: pq.StringArray(s.GetScheduleAllowlist()),
		CreatedAt:         sql.NullInt64{Int64: s.GetCreatedAt(), Valid: true},
		UpdatedAt:         sql.NullInt64{Int64: s.GetUpdatedAt(), Valid: true},
		UpdatedBy:         sql.NullString{String: s.GetUpdatedBy(), Valid: true},
	}

	return settings.Nullify()
}
