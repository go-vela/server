// SPDX-License-Identifier: Apache-2.0
//
//nolint:dupl // similar code to push.go
package actions

import "github.com/go-vela/types/constants"

// Pull is the API representation of the various actions associated
// with the pull_request event webhook from the SCM.
type Pull struct {
	Opened      *bool `json:"opened"`
	Edited      *bool `json:"edited"`
	Synchronize *bool `json:"synchronize"`
	Reopened    *bool `json:"reopened"`
}

// FromMask returns the Pull type resulting from the provided integer mask.
func (a *Pull) FromMask(mask int64) *Pull {
	a.SetOpened(mask&constants.AllowPullOpen > 0)
	a.SetSynchronize(mask&constants.AllowPullSync > 0)
	a.SetEdited(mask&constants.AllowPullEdit > 0)
	a.SetReopened(mask&constants.AllowPullReopen > 0)

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
