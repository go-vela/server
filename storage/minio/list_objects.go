// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"

	api "github.com/go-vela/server/api/types"
)

// ListBuildObjectNames lists the names of objects in a bucket for a specific build.
func (c *Client) ListBuildObjectNames(ctx context.Context, org, repo, build string) (map[string]string, error) {
	objectsWithURLs := make(map[string]string)
	// Construct the prefix path for filtering
	prefix := org + "/" + repo + "/" + build + "/"

	c.Logger.Tracef("listing object names in bucket %s with prefix %s", c.config.Bucket, prefix)

	b := api.Bucket{
		BucketName: c.GetBucket(),
		ListObjectsOptions: minio.ListObjectsOptions{
			Prefix:    prefix,
			Recursive: true,
		},
	}

	objectCh := c.client.ListObjects(ctx, c.config.Bucket, b.ListObjectsOptions)

	var objectNames []string

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		objectNames = append(objectNames, object.Key)
		// Generate presigned URL for each object
		obj := &api.Object{
			ObjectName: object.Key,
			Bucket:     b,
		}

		url, err := c.PresignedGetObject(ctx, obj)
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned URL for object %s: %w", object.Key, err)
		}

		objectsWithURLs[object.Key] = url
	}

	return objectsWithURLs, nil
}
