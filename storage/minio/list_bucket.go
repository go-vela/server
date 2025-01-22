package minio

import (
	"context"
)

// ListBuckets lists all buckets in MinIO.
func (c *MinioClient) ListBuckets(ctx context.Context) ([]string, error) {
	c.Logger.Trace("listing all buckets")

	buckets, err := c.client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	bucketNames := make([]string, len(buckets))
	for i, bucket := range buckets {
		bucketNames[i] = bucket.Name
	}
	return bucketNames, nil
}
