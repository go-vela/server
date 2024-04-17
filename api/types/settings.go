// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

// Settings is the API representation of platform settings.
//
// swagger:model Settings
type Settings struct {
	ID                *int64    `json:"id,omitempty"`
	CloneImage        *string   `json:"clone_image,omitempty"`
	TemplateDepth     *int      `json:"template_depth,omitempty"`
	StarlarkExecLimit *uint64   `json:"starlark_exec_limit,omitempty"`
	QueueRoutes       *[]string `json:"queue_routes,omitempty"`
}

// NewSettings returns a new Settings record.
func NewSettings(c *cli.Context) *Settings {
	s := new(Settings)

	// singleton record ID should always be 1
	s.SetID(1)

	// set the clone image to use for the injected clone step
	s.SetCloneImage(c.String("clone-image"))

	// set the queue routes (channels) to use for builds
	s.SetQueueRoutes(c.StringSlice("queue.routes"))

	// set the queue routes (channels) to use for builds
	s.SetQueueRoutes(c.StringSlice("queue.routes"))

	return s
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

// GetCloneImage returns the CloneImage field.
//
// When the provided Settings type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Settings) GetCloneImage() string {
	// return zero value if Settings type or CloneImage field is nil
	if s == nil || s.CloneImage == nil {
		return ""
	}

	return *s.CloneImage
}

// GetTemplateDepth returns the TemplateDepth field.
//
// When the provided Settings type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Settings) GetTemplateDepth() int {
	// return zero value if Settings type or TemplateDepth field is nil
	if s == nil || s.TemplateDepth == nil {
		return 0
	}

	return *s.TemplateDepth
}

// GetStarlarkExecLimit returns the StarlarkExecLimit field.
//
// When the provided Settings type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Settings) GetStarlarkExecLimit() uint64 {
	// return zero value if Settings type or StarlarkExecLimit field is nil
	if s == nil || s.StarlarkExecLimit == nil {
		return 0
	}

	return *s.StarlarkExecLimit
}

// GetQueueRoutes returns the QueueRoutes field.
//
// When the provided Settings type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Settings) GetQueueRoutes() []string {
	// return zero value if Settings type or QueueRoutes field is nil
	if s == nil || s.QueueRoutes == nil {
		return []string{}
	}

	return *s.QueueRoutes
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

// SetCloneImage sets the CloneImage field.
//
// When the provided Settings type is nil, it
// will set nothing and immediately return.
func (s *Settings) SetCloneImage(v string) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.CloneImage = &v
}

// SetTemplateDepth sets the TemplateDepth field.
//
// When the provided Settings type is nil, it
// will set nothing and immediately return.
func (s *Settings) SetTemplateDepth(v int) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.TemplateDepth = &v
}

// SetStarlarkExecLimit sets the StarlarkExecLimit field.
//
// When the provided Settings type is nil, it
// will set nothing and immediately return.
func (s *Settings) SetStarlarkExecLimit(v uint64) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.StarlarkExecLimit = &v
}

// SetQueueRoutes sets the QueueRoutes field.
//
// When the provided Settings type is nil, it
// will set nothing and immediately return.
func (s *Settings) SetQueueRoutes(v []string) {
	// return if Settings type is nil
	if s == nil {
		return
	}

	s.QueueRoutes = &v
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
