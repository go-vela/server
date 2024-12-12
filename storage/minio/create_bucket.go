package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
)

// CreateBucket creates a new bucket in MinIO.
func (c *MinioClient) CreateBucket(ctx context.Context, bucketName string, opts *minio.MakeBucketOptions) error {
	c.Logger.Tracef("create new bucket: %s", bucketName)
	if opts == nil {
		opts = &minio.MakeBucketOptions{}
		c.Logger.Trace("Using US Standard Region as location default")
	}
	err := c.client.MakeBucket(ctx, bucketName, *opts)
	if err != nil {
		exists, errBucketExists := c.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			c.Logger.Tracef("Bucket %s already exists", bucketName)
			return nil
		}
		return err
	}
	return nil
}
