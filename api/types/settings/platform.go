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
}`,
		s.GetID(),
		cs.String(),
		qs.String(),
		s.GetRepoAllowlist(),
		s.GetScheduleAllowlist(),
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
