// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
)

// RepoPartial is an API type that holds all relevant information
// for a repository attached to a dashboard.
type RepoPartial struct {
	Org     string         `json:"org,omitempty"`
	Name    string         `json:"name,omitempty"`
	Counter int            `json:"counter,omitempty"`
	Builds  []BuildPartial `json:"builds,omitempty"`
}

// BuildPartial is an API type that holds all relevant information
// for a build attached to a RepoPartial.
type BuildPartial struct {
	Number   int    `json:"number,omitempty"`
	Started  int64  `json:"started,omitempty"`
	Finished int64  `json:"finished,omitempty"`
	Sender   string `json:"sender,omitempty"`
	Status   string `json:"status,omitempty"`
	Event    string `json:"event,omitempty"`
	Branch   string `json:"branch,omitempty"`
	Link     string `json:"link,omitempty"`
}

// DashCard is an API type that holds the dashboard information as
// well as a list of RepoPartials attached to the dashboard.
type DashCard struct {
	Dashboard *Dashboard    `json:"dashboard,omitempty"`
	Repos     []RepoPartial `json:"repos,omitempty"`
}

// Dashboard is the library representation of a dashboard.
//
// swagger:model Dashboard
type Dashboard struct {
	ID        *string          `json:"id,omitempty"`
	Name      *string          `json:"name,omitempty"`
	CreatedAt *int64           `json:"created_at,omitempty"`
	CreatedBy *string          `json:"created_by,omitempty"`
	UpdatedAt *int64           `json:"updated_at,omitempty"`
	UpdatedBy *string          `json:"updated_by,omitempty"`
	Admins    *[]string        `json:"admins,omitempty"`
	Repos     []*DashboardRepo `json:"repos,omitempty"`
}

// GetID returns the ID field.
//
// When the provided Dashboard type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Dashboard) GetID() string {
	// return zero value if Dashboard type or ID field is nil
	if d == nil || d.ID == nil {
		return ""
	}

	return *d.ID
}

// GetName returns the Name field.
//
// When the provided Dashboard type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Dashboard) GetName() string {
	// return zero value if Dashboard type or Name field is nil
	if d == nil || d.Name == nil {
		return ""
	}

	return *d.Name
}

// GetCreatedAt returns the CreatedAt field.
//
// When the provided Dashboard type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Dashboard) GetCreatedAt() int64 {
	// return zero value if Dashboard type or CreatedAt field is nil
	if d == nil || d.CreatedAt == nil {
		return 0
	}

	return *d.CreatedAt
}

// GetCreatedBy returns the CreatedBy field.
//
// When the provided Dashboard type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Dashboard) GetCreatedBy() string {
	// return zero value if Dashboard type or CreatedBy field is nil
	if d == nil || d.CreatedBy == nil {
		return ""
	}

	return *d.CreatedBy
}

// GetUpdatedAt returns the UpdatedAt field.
//
// When the provided Dashboard type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Dashboard) GetUpdatedAt() int64 {
	// return zero value if Dashboard type or UpdatedAt field is nil
	if d == nil || d.UpdatedAt == nil {
		return 0
	}

	return *d.UpdatedAt
}

// GetUpdatedBy returns the UpdatedBy field.
//
// When the provided Dashboard type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Dashboard) GetUpdatedBy() string {
	// return zero value if Dashboard type or UpdatedBy field is nil
	if d == nil || d.UpdatedBy == nil {
		return ""
	}

	return *d.UpdatedBy
}

// GetAdmins returns the Admins field.
//
// When the provided Dashboard type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Dashboard) GetAdmins() []string {
	// return zero value if Dashboard type or Admins field is nil
	if d == nil || d.Admins == nil {
		return []string{}
	}

	return *d.Admins
}

// GetRepos returns the Repos field.
//
// When the provided Dashboard type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Dashboard) GetRepos() []*DashboardRepo {
	// return zero value if Dashboard type or Repos field is nil
	if d == nil || d.Repos == nil {
		return []*DashboardRepo{}
	}

	return d.Repos
}

// SetID sets the ID field.
//
// When the provided Dashboard type is nil, it
// will set nothing and immediately return.
func (d *Dashboard) SetID(v string) {
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
func (d *Dashboard) SetName(v string) {
	// return if Dashboard type is nil
	if d == nil {
		return
	}

	d.Name = &v
}

// SetCreatedAt sets the CreatedAt field.
//
// When the provided Dashboard type is nil, it
// will set nothing and immediately return.
func (d *Dashboard) SetCreatedAt(v int64) {
	// return if Dashboard type is nil
	if d == nil {
		return
	}

	d.CreatedAt = &v
}

// SetCreatedBy sets the CreatedBy field.
//
// When the provided Dashboard type is nil, it
// will set nothing and immediately return.
func (d *Dashboard) SetCreatedBy(v string) {
	// return if Dashboard type is nil
	if d == nil {
		return
	}

	d.CreatedBy = &v
}

// SetUpdatedAt sets the UpdatedAt field.
//
// When the provided Dashboard type is nil, it
// will set nothing and immediately return.
func (d *Dashboard) SetUpdatedAt(v int64) {
	// return if Dashboard type is nil
	if d == nil {
		return
	}

	d.UpdatedAt = &v
}

// SetUpdatedBy sets the UpdatedBy field.
//
// When the provided Dashboard type is nil, it
// will set nothing and immediately return.
func (d *Dashboard) SetUpdatedBy(v string) {
	// return if Dashboard type is nil
	if d == nil {
		return
	}

	d.UpdatedBy = &v
}

// SetAdmins sets the Admins field.
//
// When the provided Dashboard type is nil, it
// will set nothing and immediately return.
func (d *Dashboard) SetAdmins(v []string) {
	// return if Dashboard type is nil
	if d == nil {
		return
	}

	d.Admins = &v
}

// SetRepos sets the Repos field.
//
// When the provided Dashboard type is nil, it
// will set nothing and immediately return.
func (d *Dashboard) SetRepos(v []*DashboardRepo) {
	// return if Dashboard type is nil
	if d == nil {
		return
	}

	d.Repos = v
}

// String implements the Stringer interface for the Dashboard type.
func (d *Dashboard) String() string {
	return fmt.Sprintf(`{
  Name: %s,
  ID: %s,
  Admins: %v,
  CreatedAt: %d,
  CreatedBy: %s,
  UpdatedAt: %d,
  UpdatedBy: %s,
  Repos: %v,
}`,
		d.GetName(),
		d.GetID(),
		d.GetAdmins(),
		d.GetCreatedAt(),
		d.GetCreatedBy(),
		d.GetUpdatedAt(),
		d.GetUpdatedBy(),
		d.GetRepos(),
	)
}
