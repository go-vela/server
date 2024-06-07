// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/executable builds GetBuildExecutable
//
// Get a build executable
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the organization
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repository
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the build executable
//     type: json
//     schema:
//       "$ref": "#/definitions/BuildExecutable"
//   '400':
//     description: Invalid request payload or path
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// GetBuildExecutable represents the API handler to get
// a build executable for a repository.
func GetBuildExecutable(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	cl := claims.Retrieve(c)
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build":   b.GetNumber(),
		"org":     o,
		"repo":    r.GetName(),
		"subject": cl.Subject,
	}).Infof("reading build executable %s/%d", r.GetFullName(), b.GetNumber())

	// send database call to pop the requested build executable from the table
	bExecutable, err := database.FromContext(c).PopBuildExecutable(ctx, b.GetID())
	if err != nil {
		retErr := fmt.Errorf("unable to pop build executable: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, bExecutable)
}

// PublishBuildExecutable marshals a pipeline.Build into bytes and pushes that data to the build_executables table to be
// requested by a worker whenever the build has been picked up.
func PublishBuildExecutable(ctx context.Context, db database.Interface, p *pipeline.Build, b *types.Build) error {
	// marshal pipeline build into byte data to add to the build executable object
	byteExecutable, err := json.Marshal(p)
	if err != nil {
		logrus.Errorf("Failed to marshal build executable: %v", err)

		// error out the build
		CleanBuild(ctx, db, b, nil, nil, err)

		return err
	}

	// create build executable to push to database
	bExecutable := new(library.BuildExecutable)
	bExecutable.SetBuildID(b.GetID())
	bExecutable.SetData(byteExecutable)

	// send database call to create a build executable
	err = db.CreateBuildExecutable(ctx, bExecutable)
	if err != nil {
		logrus.Errorf("Failed to publish build executable to database: %v", err)

		// error out the build
		CleanBuild(ctx, db, b, nil, nil, err)

		return err
	}

	return nil
}
