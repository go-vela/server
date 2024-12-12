package storage

import "context"

// Storage defines the service interface for object storage operations.
type Storage interface {
	// Bucket Management
	CreateBucket(ctx context.Context, bucketName string) error
	DeleteBucket(ctx context.Context, bucketName string) error
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	ListBuckets(ctx context.Context) ([]string, error)
	// Object Operations
	Upload(ctx context.Context, bucketName string, objectName string, data []byte, contentType string) error
	Download(ctx context.Context, bucketName string, objectName string) ([]byte, error)
	Delete(ctx context.Context, bucketName string, objectName string) error
	ListObjects(ctx context.Context, bucketName string, prefix string) ([]string, error)
	//// Presigned URLs
	//GeneratePresignedURL(ctx context.Context, bucket string, key string, expiry int64) (string, error)
	// Object Lifecycle
	SetBucketLifecycle(ctx context.Context, bucketName string, lifecycleConfig string) error
	GetBucketLifecycle(ctx context.Context, bucketName string) (string, error)
	//// Workflow-Specific Operations
	//UploadArtifact(ctx context.Context, workflowID, artifactName string, data []byte) error
	//DownloadArtifact(ctx context.Context, workflowID, artifactName string) ([]byte, error)
	//UploadCache(ctx context.Context, key string, data []byte) error
	//DownloadCache(ctx context.Context, key string) ([]byte, error)
	//DeleteCache(ctx context.Context, key string) error
}
