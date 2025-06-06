// SPDX-License-Identifier: Apache-2.0

package scm

import (
	"context"

	"github.com/gin-gonic/gin"
)

// key defines the key type for storing
// the scm Service in the context.
const key = "scm"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the scm Service
// associated with this context.
func FromContext(c context.Context) Service {
	// get scm value from context
	v := c.Value(key)
	if v == nil {
		return nil
	}

	// cast scm value to expected Service type
	s, ok := v.(Service)
	if !ok {
		return nil
	}

	return s
}

// ToContext adds the scm Service to this
// context if it supports the Setter interface.
func ToContext(c Setter, s Service) {
	c.Set(key, s)
}

// WithGinContext inserts the scm Service into the gin.Context.
func WithGinContext(c *gin.Context, s Service) {
	// set the scm Service in the gin.Context
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Context.Set
	c.Set(key, s)
}
