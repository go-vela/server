package minio

import "context"

// DeleteBucket deletes a bucket in MinIO.
func (c *MinioClient) DeleteBucket(ctx context.Context, bucketName string) error {
	c.Logger.Tracef("deleting bucketName: %s", bucketName)

	err := c.client.RemoveBucket(ctx, bucketName)
	if err != nil {
		return err
	}
	return nil
}
