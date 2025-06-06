// SPDX-License-Identifier: Apache-2.0

package queue

import (
	"context"

	"github.com/gin-gonic/gin"
)

// key defines the key type for storing
// the queue Service in the context.
const key = "queue"

// FromContext retrieves the queue Service from the context.Context.
func FromContext(c context.Context) Service {
	// get queue value from context.Context
	v := c.Value(key)
	if v == nil {
		return nil
	}

	// cast queue value to expected Service type
	s, ok := v.(Service)
	if !ok {
		return nil
	}

	return s
}

// FromGinContext retrieves the queue Service from the gin.Context.
func FromGinContext(c *gin.Context) Service {
	// get queue value from gin.Context
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Context.Get
	v, ok := c.Get(key)
	if !ok {
		return nil
	}

	// cast queue value to expected Service type
	s, ok := v.(Service)
	if !ok {
		return nil
	}

	return s
}

// WithGinContext inserts the queue Service into the gin.Context.
func WithGinContext(c *gin.Context, s Service) {
	// set the queue Service in the gin.Context
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Context.Set
	c.Set(key, s)
}
