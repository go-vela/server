// SPDX-License-Identifier: Apache-2.0

package settings

type Queue struct {
	Routes *[]string `json:"routes,omitempty"`
}

// GetRoutes returns the Routes field.
//
// When the provided QueueSettings type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (qs *Queue) GetRoutes() []string {
	// return zero value if Settings type or Routes field is nil
	if qs == nil || qs.Routes == nil {
		return []string{}
	}

	return *qs.Routes
}

// SetRoutes sets the Routes field.
//
// When the provided Settings type is nil, it
// will set nothing and immediately return.
func (qs *Queue) SetRoutes(v []string) {
	// return if Settings type is nil
	if qs == nil {
		return
	}

	qs.Routes = &v
}

func QueueMockEmpty() Queue {
	qs := Queue{}
	qs.SetRoutes([]string{})

	return qs
}
