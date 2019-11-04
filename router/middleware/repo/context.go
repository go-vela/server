// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"context"

	"github.com/go-vela/types/library"
)

const key = "repo"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Repo associated with this context.
func FromContext(c context.Context) *library.Repo {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	r, ok := value.(*library.Repo)
	if !ok {
		return nil
	}

	return r
}

// ToContext adds the Repo to this context if it supports
// the Setter interface.
func ToContext(c Setter, r *library.Repo) {
	c.Set(key, r)
}
