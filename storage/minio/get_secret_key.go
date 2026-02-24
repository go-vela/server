// SPDX-License-Identifier: Apache-2.0

package minio

// GetSecretKey returns the secret key for the MinIO client.
func (c *Client) GetSecretKey() string {
	if c == nil || c.config == nil {
		return ""
	}

	return c.config.SecretKey
}
