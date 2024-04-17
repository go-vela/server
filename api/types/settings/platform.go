// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"fmt"
)

// Platform is the API representation of platform settings.
//
// swagger:model Platform
type Platform struct {
	ID *int64 `json:"id"`

	*Queue    `json:"queue"`
	*Compiler `json:"compiler"`

	// misc
	RepoAllowlist     *[]string `json:"repo_allowlist"`
	ScheduleAllowlist *[]string `json:"schedule_allowlist"`
}

// GetID returns the ID field.
//
// When the provided Settings type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Platform) GetID() int64 {
	// return zero value if Settings type or ID field is nil
	if s == nil || s.ID == nil {
		return 0
	}

	return *s.ID
}

// SetID sets the ID field.
//
// When the provided Settings type is nil, it
// will set nothing and immediately return.
func (s *Platform) SetID(v int64) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.ID = &v
}

// SetCompilerSettings sets the CompilerSettings field.
//
// When the provided CompilerSettings type is nil, it
// will set nothing and immediately return.
func (s *Platform) SetCompilerSettings(cs Compiler) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.Compiler = &cs
}

// SetQueueSettings sets the QueueSettings field.
//
// When the provided QueueSettings type is nil, it
// will set nothing and immediately return.
func (s *Platform) SetQueueSettings(qs Queue) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.Queue = &qs
}

// GetRepoAllowlist returns the RepoAllowlist field.
//
// When the provided Settings type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Platform) GetRepoAllowlist() []string {
	// return zero value if Settings type or RepoAllowlist field is nil
	if s == nil || s.RepoAllowlist == nil {
		return []string{}
	}

	return *s.RepoAllowlist
}

// SetRepoAllowlist sets the RepoAllowlist field.
//
// When the provided Settings type is nil, it
// will set nothing and immediately return.
func (s *Platform) SetRepoAllowlist(v []string) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.RepoAllowlist = &v
}

// GetScheduleAllowlist returns the ScheduleAllowlist field.
//
// When the provided Settings type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Platform) GetScheduleAllowlist() []string {
	// return zero value if Settings type or ScheduleAllowlist field is nil
	if s == nil || s.ScheduleAllowlist == nil {
		return []string{}
	}

	return *s.ScheduleAllowlist
}

// SetScheduleAllowlist sets the RepoAllowlist field.
//
// When the provided Settings type is nil, it
// will set nothing and immediately return.
func (s *Platform) SetScheduleAllowlist(v []string) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.ScheduleAllowlist = &v
}

// String implements the Stringer interface for the Settings type.
func (s *Platform) String() string {
	return fmt.Sprintf(`{
  ID: %d,
  CloneImage: %s,
  QueueRoutes: %v,
  other stuff: %v,
}`,
		s.GetID(),
		s.GetCloneImage(),
		s.GetRoutes(),
		s.GetRoutes(),
	)
}

// ToEnv converts the Settings type to a string format compatible with standard posix environments.
func (s *Platform) ToEnv() string {
	return fmt.Sprintf(`VELA_CLONE_IMAGE='%s'
VELA_QUEUE_ROUTES='%v'
`,
		s.GetCloneImage(),
		s.GetRoutes(),
	)
}

// ToYAML converts the Settings type to a YAML string.
func (s *Platform) ToYAML() string {
	return fmt.Sprintf(`VELA_CLONE_IMAGE: '%s'
VELA_QUEUE_ROUTES: '%s'
`,
		s.GetCloneImage(),
		s.GetRoutes(),
	)
}
