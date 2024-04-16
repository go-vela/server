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
	RepoID      *int64  `json:"repo_id,omitempty"`
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
}

// GetID returns the ID field from the provided Schedule. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (s *Schedule) GetID() int64 {
	// return zero value if Schedule type or ID field is nil
	if s == nil || s.ID == nil {
		return 0
	}

	return *s.ID
}

// GetRepoID returns the RepoID field from the provided Schedule. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (s *Schedule) GetRepoID() int64 {
	// return zero value if Schedule type or RepoID field is nil
	if s == nil || s.RepoID == nil {
		return 0
	}

	return *s.RepoID
}

// GetActive returns the Active field from the provided Schedule. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (s *Schedule) GetActive() bool {
	// return zero value if Schedule type or Active field is nil
	if s == nil || s.Active == nil {
		return false
	}

	return *s.Active
}

// GetName returns the Name field from the provided Schedule. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (s *Schedule) GetName() string {
	// return zero value if Schedule type or Name field is nil
	if s == nil || s.Name == nil {
		return ""
	}

	return *s.Name
}

// GetEntry returns the Entry field from the provided Schedule. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (s *Schedule) GetEntry() string {
	// return zero value if Schedule type or Entry field is nil
	if s == nil || s.Entry == nil {
		return ""
	}

	return *s.Entry
}

// GetCreatedAt returns the CreatedAt field from the provided Schedule. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (s *Schedule) GetCreatedAt() int64 {
	// return zero value if Schedule type or CreatedAt field is nil
	if s == nil || s.CreatedAt == nil {
		return 0
	}

	return *s.CreatedAt
}

// GetCreatedBy returns the CreatedBy field from the provided Schedule. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (s *Schedule) GetCreatedBy() string {
	// return zero value if Schedule type or CreatedBy field is nil
	if s == nil || s.CreatedBy == nil {
		return ""
	}

	return *s.CreatedBy
}

// GetUpdatedAt returns the UpdatedAt field from the provided Schedule. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (s *Schedule) GetUpdatedAt() int64 {
	// return zero value if Schedule type or UpdatedAt field is nil
	if s == nil || s.UpdatedAt == nil {
		return 0
	}

	return *s.UpdatedAt
}

// GetUpdatedBy returns the UpdatedBy field from the provided Schedule. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (s *Schedule) GetUpdatedBy() string {
	// return zero value if Schedule type or UpdatedBy field is nil
	if s == nil || s.UpdatedBy == nil {
		return ""
	}

	return *s.UpdatedBy
}

// GetScheduledAt returns the ScheduledAt field from the provided Schedule. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (s *Schedule) GetScheduledAt() int64 {
	// return zero value if Schedule type or ScheduledAt field is nil
	if s == nil || s.ScheduledAt == nil {
		return 0
	}

	return *s.ScheduledAt
}

// GetBranch returns the Branch field from the provided Schedule. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (s *Schedule) GetBranch() string {
	// return zero value if Schedule type or ScheduledAt field is nil
	if s == nil || s.Branch == nil {
		return ""
	}

	return *s.Branch
}

// GetError returns the Error field from the provided Schedule. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (s *Schedule) GetError() string {
	// return zero value if Schedule type or Error field is nil
	if s == nil || s.Error == nil {
		return ""
	}

	return *s.Error
}

// SetID sets the ID field in the provided Schedule. If the object is nil,
// it will set nothing and immediately return making this a no-op.
func (s *Schedule) SetID(id int64) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.ID = &id
}

// SetRepoID sets the RepoID field in the provided Schedule. If the object is nil,
// it will set nothing and immediately return making this a no-op.
func (s *Schedule) SetRepoID(repoID int64) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.RepoID = &repoID
}

// SetActive sets the Active field in the provided Schedule. If the object is nil,
// it will set nothing and immediately return making this a no-op.
func (s *Schedule) SetActive(active bool) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.Active = &active
}

// SetName sets the Name field in the provided Schedule. If the object is nil,
// it will set nothing and immediately return making this a no-op.
func (s *Schedule) SetName(name string) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.Name = &name
}

// SetEntry sets the Entry field in the provided Schedule. If the object is nil,
// it will set nothing and immediately return making this a no-op.
func (s *Schedule) SetEntry(entry string) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.Entry = &entry
}

// SetCreatedAt sets the CreatedAt field in the provided Schedule. If the object is nil,
// it will set nothing and immediately return making this a no-op.
func (s *Schedule) SetCreatedAt(createdAt int64) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.CreatedAt = &createdAt
}

// SetCreatedBy sets the CreatedBy field in the provided Schedule. If the object is nil,
// it will set nothing and immediately return making this a no-op.
func (s *Schedule) SetCreatedBy(createdBy string) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.CreatedBy = &createdBy
}

// SetUpdatedAt sets the UpdatedAt field in the provided Schedule. If the object is nil,
// it will set nothing and immediately return making this a no-op.
func (s *Schedule) SetUpdatedAt(updatedAt int64) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.UpdatedAt = &updatedAt
}

// SetUpdatedBy sets the UpdatedBy field in the provided Schedule. If the object is nil,
// it will set nothing and immediately return making this a no-op.
func (s *Schedule) SetUpdatedBy(updatedBy string) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.UpdatedBy = &updatedBy
}

// SetScheduledAt sets the ScheduledAt field in the provided Schedule. If the object is nil,
// it will set nothing and immediately return making this a no-op.
func (s *Schedule) SetScheduledAt(scheduledAt int64) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.ScheduledAt = &scheduledAt
}

// SetBranch sets the Branch field in the provided Schedule. If the object is nil,
// it will set nothing and immediately return making this a no-op.
func (s *Schedule) SetBranch(branch string) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.Branch = &branch
}

// SetError sets the Error field in the provided Schedule. If the object is nil,
// it will set nothing and immediately return making this a no-op.
func (s *Schedule) SetError(err string) {
	// return if Schedule type is nil
	if s == nil {
		return
	}

	s.Error = &err
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
  RepoID: %d,
  ScheduledAt: %d,
  UpdatedAt: %d,
  UpdatedBy: %s,
  Branch: %s,
  Error: %s,
}`,
		s.GetActive(),
		s.GetCreatedAt(),
		s.GetCreatedBy(),
		s.GetEntry(),
		s.GetID(),
		s.GetName(),
		s.GetRepoID(),
		s.GetScheduledAt(),
		s.GetUpdatedAt(),
		s.GetUpdatedBy(),
		s.GetBranch(),
		s.GetError(),
	)
}
