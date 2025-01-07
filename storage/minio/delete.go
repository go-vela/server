package minio

import (
	"context"
	api "github.com/go-vela/server/api/types"
	"github.com/minio/minio-go/v7"
)

// Delete deletes an object in a bucket in MinIO.
func (c *MinioClient) Delete(ctx context.Context, object *api.Object) error {
	c.Logger.Tracef("deleting objectName: %s from bucketName: %s", object.ObjectName, object.BucketName)

	err := c.client.RemoveObject(ctx, object.BucketName, object.ObjectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}
