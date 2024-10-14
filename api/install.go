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
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/util"
	"github.com/google/go-github/v62/github"
	"github.com/sirupsen/logrus"
)

// HandleInstallCallback represents the API handler to
// process an SCM app installation for Vela.
func HandleInstallCallback(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*internal.Metadata)
	l := c.MustGet("logger").(*logrus.Entry)

	ctx := c.Request.Context()

	// GitHub App and OAuth share the same callback URL,
	// so we need to differentiate between the two using setup_action
	setupAction := c.Request.FormValue("setup_action")
	switch setupAction {
	case "install":
	case "update":
		installID := c.Request.FormValue("installation_id")
		if len(installID) == 0 {
			retErr := errors.New("setup_action is install but installation_id is missing")

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		// todo: if the repo is already added, then redirecting to the install url will try to add ALL repos...

		// todo: on "install" we also need to check if it was just a regular github ui manual installation
		// todo: on "update" this might just be a regular ui update to the github app
		// todo: we need to capture the installation ID and sync all the vela repos for that installation
		redirect, err := GetAppInstallRedirectURL(ctx, l, m, c.Request.URL.Query())
		if err != nil {
			retErr := fmt.Errorf("unable to get app install redirect URL: %w", err)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}

		if len(redirect) == 0 {
			c.JSON(http.StatusOK, "installation completed")

			return
		}

		c.Redirect(http.StatusTemporaryRedirect, redirect)

		return
	case "":
		break
	}
}

// Install represents the API handler to
// process an SCM app installation for Vela from
// the API or UI.
func Install(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	// scm := scm.FromContext(c)

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
	// todo: this would need a github client
	redirectURL, err := GetRepoInstallURL(c.Request.Context(), nil, ri)
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

// GetInstallInfo represents the API handler to retrieve the
// SCM installation HTML URL for a particular repo and Vela server.
func GetInstallInfo(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*internal.Metadata)
	l := c.MustGet("logger").(*logrus.Entry)
	u := user.Retrieve(c)
	r := repo.Retrieve(c)
	// scm := scm.FromContext(c)

	l.Debug("retrieving repo install information")

	// todo: this would need github clients
	ri, err := GetRepoInstallInfo(c.Request.Context(), nil, nil, u, r)
	if err != nil {
		retErr := fmt.Errorf("unable to get repo scm install info %s: %w", u.GetName(), err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// todo: use url.values etc
	ri.InstallURL = fmt.Sprintf(
		"%s/install?org_scm_id=%d&repo_scm_id=%d",
		m.Vela.Address,
		ri.OrgSCMID, ri.RepoSCMID,
	)

	c.JSON(http.StatusOK, ri)
}

// GetRepoInstallInfo retrieves the repo information required for installation, such as org and repo ID for the given org and repo name.
func GetRepoInstallInfo(ctx context.Context, userClient *github.Client, appClient *github.Client, u *types.User, r *types.Repo) (*types.RepoInstall, error) {
	// client := c.newClientToken(ctx, u.GetToken())

	// send an API call to get the org info
	repoInfo, resp, err := userClient.Repositories.Get(ctx, r.GetOrg(), r.GetName())

	orgID := repoInfo.GetOwner().GetID()

	// if org is not found, return the personal org
	if resp.StatusCode == http.StatusNotFound {
		user, _, err := userClient.Users.Get(ctx, "")
		if err != nil {
			return nil, err
		}
		orgID = user.GetID()
	} else if err != nil {
		return nil, err
	}

	ri := &types.RepoInstall{
		OrgSCMID:  orgID,
		RepoSCMID: repoInfo.GetID(),
	}

	// todo: pagination...
	installations, resp, err := appClient.Apps.ListInstallations(ctx, &github.ListOptions{})
	if err != nil && (resp == nil || resp.StatusCode != http.StatusNotFound) {
		return nil, err
	}

	// check if the app is installed on the org
	var id int64
	for _, installation := range installations {
		// app is installed to the org
		if installation.GetAccount().GetID() == orgID {
			ri.AppInstalled = true
			ri.InstallID = installation.GetID()
		}
	}

	// todo: remove all this, it doesnt work without a PAT, lol
	_, _, err = appClient.Apps.AddRepository(ctx, id, repoInfo.GetID())
	if err != nil {
		return nil, err
	}

	return ri, nil
}

// GetRepoInstallURL takes RepoInstall configurations and returns the SCM URL for installing the application.
func GetRepoInstallURL(ctx context.Context, appClient *github.Client, ri *types.RepoInstall) (string, error) {
	// retrieve the authenticated app information
	// required for slug and HTML URL
	app, _, err := appClient.Apps.Get(ctx, "")
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf(
		"%s/installations/new/permissions",
		app.GetHTMLURL())

	// stored as state to retrieve from the post-install callback
	state := fmt.Sprintf("type=%s,port=%s", ri.Type, ri.Port)

	v := &url.Values{}
	v.Set("state", state)
	v.Set("suggested_target_id", strconv.FormatInt(ri.OrgSCMID, 10))
	v.Set("repository_ids", strconv.FormatInt(ri.RepoSCMID, 10))

	return fmt.Sprintf("%s?%s", path, v.Encode()), nil
}
