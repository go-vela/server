// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package repo

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/repos repos ListRepos
//
// Get all repos in the configured backend
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// parameters:
// - in: query
//   name: page
//   description: The page of results to retrieve
//   type: integer
//   default: 1
// - in: query
//   name: per_page
//   description: How many results per page to return
//   type: integer
//   maximum: 100
//   default: 10
// responses:
//   '200':
//     description: Successfully retrieved the repo
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Repo"
//     headers:
//       X-Total-Count:
//         description: Total number of results
//         type: integer
//       Link:
//         description: see https://tools.ietf.org/html/rfc5988
//         type: string
//   '400':
//     description: Unable to retrieve the repo
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unable to retrieve the repo
//     schema:
//       "$ref": "#/definitions/Error"

// ListRepos represents the API handler to capture a list
// of repos for a user from the configured backend.
func ListRepos(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("listing repos for user %s", u.GetName())

	// capture page query parameter if present
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert page query parameter for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// capture per_page query parameter if present
	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		retErr := fmt.Errorf("unable to convert per_page query parameter for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// ensure per_page isn't above or below allowed values
	perPage = util.MaxInt(1, util.MinInt(100, perPage))

	// capture the sort_by query parameter if present
	sortBy := util.QueryParameter(c, "sort_by", "name")

	// capture the query parameters if present:
	//
	// * active
	filters := map[string]interface{}{
		"active": util.QueryParameter(c, "active", "true"),
	}

	// send API call to capture the list of repos for the user
	r, t, err := database.FromContext(c).ListReposForUser(ctx, u, sortBy, filters, page, perPage)
	if err != nil {
		retErr := fmt.Errorf("unable to get repos for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// create pagination object
	pagination := api.Pagination{
		Page:    page,
		PerPage: perPage,
		Total:   t,
	}
	// set pagination headers
	pagination.SetHeaderLink(c)

	c.JSON(http.StatusOK, r)
}
