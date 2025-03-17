// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
)

func (c *Client) GetBucket(context.Context) string {
	// GetBucket returns the bucket name for the MinIO client.
	return c.config.Bucket
}
