package storage

import (
	"context"
	api "github.com/go-vela/server/api/types"
)

// Storage defines the service interface for object storage operations.
type Storage interface {
	// Bucket Management
	CreateBucket(ctx context.Context, bucket *api.Bucket) error
	DeleteBucket(ctx context.Context, bucket *api.Bucket) error
	BucketExists(ctx context.Context, bucket *api.Bucket) (bool, error)
	ListBuckets(ctx context.Context) ([]string, error)
	// Object Operations
	Upload(ctx context.Context, object *api.Object) error
	Download(ctx context.Context, object *api.Object) error
	Delete(ctx context.Context, object *api.Object) error
	ListObjects(ctx context.Context, bucketName string) ([]string, error)
	// Presigned URLs
	//GeneratePresignedURL(ctx context.Context, bucket string, key string, expiry int64) (string, error)
	PresignedGetObject(ctx context.Context, object *api.Object) (string, error)
	// Object Lifecycle
	SetBucketLifecycle(ctx context.Context, bucketName *api.Bucket) error
	GetBucketLifecycle(ctx context.Context, bucket *api.Bucket) (*api.Bucket, error)
	//// Workflow-Specific Operations
	//UploadArtifact(ctx context.Context, workflowID, artifactName string, data []byte) error
	//DownloadArtifact(ctx context.Context, workflowID, artifactName string) ([]byte, error)
	//UploadCache(ctx context.Context, key string, data []byte) error
	//DownloadCache(ctx context.Context, key string) ([]byte, error)
	//DeleteCache(ctx context.Context, key string) error
}
