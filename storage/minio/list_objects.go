// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"

	"github.com/minio/minio-go/v7"
)

// ListObjects lists the objects in a bucket.
func (c *Client) ListObjects(ctx context.Context, bucketName string) ([]string, error) {
	c.Logger.Tracef("listing objects in bucket %s", bucketName)

	objectCh := c.client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{})

	var objects []string

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		objects = append(objects, object.Key)
	}

	return objects, nil
}
