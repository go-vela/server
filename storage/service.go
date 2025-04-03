// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"
	"io"

	"github.com/minio/minio-go/v7"

	api "github.com/go-vela/server/api/types"
)

// Storage defines the service interface for object storage operations.
type Storage interface {
	// Bucket Management
	CreateBucket(ctx context.Context, bucket *api.Bucket) error
	BucketExists(ctx context.Context, bucket *api.Bucket) (bool, error)
	ListBuckets(ctx context.Context) ([]string, error)
	GetBucket(ctx context.Context) string
	// Object Operations
	StatObject(ctx context.Context, object *api.Object) (*api.Object, error)
	Upload(ctx context.Context, object *api.Object) error
	UploadObject(ctx context.Context, object *api.Object, reader io.Reader, size int64) error
	//Download(ctx context.Context, object *api.Object) error
	ListObjects(ctx context.Context, bucket *api.Bucket) ([]minio.ObjectInfo, error)
	// Presigned URLs
	PresignedGetObject(ctx context.Context, object *api.Object) (string, error)
}
