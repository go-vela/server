// SPDX-License-Identifier: Apache-2.0

package minio

import "net/url"

// GetAddress returns the endpoint address for the MinIO client.
func (c *Client) GetAddress() string {
	if c == nil || c.config == nil {
		return ""
	}
	// Parse the configured endpoint to extract just the host:port
	if c.config.Endpoint == "" {
		return ""
	}

	u, err := url.Parse(c.config.Endpoint)
	if err != nil {
		// If parsing fails, return the endpoint as-is
		return c.config.Endpoint
	}

	return u.Host
}

// GetEndpoint returns the configured endpoint.
func (c *Client) GetEndpoint() string {
	if c == nil || c.config == nil {
		return ""
	}
	return c.config.Endpoint
}
