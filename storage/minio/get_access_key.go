// SPDX-License-Identifier: Apache-2.0

package minio

// GetAccessKey returns the configured access key.
func (c *Client) GetAccessKey() string {
	if c == nil || c.config == nil {
		return ""
	}
	return c.config.AccessKey
}
