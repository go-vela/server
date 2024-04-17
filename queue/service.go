// SPDX-License-Identifier: Apache-2.0

package queue

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/queue/models"
	"github.com/go-vela/types/pipeline"
)

// Service represents the interface for Vela integrating
// with the different supported Queue backends.
type Service interface {
	// Service Interface Functions

	// Driver defines a function that outputs
	// the configured queue driver.
	Driver() string

	// Length defines a function that outputs
	// the length of a queue channel
	Length(context.Context) (int64, error)

	// Pop defines a function that grabs an
	// item off the queue.
	Pop(context.Context, []string) (*models.Item, error)

	// Push defines a function that publishes an
	// item to the specified route in the queue.
	Push(context.Context, string, []byte) error

	// Ping defines a function that checks the
	// connection to the queue.
	Ping(context.Context) error

	// Route defines a function that decides which
	// channel a build gets placed within the queue.
	Route(*pipeline.Worker) (string, error)

	// GetSettings defines a function that returns
	// queue settings.
	GetSettings() *api.QueueSettings

	// SetSettings defines a function that takes api settings
	// and updates the compiler Engine.
	SetSettings(*api.Settings)
}
