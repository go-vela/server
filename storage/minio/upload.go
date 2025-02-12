package minio

import (
	"context"
	api "github.com/go-vela/server/api/types"
	"github.com/minio/minio-go/v7"
	"io"
)

// Upload uploads an object to a bucket in MinIO.ts
func (c *MinioClient) Upload(ctx context.Context, object *api.Object) error {
	c.Logger.Tracef("uploading data to bucket %s", object.Bucket.BucketName)
	_, err := c.client.FPutObject(ctx, object.Bucket.BucketName, object.ObjectName, object.FilePath, minio.PutObjectOptions{})

	return err
}

// UploadObject uploads an object to a bucket in MinIO.ts
func (c *MinioClient) UploadObject(ctx context.Context, object *api.Object, reader io.Reader, size int64) error {
	c.Logger.Tracef("uploading data to bucket %s", object.Bucket.BucketName)
	//_, err := c.client.FPutObject(ctx, object.Bucket.BucketName, object.ObjectName, object.FilePath, minio.PutObjectOptions{})
	_, err := c.client.PutObject(ctx, object.Bucket.BucketName, object.ObjectName, reader, size, minio.PutObjectOptions{})

	return err
}
