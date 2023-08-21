// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package user

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /api/v1/user/source/repos users GetSourceRepos
//
// Retrieve a list of repos for the current authenticated user
//
// ---
// produces:
// - application/json
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved a list of repos for the current user
//     schema:
//       "$ref": "#/definitions/Repo"
//   '404':
//     description: Unable to retrieve a list of repos for the current user
//     schema:
//       "$ref": "#/definitions/Error"

// GetSourceRepos represents the API handler to capture
// the list of repos for a user from the configured backend.
func GetSourceRepos(c *gin.Context) {
	// capture middleware values
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Infof("reading available SCM repos for user %s", u.GetName())

	// variables to capture requested data
	dbRepos := []*library.Repo{}
	output := make(map[string][]library.Repo)

	// send API call to capture the list of repos for the user
	srcRepos, err := scm.FromContext(c).ListUserRepos(u)
	if err != nil {
		retErr := fmt.Errorf("unable to get SCM repos for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// create a map
	// TODO: clean this up
	for _, srepo := range srcRepos {
		// local variables to avoid bad memory address de-referencing
		// initialize active to false
		org := srepo.Org
		name := srepo.Name
		active := false

		// library struct to omit optional fields
		repo := library.Repo{
			Org:    org,
			Name:   name,
			Active: &active,
		}
		output[srepo.GetOrg()] = append(output[srepo.GetOrg()], repo)
	}

	for org := range output {
		// capture source repos from the database backend, grouped by org
		page := 1
		filters := map[string]interface{}{}

		for page > 0 {
			// send API call to capture the list of repos for the org
			dbReposPart, _, err := database.FromContext(c).ListReposForOrg(ctx, org, "name", filters, page, 100)
			if err != nil {
				retErr := fmt.Errorf("unable to get repos for org %s: %w", org, err)

				util.HandleError(c, http.StatusNotFound, retErr)

				return
			}

			// add repos to list of database org repos
			dbRepos = append(dbRepos, dbReposPart...)

			// assume no more pages exist if under 100 results are returned
			if len(dbReposPart) < 100 {
				page = 0
			} else {
				page++
			}
		}

		// apply org repos active status to output map
		for _, dbRepo := range dbRepos {
			if orgRepos, ok := output[dbRepo.GetOrg()]; ok {
				for i := range orgRepos {
					if orgRepos[i].GetName() == dbRepo.GetName() {
						active := dbRepo.GetActive()
						(&orgRepos[i]).Active = &active
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, output)
}
