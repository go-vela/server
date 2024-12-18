package minio

import (
	"context"
	api "github.com/go-vela/server/api/types"
	"time"
)

// PresignedGetObject generates a presigned URL for downloading an object.
func (c *MinioClient) PresignedGetObject(ctx context.Context, object *api.Object) (string, error) {
	// Generate presigned URL for downloading the object.
	// The URL is valid for 7 days.
	presignedURL, err := c.client.PresignedGetObject(ctx, object.BucketName, object.ObjectName, 7*24*time.Hour, nil)
	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}
