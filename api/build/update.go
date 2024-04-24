// SPDX-License-Identifier: Apache-2.0

package build

import (
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
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// swagger:operation PUT /api/v1/repos/{org}/{repo}/builds/{build} builds UpdateBuild
//
// Updates a build in the configured backend
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
//   name: build
//   description: Build number to update
//   required: true
//   type: integer
// - in: body
//   name: body
//   description: Payload containing the build to update
//   required: true
//   schema:
//     "$ref": "#/definitions/Build"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully updated the build
//     schema:
//       "$ref": "#/definitions/Build"
//   '404':
//     description: Unable to update the build
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to update the build
//     schema:
//       "$ref": "#/definitions/Error"

// UpdateBuild represents the API handler to update
// a build for a repo in the configured backend.
func UpdateBuild(c *gin.Context) {
	// capture middleware values
	cl := claims.Retrieve(c)
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  cl.Subject,
	}).Infof("updating build %s", entry)

	// capture body from API request
	input := new(types.Build)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for build %s: %w", entry, err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// update build fields if provided
	if len(input.GetStatus()) > 0 {
		// update status if set
		b.SetStatus(input.GetStatus())
	}

	if len(input.GetError()) > 0 {
		// update error if set
		b.SetError(input.GetError())
	}

	if input.GetEnqueued() > 0 {
		// update enqueued if set
		b.SetEnqueued(input.GetEnqueued())
	}

	if input.GetStarted() > 0 {
		// update started if set
		b.SetStarted(input.GetStarted())
	}

	if input.GetFinished() > 0 {
		// update finished if set
		b.SetFinished(input.GetFinished())
	}

	if len(input.GetTitle()) > 0 {
		// update title if set
		b.SetTitle(input.GetTitle())
	}

	if len(input.GetMessage()) > 0 {
		// update message if set
		b.SetMessage(input.GetMessage())
	}

	if len(input.GetHost()) > 0 {
		// update host if set
		b.SetHost(input.GetHost())
	}

	if len(input.GetRuntime()) > 0 {
		// update runtime if set
		b.SetRuntime(input.GetRuntime())
	}

	if len(input.GetDistribution()) > 0 {
		// update distribution if set
		b.SetDistribution(input.GetDistribution())
	}

	// send API call to update the build
	b, err = database.FromContext(c).UpdateBuild(ctx, b)
	if err != nil {
		retErr := fmt.Errorf("unable to update build %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, b)

	// check if the build is in a "final" state
	// and if build is not a scheduled event
	if (b.GetStatus() == constants.StatusSuccess ||
		b.GetStatus() == constants.StatusFailure ||
		b.GetStatus() == constants.StatusCanceled ||
		b.GetStatus() == constants.StatusKilled ||
		b.GetStatus() == constants.StatusError) && b.GetEvent() != constants.EventSchedule {
		// send API call to set the status on the commit
		err = scm.FromContext(c).Status(ctx, r.GetOwner(), b, r.GetOrg(), r.GetName())
		if err != nil {
			logrus.Errorf("unable to set commit status for build %s: %v", entry, err)
		}
	}
}

// UpdateComponentStatuses updates all components (steps and services) for a build to a given status.
func UpdateComponentStatuses(c *gin.Context, b *types.Build, status string) error {
	ctx := c.Request.Context()

	// retrieve the steps for the build from the step table
	steps := []*library.Step{}
	page := 1
	perPage := 100

	for page > 0 {
		// retrieve build steps (per page) from the database
		stepsPart, _, err := database.FromContext(c).ListStepsForBuild(ctx, b, map[string]interface{}{}, page, perPage)
		if err != nil {
			return err
		}

		// add page of steps to list steps
		steps = append(steps, stepsPart...)

		// assume no more pages exist if under 100 results are returned
		if len(stepsPart) < 100 {
			page = 0
		} else {
			page++
		}
	}

	// iterate over each step for the build
	// setting status
	for _, step := range steps {
		step.SetStatus(status)

		_, err := database.FromContext(c).UpdateStep(ctx, step)
		if err != nil {
			return err
		}
	}

	// retrieve the services for the build from the service table
	services := []*library.Service{}
	page = 1

	for page > 0 {
		// retrieve build services (per page) from the database
		servicesPart, _, err := database.FromContext(c).ListServicesForBuild(ctx, b, map[string]interface{}{}, page, perPage)
		if err != nil {
			return err
		}

		// add page of services to the list of services
		services = append(services, servicesPart...)

		// assume no more pages exist if under 100 results are returned
		if len(servicesPart) < 100 {
			page = 0
		} else {
			page++
		}
	}

	// iterate over each service for the build
	// setting status
	for _, service := range services {
		service.SetStatus(status)

		_, err := database.FromContext(c).UpdateService(ctx, service)
		if err != nil {
			return err
		}
	}

	return nil
}
