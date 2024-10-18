// SPDX-License-Identifier: Apache-2.0
//
//nolint:dupl // similar code to schedule.go
package actions

import "github.com/go-vela/server/constants"

// Deploy is the API representation of the various actions associated
// with the deploy event webhook from the SCM.
type Deploy struct {
	Created *bool `json:"created"`
}

// FromMask returns the Deploy type resulting from the provided integer mask.
func (a *Deploy) FromMask(mask int64) *Deploy {
	a.SetCreated(mask&constants.AllowDeployCreate > 0)

	return a
}

// ToMask returns the integer mask of the values for the Deploy set.
func (a *Deploy) ToMask() int64 {
	mask := int64(0)

	if a.GetCreated() {
		mask = mask | constants.AllowDeployCreate
	}

	return mask
}

// GetCreated returns the Created field from the provided Deploy. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (a *Deploy) GetCreated() bool {
	// return zero value if Deploy type or Created field is nil
	if a == nil || a.Created == nil {
		return false
	}

	return *a.Created
}

// SetCreated sets the Deploy Created field.
//
// When the provided Deploy type is nil, it
// will set nothing and immediately return.
func (a *Deploy) SetCreated(v bool) {
	// return if Deploy type is nil
	if a == nil {
		return
	}

	a.Created = &v
}
