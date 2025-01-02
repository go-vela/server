package minio

import (
	"context"
	api "github.com/go-vela/server/api/types"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

// Upload uploads an object to a bucket in MinIO.ts
func (c *MinioClient) Upload(ctx context.Context, object *api.Object) error {
	c.Logger.WithFields(logrus.Fields{
		"bucket": object.BucketName,
		"object": object.ObjectName,
	}).Tracef("uploading data to bucket %s", object.BucketName)
	_, err := c.client.FPutObject(ctx, object.BucketName, object.ObjectName, object.FilePath, minio.PutObjectOptions{})
	return err
}
