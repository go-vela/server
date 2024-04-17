// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
)

// Settings is the API representation of platform settings.
//
// swagger:model Settings
type Settings struct {
	ID                *int64    `json:"id,omitempty"`
	RepoAllowlist     *[]string `json:"repo_allowlist,omitempty"`
	*QueueSettings    `json:"queue,omitempty"`
	*CompilerSettings `json:"compiler,omitempty"`
}

// GetID returns the ID field.
//
// When the provided Settings type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Settings) GetID() int64 {
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
func (s *Settings) SetID(v int64) {
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
func (s *Settings) SetCompilerSettings(cs *CompilerSettings) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.CompilerSettings = cs
}

// SetQueueSettings sets the QueueSettings field.
//
// When the provided QueueSettings type is nil, it
// will set nothing and immediately return.
func (s *Settings) SetQueueSettings(qs *QueueSettings) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.QueueSettings = qs
}

// GetRepoAllowlist returns the RepoAllowlist field.
//
// When the provided Settings type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Settings) GetRepoAllowlist() []string {
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
func (s *Settings) SetRepoAllowlist(v []string) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.RepoAllowlist = &v
}

// String implements the Stringer interface for the Settings type.
func (s *Settings) String() string {
	return fmt.Sprintf(`{
  ID: %d,
  CloneImage: %s,
  QueueRoutes: %v,
}`,
		s.GetID(),
		s.GetCloneImage(),
		s.GetQueueRoutes(),
	)
}

// ToEnv converts the Settings type to a string format compatible with standard posix environments.
func (s *Settings) ToEnv() string {
	return fmt.Sprintf(`VELA_CLONE_IMAGE='%s'
VELA_QUEUE_ROUTES='%v'
`,
		s.GetCloneImage(),
		s.GetQueueRoutes(),
	)
}

// ToYAML converts the Settings type to a YAML string.
func (s *Settings) ToYAML() string {
	return fmt.Sprintf(`VELA_CLONE_IMAGE: '%s'
VELA_QUEUE_ROUTES: '%s'
`,
		s.GetCloneImage(),
		s.GetQueueRoutes(),
	)
}

type CompilerSettings struct {
	CloneImage        *string `json:"clone_image,omitempty"`
	TemplateDepth     *int    `json:"template_depth,omitempty"`
	StarlarkExecLimit *uint64 `json:"starlark_exec_limit,omitempty"`
}

// GetCloneImage returns the CloneImage field.
//
// When the provided CompilerSettings type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *CompilerSettings) GetCloneImage() string {
	// return zero value if Settings type or CloneImage field is nil
	if s == nil || s.CloneImage == nil {
		return ""
	}

	return *s.CloneImage
}

// GetTemplateDepth returns the TemplateDepth field.
//
// When the provided CompilerSettings type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *CompilerSettings) GetTemplateDepth() int {
	// return zero value if Settings type or TemplateDepth field is nil
	if s == nil || s.TemplateDepth == nil {
		return 0
	}

	return *s.TemplateDepth
}

// GetStarlarkExecLimit returns the StarlarkExecLimit field.
//
// When the provided CompilerSettings type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *CompilerSettings) GetStarlarkExecLimit() uint64 {
	// return zero value if Settings type or StarlarkExecLimit field is nil
	if s == nil || s.StarlarkExecLimit == nil {
		return 0
	}

	return *s.StarlarkExecLimit
}

// SetCloneImage sets the CloneImage field.
//
// When the provided CompilerSettings type is nil, it
// will set nothing and immediately return.
func (s *CompilerSettings) SetCloneImage(v string) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.CloneImage = &v
}

// SetTemplateDepth sets the TemplateDepth field.
//
// When the provided CompilerSettings type is nil, it
// will set nothing and immediately return.
func (s *CompilerSettings) SetTemplateDepth(v int) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.TemplateDepth = &v
}

// SetStarlarkExecLimit sets the StarlarkExecLimit field.
//
// When the provided CompilerSettings type is nil, it
// will set nothing and immediately return.
func (s *CompilerSettings) SetStarlarkExecLimit(v uint64) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.StarlarkExecLimit = &v
}

type QueueSettings struct {
	Routes *[]string `json:"routes,omitempty"`
}

// GetQueueRoutes returns the QueueRoutes field.
//
// When the provided QueueSettings type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *QueueSettings) GetQueueRoutes() []string {
	// return zero value if Settings type or QueueRoutes field is nil
	if s == nil || s.Routes == nil {
		return []string{}
	}

	return *s.Routes
}

// SetQueueRoutes sets the QueueRoutes field.
//
// When the provided Settings type is nil, it
// will set nothing and immediately return.
func (s *QueueSettings) SetQueueRoutes(v []string) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.Routes = &v
}
