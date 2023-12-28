// SPDX-License-Identifier: Apache-2.0

//nolint:dupl // ignore similar code with step
package log

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/service"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/services/{service}/logs services GetServiceLog
//
// Retrieve the logs for a service
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
//   description: Build number
//   required: true
//   type: integer
// - in: path
//   name: service
//   description: Service number
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved the service logs
//     schema:
//       "$ref": "#/definitions/Log"
//   '500':
//     description: Unable to retrieve the service logs
//     schema:
//       "$ref": "#/definitions/Error"

// GetServiceLog represents the API handler to capture
// the logs for a service from the configured backend.
func GetServiceLog(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	s := service.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d/%d", r.GetFullName(), b.GetNumber(), s.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build":   b.GetNumber(),
		"org":     o,
		"repo":    r.GetName(),
		"service": s.GetNumber(),
		"user":    u.GetName(),
	}).Infof("reading logs for service %s", entry)

	// send API call to capture the service logs
	l, err := database.FromContext(c).GetLogForService(ctx, s)
	if err != nil {
		retErr := fmt.Errorf("unable to get logs for service %s: %w", entry, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	c.JSON(http.StatusOK, l)
}
