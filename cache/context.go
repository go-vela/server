// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"context"

	"github.com/gin-gonic/gin"
)

// key defines the key type for storing
// the cache Service in the context.
const key = "cache"

// FromContext retrieves the cache Service from the context.Context.
func FromContext(c context.Context) Service {
	// get cache value from context.Context
	v := c.Value(key)
	if v == nil {
		return nil
	}

	// cast cache value to expected Service type
	s, ok := v.(Service)
	if !ok {
		return nil
	}

	return s
}

// FromGinContext retrieves the cache Service from the gin.Context.
func FromGinContext(c *gin.Context) Service {
	// get cache value from gin.Context
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Context.Get
	v, ok := c.Get(key)
	if !ok {
		return nil
	}

	// cast cache value to expected Service type
	s, ok := v.(Service)
	if !ok {
		return nil
	}

	return s
}

// WithGinContext inserts the cache Service into the gin.Context.
func WithGinContext(c *gin.Context, s Service) {
	// set the cache Service in the gin.Context
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Context.Set
	c.Set(key, s)
}
