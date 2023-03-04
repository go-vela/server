// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package claims

import (
	"context"

	"github.com/go-vela/server/internal/token"
)

const key = "claims"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Claims associated with this context.
func FromContext(c context.Context) *token.Claims {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	cl, ok := value.(*token.Claims)
	if !ok {
		return nil
	}

	return cl
}

// ToContext adds the Claims to this context if it supports
// the Setter interface.
func ToContext(c Setter, cl *token.Claims) {
	c.Set(key, cl)
}
