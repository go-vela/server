// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"fmt"
)

// Platform is the API representation of platform settings.
//
// swagger:model Platform
type Platform struct {
	ID                *int64 `json:"id"`
	*Queue            `json:"queue"`
	*Compiler         `json:"compiler"`
	RepoAllowlist     *[]string `json:"repo_allowlist"`
	ScheduleAllowlist *[]string `json:"schedule_allowlist"`
	CreatedAt         *int64    `json:"created_at,omitempty"`
	UpdatedAt         *int64    `json:"updated_at,omitempty"`
	UpdatedBy         *string   `json:"updated_by,omitempty"`
}

// GetID returns the ID field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Platform) GetID() int64 {
	// return zero value if Platform type or ID field is nil
	if s == nil || s.ID == nil {
		return 0
	}

	return *s.ID
}

// GetCompiler returns the Compiler field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Platform) GetCompiler() Compiler {
	// return zero value if Platform type or Compiler field is nil
	if s == nil || s.Compiler == nil {
		return Compiler{}
	}

	return *s.Compiler
}

// GetQueue returns the Queue field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Platform) GetQueue() Queue {
	// return zero value if Platform type or Queue field is nil
	if s == nil || s.Queue == nil {
		return Queue{}
	}

	return *s.Queue
}

// GetRepoAllowlist returns the RepoAllowlist field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Platform) GetRepoAllowlist() []string {
	// return zero value if Platform type or RepoAllowlist field is nil
	if s == nil || s.RepoAllowlist == nil {
		return []string{}
	}

	return *s.RepoAllowlist
}

// GetScheduleAllowlist returns the ScheduleAllowlist field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Platform) GetScheduleAllowlist() []string {
	// return zero value if Platform type or ScheduleAllowlist field is nil
	if s == nil || s.ScheduleAllowlist == nil {
		return []string{}
	}

	return *s.ScheduleAllowlist
}

// GetCreatedAt returns the CreatedAt field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Platform) GetCreatedAt() int64 {
	// return zero value if Platform type or CreatedAt field is nil
	if s == nil || s.CreatedAt == nil {
		return 0
	}

	return *s.CreatedAt
}

// GetUpdatedAt returns the UpdatedAt field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Platform) GetUpdatedAt() int64 {
	// return zero value if Platform type or UpdatedAt field is nil
	if s == nil || s.UpdatedAt == nil {
		return 0
	}

	return *s.UpdatedAt
}

// GetUpdatedBy returns the UpdatedBy field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Platform) GetUpdatedBy() string {
	// return zero value if Platform type or UpdatedBy field is nil
	if s == nil || s.UpdatedBy == nil {
		return ""
	}

	return *s.UpdatedBy
}

// SetID sets the ID field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (s *Platform) SetID(v int64) {
	// return if Platform type is nil
	if s == nil {
		return
	}

	s.ID = &v
}

// SetCompiler sets the Compiler field.
//
// When the provided Compiler type is nil, it
// will set nothing and immediately return.
func (s *Platform) SetCompiler(cs Compiler) {
	// return if Platform type is nil
	if s == nil {
		return
	}

	s.Compiler = &cs
}

// SetQueue sets the Queue field.
//
// When the provided Queue type is nil, it
// will set nothing and immediately return.
func (s *Platform) SetQueue(qs Queue) {
	// return if Platform type is nil
	if s == nil {
		return
	}

	s.Queue = &qs
}

// SetRepoAllowlist sets the RepoAllowlist field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (s *Platform) SetRepoAllowlist(v []string) {
	// return if Platform type is nil
	if s == nil {
		return
	}

	s.RepoAllowlist = &v
}

// SetScheduleAllowlist sets the RepoAllowlist field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (s *Platform) SetScheduleAllowlist(v []string) {
	// return if Platform type is nil
	if s == nil {
		return
	}

	s.ScheduleAllowlist = &v
}

// SetCreatedAt sets the CreatedAt field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (s *Platform) SetCreatedAt(v int64) {
	// return if Platform type is nil
	if s == nil {
		return
	}

	s.CreatedAt = &v
}

// SetUpdatedAt sets the UpdatedAt field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (s *Platform) SetUpdatedAt(v int64) {
	// return if Platform type is nil
	if s == nil {
		return
	}

	s.UpdatedAt = &v
}

// SetUpdatedBy sets the UpdatedBy field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (s *Platform) SetUpdatedBy(v string) {
	// return if Platform type is nil
	if s == nil {
		return
	}

	s.UpdatedBy = &v
}

// Update takes another settings record and updates the internal fields, intended
// to be used when the refreshing settings record shared across the server.
func (s *Platform) Update(s_ *Platform) {
	if s == nil {
		return
	}

	if s_ == nil {
		return
	}

	s.SetCompiler(s_.GetCompiler())
	s.SetQueue(s_.GetQueue())
	s.SetRepoAllowlist(s_.GetRepoAllowlist())
	s.SetScheduleAllowlist(s_.GetScheduleAllowlist())
}

// String implements the Stringer interface for the Platform type.
func (s *Platform) String() string {
	cs := s.GetCompiler()
	qs := s.GetQueue()

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
		s.GetID(),
		cs.String(),
		qs.String(),
		s.GetRepoAllowlist(),
		s.GetScheduleAllowlist(),
		s.GetCreatedAt(),
		s.GetUpdatedAt(),
		s.GetUpdatedBy(),
	)
}

// PlatformMockEmpty returns an empty Platform type.
func PlatformMockEmpty() Platform {
	s := Platform{}

	s.SetCompiler(CompilerMockEmpty())
	s.SetQueue(QueueMockEmpty())

	s.SetRepoAllowlist([]string{})
	s.SetScheduleAllowlist([]string{})

	return s
}
