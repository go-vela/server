// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dustin/go-humanize"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/archiver/v3"
	api "github.com/go-vela/server/api/types"
)

func (c *Client) Download(ctx context.Context, object *api.Object) error {
	// Temporary file to store the object
	filename := "/"

	logrus.Debugf("getting object info on bucket %s from path: %s", object.Bucket.BucketName, object.ObjectName)

	// collect metadata on the object
	objInfo, err := c.client.StatObject(ctx, object.Bucket.BucketName, object.ObjectName, minio.StatObjectOptions{})
	if objInfo.Key == "" {
		logrus.Error(err)
		return nil
	}

	// retrieve the object in specified path of the bucket
	err = c.client.FGetObject(ctx, object.Bucket.BucketName, object.ObjectName, filename, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	safeDir := "/safe/directory"
	absFilePath, err := filepath.Abs(filepath.Join(safeDir, object.FilePath))
	if err != nil || !strings.HasPrefix(absFilePath, safeDir) {
		return fmt.Errorf("invalid file path")
	}
	stat, err := os.Stat(absFilePath)
	if err != nil {
		return err
	}

	logrus.Infof("downloaded %s to %s on local filesystem", humanize.Bytes(uint64(stat.Size())), filename)

	logrus.Debug("getting current working directory")

	// grab the current working directory for unpacking the object
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	logrus.Debugf("unarchiving file %s into directory %s", filename, pwd)

	// expand the object back onto the filesystem
	err = archiver.Unarchive(absFilePath, pwd)
	if err != nil {
		return err
	}

	logrus.Infof("successfully unpacked file %s", object.FilePath)

	// delete the temporary archive file
	err = os.Remove(absFilePath)
	if err != nil {
		logrus.Infof("delete of file %s unsuccessful", filename)
	} else {
		logrus.Infof("file archive %s successfully deleted", filename)
	}

	logrus.Infof("object downloaded successfully")

	return nil
}
