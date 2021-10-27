// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package queue

import (
	"context"

	"github.com/go-vela/types"
	"github.com/go-vela/types/pipeline"
)

// Service represents the interface for Vela integrating
// with the different supported Queue backends.
type Service interface {
	// Service Interface Functions

	// Driver defines a function that outputs
	// the configured queue driver.
	Driver() string

	// Pop defines a function that grabs an
	// item off the queue.
	Pop(context.Context) (*types.Item, error)

	// Push defines a function that publishes an
	// item to the specified route in the queue.
	Push(context.Context, string, []byte) error

	// Route defines a function that decides which
	// channel a build gets placed within the queue.
	Route(*pipeline.Worker) (string, error)
}
