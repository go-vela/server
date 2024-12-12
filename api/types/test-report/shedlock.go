package test_report

// Shedlock is the API types representation of a shedlock lock.
//
// swagger:model Shedlock
type Shedlock struct {
	Name      *string `json:"name,omitempty"`
	LockUntil *int64  `json:"lock_until,omitempty"`
	LockAt    *int64  `json:"lock_at,omitempty"`
	LockedBy  *string `json:"locked_by,omitempty"`
}

// SetName sets the Name field.
//
// When the provided Shedlock type is nil, it
// will set nothing and immediately return.
func (s *Shedlock) SetName(v string) {
	// return if Shedlock type is nil
	if s == nil {
		return
	}

	s.Name = &v
}

// GetName returns the Name field.
//
// When the provided Shedlock type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Shedlock) GetName() string {
	// return zero value if Shedlock type or Name field is nil
	if s == nil || s.Name == nil {
		return ""
	}

	return *s.Name
}

// SetLockUntil sets the LockUntil field.
//
// When the provided Shedlock type is nil, it
// will set nothing and immediately return.
func (s *Shedlock) SetLockUntil(v int64) {
	// return if Shedlock type is nil
	if s == nil {
		return
	}

	s.LockUntil = &v
}

// GetLockUntil returns the LockUntil field.
//
// When the provided Shedlock type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Shedlock) GetLockUntil() int64 {
	// return zero value if Shedlock type or LockUntil field is nil
	if s == nil || s.LockUntil == nil {
		return 0
	}

	return *s.LockUntil
}

// SetLockAt sets the LockAt field.
//
// When the provided Shedlock type is nil, it
// will set nothing and immediately return.
func (s *Shedlock) SetLockAt(v int64) {
	// return if Shedlock type is nil
	if s == nil {
		return
	}

	s.LockAt = &v
}

// GetLockAt returns the LockAt field.
//
// When the provided Shedlock type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Shedlock) GetLockAt() int64 {
	// return zero value if Shedlock type or LockAt field is nil
	if s == nil || s.LockAt == nil {
		return 0
	}

	return *s.LockAt
}

// SetLockedBy sets the LockedBy field.
//
// When the provided Shedlock type is nil, it
// will set nothing and immediately return.
func (s *Shedlock) SetLockedBy(v string) {
	// return if Shedlock type is nil
	if s == nil {
		return
	}

	s.LockedBy = &v
}

// GetLockedBy returns the LockedBy field.
//
// When the provided Shedlock type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Shedlock) GetLockedBy() string {
	// return zero value if Shedlock type or LockedBy field is nil
	if s == nil || s.LockedBy == nil {
		return ""
	}

	return *s.LockedBy
}
