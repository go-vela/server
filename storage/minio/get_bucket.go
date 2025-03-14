package minio

import (
	"context"
)

func (c *MinioClient) GetBucket(ctx context.Context) string {
	// GetBucket returns the bucket name for the MinIO client.
	return c.config.Bucket
}
