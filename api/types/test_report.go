// SPDX-License-Identifier: Apache-2.0

package types

// TestReport is the API representation of a test report for a pipeline.
//
// swagger:model TestReport
type TestReport struct {
	ID        *int64 `json:"id,omitempty"`
	BuildID   *int64 `json:"build_id,omitempty"`
	CreatedAt *int64 `json:"created_at,omitempty"`
}

// GetID returns the ID field.
//
// When the provided TestReport type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (tr *TestReport) GetID() int64 {
	// return zero value if Artifacts type or ID field is nil
	if tr == nil || tr.ID == nil {
		return 0
	}

	return *tr.ID
}

// GetBuildID returns the BuildID field.
//
// When the provided TestReport type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (tr *TestReport) GetBuildID() int64 {
	// return zero value if Artifacts type or BuildID field is nil
	if tr == nil || tr.BuildID == nil {
		return 0
	}

	return *tr.BuildID
}

// GetCreatedAt returns the CreatedAt field.
//
// When the provided TestReport type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (tr *TestReport) GetCreatedAt() int64 {
	// return zero value if Artifacts type or CreatedAt field is nil
	if tr == nil || tr.CreatedAt == nil {
		return 0
	}

	return *tr.CreatedAt
}

// SetID sets the ID field.
func (tr *TestReport) SetID(v int64) {
	// return if Artifacts type is nil
	if tr == nil {
		return
	}
	// set the ID field
	tr.ID = &v
}

// SetBuildID sets the BuildID field.
func (tr *TestReport) SetBuildID(v int64) {
	// return if Artifacts type is nil
	if tr == nil {
		return
	}
	// set the BuildID field
	tr.BuildID = &v
}

// SetCreatedAt sets the CreatedAt field.
func (tr *TestReport) SetCreatedAt(v int64) {
	// return if Artifacts type is nil
	if tr == nil {
		return
	}
	// set the CreatedAt field
	tr.CreatedAt = &v
}
