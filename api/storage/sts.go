// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/storage"
)

func GetSTSCreds(c *gin.Context) {
	l := c.MustGet("logger").(*logrus.Entry)
	enabled := c.MustGet("storage-enable").(bool)
	if !enabled {
		l.Info("storage is not enabled, skipping credentials request")
		c.JSON(http.StatusForbidden, gin.H{"error": "storage is not enabled"})

		return
	}

	org := c.Param("org")
	repo := c.Param("repo")
	build := c.Param("build")

	bucket := c.MustGet("storage-bucket").(string)

	prefix := fmt.Sprintf("%s/%s/%s/", org, repo, build)

	policy, _ := buildPutOnlyPolicy(bucket, prefix)
	sessionName := fmt.Sprintf("vela-%s-%s-%s", org, repo, build)

	creds, err := storage.FromGinContext(c).AssumeRole(c, int(1*time.Hour/time.Second), policy, sessionName)
	if creds == nil {
		l.Error("unable to assume role and generate temporary credentials")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to assume role and generate temporary credentials"})
		return
	}
	if err != nil {
		l.Errorf("unable to assume role: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, creds)
}

type policyDoc struct {
	Version   string      `json:"Version"`
	Statement []statement `json:"Statement"`
}

type statement struct {
	Effect   string   `json:"Effect"`
	Action   []string `json:"Action"`
	Resource []string `json:"Resource"`
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
