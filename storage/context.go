// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"

	"github.com/gin-gonic/gin"
)

// key is the key used to store minio service in context.
const key = "minio"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext retrieves minio service from the context.
func FromContext(ctx context.Context) Storage {
	// get minio value from context.Context
	v := ctx.Value(key)
	if v == nil {
		return nil
	}

	// cast minio value to expected Storage type
	s, ok := v.(Storage)
	if !ok {
		return nil
	}

	return s
}

// FromGinContext retrieves the S3 Service from the gin.Context.
func FromGinContext(c *gin.Context) Storage {
	// get minio value from gin.Context
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Context.Get
	v, ok := c.Get(key)
	if !ok {
		return nil
	}

	// cast minio value to expected Service type
	s, ok := v.(Storage)
	if !ok {
		return nil
	}

	return s
}

// ToContext adds the secret Service to this
// context if it supports the Setter interface.
func ToContext(c Setter, s Storage) {
	c.Set(key, s)
}

// WithContext adds the minio Storage to the context.
func WithContext(ctx context.Context, s Storage) context.Context {
	// set the storage Service in the context.Context
	//
	// https://pkg.go.dev/context?tab=doc#WithValue
	//
	return context.WithValue(ctx, key, s)
}

// WithGinContext inserts the minio Storage into the gin.Context.
func WithGinContext(c *gin.Context, s Storage) {
	// set the minio Storage in the gin.Context
	//
	// https://pkg.go.dev/github.com/gin-gonic/gin?tab=doc#Context.Set
	c.Set(key, s)
}
