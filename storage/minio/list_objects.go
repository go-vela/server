// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"

	"github.com/minio/minio-go/v7"

	api "github.com/go-vela/server/api/types"
)

// ListObjects lists the objects in a bucket.
func (c *Client) ListObjects(ctx context.Context, b *api.Bucket) ([]minio.ObjectInfo, error) {
	c.Logger.Tracef("listing objects in bucket %s", b.BucketName)

	objectCh := c.client.ListObjects(ctx, b.BucketName, minio.ListObjectsOptions{})

	var objects []minio.ObjectInfo

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		objects = append(objects, object)
	}

	return objects, nil
}
