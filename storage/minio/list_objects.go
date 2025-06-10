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

	opts := minio.ListObjectsOptions{
		Recursive: b.Recursive,
	}

	objectCh := c.client.ListObjects(ctx, b.BucketName, opts)

	var objects []minio.ObjectInfo

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		objects = append(objects, object)
	}

	return objects, nil
}

// ListObjectNames lists only the names of objects in a bucket.
func (c *Client) ListObjectNames(ctx context.Context, b *api.Bucket) ([]string, error) {
	c.Logger.Tracef("listing object names in bucket %s", b.BucketName)

	// Set ListObjectsOptions with Recursive flag from the Bucket type
	opts := minio.ListObjectsOptions{
		Recursive: b.Recursive,
	}

	objectCh := c.client.ListObjects(ctx, b.BucketName, opts)

	var objectNames []string

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		objectNames = append(objectNames, object.Key)
	}

	return objectNames, nil
}

// ListBuildObjectNames lists the names of objects in a bucket for a specific build.
func (c *Client) ListBuildObjectNames(ctx context.Context, b *api.Bucket, org, repo, build string) ([]string, error) {
	// Construct the prefix path for filtering
	prefix := org + "/" + repo + "/" + build + "/"

	c.Logger.Tracef("listing object names in bucket %s with prefix %s", b.BucketName, prefix)

	// Set ListObjectsOptions with Recursive flag and prefix
	opts := minio.ListObjectsOptions{
		Recursive: b.Recursive,
		Prefix:    prefix,
	}

	objectCh := c.client.ListObjects(ctx, b.BucketName, opts)

	var objectNames []string

	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}

		objectNames = append(objectNames, object.Key)
	}

	return objectNames, nil
}
