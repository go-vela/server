// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"fmt"

	api "github.com/go-vela/server/api/types"
)

// CreateBucket creates a new bucket in MinIO.
func (c *Client) CreateBucket(ctx context.Context, bucket *api.Bucket) error {
	c.Logger.Tracef("create new bucket: %s", bucket.BucketName)

	exists, errBucketExists := c.BucketExists(ctx, bucket)
	if errBucketExists != nil && exists {
		c.Logger.Tracef("Bucket %s already exists", bucket.BucketName)

		return fmt.Errorf("bucket %s already exists", bucket.BucketName)
	}

	err := c.client.MakeBucket(ctx, bucket.BucketName, bucket.MakeBucketOptions)
	if err != nil {
		c.Logger.Errorf("unable to create bucket %s: %v", bucket.BucketName, err)
		return err
	}

	return nil
}
