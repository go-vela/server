// SPDX-License-Identifier: Apache-2.0

package settings

import "fmt"

type Queue struct {
	Routes *[]string `json:"routes,omitempty" yaml:"routes,omitempty"`
}

// GetRoutes returns the Routes field.
//
// When the provided Queue type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (qs *Queue) GetRoutes() []string {
	// return zero value if Queue type or Routes field is nil
	if qs == nil || qs.Routes == nil {
		return []string{}
	}

	return *qs.Routes
}

// SetRoutes sets the Routes field.
//
// When the provided Queue type is nil, it
// will set nothing and immediately return.
func (qs *Queue) SetRoutes(v []string) {
	// return if Queue type is nil
	if qs == nil {
		return
	}

	qs.Routes = &v
}

// String implements the Stringer interface for the Queue type.
func (qs *Queue) String() string {
	return fmt.Sprintf(`{
  Routes: %v,
}`,
		qs.GetRoutes(),
	)
}

// QueueMockEmpty returns an empty Queue type.
func QueueMockEmpty() Queue {
	qs := Queue{}
	qs.SetRoutes([]string{})

	return qs
}
