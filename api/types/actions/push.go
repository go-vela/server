// SPDX-License-Identifier: Apache-2.0

package actions

import "github.com/go-vela/server/constants"

// Push is the API representation of the various actions associated
// with the push event webhook from the SCM.
type Push struct {
	Branch       *bool `json:"branch"`
	Tag          *bool `json:"tag"`
	DeleteBranch *bool `json:"delete_branch"`
	DeleteTag    *bool `json:"delete_tag"`
}

// FromMask returns the Push type resulting from the provided integer mask.
func (a *Push) FromMask(mask int64) *Push {
	a.SetBranch(mask&constants.AllowPushBranch > 0)
	a.SetTag(mask&constants.AllowPushTag > 0)
	a.SetDeleteBranch(mask&constants.AllowPushDeleteBranch > 0)
	a.SetDeleteTag(mask&constants.AllowPushDeleteTag > 0)

	return a
}

// ToMask returns the integer mask of the values for the Push set.
func (a *Push) ToMask() int64 {
	mask := int64(0)

	if a.GetBranch() {
		mask = mask | constants.AllowPushBranch
	}

	if a.GetTag() {
		mask = mask | constants.AllowPushTag
	}

	if a.GetDeleteBranch() {
		mask = mask | constants.AllowPushDeleteBranch
	}

	if a.GetDeleteTag() {
		mask = mask | constants.AllowPushDeleteTag
	}

	return mask
}

// GetBranch returns the Branch field from the provided Push. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Push) GetBranch() bool {
	// return zero value if Push type or Branch field is nil
	if a == nil || a.Branch == nil {
		return false
	}

	return *a.Branch
}

// GetTag returns the Tag field from the provided Push. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Push) GetTag() bool {
	// return zero value if Push type or Tag field is nil
	if a == nil || a.Tag == nil {
		return false
	}

	return *a.Tag
}

// GetDeleteBranch returns the DeleteBranch field from the provided Push. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Push) GetDeleteBranch() bool {
	// return zero value if Push type or DeleteBranch field is nil
	if a == nil || a.DeleteBranch == nil {
		return false
	}

	return *a.DeleteBranch
}

// GetDeleteTag returns the DeleteTag field from the provided Push. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Push) GetDeleteTag() bool {
	// return zero value if Push type or DeleteTag field is nil
	if a == nil || a.DeleteTag == nil {
		return false
	}

	return *a.DeleteTag
}

// SetBranch sets the Push Branch field.
//
// When the provided Push type is nil, it
// will set nothing and immediately return.
func (a *Push) SetBranch(v bool) {
	// return if Events type is nil
	if a == nil {
		return
	}

	a.Branch = &v
}

// SetTag sets the Push Tag field.
//
// When the provided Push type is nil, it
// will set nothing and immediately return.
func (a *Push) SetTag(v bool) {
	// return if Events type is nil
	if a == nil {
		return
	}

	a.Tag = &v
}

// SetDeleteBranch sets the Push DeleteBranch field.
//
// When the provided Push type is nil, it
// will set nothing and immediately return.
func (a *Push) SetDeleteBranch(v bool) {
	// return if Events type is nil
	if a == nil {
		return
	}

	a.DeleteBranch = &v
}

// SetDeleteTag sets the Push DeleteTag field.
//
// When the provided Push type is nil, it
// will set nothing and immediately return.
func (a *Push) SetDeleteTag(v bool) {
	// return if Events type is nil
	if a == nil {
		return
	}

	a.DeleteTag = &v
}
