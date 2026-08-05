// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"

	api "github.com/go-vela/server/api/types"
)

// ListBuildObjectNames lists all artifact objects for a build, returning a map of
// object key to download URL. Objects stored under the standard prefix receive a
// 2-minute presigned GET URL (authenticated). Objects stored under the public/ prefix
// receive a direct, non-presigned URL that is accessible without authentication,
// provided the bucket has a public-read policy on the public/* prefix.
func (c *Client) ListBuildObjectNames(ctx context.Context, org, repo, build string) (map[string]string, error) {
	objectsWithURLs := make(map[string]string)

	type prefixEntry struct {
		prefix string
		public bool
	}

	prefixes := []prefixEntry{
		{prefix: fmt.Sprintf("%s/%s/%s/", org, repo, build), public: false},
		{prefix: fmt.Sprintf("public/%s/%s/%s/", org, repo, build), public: true},
	}

	for _, p := range prefixes {
		c.Logger.Tracef("listing object names in bucket %s with prefix %s", c.config.Bucket, p.prefix)

		opts := minio.ListObjectsOptions{
			Prefix:    p.prefix,
			Recursive: true,
		}

		for object := range c.client.ListObjects(ctx, c.config.Bucket, opts) {
			if object.Err != nil {
				return nil, object.Err
			}

			var (
				url string
				err error
			)

			if p.public {
				url = c.DirectObjectURL(object.Key)
			} else {
				obj := &api.Object{
					ObjectName: object.Key,
					Bucket: api.Bucket{
						BucketName: c.config.Bucket,
					},
				}

				url, err = c.PresignedGetObject(ctx, obj)
				if err != nil {
					return nil, fmt.Errorf("failed to generate presigned URL for object %s: %w", object.Key, err)
				}
			}

			objectsWithURLs[object.Key] = url
		}
	}

	return objectsWithURLs, nil
}
