package minio

import (
	"context"
	api "github.com/go-vela/server/api/types"
)

// GetBucketLifecycle retrieves the lifecycle configuration for a bucket.
func (c *Client) GetBucketLifecycle(ctx context.Context, bucket *api.Bucket) (*api.Bucket, error) {
	c.Logger.Tracef("getting lifecycle configuration for bucket %s", bucket.BucketName)

	var lifecycleConfig *api.Bucket
	lifeCycle, err := c.client.GetBucketLifecycle(ctx, bucket.BucketName)
	if err != nil {
		return lifecycleConfig, err
	}

	lifecycleConfig = &api.Bucket{BucketName: bucket.BucketName, LifecycleConfig: *lifeCycle}

	return lifecycleConfig, nil
}
