// SPDX-License-Identifier: Apache-2.0

package repo

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/install/html_url repos GetInstallHTMLURL
//
// Repair a hook for a repository in Vela and the configured SCM
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
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully constructed the repo installation HTML URL
//     schema:
//       type: string
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

// GetInstallHTMLURL represents the API handler to retrieve the
// SCM installation HTML URL for a particular repo and Vela server.
func GetInstallHTMLURL(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*internal.Metadata)
	l := c.MustGet("logger").(*logrus.Entry)
	u := user.Retrieve(c)
	r := repo.Retrieve(c)
	scm := scm.FromContext(c)

	l.Debug("constructing repo install url")

	ri, err := scm.GetRepoInstallInfo(c.Request.Context(), u, r.GetOrg(), r.GetName())
	if err != nil {
		retErr := fmt.Errorf("unable to get repo scm install info %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// todo: use url.values etc
	appInstallURL := fmt.Sprintf(
		"%s/install?org_scm_id=%d&repo_scm_id=%d",
		m.Vela.Address,
		ri.OrgSCMID, ri.RepoSCMID,
	)

	c.JSON(http.StatusOK, fmt.Sprintf("%s", appInstallURL))
}
