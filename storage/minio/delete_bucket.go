package minio

import (
	"context"
	api "github.com/go-vela/server/api/types"
	"github.com/sirupsen/logrus"
)

// DeleteBucket deletes a bucket in MinIO.
func (c *MinioClient) DeleteBucket(ctx context.Context, bucket *api.Bucket) error {
	c.Logger.WithFields(logrus.Fields{
		"bucket": bucket.BucketName,
	}).Tracef("deleting bucketName: %s", bucket.BucketName)

	err := c.client.RemoveBucket(ctx, bucket.BucketName)
	if err != nil {
		return err
	}
	return nil
}
