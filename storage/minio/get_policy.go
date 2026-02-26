// SPDX-License-Identifier: Apache-2.0

package minio

import (
	"encoding/json"
	"strings"
)

type policyDoc struct {
	Version   string      `json:"Version"`
	Statement []statement `json:"Statement"`
}

type statement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource []string `json:"Resource"`
}

func (c *Client) GetPolicy(prefix string) string {
	policy, err := buildPutOnlyPolicy(c.config.Bucket, prefix)
	if err != nil {
		c.Logger.Debugf("failed to build policy: %v", err)
		return ""
	}

	return policy
}

func buildPutOnlyPolicy(bucket, prefix string) (string, error) {
	// Normalize prefix
	prefix = strings.TrimPrefix(prefix, "/")
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}

	doc := policyDoc{
		Version: "2012-10-17",
		Statement: []statement{
			{
				Effect: "Allow",
				Action: []string{
					"s3:PutObject",
					"s3:AbortMultipartUpload",
					"s3:ListMultipartUploadParts",
				},
				Resource: []string{
					"arn:aws:s3:::" + bucket + "/" + prefix + "*",
				},
			},
		},
	}

	b, err := json.Marshal(doc)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
