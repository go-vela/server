// SPDX-License-Identifier: Apache-2.0

package actions

import "github.com/go-vela/server/constants"

// Comment is the API representation of the various actions associated
// with the comment event webhook from the SCM.
type Comment struct {
	Created *bool `json:"created"`
	Edited  *bool `json:"edited"`
}

// FromMask returns the Comment type resulting from the provided integer mask.
func (a *Comment) FromMask(mask int64) *Comment {
	a.SetCreated(mask&constants.AllowCommentCreate > 0)
	a.SetEdited(mask&constants.AllowCommentEdit > 0)

	return a
}

// ToMask returns the integer mask of the values for the Comment set.
func (a *Comment) ToMask() int64 {
	mask := int64(0)

	if a.GetCreated() {
		mask = mask | constants.AllowCommentCreate
	}

	if a.GetEdited() {
		mask = mask | constants.AllowCommentEdit
	}

	return mask
}

// GetCreated returns the Created field from the provided Comment. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Comment) GetCreated() bool {
	// return zero value if Events type or Created field is nil
	if a == nil || a.Created == nil {
		return false
	}

	return *a.Created
}

// GetEdited returns the Edited field from the provided Comment. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Comment) GetEdited() bool {
	// return zero value if Events type or Edited field is nil
	if a == nil || a.Edited == nil {
		return false
	}

	return *a.Edited
}

// SetCreated sets the Comment Created field.
//
// When the provided Events type is nil, it
// will set nothing and immediately return.
func (a *Comment) SetCreated(v bool) {
	// return if Events type is nil
	if a == nil {
		return
	}

	a.Created = &v
}

// SetEdited sets the Comment Edited field.
//
// When the provided Events type is nil, it
// will set nothing and immediately return.
func (a *Comment) SetEdited(v bool) {
	// return if Events type is nil
	if a == nil {
		return
	}

	a.Edited = &v
}
