package minio

import (
	"fmt"
)

// ClientOpt represents a configuration option to initialize the MinIO client.
type ClientOpt func(client *MinioClient) error

// WithAccessKey sets the access key in the MinIO client.
func WithAccessKey(accessKey string) ClientOpt {
	return func(c *MinioClient) error {
		c.Logger.Trace("configuring access key in minio client")

		// check if the access key provided is empty
		if len(accessKey) == 0 {
			return fmt.Errorf("no MinIO access key provided")
		}

		// set the access key in the minio client
		c.config.AccessKey = accessKey

		return nil
	}
}

// WithSecretKey sets the secret key in the MinIO client.
func WithSecretKey(secretKey string) ClientOpt {
	return func(c *MinioClient) error {
		c.Logger.Trace("configuring secret key in minio client")

		// check if the secret key provided is empty
		if len(secretKey) == 0 {
			return fmt.Errorf("no MinIO secret key provided")
		}

		// set the secret key in the minio client
		c.config.SecretKey = secretKey

		return nil
	}
}

// WithSecure sets the secure connection mode in the MinIO client.
func WithSecure(secure bool) ClientOpt {
	return func(c *MinioClient) error {
		c.Logger.Trace("configuring secure connection mode in minio client")

		// set the secure connection mode in the minio client
		c.config.Secure = secure

		return nil
	}
}
