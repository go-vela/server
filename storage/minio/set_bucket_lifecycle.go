package minio

import (
	"context"
	"encoding/xml"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

// SetBucketLifecycle sets the lifecycle configuration for a bucket.
func (c *MinioClient) SetBucketLifecycle(ctx context.Context, bucketName string, lifecycleConfig string) error {
	c.Logger.Tracef("setting lifecycle configuration for bucket %s", bucketName)

	var config lifecycle.Configuration
	if err := xml.Unmarshal([]byte(lifecycleConfig), &config); err != nil {
		c.Logger.Errorf("failed to unmarshal lifecycle configuration: %s", err)
		return err
	}

	return c.client.SetBucketLifecycle(ctx, bucketName, &config)
}
