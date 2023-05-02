// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"context"
	"github.com/go-vela/server/api/types"
)

const key = "schedule"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Schedule associated with this context.
func FromContext(c context.Context) *types.Schedule {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	s, ok := value.(*types.Schedule)
	if !ok {
		return nil
	}

	return s
}

// ToContext adds the Schedule to this context if it supports
// the Setter interface.
func ToContext(c Setter, s *types.Schedule) {
	c.Set(key, s)
}
