package testreport

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

const key = "testreport"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the TestReport associated with this context.
func FromContext(c context.Context) *api.TestReport {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	tr, ok := value.(*api.TestReport)
	if !ok {
		return nil
	}

	return tr
}

// ToContext adds the TestReport to this context if it supports
// the Setter interface.
func ToContext(c Setter, tr *api.TestReport) {
	c.Set(key, tr)
}
