// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/pipeline"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// swagger:operation DELETE /api/v1/pipelines/{org}/{repo}/{pipeline} pipelines DeletePipeline
//
// Delete a pipeline from the configured backend
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
//   description: Commit SHA for pipeline to delete
//   required: true
//   type: string
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully deleted the pipeline
//     schema:
//       type: string
//   '400':
//     description: Unable to delete the pipeline
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to delete the pipeline
//     schema:
//       "$ref": "#/definitions/Error"

// DeletePipeline represents the API handler to remove
// a pipeline for a repo from the configured backend.
func DeletePipeline(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	p := pipeline.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), p.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"org":      o,
		"pipeline": p.GetNumber(),
		"repo":     r.GetName(),
		"user":     u.GetName(),
	}).Infof("deleting pipeline %s", entry)

	// send API call to remove the build
	err := database.FromContext(c).DeletePipeline(p)
	if err != nil {
		retErr := fmt.Errorf("unable to delete pipeline %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, fmt.Sprintf("pipeline %s deleted", entry))
}
