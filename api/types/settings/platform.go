// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"fmt"
)

// Platform is the API representation of platform settingps.
//
// swagger:model Platform
type Platform struct {
	ID                *int64 `json:"id"`
	*Compiler         `json:"compiler,omitempty" yaml:"compiler,omitempty"`
	*Queue            `json:"queue,omitempty" yaml:"queue,omitempty"`
	RepoAllowlist     *[]string `json:"repo_allowlist,omitempty" yaml:"repo_allowlist,omitempty"`
	ScheduleAllowlist *[]string `json:"schedule_allowlist,omitempty" yaml:"schedule_allowlist,omitempty"`
	CreatedAt         *int64    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt         *int64    `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	UpdatedBy         *string   `json:"updated_by,omitempty" yaml:"updated_by,omitempty"`
}

// GetID returns the ID field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetID() int64 {
	// return zero value if Platform type or ID field is nil
	if ps == nil || ps.ID == nil {
		return 0
	}

	return *ps.ID
}

// GetCompiler returns the Compiler field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetCompiler() Compiler {
	// return zero value if Platform type or Compiler field is nil
	if ps == nil || ps.Compiler == nil {
		return Compiler{}
	}

	return *ps.Compiler
}

// GetQueue returns the Queue field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetQueue() Queue {
	// return zero value if Platform type or Queue field is nil
	if ps == nil || ps.Queue == nil {
		return Queue{}
	}

	return *ps.Queue
}

// GetRepoAllowlist returns the RepoAllowlist field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetRepoAllowlist() []string {
	// return zero value if Platform type or RepoAllowlist field is nil
	if ps == nil || ps.RepoAllowlist == nil {
		return []string{}
	}

	return *ps.RepoAllowlist
}

// GetScheduleAllowlist returns the ScheduleAllowlist field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetScheduleAllowlist() []string {
	// return zero value if Platform type or ScheduleAllowlist field is nil
	if ps == nil || ps.ScheduleAllowlist == nil {
		return []string{}
	}

	return *ps.ScheduleAllowlist
}

// GetCreatedAt returns the CreatedAt field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetCreatedAt() int64 {
	// return zero value if Platform type or CreatedAt field is nil
	if ps == nil || ps.CreatedAt == nil {
		return 0
	}

	return *ps.CreatedAt
}

// GetUpdatedAt returns the UpdatedAt field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetUpdatedAt() int64 {
	// return zero value if Platform type or UpdatedAt field is nil
	if ps == nil || ps.UpdatedAt == nil {
		return 0
	}

	return *ps.UpdatedAt
}

// GetUpdatedBy returns the UpdatedBy field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetUpdatedBy() string {
	// return zero value if Platform type or UpdatedBy field is nil
	if ps == nil || ps.UpdatedBy == nil {
		return ""
	}

	return *ps.UpdatedBy
}

// SetID sets the ID field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetID(v int64) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.ID = &v
}

// SetCompiler sets the Compiler field.
//
// When the provided Compiler type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetCompiler(cs Compiler) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.Compiler = &cs
}

// SetQueue sets the Queue field.
//
// When the provided Queue type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetQueue(qs Queue) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.Queue = &qs
}

// SetRepoAllowlist sets the RepoAllowlist field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetRepoAllowlist(v []string) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.RepoAllowlist = &v
}

// SetScheduleAllowlist sets the RepoAllowlist field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetScheduleAllowlist(v []string) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.ScheduleAllowlist = &v
}

// SetCreatedAt sets the CreatedAt field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetCreatedAt(v int64) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.CreatedAt = &v
}

// SetUpdatedAt sets the UpdatedAt field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetUpdatedAt(v int64) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.UpdatedAt = &v
}

// SetUpdatedBy sets the UpdatedBy field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetUpdatedBy(v string) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.UpdatedBy = &v
}

// Update takes another settings record and updates the internal fields, intended
// to be used when the refreshing settings record shared across the server.
func (ps *Platform) Update(newSettingps *Platform) {
	if ps == nil {
		return
	}

	if newSettingps == nil {
		return
	}

	ps.SetCompiler(newSettingps.GetCompiler())
	ps.SetQueue(newSettingps.GetQueue())
	ps.SetRepoAllowlist(newSettingps.GetRepoAllowlist())
	ps.SetScheduleAllowlist(newSettingps.GetScheduleAllowlist())
}

// String implements the Stringer interface for the Platform type.
func (ps *Platform) String() string {
	cs := ps.GetCompiler()
	qs := ps.GetQueue()

	return fmt.Sprintf(`{
  ID: %d,
  Compiler: %v,
  Queue: %v,
  RepoAllowlist: %v,
  ScheduleAllowlist: %v,
  CreatedAt: %d,
  UpdatedAt: %d,
  UpdatedBy: %s,
}`,
		ps.GetID(),
		cs.String(),
		qs.String(),
		ps.GetRepoAllowlist(),
		ps.GetScheduleAllowlist(),
		ps.GetCreatedAt(),
		ps.GetUpdatedAt(),
		ps.GetUpdatedBy(),
	)
}

// PlatformMockEmpty returns an empty Platform type.
func PlatformMockEmpty() Platform {
	ps := Platform{}

	ps.SetCompiler(CompilerMockEmpty())
	ps.SetQueue(QueueMockEmpty())

	ps.SetRepoAllowlist([]string{})
	ps.SetScheduleAllowlist([]string{})

	return ps
}
