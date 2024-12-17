package minio

import (
	"context"
	api "github.com/go-vela/server/api/types"
	"github.com/minio/minio-go/v7"
)

// CreateBucket creates a new bucket in MinIO.
func (c *MinioClient) CreateBucket(ctx context.Context, bucket *api.Bucket) error {
	c.Logger.Tracef("create new bucket: %s", bucket.BucketName)
	var opts minio.MakeBucketOptions
	if &bucket.Options == nil {
		c.Logger.Trace("Using US Standard Region as location default")
		opts = minio.MakeBucketOptions{}
	} else {
		opts = minio.MakeBucketOptions{
			Region:        bucket.Options.Region,
			ObjectLocking: bucket.Options.ObjectLocking,
		}
	}
	err := c.client.MakeBucket(ctx, bucket.BucketName, opts)
	if err != nil {
		exists, errBucketExists := c.BucketExists(ctx, bucket)
		if errBucketExists == nil && exists {
			c.Logger.Tracef("Bucket %s already exists", bucket.BucketName)
			return nil
		}
		return err
	}
	return nil
}
