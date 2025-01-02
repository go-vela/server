package minio

import (
	"context"
	api "github.com/go-vela/server/api/types"
	"github.com/sirupsen/logrus"
)

// GetBucketLifecycle retrieves the lifecycle configuration for a bucket.
func (c *MinioClient) GetBucketLifecycle(ctx context.Context, bucket *api.Bucket) (*api.Bucket, error) {
	c.Logger.WithFields(logrus.Fields{
		"bucket": bucket.BucketName,
	}).Tracef("getting lifecycle configuration for bucket %s", bucket.BucketName)

	var lifecycleConfig *api.Bucket
	lifeCycle, err := c.client.GetBucketLifecycle(ctx, bucket.BucketName)
	if err != nil {
		return lifecycleConfig, err
	}

	lifecycleConfig = &api.Bucket{BucketName: bucket.BucketName, LifecycleConfig: *lifeCycle}

	return lifecycleConfig, nil
}
