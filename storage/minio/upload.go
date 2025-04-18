// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"io"
	"mime"
	"path/filepath"

	"github.com/minio/minio-go/v7"

	api "github.com/go-vela/server/api/types"
)

// Upload uploads an object to a bucket in MinIO.ts.
func (c *Client) Upload(ctx context.Context, object *api.Object) error {
	c.Logger.Tracef("uploading data to bucket %s", object.Bucket.BucketName)
	info, err := c.client.FPutObject(ctx, object.Bucket.BucketName, object.ObjectName, object.FilePath, minio.PutObjectOptions{})

	c.Logger.Infof("uploaded object %v with size %d", info, info.Size)

	return err
}

// UploadObject uploads an object to a bucket in MinIO.ts.
func (c *Client) UploadObject(ctx context.Context, object *api.Object, reader io.Reader, size int64) error {
	c.Logger.Infof("uploading data to bucket %s", object.Bucket.BucketName)
	ext := filepath.Ext(object.FilePath)
	contentType := mime.TypeByExtension(ext)

	c.Logger.Infof("uploading object %s with content type %s", object.ObjectName, contentType)
	// TODO - better way to get bucket name
	info, err := c.client.PutObject(ctx, object.Bucket.BucketName, object.ObjectName, reader, size,
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		c.Logger.Errorf("unable to upload object %s: %v", object.ObjectName, err)
		return err
	}

	c.Logger.Infof("uploaded object %v with size %d", info, info.Size)

	return nil
}
