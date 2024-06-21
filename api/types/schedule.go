// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
)

// Schedule is the API representation of a schedule for a repo.
//
// swagger:model Schedule
type Schedule struct {
	ID          *int64  `json:"id,omitempty"`
	Repo        *Repo   `json:"repo,omitempty"`
	Active      *bool   `json:"active,omitempty"`
	Name        *string `json:"name,omitempty"`
	Entry       *string `json:"entry,omitempty"`
	CreatedAt   *int64  `json:"created_at,omitempty"`
	CreatedBy   *string `json:"created_by,omitempty"`
	UpdatedAt   *int64  `json:"updated_at,omitempty"`
	UpdatedBy   *string `json:"updated_by,omitempty"`
	ScheduledAt *int64  `json:"scheduled_at,omitempty"`
	Branch      *string `json:"branch,omitempty"`
	Error       *string `json:"error,omitempty"`
	NextRun     *int64  `json:"next_run,omitempty"`
}

// GetID returns the ID field.
//
// When the provided Schedule type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Schedule) GetID() int64 {
	// return zero value if Schedule type or ID field is nil
	if s == nil || s.ID == nil {
		return 0
	}

	return *s.ID
}

// GetRepo returns the Repo field.
//
// When the provided Schedule type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Schedule) GetRepo() *Repo {
	// return zero value if Schedule type or RepoID field is nil
	if s == nil || s.Repo == nil {
		return new(Repo)
	}

	return s.Repo
}

// GetActive returns the Active field.
//
// When the provided Schedule type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Schedule) GetActive() bool {
	// return zero value if Schedule type or Active field is nil
	if s == nil || s.Active == nil {
		return false
	}

	return *s.Active
}

// GetName returns the Name field.
//
// When the provided Schedule type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Schedule) GetName() string {
	// return zero value if Schedule type or Name field is nil
	if s == nil || s.Name == nil {
		return ""
	}

	return *s.Name
}

// GetEntry returns the Entry field.
//
// When the provided Schedule type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Schedule) GetEntry() string {
	// return zero value if Schedule type or Entry field is nil
	if s == nil || s.Entry == nil {
		return ""
	}

	return *s.Entry
}

// GetCreatedAt returns the CreatedAt field.
//
// When the provided Schedule type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Schedule) GetCreatedAt() int64 {
	// return zero value if Schedule type or CreatedAt field is nil
	if s == nil || s.CreatedAt == nil {
		return 0
	}

	return *s.CreatedAt
}

// GetCreatedBy returns the CreatedBy field.
//
// When the provided Schedule type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Schedule) GetCreatedBy() string {
	// return zero value if Schedule type or CreatedBy field is nil
	if s == nil || s.CreatedBy == nil {
		return ""
	}

	return *s.CreatedBy
}

// GetUpdatedAt returns the UpdatedAt field.
//
// When the provided Schedule type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Schedule) GetUpdatedAt() int64 {
	// return zero value if Schedule type or UpdatedAt field is nil
	if s == nil || s.UpdatedAt == nil {
		return 0
	}

	return *s.UpdatedAt
}

// GetUpdatedBy returns the UpdatedBy field.
//
// When the provided Schedule type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Schedule) GetUpdatedBy() string {
	// return zero value if Schedule type or UpdatedBy field is nil
	if s == nil || s.UpdatedBy == nil {
		return ""
	}

	return *s.UpdatedBy
}

// GetScheduledAt returns the ScheduledAt field.
//
// When the provided Schedule type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Schedule) GetScheduledAt() int64 {
	// return zero value if Schedule type or ScheduledAt field is nil
	if s == nil || s.ScheduledAt == nil {
		return 0
	}

	return *s.ScheduledAt
}

// GetBranch returns the Branch field.
//
// When the provided Schedule type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Schedule) GetBranch() string {
	// return zero value if Schedule type or ScheduledAt field is nil
	if s == nil || s.Branch == nil {
		return ""
	}

	return *s.Branch
}

// GetError returns the Error field.
//
// When the provided Schedule type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Schedule) GetError() string {
	// return zero value if Schedule type or Error field is nil
	if s == nil || s.Error == nil {
		return ""
	}

	return *s.Error
}

// GetNextRun returns the NextRun field.
//
// When the provided Schedule type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Schedule) GetNextRun() int64 {
	// return zero value if Schedule type or NextRun field is nil
	if s == nil || s.NextRun == nil {
		return 0
	}

	return *s.NextRun
}

// SetID sets the ID field.
//
// When the provided Schedule type is nil, it
// will set nothing and immediately return.
func (s *Schedule) SetID(id int64) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.ID = &id
}

// SetRepo sets the Repo field.
//
// When the provided Schedule type is nil, it
// will set nothing and immediately return.
func (s *Schedule) SetRepo(v *Repo) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.Repo = v
}

// SetActive sets the Active field.
//
// When the provided Schedule type is nil, it
// will set nothing and immediately return.
func (s *Schedule) SetActive(active bool) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.Active = &active
}

// SetName sets the Name field.
//
// When the provided Schedule type is nil, it
// will set nothing and immediately return.
func (s *Schedule) SetName(name string) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.Name = &name
}

// SetEntry sets the Entry field.
//
// When the provided Schedule type is nil, it
// will set nothing and immediately return.
func (s *Schedule) SetEntry(entry string) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.Entry = &entry
}

// SetCreatedAt sets the CreatedAt field.
//
// When the provided Schedule type is nil, it
// will set nothing and immediately return.
func (s *Schedule) SetCreatedAt(createdAt int64) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.CreatedAt = &createdAt
}

// SetCreatedBy sets the CreatedBy field.
//
// When the provided Schedule type is nil, it
// will set nothing and immediately return.
func (s *Schedule) SetCreatedBy(createdBy string) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.CreatedBy = &createdBy
}

// SetUpdatedAt sets the UpdatedAt field.
//
// When the provided Schedule type is nil, it
// will set nothing and immediately return.
func (s *Schedule) SetUpdatedAt(updatedAt int64) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.UpdatedAt = &updatedAt
}

// SetUpdatedBy sets the UpdatedBy field.
//
// When the provided Schedule type is nil, it
// will set nothing and immediately return.
func (s *Schedule) SetUpdatedBy(updatedBy string) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.UpdatedBy = &updatedBy
}

// SetScheduledAt sets the ScheduledAt field.
//
// When the provided Schedule type is nil, it
// will set nothing and immediately return.
func (s *Schedule) SetScheduledAt(scheduledAt int64) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.ScheduledAt = &scheduledAt
}

// SetBranch sets the Branch field.
//
// When the provided Schedule type is nil, it
// will set nothing and immediately return.
func (s *Schedule) SetBranch(branch string) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.Branch = &branch
}

// SetError sets the Error field.
//
// When the provided Schedule type is nil, it
// will set nothing and immediately return.
func (s *Schedule) SetError(err string) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.Error = &err
}

// SetNextRun sets the NextRun field.
//
// When the provided Schedule type is nil, it
// will set nothing and immediately return.
func (s *Schedule) SetNextRun(nextRun int64) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.NextRun = &nextRun
}

// String implements the Stringer interface for the Schedule type.
func (s *Schedule) String() string {
	return fmt.Sprintf(`{
  Active: %t,
  CreatedAt: %d,
  CreatedBy: %s,
  Entry: %s,
  ID: %d,
  Name: %s,
  Repo: %v,
  ScheduledAt: %d,
  UpdatedAt: %d,
  UpdatedBy: %s,
  Branch: %s,
  Error: %s,
  NextRun: %d,
}`,
		s.GetActive(),
		s.GetCreatedAt(),
		s.GetCreatedBy(),
		s.GetEntry(),
		s.GetID(),
		s.GetName(),
		s.GetRepo(),
		s.GetScheduledAt(),
		s.GetUpdatedAt(),
		s.GetUpdatedBy(),
		s.GetBranch(),
		s.GetError(),
		s.GetNextRun(),
	)
}
