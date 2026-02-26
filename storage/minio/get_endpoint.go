// SPDX-License-Identifier: Apache-2.0

package minio

// GetEndpoint returns the configured endpoint.
func (c *Client) GetEndpoint() string {
	if c == nil || c.config == nil {
		return ""
	}

	return c.config.Endpoint
}
