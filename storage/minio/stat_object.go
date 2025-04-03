// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"

	"github.com/go-vela/server/api/types"
)

// StatObject retrieves the metadata of an object from the MinIO storage.
func (c *Client) StatObject(ctx context.Context, object *types.Object) (*types.Object, error) {
	c.Logger.Tracef("retrieving metadata for object %s from bucket %s", object.ObjectName, object.Bucket.BucketName)

	// Get object info
	info, err := c.client.StatObject(ctx, object.Bucket.BucketName, object.ObjectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get object info %s from bucket %s: %w", object.ObjectName, object.Bucket.BucketName, err)
	}

	// Map MinIO object info to API object
	return &types.Object{
		ObjectName: info.Key,
	}, nil
}
