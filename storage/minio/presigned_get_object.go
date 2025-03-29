// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
)

// TODO hide URL behind a different name
// PresignedGetObject generates a presigned URL for downloading an object.
func (c *Client) PresignedGetObject(ctx context.Context, object *api.Object) (string, error) {
	c.Logger.Tracef("generating presigned URL for object %s in bucket %s", object.ObjectName, object.Bucket.BucketName)

	// collect metadata on the object
	// make sure the object exists before generating the presigned URL
	objInfo, err := c.client.StatObject(ctx, object.Bucket.BucketName, object.ObjectName, minio.StatObjectOptions{})
	if objInfo.Key == "" {
		logrus.Errorf("unable to get object info %s from bucket %s: %v", object.ObjectName, object.Bucket.BucketName, err)
		return "", err
	}

	// Generate presigned URL for downloading the object.
	// The URL is valid for 7 days.
	presignedURL, err := c.client.PresignedGetObject(ctx, object.Bucket.BucketName, object.ObjectName, 7*24*time.Hour, nil)
	if err != nil {
		return "", err
	}
	//presignedURL.RequestURI()
	return presignedURL.String(), nil
}
