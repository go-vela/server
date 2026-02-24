// SPDX-License-Identifier: Apache-2.0

package minio

// GetBucket returns the configured bucket name.
func (c *Client) GetBucket() string {
	if c == nil || c.config == nil {
		return ""
	}

	return c.config.Bucket
}
