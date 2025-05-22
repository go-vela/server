// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
)

// BuildExecutable is the API representation of a BuildExecutable.
//
// swagger:model BuildExecutable
type BuildExecutable struct {
	ID    *int64 `json:"id,omitempty"`
	Build *Build `json:"build,omitempty"`
	// swagger:strfmt base64
	Data *[]byte `json:"data,omitempty"`
}

// GetID returns the ID field.
//
// When the provided BuildExecutable type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *BuildExecutable) GetID() int64 {
	// return zero value if BuildExecutable type or ID field is nil
	if b == nil || b.ID == nil {
		return 0
	}

	return *b.ID
}

// GetBuild returns the Build field.
//
// When the provided BuildExecutable type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *BuildExecutable) GetBuild() *Build {
	// return zero value if BuildExecutable type or BuildID field is nil
	if b == nil || b.Build == nil {
		return nil
	}

	return b.Build
}

// GetData returns the Data field.
//
// When the provided BuildExecutable type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *BuildExecutable) GetData() []byte {
	// return zero value if BuildExecutable type or Data field is nil
	if b == nil || b.Data == nil {
		return []byte{}
	}

	return *b.Data
}

// SetID sets the ID field.
//
// When the provided BuildExecutable type is nil, it
// will set nothing and immediately return.
func (b *BuildExecutable) SetID(v int64) {
	// return if BuildExecutable type is nil
	if b == nil {
		return
	}

	b.ID = &v
}

// SetBuildID sets the BuildID field.
//
// When the provided BuildExecutable type is nil, it
// will set nothing and immediately return.
func (b *BuildExecutable) SetBuild(v *Build) {
	// return if BuildExecutable type is nil
	if b == nil {
		return
	}

	b.Build = v
}

// SetData sets the Data field.
//
// When the provided BuildExecutable type is nil, it
// will set nothing and immediately return.
func (b *BuildExecutable) SetData(v []byte) {
	// return if Log type is nil
	if b == nil {
		return
	}

	b.Data = &v
}

// String implements the Stringer interface for the BuildExecutable type.
func (b *BuildExecutable) String() string {
	return fmt.Sprintf(`{
  ID: %d,
  Build: %v,
  Data: %s,
}`,
		b.GetID(),
		b.GetBuild(),
		b.GetData(),
	)
}
