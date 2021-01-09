// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package executors

import (
	"context"

	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

const key = "executors"

// Setter defines a context that enables setting values.
type Setter interface {
	Set(string, interface{})
}

// FromContext returns the executors associated with this context.
func FromContext(c context.Context) []library.Executor {
	value := c.Value(key)
	if value == nil {
		return nil
	}

	e, ok := value.([]library.Executor)
	if !ok {
		return nil
	}

	logrus.Infof("FromContext: %v", e)
	return e
}

// ToContext adds the executor to this context if it supports
// the Setter interface.
func ToContext(c Setter, e []library.Executor) {
	logrus.Infof("ToContext: %v", e)
	c.Set(key, e)
}
