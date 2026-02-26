// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"
	"io"

	api "github.com/go-vela/server/api/types"
)

// Storage defines the service interface for object storage operations.
type Storage interface {
	GetEndpoint() string
	GetBucket() string
	GetPolicy(string) string
	StatObject(context.Context, *api.Object) (*api.Object, error)
	ListBuildObjectNames(context.Context, string, string, string) (map[string]string, error)
	PresignedGetObject(context.Context, *api.Object) (string, error)
	AssumeRole(ctx context.Context, durationSeconds int, prefix, sessionName string) (*api.STSCreds, error)
	UploadObject(ctx context.Context, object *api.Object, reader io.Reader, size int64) error
}
