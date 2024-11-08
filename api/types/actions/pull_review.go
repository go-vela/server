// SPDX-License-Identifier: Apache-2.0

package actions

import "github.com/go-vela/server/constants"

// PullReview is the API representation of the various actions associated
// with the pull_request_review event webhook from the SCM.
type PullReview struct {
	Submitted *bool `json:"submitted"`
	Edited    *bool `json:"edited"`
	Dismissed *bool `json:"dismissed"`
}

// FromMask returns the PullReview type resulting from the provided integer mask.
func (a *PullReview) FromMask(mask int64) *PullReview {
	a.SetSubmitted(mask&constants.AllowPullReviewSubmit > 0)
	a.SetEdited(mask&constants.AllowPullReviewEdit > 0)
	a.SetDismissed(mask&constants.AllowPullReviewDismiss > 0)

	return a
}

// ToMask returns the integer mask of the values for the PullReview set.
func (a *PullReview) ToMask() int64 {
	mask := int64(0)

	if a.GetSubmitted() {
		mask = mask | constants.AllowPullReviewSubmit
	}

	if a.GetEdited() {
		mask = mask | constants.AllowPullReviewEdit
	}

	if a.GetDismissed() {
		mask = mask | constants.AllowPullReviewDismiss
	}

	return mask
}

// GetSubmitted returns the Submitted field from the provided PullReview. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *PullReview) GetSubmitted() bool {
	// return zero value if PullReview type or Submitted field is nil
	if a == nil || a.Submitted == nil {
		return false
	}

	return *a.Submitted
}

// GetEdited returns the Edited field from the provided PullReview. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *PullReview) GetEdited() bool {
	// return zero value if PullReview type or Edited field is nil
	if a == nil || a.Edited == nil {
		return false
	}

	return *a.Edited
}

// GetDismissed returns the Dismissed field from the provided PullReview. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *PullReview) GetDismissed() bool {
	// return zero value if PullReview type or Dismissed field is nil
	if a == nil || a.Dismissed == nil {
		return false
	}

	return *a.Dismissed
}

// SetSubmitted sets the PullReview Submitted field.
//
// When the provided PullReview type is nil, it
// will set nothing and immediately return.
func (a *PullReview) SetSubmitted(v bool) {
	// return if PullReview type is nil
	if a == nil {
		return
	}

	a.Submitted = &v
}

// SetEdited sets the PullReview Edited field.
//
// When the provided PullReview type is nil, it
// will set nothing and immediately return.
func (a *PullReview) SetEdited(v bool) {
	// return if PullReview type is nil
	if a == nil {
		return
	}

	a.Edited = &v
}

// SetDismissed sets the PullReview Dismissed field.
//
// When the provided PullReview type is nil, it
// will set nothing and immediately return.
func (a *PullReview) SetDismissed(v bool) {
	// return if PullReview type is nil
	if a == nil {
		return
	}

	a.Dismissed = &v
}
