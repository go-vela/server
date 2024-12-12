package minio

import (
	"context"
)

// BucketExists checks if a bucket exists in MinIO.
func (c *MinioClient) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	c.Logger.Tracef("checking if bucket %s exists", bucketName)

	exists, err := c.client.BucketExists(ctx, bucketName)
	if err != nil {
		return false, err
	}
	return exists, nil
}
