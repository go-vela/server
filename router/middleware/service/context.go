// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"context"

	"github.com/go-vela/types/library"
)

const key = "service"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the Service associated with this context.
func FromContext(c context.Context) *library.Service {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	s, ok := value.(*library.Service)
	if !ok {
		return nil
	}

	return s
}

// ToContext adds the Service to this context if it supports
// the Setter interface.
func ToContext(c Setter, s *library.Service) {
	c.Set(key, s)
}
