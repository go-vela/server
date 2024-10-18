// SPDX-License-Identifier: Apache-2.0
//
//nolint:dupl // similar code to deploy.go
package actions

import "github.com/go-vela/server/constants"

// Schedule is the API representation of the various actions associated
// with the schedule event.
type Schedule struct {
	Run *bool `json:"run"`
}

// FromMask returns the Schedule type resulting from the provided integer mask.
func (a *Schedule) FromMask(mask int64) *Schedule {
	a.SetRun(mask&constants.AllowSchedule > 0)

	return a
}

// ToMask returns the integer mask of the values for the Schedule set.
func (a *Schedule) ToMask() int64 {
	mask := int64(0)

	if a.GetRun() {
		mask = mask | constants.AllowSchedule
	}

	return mask
}

// GetRun returns the Run field from the provided Schedule. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Schedule) GetRun() bool {
	// return zero value if Schedule type or Run field is nil
	if a == nil || a.Run == nil {
		return false
	}

	return *a.Run
}

// SetRun sets the Schedule Run field.
//
// When the provided Schedule type is nil, it
// will set nothing and immediately return.
func (a *Schedule) SetRun(v bool) {
	// return if Schedule type is nil
	if a == nil {
		return
	}

	a.Run = &v
}
