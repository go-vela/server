package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
)

// Delete deletes an object in a bucket in MinIO.
func (c *MinioClient) Delete(ctx context.Context, bucketName string, objectName string) error {
	c.Logger.Tracef("deleting objectName: %s from bucketName: %s", objectName, bucketName)

	err := c.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
