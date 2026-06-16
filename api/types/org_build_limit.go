// SPDX-License-Identifier: Apache-2.0

package types

import "fmt"

// OrgBuildLimit is the API representation of an
// organization concurrent build limit.
//
// swagger:model OrgBuildLimit
type OrgBuildLimit struct {
	ID         *int64  `json:"id"`
	Org        *string `json:"org,omitempty"`
	BuildLimit *int32  `json:"build_limit,omitempty"`
	CreatedAt  *int64  `json:"created_at,omitempty"`
	UpdatedAt  *int64  `json:"updated_at,omitempty"`
	UpdatedBy  *string `json:"updated_by,omitempty"`
}

// GetID returns the ID field.
//
// When the provided OrgBuildLimit type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (o *OrgBuildLimit) GetID() int64 {
	// return zero value if OrgBuildLimit type or ID field is nil
	if o == nil || o.ID == nil {
		return 0
	}

	return *o.ID
}

// GetOrg returns the Org field.
//
// When the provided OrgBuildLimit type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (o *OrgBuildLimit) GetOrg() string {
	// return zero value if OrgBuildLimit type or Org field is nil
	if o == nil || o.Org == nil {
		return ""
	}

	return *o.Org
}

// GetBuildLimit returns the BuildLimit field.
//
// When the provided OrgBuildLimit type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (o *OrgBuildLimit) GetBuildLimit() int32 {
	// return zero value if OrgBuildLimit type or BuildLimit field is nil
	if o == nil || o.BuildLimit == nil {
		return 0
	}

	return *o.BuildLimit
}

// GetCreatedAt returns the CreatedAt field.
//
// When the provided OrgBuildLimit type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (o *OrgBuildLimit) GetCreatedAt() int64 {
	// return zero value if OrgBuildLimit type or CreatedAt field is nil
	if o == nil || o.CreatedAt == nil {
		return 0
	}

	return *o.CreatedAt
}

// GetUpdatedAt returns the UpdatedAt field.
//
// When the provided OrgBuildLimit type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (o *OrgBuildLimit) GetUpdatedAt() int64 {
	// return zero value if OrgBuildLimit type or UpdatedAt field is nil
	if o == nil || o.UpdatedAt == nil {
		return 0
	}

	return *o.UpdatedAt
}

// GetUpdatedBy returns the UpdatedBy field.
//
// When the provided OrgBuildLimit type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (o *OrgBuildLimit) GetUpdatedBy() string {
	// return zero value if OrgBuildLimit type or UpdatedBy field is nil
	if o == nil || o.UpdatedBy == nil {
		return ""
	}

	return *o.UpdatedBy
}

// SetID sets the ID field.
//
// When the provided OrgBuildLimit type is nil, it
// will set nothing and immediately return.
func (o *OrgBuildLimit) SetID(v int64) {
	// return if OrgBuildLimit type is nil
	if o == nil {
		return
	}

	o.ID = &v
}

// SetOrg sets the Org field.
//
// When the provided OrgBuildLimit type is nil, it
// will set nothing and immediately return.
func (o *OrgBuildLimit) SetOrg(v string) {
	// return if OrgBuildLimit type is nil
	if o == nil {
		return
	}

	o.Org = &v
}

// SetBuildLimit sets the BuildLimit field.
//
// When the provided OrgBuildLimit type is nil, it
// will set nothing and immediately return.
func (o *OrgBuildLimit) SetBuildLimit(v int32) {
	// return if OrgBuildLimit type is nil
	if o == nil {
		return
	}

	o.BuildLimit = &v
}

// SetCreatedAt sets the CreatedAt field.
//
// When the provided OrgBuildLimit type is nil, it
// will set nothing and immediately return.
func (o *OrgBuildLimit) SetCreatedAt(v int64) {
	// return if OrgBuildLimit type is nil
	if o == nil {
		return
	}

	o.CreatedAt = &v
}

// SetUpdatedAt sets the UpdatedAt field.
//
// When the provided OrgBuildLimit type is nil, it
// will set nothing and immediately return.
func (o *OrgBuildLimit) SetUpdatedAt(v int64) {
	// return if OrgBuildLimit type is nil
	if o == nil {
		return
	}

	o.UpdatedAt = &v
}

// SetUpdatedBy sets the UpdatedBy field.
//
// When the provided OrgBuildLimit type is nil, it
// will set nothing and immediately return.
func (o *OrgBuildLimit) SetUpdatedBy(v string) {
	// return if OrgBuildLimit type is nil
	if o == nil {
		return
	}

	o.UpdatedBy = &v
}

// String implements the Stringer interface for the OrgBuildLimit type.
func (o *OrgBuildLimit) String() string {
	return fmt.Sprintf(`{
  ID: %d,
  Org: %s,
  BuildLimit: %d,
  CreatedAt: %d,
  UpdatedAt: %d,
  UpdatedBy: %s,
}`,
		o.GetID(),
		o.GetOrg(),
		o.GetBuildLimit(),
		o.GetCreatedAt(),
		o.GetUpdatedAt(),
		o.GetUpdatedBy(),
	)
}
