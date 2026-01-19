// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

// ArtifactResp is the mock response for creating an artifact.
const ArtifactResp = `{
	"id": 1,
	"build_id": 1,
	"file_name": "test-results.xml",
	"object_path": "builds/1/test-results.xml",
	"file_size": 1024,
	"file_type": "xml",
	"presigned_url": "https://storage.example.com/builds/1/test-results.xml",
	"created_at": 1750710551
}`

func addArtifact(c *gin.Context) {
	data := []byte(ArtifactResp)

	var body api.Artifact

	_ = json.Unmarshal(data, &body)
	c.JSON(http.StatusCreated, body)
}
