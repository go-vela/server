// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// Storage defines the service interface for object storage operations.
type Storage interface {
	GetAddress() string
	GetBucket() string
	GetPolicy(string) string
	// Object Operations
	StatObject(context.Context, *api.Object) (*api.Object, error)
	ListBuildObjectNames(context.Context, string, string, string) (map[string]string, error)
	// Presigned URLs
	PresignedGetObject(context.Context, *api.Object) (string, error)
	// Storage info
	StorageEnable() bool
	AssumeRole(ctx context.Context, durationSeconds int, prefix, sessionName string) (*api.STSCreds, error)
}
