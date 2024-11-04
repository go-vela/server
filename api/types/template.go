// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
)

// Template is the API representation of a template for a pipeline.
//
// swagger:model Template
type Template struct {
	Link   *string `json:"link,omitempty"`
	Name   *string `json:"name,omitempty"`
	Source *string `json:"source,omitempty"`
	Type   *string `json:"type,omitempty"`
}

// GetLink returns the Link field.
//
// When the provided Template type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *Template) GetLink() string {
	// return zero value if Template type or Link field is nil
	if t == nil || t.Link == nil {
		return ""
	}

	return *t.Link
}

// GetName returns the Name field.
//
// When the provided Template type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *Template) GetName() string {
	// return zero value if Template type or Name field is nil
	if t == nil || t.Name == nil {
		return ""
	}

	return *t.Name
}

// GetSource returns the Source field.
//
// When the provided Template type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *Template) GetSource() string {
	// return zero value if Template type or Source field is nil
	if t == nil || t.Source == nil {
		return ""
	}

	return *t.Source
}

// GetType returns the Type field.
//
// When the provided Template type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *Template) GetType() string {
	// return zero value if Template type or Type field is nil
	if t == nil || t.Type == nil {
		return ""
	}

	return *t.Type
}

// SetLink sets the Link field.
//
// When the provided Template type is nil, it
// will set nothing and immediately return.
func (t *Template) SetLink(v string) {
	// return if Template type is nil
	if t == nil {
		return
	}

	t.Link = &v
}

// SetName sets the Name field.
//
// When the provided Template type is nil, it
// will set nothing and immediately return.
func (t *Template) SetName(v string) {
	// return if Template type is nil
	if t == nil {
		return
	}

	t.Name = &v
}

// SetSource sets the Source field.
//
// When the provided Template type is nil, it
// will set nothing and immediately return.
func (t *Template) SetSource(v string) {
	// return if Template type is nil
	if t == nil {
		return
	}

	t.Source = &v
}

// SetType sets the Type field.
//
// When the provided Template type is nil, it
// will set nothing and immediately return.
func (t *Template) SetType(v string) {
	// return if Template type is nil
	if t == nil {
		return
	}

	t.Type = &v
}

// String implements the Stringer interface for the Template type.
func (t *Template) String() string {
	return fmt.Sprintf(`{
  Link: %s,
  Name: %s,
  Source: %s,
  Type: %s,
}`,
		t.GetLink(),
		t.GetName(),
		t.GetSource(),
		t.GetType(),
	)
}
