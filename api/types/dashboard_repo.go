// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
)

// DashboardRepo is the API representation of a repo belonging to a Dashboard.
//
// swagger:model DashboardRepo
type DashboardRepo struct {
	ID       *int64    `json:"id,omitempty"`
	Name     *string   `json:"name,omitempty"`
	Branches *[]string `json:"branches,omitempty"`
	Events   *[]string `json:"events,omitempty"`
}

// GetID returns the ID field.
//
// When the provided Dashboard type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *DashboardRepo) GetID() int64 {
	// return zero value if Dashboard type or ID field is nil
	if d == nil || d.ID == nil {
		return 0
	}

	return *d.ID
}

// GetName returns the Name field.
//
// When the provided Dashboard type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *DashboardRepo) GetName() string {
	// return zero value if Dashboard type or ID field is nil
	if d == nil || d.Name == nil {
		return ""
	}

	return *d.Name
}

// GetBranches returns the Branches field.
//
// When the provided Dashboard type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *DashboardRepo) GetBranches() []string {
	// return zero value if Dashboard type or Branches field is nil
	if d == nil || d.Branches == nil {
		return []string{}
	}

	return *d.Branches
}

// GetEvents returns the Events field.
//
// When the provided Dashboard type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *DashboardRepo) GetEvents() []string {
	// return zero value if Dashboard type or Events field is nil
	if d == nil || d.Events == nil {
		return []string{}
	}

	return *d.Events
}

// SetID sets the ID field.
//
// When the provided Dashboard type is nil, it
// will set nothing and immediately return.
func (d *DashboardRepo) SetID(v int64) {
	// return if Dashboard type is nil
	if d == nil {
		return
	}

	d.ID = &v
}

// SetName sets the Name field.
//
// When the provided Dashboard type is nil, it
// will set nothing and immediately return.
func (d *DashboardRepo) SetName(v string) {
	// return if Dashboard type is nil
	if d == nil {
		return
	}

	d.Name = &v
}

// SetBranches sets the Branches field.
//
// When the provided Dashboard type is nil, it
// will set nothing and immediately return.
func (d *DashboardRepo) SetBranches(v []string) {
	// return if Dashboard type is nil
	if d == nil {
		return
	}

	d.Branches = &v
}

// SetEvents sets the Events field.
//
// When the provided Dashboard type is nil, it
// will set nothing and immediately return.
func (d *DashboardRepo) SetEvents(v []string) {
	// return if Dashboard type is nil
	if d == nil {
		return
	}

	d.Events = &v
}

// String implements the Stringer interface for the Dashboard type.
func (d *DashboardRepo) String() string {
	return fmt.Sprintf(`{
  Name: %s,
  ID: %d,
  Branches: %v,
  Events: %v,
}`,
		d.GetName(),
		d.GetID(),
		d.GetBranches(),
		d.GetEvents(),
	)
}
