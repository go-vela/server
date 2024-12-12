package minio

import (
	"context"
	"encoding/xml"
)

// GetBucketLifecycle retrieves the lifecycle configuration for a bucket.
func (c *MinioClient) GetBucketLifecycle(ctx context.Context, bucketName string) (string, error) {
	c.Logger.Tracef("getting lifecycle configuration for bucket %s", bucketName)

	lifecycleConfig, err := c.client.GetBucketLifecycle(ctx, bucketName)
	if err != nil {
		return "", err
	}

	lifecycleBytes, err := xml.MarshalIndent(lifecycleConfig, "", "  ")

	return string(lifecycleBytes), nil
}
