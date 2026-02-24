// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
)

// PresignedGetObject generates a presigned URL for downloading an object.
func (c *Client) PresignedGetObject(ctx context.Context, object *api.Object) (string, error) {
	c.Logger.Tracef("generating presigned URL for object %s in bucket %s", object.ObjectName, object.Bucket.BucketName)

	var url string
	// collect metadata on the object
	// make sure the object exists before generating the presigned URL
	objInfo, err := c.client.StatObject(ctx, object.Bucket.BucketName, object.ObjectName, minio.StatObjectOptions{})
	if objInfo.Key == "" {
		logrus.Errorf("unable to get object info %s from bucket %s: %v", object.ObjectName, object.Bucket.BucketName, err)
		return "", err
	}

	// Generate presigned URL for downloading the object.
	// The URL is valid for 2 minutes.
	presignedURL, err := c.client.PresignedGetObject(ctx, object.Bucket.BucketName, object.ObjectName, 2*time.Minute, nil)
	if err != nil {
		return fmt.Sprintf("Unable to generate presigned URL for object %s", object.ObjectName), err
	}

	url = presignedURL.String()

	return url, nil
}
