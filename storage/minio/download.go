package minio

import (
	"context"
	"github.com/dustin/go-humanize"
	"github.com/go-vela/archiver/v3"
	api "github.com/go-vela/server/api/types"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"os"
)

func (c *MinioClient) Download(ctx context.Context, object *api.Object) error {
	// Temporary file to store the object
	filename := "/"
	// Check if the directory exists
	//_, err := os.Stat(object.FilePath)
	//if os.IsNotExist(err) {
	//	// Create the directory if it does not exist
	//	err = os.MkdirAll(object.FilePath, 0755)
	//	if err != nil {
	//		return fmt.Errorf("failed to create directory: %w", err)
	//	}
	//} else if err != nil {
	//	return fmt.Errorf("failed to check directory: %w", err)
	//}
	//err := c.client.FGetObject(ctx, object.BucketName, object.ObjectName, object.FilePath, minio.GetObjectOptions{})
	//if err != nil {
	//	c.Logger.Errorf("unable to retrive object %s", object.ObjectName)
	//	return err
	//}
	//
	//c.Logger.Tracef("successfully downloaded object %s to %s", object.ObjectName, object.FilePath)
	//return nil
	logrus.Debugf("getting object info on bucket %s from path: %s", object.Bucket.BucketName, object.ObjectName)

	// set a timeout on the request to the cache provider
	//ctx, cancel := context.WithTimeout(context.Background(), r.Timeout)
	//defer cancel()

	// collect metadata on the object
	objInfo, err := c.client.StatObject(ctx, object.Bucket.BucketName, object.ObjectName, minio.StatObjectOptions{})
	if objInfo.Key == "" {
		logrus.Error(err)
		return nil
	}

	logrus.Debugf("getting object in bucket %s from path: %s", object.Bucket.BucketName, object.ObjectName)

	logrus.Infof("%s to download", humanize.Bytes(uint64(objInfo.Size)))

	// retrieve the object in specified path of the bucket
	err = c.client.FGetObject(ctx, object.Bucket.BucketName, object.ObjectName, filename, minio.GetObjectOptions{})
	if err != nil {
		return err
	}

	stat, err := os.Stat(object.FilePath)
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
	err = archiver.Unarchive(object.FilePath, pwd)
	if err != nil {
		return err
	}

	logrus.Infof("successfully unpacked file %s", object.FilePath)

	// delete the temporary archive file
	err = os.Remove(filename)
	if err != nil {
		logrus.Infof("delete of file %s unsuccessful", filename)
	} else {
		logrus.Infof("file archive %s successfully deleted", filename)
	}

	logrus.Infof("object downloaded successfully")

	return nil

}
