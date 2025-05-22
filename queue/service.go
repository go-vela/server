// SPDX-License-Identifier: Apache-2.0

package queue

import (
	"context"

	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/compiler/types/pipeline"

	api "github.com/go-vela/server/api/types"
)

// Service represents the interface for Vela integrating
// with the different supported Queue backends.
type Service interface {
	// Service Interface Functions

	// Driver defines a function that outputs
	// the configured queue driver.
	Driver() string

	// Length defines a function that outputs
	// the length of all queue channels
	Length(context.Context) (int64, error)

	// RouteLength defines a function that outputs
	// the length of a defined queue route
	RouteLength(context.Context, string) (int64, error)

	// Pop defines a function that grabs an
	// item off the queue.
	Pop(context.Context, []string) (int64, error)

	// Position defines a function that returns
	// the position of a build in the queue.
	Position(context.Context, *api.Build) int64

	// Push defines a function that publishes an
	// item to the specified route in the queue.
	Push(context.Context, string, int64) error

	// Ping defines a function that checks the
	// connection to the queue.
	Ping(context.Context) error

	// Route defines a function that decides which
	// channel a build gets placed within the queue.
	Route(*pipeline.Worker) (string, error)

	// GetSettings defines a function that returns
	// queue settings.
	GetSettings() settings.Queue

	// SetSettings defines a function that takes api settings
	// and updates the compiler Engine.
	SetSettings(*settings.Platform)
}
