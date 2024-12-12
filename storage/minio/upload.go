package minio

import (
	"bytes"
	"context"
	"github.com/minio/minio-go/v7"
)

// Helper methods for uploading objects
func (c *MinioClient) Upload(ctx context.Context, bucketName, objectName string, data []byte, contentType string) error {
	c.Logger.Tracef("uploading data to bucket %s", bucketName)

	reader := bytes.NewReader(data)
	_, err := c.client.PutObject(ctx, bucketName, objectName, reader, int64(len(data)), minio.PutObjectOptions{ContentType: contentType})
	return err
}
