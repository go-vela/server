// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"fmt"
	"net/url"
)

func (c *Client) GetAddress() string {
	// GetAddress returns the endpoint address for the MinIO client.
	urlEndpoint, err := url.Parse(c.config.Endpoint)
	if err != nil {
		return fmt.Sprintf("invalid server %s: must to be a HTTP URI", c.config.Endpoint)
	}

	return urlEndpoint.Host
}
