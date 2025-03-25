// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"time"

	api "github.com/go-vela/server/api/types"
)

// TODO hide URL behind a different name
// PresignedGetObject generates a presigned URL for downloading an object.
func (c *Client) PresignedGetObject(ctx context.Context, object *api.Object) (string, error) {
	c.Logger.Tracef("generating presigned URL for object %s in bucket %s", object.ObjectName, object.Bucket.BucketName)

	// Generate presigned URL for downloading the object.
	// The URL is valid for 7 days.
	presignedURL, err := c.client.PresignedGetObject(ctx, object.Bucket.BucketName, object.ObjectName, 7*24*time.Hour, nil)
	if err != nil {
		return "", err
	}
	//presignedURL.RequestURI()
	return presignedURL.String(), nil
}
