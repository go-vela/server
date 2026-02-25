// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
)

// Favorite is the API representation of a user's favorite.
//
// swagger:model Favorite
type Favorite struct {
	Position *int64  `json:"position,omitempty"`
	Repo     *string `json:"repo,omitempty"`
}

// GetPosition returns the Position field.
//
// When the provided Favorite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (f *Favorite) GetPosition() int64 {
	// return zero value if Favorite type or Position field is nil
	if f == nil || f.Position == nil {
		return 0
	}

	return *f.Position
}

// GetRepo returns the Repo field.
//
// When the provided Favorite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (f *Favorite) GetRepo() string {
	// return zero value if Favorite type or Repo field is nil
	if f == nil || f.Repo == nil {
		return ""
	}

	return *f.Repo
}

// SetPosition sets the Position field.
//
// When the provided Favorite type is nil, it
// will set nothing and immediately return.
func (f *Favorite) SetPosition(v int64) {
	// return if Favorite type is nil
	if f == nil {
		return
	}

	f.Position = &v
}

// SetRepo sets the Repo field.
//
// When the provided Favorite type is nil, it
// will set nothing and immediately return.
func (f *Favorite) SetRepo(v string) {
	// return if Favorite type is nil
	if f == nil {
		return
	}

	f.Repo = &v
}

// String implements the Stringer interface for the Favorite type.
func (f *Favorite) String() string {
	return fmt.Sprintf(`{
  Position: %d,
  Repo: %s,
}`,
		f.GetPosition(),
		f.GetRepo(),
	)
}
