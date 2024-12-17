package minio

import (
	"context"
	api "github.com/go-vela/server/api/types"
)

// SetBucketLifecycle sets the lifecycle configuration for a bucket.
func (c *MinioClient) SetBucketLifecycle(ctx context.Context, bucket *api.Bucket) error {
	c.Logger.Tracef("setting lifecycle configuration for bucket %s", bucket.BucketName)

	return c.client.SetBucketLifecycle(ctx, bucket.BucketName, &bucket.LifecycleConfig)
}
