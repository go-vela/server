package minio

import (
	"context"
	api "github.com/go-vela/server/api/types"
)

// BucketExists checks if a bucket exists in MinIO.
func (c *MinioClient) BucketExists(ctx context.Context, bucket *api.Bucket) (bool, error) {
	c.Logger.Tracef("checking if bucket %s exists", bucket.BucketName)

	exists, err := c.client.BucketExists(ctx, bucket.BucketName)
	if err != nil {
		return false, err
	}
	return exists, nil
}
