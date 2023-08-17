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
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation POST /api/v1/pipelines/{org}/{repo} pipelines CreatePipeline
//
// Create a pipeline in the configured backend
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
// - in: body
//   name: body
//   description: Payload containing the pipeline to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Pipeline"
// security:
//   - ApiKeyAuth: []
// responses:
//   '201':
//     description: Successfully created the pipeline
//     type: json
//     schema:
//       "$ref": "#/definitions/Pipeline"
//   '400':
//     description: Unable to create the pipeline
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to create the pipeline
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to create the pipeline
//     schema:
//       "$ref": "#/definitions/Error"

// CreatePipeline represents the API handler to
// create a pipeline in the configured backend.
func CreatePipeline(c *gin.Context) {
	// capture middleware values
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logger := logrus.WithFields(logrus.Fields{
		"org":  o,
		"repo": r.GetName(),
		"user": u.GetName(),
	})

	logger.Infof("creating new pipeline for repo %s", r.GetFullName())

	// capture body from API request
	input := new(library.Pipeline)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for new build for repo %s: %w", r.GetFullName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in pipeline object
	input.SetRepoID(r.GetID())

	// send API call to create the pipeline
	p, err := database.FromContext(c).CreatePipeline(ctx, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create pipeline %s/%s: %w", r.GetFullName(), input.GetCommit(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusCreated, p)
}
