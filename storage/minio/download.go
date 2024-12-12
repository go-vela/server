package minio

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
)

func (c *MinioClient) Download(ctx context.Context, bucketName, key string) ([]byte, error) {
	object, err := c.client.GetObject(ctx, bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()
	return io.ReadAll(object)
}
