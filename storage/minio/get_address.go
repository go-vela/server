// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"fmt"
	"strings"
)

func (c *Client) GetAddress() string {
	// GetAddress returns the endpoint address for the MinIO client.
	var urlEndpoint string

	if len(c.config.Endpoint) > 0 {
		useSSL := strings.HasPrefix(c.config.Endpoint, "https://")

		if !useSSL {
			if !strings.HasPrefix(c.config.Endpoint, "http://") {
				return fmt.Sprintf("invalid server %s: must to be a HTTP URI", c.config.Endpoint)
			}

			urlEndpoint = c.config.Endpoint[7:]
		} else {
			urlEndpoint = c.config.Endpoint[8:]
		}
	}

	return urlEndpoint
}
