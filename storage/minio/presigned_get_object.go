package minio

import (
	"context"
	api "github.com/go-vela/server/api/types"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"time"
)

// TODO hide URL behind a different name
// PresignedGetObject generates a presigned URL for downloading an object.
func (c *MinioClient) PresignedGetObject(ctx context.Context, object *api.Object) (string, error) {
	c.Logger.Tracef("generating presigned URL for object %s in bucket %s", object.ObjectName, object.Bucket.BucketName)

	// collect metadata on the object
	objInfo, err := c.client.StatObject(ctx, object.Bucket.BucketName, object.ObjectName, minio.StatObjectOptions{})
	if objInfo.Key == "" {
		logrus.Errorf("unable to get object info %s from bucket %s: %v", object.ObjectName, object.Bucket.BucketName, err)
		return "", err
	}

	_, err = c.client.BucketExists(ctx, object.Bucket.BucketName)
	if err != nil {
		logrus.Errorf("unable to check if bucket %s exists: %v", object.Bucket.BucketName, err)
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
