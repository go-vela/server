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
	ID         *int64  `json:"id,omitempty"`
	CloneImage *string `json:"clone_image,omitempty"`
}

// NewSettings returns a new Settings record.
func NewSettings(c *cli.Context) *Settings {
	s := new(Settings)

	// singleton record ID should always be 1
	s.SetID(1)

	// set the clone image to use for the injected clone step
	s.SetCloneImage(c.String("clone-image"))

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

// String implements the Stringer interface for the Settings type.
func (s *Settings) String() string {
	return fmt.Sprintf(`{
  ID: %d,
  CloneImage: %s,
}`,
		s.GetID(),
		s.GetCloneImage(),
	)
}

// ToEnv converts the Settings type to a string format compatible with standard posix environments.
func (s *Settings) ToEnv() string {
	return fmt.Sprintf(`CloneImage='%s'
FooBar='%s'
`,
		s.GetCloneImage(),
		"something-cool",
	)
}

// ToYAML converts the Settings type to a YAML string.
func (s *Settings) ToYAML() string {
	return fmt.Sprintf(`CloneImage: '%s'
FooBar: '%s'
`,
		s.GetCloneImage(),
		"something-cool",
	)
}
