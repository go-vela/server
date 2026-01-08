// SPDX-License-Identifier: Apache-2.0

package actions

import "github.com/go-vela/server/constants"

// Pull is the API representation of the various actions associated
// with the pull_request event webhook from the SCM.
type Pull struct {
	Opened      *bool `json:"opened"`
	Edited      *bool `json:"edited"`
	Synchronize *bool `json:"synchronize"`
	Reopened    *bool `json:"reopened"`
	Labeled     *bool `json:"labeled"`
	Unlabeled   *bool `json:"unlabeled"`
	Merged      *bool `json:"merged"`
	Closed      *bool `json:"closed"`
}

// FromMask returns the Pull type resulting from the provided integer mask.
func (a *Pull) FromMask(mask int64) *Pull {
	a.SetOpened(mask&constants.AllowPullOpen > 0)
	a.SetSynchronize(mask&constants.AllowPullSync > 0)
	a.SetEdited(mask&constants.AllowPullEdit > 0)
	a.SetReopened(mask&constants.AllowPullReopen > 0)
	a.SetLabeled(mask&constants.AllowPullLabel > 0)
	a.SetUnlabeled(mask&constants.AllowPullUnlabel > 0)
	a.SetMerged(mask&constants.AllowPullMerged > 0)
	a.SetClosed(mask&constants.AllowPullClosedUnmerged > 0)

	return a
}

// ToMask returns the integer mask of the values for the Pull set.
func (a *Pull) ToMask() int64 {
	mask := int64(0)

	if a.GetOpened() {
		mask = mask | constants.AllowPullOpen
	}

	if a.GetSynchronize() {
		mask = mask | constants.AllowPullSync
	}

	if a.GetEdited() {
		mask = mask | constants.AllowPullEdit
	}

	if a.GetReopened() {
		mask = mask | constants.AllowPullReopen
	}

	if a.GetLabeled() {
		mask = mask | constants.AllowPullLabel
	}

	if a.GetUnlabeled() {
		mask = mask | constants.AllowPullUnlabel
	}

	if a.GetMerged() {
		mask = mask | constants.AllowPullMerged
	}

	if a.GetClosed() {
		mask = mask | constants.AllowPullClosedUnmerged
	}

	return mask
}

// GetOpened returns the Opened field from the provided Pull. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Pull) GetOpened() bool {
	// return zero value if Pull type or Opened field is nil
	if a == nil || a.Opened == nil {
		return false
	}

	return *a.Opened
}

// GetSynchronize returns the Synchronize field from the provided Pull. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Pull) GetSynchronize() bool {
	// return zero value if Pull type or Synchronize field is nil
	if a == nil || a.Synchronize == nil {
		return false
	}

	return *a.Synchronize
}

// GetEdited returns the Edited field from the provided Pull. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Pull) GetEdited() bool {
	// return zero value if Pull type or Edited field is nil
	if a == nil || a.Edited == nil {
		return false
	}

	return *a.Edited
}

// GetReopened returns the Reopened field from the provided Pull. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Pull) GetReopened() bool {
	// return zero value if Pull type or Reopened field is nil
	if a == nil || a.Reopened == nil {
		return false
	}

	return *a.Reopened
}

// GetLabeled returns the Labeled field from the provided Pull. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Pull) GetLabeled() bool {
	// return zero value if Pull type or Labeled field is nil
	if a == nil || a.Labeled == nil {
		return false
	}

	return *a.Labeled
}

// GetUnlabeled returns the Unlabeled field from the provided Pull. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Pull) GetUnlabeled() bool {
	// return zero value if Pull type or Unlabeled field is nil
	if a == nil || a.Unlabeled == nil {
		return false
	}

	return *a.Unlabeled
}

// GetMerged returns the Merged field from the provided Pull. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Pull) GetMerged() bool {
	// return zero value if Pull type or Merged field is nil
	if a == nil || a.Merged == nil {
		return false
	}

	return *a.Merged
}

// GetClosed returns the Closed field from the provided Pull. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Pull) GetClosed() bool {
	// return zero value if Pull type or Closed field is nil
	if a == nil || a.Closed == nil {
		return false
	}

	return *a.Closed
}

// SetOpened sets the Pull Opened field.
//
// When the provided Pull type is nil, it
// will set nothing and immediately return.
func (a *Pull) SetOpened(v bool) {
	// return if Pull type is nil
	if a == nil {
		return
	}

	a.Opened = &v
}

// SetSynchronize sets the Pull Synchronize field.
//
// When the provided Pull type is nil, it
// will set nothing and immediately return.
func (a *Pull) SetSynchronize(v bool) {
	// return if Pull type is nil
	if a == nil {
		return
	}

	a.Synchronize = &v
}

// SetEdited sets the Pull Edited field.
//
// When the provided Pull type is nil, it
// will set nothing and immediately return.
func (a *Pull) SetEdited(v bool) {
	// return if Pull type is nil
	if a == nil {
		return
	}

	a.Edited = &v
}

// SetReopened sets the Pull Reopened field.
//
// When the provided Pull type is nil, it
// will set nothing and immediately return.
func (a *Pull) SetReopened(v bool) {
	// return if Pull type is nil
	if a == nil {
		return
	}

	a.Reopened = &v
}

// SetLabeled sets the Pull Labeled field.
//
// When the provided Pull type is nil, it
// will set nothing and immediately return.
func (a *Pull) SetLabeled(v bool) {
	// return if Pull type is nil
	if a == nil {
		return
	}

	a.Labeled = &v
}

// SetUnlabeled sets the Pull Unlabeled field.
//
// When the provided Pull type is nil, it
// will set nothing and immediately return.
func (a *Pull) SetUnlabeled(v bool) {
	// return if Pull type is nil
	if a == nil {
		return
	}

	a.Unlabeled = &v
}

// SetMerged sets the Pull Merged field.
//
// When the provided Pull type is nil, it
// will set nothing and immediately return.
func (a *Pull) SetMerged(v bool) {
	// return if Pull type is nil
	if a == nil {
		return
	}

	a.Merged = &v
}

// SetClosed sets the Pull Closed field.
//
// When the provided Pull type is nil, it
// will set nothing and immediately return.
func (a *Pull) SetClosed(v bool) {
	// return if Pull type is nil
	if a == nil {
		return
	}

	a.Closed = &v
}
