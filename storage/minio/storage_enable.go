// SPDX-License-Identifier: Apache-2.0

package minio

func (c *Client) StorageEnable() bool {
	// GetBucket returns the bucket name for the MinIO client.
	return c.config.Enable
}
