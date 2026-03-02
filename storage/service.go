// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"
	"time"

	api "github.com/go-vela/server/api/types"
)

// Storage defines the service interface for object storage operations.
type Storage interface {
	StatObject(context.Context, *api.Object) (*api.Object, error)
	ListBuildObjectNames(context.Context, string, string, string) (map[string]string, error)
	PresignedGetObject(context.Context, *api.Object) (string, error)
	PresignedPutObject(context.Context, string, time.Duration) (string, error)
}
