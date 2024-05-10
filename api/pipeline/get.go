// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"bytes"
	"net/http"

	"github.com/alecthomas/chroma/v2/quick"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
)

// swagger:operation GET /api/v1/pipelines/{org}/{repo}/{pipeline} pipelines GetPipeline
//
// Get a pipeline from the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: pipeline
//   description: Commit SHA for pipeline to retrieve
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the pipeline
//     type: json
//     schema:
//       "$ref": "#/definitions/Pipeline"

// GetPipeline represents the API handler to capture
// a pipeline for a repo from the configured backend.
func GetPipeline(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	p := pipeline.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":      o,
		"pipeline": p.GetCommit(),
		"repo":     r.GetName(),
		"user":     u.GetName(),
	}).Infof("reading pipeline %s/%s", r.GetFullName(), p.GetCommit())

	buf := new(bytes.Buffer)
	err := quick.Highlight(buf, string(p.GetData()), "yaml", "terminal16", "monokai")
	if err == nil {
		p.SetData(buf.Bytes())
	}
	p.SetData(buf.Bytes())

	c.JSON(http.StatusOK, p)
}
