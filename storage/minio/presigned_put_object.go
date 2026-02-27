// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"fmt"
	"time"
)

func (c *Client) PresignedPutObject(ctx context.Context, path string, durationSeconds time.Duration) (string, error) {
	c.Logger.Tracef("generating presigned PUT URL for object %s in bucket %s", path, c.config.Bucket)
	// Generate presigned URL for downloading the object.
	// The URL is valid for 2 minutes.
	presignedURL, err := c.client.PresignedPutObject(ctx, c.config.Bucket, path, durationSeconds)
	if err != nil {
		return fmt.Sprintf("Unable to generate presigned URL for object %s", path), err
	}

	return presignedURL.String(), nil
}
