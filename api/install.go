// SPDX-License-Identifier: Apache-2.0

package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/sirupsen/logrus"
)

// swagger:operation GET /install install Install
//
// Start SCM app installation flow and redirect to the external SCM destination
//
// ---
// produces:
// - application/json
// parameters:
// - in: query
//   name: type
//   description: The type of installation flow, either 'cli' or 'web'
//   type: string
// - in: query
//   name: port
//   description: The local server port used during 'cli' flow
//   type: string
// - in: query
//   name: org_scm_id
//   description: The SCM org id
//   type: string
// - in: query
//   name: repo_scm_id
//   description: The SCM repo id
//   type: string
// responses:
//   '307':
//     description: Redirected for installation
//   '400':
//     description: Invalid request payload
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '503':
//     description: Service unavailable
//     schema:
//       "$ref": "#/definitions/Error"

// Install represents the API handler to
// process an SCM app installation for Vela from
// the API or UI.
func Install(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	scm := scm.FromContext(c)

	l.Debug("redirecting to SCM to complete app flow installation")

	orgSCMID, err := strconv.Atoi(util.FormParameter(c, "org_scm_id"))
	if err != nil {
		retErr := fmt.Errorf("unable to parse org_scm_id to integer: %v", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	repoSCMID, err := strconv.Atoi(util.FormParameter(c, "repo_scm_id"))
	if err != nil {
		retErr := fmt.Errorf("unable to parse repo_scm_id to integer: %v", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// first, check if the org installation exists.
	// if it does, just add the repo manually using the api and be done with it
	// if it doesn't, then we need to start the installation flow
	// but this came from the browser... it has NO auth to contact github api

	// type cannot be empty
	t := util.FormParameter(c, "type")
	if len(t) == 0 {
		retErr := errors.New("no type query provided")

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// port can be empty when using web flow
	p := util.FormParameter(c, "port")

	// capture query params
	ri := &types.RepoInstall{
		OrgSCMID:  int64(orgSCMID),
		RepoSCMID: int64(repoSCMID),
		InstallCallback: types.InstallCallback{
			Type: t,
			Port: p,
		},
	}

	// construct the repo installation url
	redirectURL, err := scm.GetRepoInstallURL(c.Request.Context(), ri)
	if err != nil {
		l.Errorf("unable to get repo install url: %v", err)

		return
	}

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// GetAppInstallRedirectURL is a helper function to generate the redirect URL for completing an app installation flow.
func GetAppInstallRedirectURL(ctx context.Context, l *logrus.Entry, m *internal.Metadata, q url.Values) (string, error) {
	// extract state that is passed along during the installation process
	pairs := strings.Split(q.Get("state"), ",")

	values := make(map[string]string)

	for _, pair := range pairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])

			value := strings.TrimSpace(parts[1])

			values[key] = value
		}
	}

	t, p := values["type"], values["port"]

	// default redirect location if a user ended up here
	// by providing an unsupported type
	// this is ignored when empty
	r := ""

	switch t {
	// cli auth flow
	case "cli":
		r = fmt.Sprintf("http://127.0.0.1:%s", p)
	// web auth flow
	case "web":
		r = fmt.Sprintf("%s%s", m.Vela.WebAddress, m.Vela.WebOauthCallbackPath)
	}

	// append the code and state values
	r = fmt.Sprintf("%s?%s", r, q.Encode())

	l.Debug("redirecting for final app installation flow")

	return r, nil
}
