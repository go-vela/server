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
		ID       sql.NullInt32 `sql:"id"`
		Compiler `json:"compiler" sql:"compiler"`
		Queue    `json:"queue"    sql:"queue"`

		RepoAllowlist     pq.StringArray `json:"repo_allowlist"      sql:"repo_allowlist"      gorm:"type:varchar(1000)"`
		ScheduleAllowlist pq.StringArray `json:"schedule_allowlist"  sql:"schedule_allowlist"  gorm:"type:varchar(1000)"`
		MaxDashboardRepos sql.NullInt32  `json:"max_dashboard_repos" sql:"max_dashboard_repos"`

		CreatedAt sql.NullInt64  `sql:"created_at"`
		UpdatedAt sql.NullInt64  `sql:"updated_at"`
		UpdatedBy sql.NullString `sql:"updated_by"`
	}

	// Compiler is the database representation of compiler settings.
	Compiler struct {
		CloneImage        sql.NullString `json:"clone_image"         sql:"clone_image"`
		TemplateDepth     sql.NullInt64  `json:"template_depth"      sql:"template_depth"`
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
func (ps *Platform) Nullify() *Platform {
	if ps == nil {
		return nil
	}

	// check if the ID field should be false
	if ps.ID.Int32 == 0 {
		ps.ID.Valid = false
	}

	// check if the CloneImage field should be false
	if len(ps.CloneImage.String) == 0 {
		ps.CloneImage.Valid = false
	}

	// check if the MaxDashboardRepos field should be false
	if ps.MaxDashboardRepos.Int32 == 0 {
		ps.MaxDashboardRepos.Valid = false
	}

	// check if the CreatedAt field should be false
	if ps.CreatedAt.Int64 < 0 {
		ps.CreatedAt.Valid = false
	}

	// check if the UpdatedAt field should be false
	if ps.UpdatedAt.Int64 < 0 {
		ps.UpdatedAt.Valid = false
	}

	return ps
}

// ToAPI converts the Settings type
// to an API Settings type.
func (ps *Platform) ToAPI() *settings.Platform {
	psAPI := new(settings.Platform)
	psAPI.SetID(ps.ID.Int32)

	psAPI.SetRepoAllowlist(ps.RepoAllowlist)
	psAPI.SetScheduleAllowlist(ps.ScheduleAllowlist)
	psAPI.SetMaxDashboardRepos(ps.MaxDashboardRepos.Int32)

	psAPI.Compiler = &settings.Compiler{}
	psAPI.SetCloneImage(ps.CloneImage.String)
	psAPI.SetTemplateDepth(int(ps.TemplateDepth.Int64))
	psAPI.SetStarlarkExecLimit(ps.StarlarkExecLimit.Int64)

	psAPI.Queue = &settings.Queue{}
	psAPI.SetRoutes(ps.Routes)

	psAPI.SetCreatedAt(ps.CreatedAt.Int64)
	psAPI.SetUpdatedAt(ps.UpdatedAt.Int64)
	psAPI.SetUpdatedBy(ps.UpdatedBy.String)

	return psAPI
}

// Validate verifies the necessary fields for
// the Settings type are populated correctly.
func (ps *Platform) Validate() error {
	// verify the CloneImage field is populated
	if len(ps.CloneImage.String) == 0 {
		return ErrEmptyCloneImage
	}

	// verify compiler settings are within limits
	if ps.TemplateDepth.Int64 <= 0 {
		return fmt.Errorf("template depth must be greater than zero, got: %d", ps.TemplateDepth.Int64)
	}

	if ps.StarlarkExecLimit.Int64 <= 0 {
		return fmt.Errorf("starlark exec limit must be greater than zero, got: %d", ps.StarlarkExecLimit.Int64)
	}

	// ensure that all Settings string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	ps.CloneImage = sql.NullString{String: util.Sanitize(ps.CloneImage.String), Valid: ps.CloneImage.Valid}

	// ensure that all Queue.Routes are sanitized
	// to avoid unsafe HTML content
	for i, v := range ps.Routes {
		ps.Routes[i] = util.Sanitize(v)
	}

	// ensure that all RepoAllowlist are sanitized
	// to avoid unsafe HTML content
	for i, v := range ps.RepoAllowlist {
		ps.RepoAllowlist[i] = util.Sanitize(v)
	}

	// ensure that all ScheduleAllowlist are sanitized
	// to avoid unsafe HTML content
	for i, v := range ps.ScheduleAllowlist {
		ps.ScheduleAllowlist[i] = util.Sanitize(v)
	}

	if ps.MaxDashboardRepos.Int32 <= 0 {
		return fmt.Errorf("max dashboard repos must be greater than zero, got: %d", ps.MaxDashboardRepos.Int32)
	}

	if ps.CreatedAt.Int64 < 0 {
		return fmt.Errorf("created_at must be greater than zero, got: %d", ps.CreatedAt.Int64)
	}

	if ps.UpdatedAt.Int64 < 0 {
		return fmt.Errorf("updated_at must be greater than zero, got: %d", ps.UpdatedAt.Int64)
	}

	return nil
}

// SettingsFromAPI converts the API Settings type
// to a database Settings type.
func SettingsFromAPI(s *settings.Platform) *Platform {
	settings := &Platform{
		ID: sql.NullInt32{Int32: s.GetID(), Valid: true},
		Compiler: Compiler{
			CloneImage:        sql.NullString{String: s.GetCloneImage(), Valid: true},
			TemplateDepth:     sql.NullInt64{Int64: int64(s.GetTemplateDepth()), Valid: true},
			StarlarkExecLimit: sql.NullInt64{Int64: s.GetStarlarkExecLimit(), Valid: true},
		},
		Queue: Queue{
			Routes: pq.StringArray(s.GetRoutes()),
		},
		RepoAllowlist:     pq.StringArray(s.GetRepoAllowlist()),
		ScheduleAllowlist: pq.StringArray(s.GetScheduleAllowlist()),
		MaxDashboardRepos: sql.NullInt32{Int32: s.GetMaxDashboardRepos(), Valid: true},
		CreatedAt:         sql.NullInt64{Int64: s.GetCreatedAt(), Valid: true},
		UpdatedAt:         sql.NullInt64{Int64: s.GetUpdatedAt(), Valid: true},
		UpdatedBy:         sql.NullString{String: s.GetUpdatedBy(), Valid: true},
	}

	return settings.Nullify()
}
