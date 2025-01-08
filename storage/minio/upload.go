package minio

import (
	"context"
	api "github.com/go-vela/server/api/types"
	"github.com/minio/minio-go/v7"
)

// Upload uploads an object to a bucket in MinIO.ts
func (c *MinioClient) Upload(ctx context.Context, object *api.Object) error {
	c.Logger.Tracef("uploading data to bucket %s", object.Bucket.BucketName)
	_, err := c.client.FPutObject(ctx, object.Bucket.BucketName, object.ObjectName, object.FilePath, minio.PutObjectOptions{})
	return err
}
