package minio

import (
	"context"
	api "github.com/go-vela/server/api/types"
	"github.com/minio/minio-go/v7"
)

func (c *MinioClient) Download(ctx context.Context, object *api.Object) error {

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
	err := c.client.FGetObject(ctx, object.BucketName, object.ObjectName, object.FilePath, minio.GetObjectOptions{})
	if err != nil {
		c.Logger.Errorf("unable to retrive object %s", object.ObjectName)
		return err
	}

	c.Logger.Tracef("successfully downloaded object %s to %s", object.ObjectName, object.FilePath)
	return nil
}
