// SPDX-License-Identifier: Apache-2.0

package minio

import "fmt"

// DirectObjectURL returns a non-presigned direct URL for an object in the bucket.
// This is used for objects stored under the public/ prefix that are accessible
// without authentication when the bucket has a public-read policy configured
// for that prefix.
func (c *Client) DirectObjectURL(objectKey string) string {
	// c.config.Endpoint is a fully-qualified URL (e.g. "http://minio:9000"),
	// so we just append the bucket and object key directly.
	return fmt.Sprintf("%s/%s/%s", c.config.Endpoint, c.config.Bucket, objectKey)
}
