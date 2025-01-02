package minio

import (
	"context"
	api "github.com/go-vela/server/api/types"
	"github.com/sirupsen/logrus"
)

// SetBucketLifecycle sets the lifecycle configuration for a bucket.
func (c *MinioClient) SetBucketLifecycle(ctx context.Context, bucket *api.Bucket) error {
	c.Logger.WithFields(logrus.Fields{
		"bucket": bucket.BucketName,
	}).Tracef("setting lifecycle configuration for bucket %s", bucket.BucketName)
	return c.client.SetBucketLifecycle(ctx, bucket.BucketName, &bucket.LifecycleConfig)
}
