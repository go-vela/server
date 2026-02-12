// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"fmt"
)

// ClientOpt represents a configuration option to initialize the MinIO client.
type ClientOpt func(client *Client) error

// WithOptions sets multiple options in the MinIO client.
func WithOptions(enable, secure bool, endpoint, accessKey, secretKey, bucket, token string) ClientOpt {
	return func(c *Client) error {
		c.Logger.Trace("configuring multiple options in minio client")

		if len(accessKey) == 0 {
			return fmt.Errorf("no MinIO access key provided")
		}
		// check if the secret key provided is empty
		if len(secretKey) == 0 {
			return fmt.Errorf("no MinIO secret key provided")
		}
		// check if the bucket name provided is empty
		if len(bucket) == 0 {
			return fmt.Errorf("no MinIO bucket name provided")
		}
		// set the enable flag in the minio client
		c.config.Enable = enable
		// set the endpoint in the minio client
		c.config.Endpoint = endpoint
		// set the secret key in the minio client
		c.config.SecretKey = secretKey
		// set the access key in the minio client
		c.config.AccessKey = accessKey
		// set the secure connection mode in the minio client
		c.config.Secure = secure
		// set the bucket name in the minio client
		c.config.Bucket = bucket
		// set the token in the minio client
		c.config.Token = token

		return nil
	}
}
