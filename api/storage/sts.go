// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
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

	r := repo.Retrieve(c)
	org := r.GetOrg()
	b := build.Retrieve(c)
	repoName := r.GetName()
	buildNum := b.GetNumber()
	ctx := c.Request.Context()

	bucket := c.MustGet("storage-bucket").(string)

	prefix := fmt.Sprintf("%s/%s/%d/", org, repoName, buildNum)

	policy, _ := buildPutOnlyPolicy(bucket, prefix)
	sessionName := fmt.Sprintf("vela-%s-%s-%d", org, repoName, buildNum)

	creds, err := storage.FromGinContext(c).AssumeRole(ctx, int(r.GetTimeout())*60, policy, sessionName)
	if creds == nil {
		l.Errorf("unable to assume role and generate temporary credentials without error %s", err)
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
