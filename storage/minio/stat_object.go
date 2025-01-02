package minio

import (
	"context"
	"fmt"
	"github.com/go-vela/server/api/types"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

// StatObject retrieves the metadata of an object from the MinIO storage.
func (c *MinioClient) StatObject(ctx context.Context, object *types.Object) (*types.Object, error) {
	c.Logger.WithFields(logrus.Fields{
		"bucket": object.BucketName,
		"object": object.ObjectName,
	}).Tracef("retrieving metadata for object %s from bucket %s", object.ObjectName, object.BucketName)

	// Get object info
	info, err := c.client.StatObject(ctx, object.BucketName, object.ObjectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("unable to get object info %s from bucket %s: %v", object.ObjectName, object.BucketName, err)
	}

	// Map MinIO object info to API object
	return &types.Object{
		ObjectName: info.Key,
	}, nil
}
