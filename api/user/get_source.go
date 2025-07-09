// SPDX-License-Identifier: Apache-2.0

package user

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/user/source/repos users GetSourceRepos
//
// Get all repos for the current authenticated user
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
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Error"

// GetSourceRepos represents the API handler to get a list of repos for a user.
func GetSourceRepos(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	l.Debugf("reading available SCM repos for user %s", u.GetName())

	// variables to capture requested data
	output := make(map[string][]types.Repo)

	// send API call to capture the list of repos for the user
	srcRepos, err := scm.FromContext(c).ListUserRepos(ctx, u)
	if err != nil {
		retErr := fmt.Errorf("unable to get SCM repos for user %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	dbRepos, err := database.FromContext(c).GetReposInList(ctx, srcRepos)
	if err != nil {
		retErr := fmt.Errorf("unable to get repos from database: %w", err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// create a lookup map to store the active status of each repo
	lookup := make(map[string]bool, len(dbRepos))

	for _, repo := range dbRepos {
		key := repo.GetFullName()
		lookup[key] = repo.GetActive()
	}

	for _, r := range srcRepos {
		// local variables to avoid bad memory address de-referencing
		// initialize active to false
		splitR := strings.Split(r, "/")

		// safety check
		if len(splitR) != 2 {
			continue
		}

		org := splitR[0]
		name := splitR[1]
		active := lookup[r]

		// API struct to omit optional fields
		repo := types.Repo{
			Org:    &org,
			Name:   &name,
			Active: &active,
		}
		output[org] = append(output[org], repo)
	}

	c.JSON(http.StatusOK, output)
}
