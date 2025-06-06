// SPDX-License-Identifier: Apache-2.0

package secret

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/api/types/actions"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/util"
)

// swagger:operation POST /api/v1/secrets/{engine}/{type}/{org}/{name} secrets CreateSecret
//
// Create a secret
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: engine
//   description: Secret engine to create a secret in, eg. "native"
//   required: true
//   type: string
// - in: path
//   name: type
//   description: Secret type to create
//   enum:
//   - org
//   - repo
//   - shared
//   required: true
//   type: string
// - in: path
//   name: org
//   description: Name of the organization
//   required: true
//   type: string
// - in: path
//   name: name
//   description: Name of the repository if a repository secret, team name if a shared secret, or '*' if an organization secret
//   required: true
//   type: string
// - in: body
//   name: body
//   description: Secret object to create
//   required: true
//   schema:
//     "$ref": "#/definitions/Secret"
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully created the secret
//     schema:
//       "$ref": "#/definitions/Secret"
//   '400':
//     description: Invalid request payload or path
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// CreateSecret represents the API handler to
// create a secret.
//
//nolint:funlen // suppress long function error
func CreateSecret(c *gin.Context) {
	// capture middleware values
	l := c.MustGet("logger").(*logrus.Entry)
	u := user.Retrieve(c)
	e := util.PathParameter(c, "engine")
	t := util.PathParameter(c, "type")
	o := util.PathParameter(c, "org")
	n := util.PathParameter(c, "name")
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%s/%s", t, o, n)

	// create log fields from API metadata
	fields := logrus.Fields{
		"secret_engine": e,
		"secret_org":    o,
		"secret_repo":   n,
		"secret_type":   t,
	}

	// check if secret is a shared secret
	if strings.EqualFold(t, constants.SecretShared) {
		// update log fields for shared secret
		delete(fields, "secret_repo")
		fields["secret_team"] = n
	}

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logger := l.WithFields(fields)

	if strings.EqualFold(t, constants.SecretOrg) {
		// retrieve org name from SCM
		//
		// SCM can be case insensitive, causing access retrieval to work
		// but Org/Repo != org/repo in Vela. So this check ensures that
		// what a user inputs matches the casing we expect in Vela since
		// the SCM will have the source of truth for casing.
		org, err := scm.FromContext(c).GetOrgName(ctx, u, o)
		if err != nil {
			retErr := fmt.Errorf("unable to retrieve organization %s", o)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		// check if casing is accurate
		if org != o {
			retErr := fmt.Errorf("unable to retrieve organization %s. Did you mean %s?", o, org)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}
	}

	if strings.EqualFold(t, constants.SecretRepo) {
		// retrieve org and repo name from SCM
		//
		// same story as org secret. SCM has accurate casing.
		scmOrg, scmRepo, err := scm.FromContext(c).GetOrgAndRepoName(ctx, u, o, n)
		if err != nil {
			retErr := fmt.Errorf("unable to retrieve repository %s/%s", o, n)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		// check if casing is accurate for org entry
		if scmOrg != o {
			retErr := fmt.Errorf("unable to retrieve org %s. Did you mean %s?", o, scmOrg)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		// check if casing is accurate for repo entry
		if scmRepo != n {
			retErr := fmt.Errorf("unable to retrieve repository %s. Did you mean %s?", n, scmRepo)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}
	}

	logger.Debugf("creating new secret %s for %s service", entry, e)

	// capture body from API request
	input := new(types.Secret)

	err := c.Bind(input)
	if err != nil {
		retErr := fmt.Errorf("unable to decode JSON for secret %s for %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	if len(input.GetRepoAllowlist()) > 0 {
		input.SetRepoAllowlist(util.Unique(input.GetRepoAllowlist()))
	}

	err = validateAllowlist(ctx, database.FromContext(c), []string{}, input)
	if err != nil {
		retErr := fmt.Errorf("invalid allowlist: %w", err)

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// reject secrets with solely whitespace characters as its value
	trimmed := strings.TrimSpace(input.GetValue())
	if len(trimmed) == 0 {
		retErr := fmt.Errorf("secret value must contain non-whitespace characters")

		util.HandleError(c, http.StatusBadRequest, retErr)

		return
	}

	// update fields in secret object
	input.SetOrg(o)
	input.SetRepo(n)
	input.SetType(t)
	input.SetCreatedAt(time.Now().UTC().Unix())
	input.SetCreatedBy(u.GetName())
	input.SetUpdatedAt(time.Now().UTC().Unix())
	input.SetUpdatedBy(u.GetName())

	if len(input.GetImages()) > 0 {
		input.SetImages(util.Unique(input.GetImages()))
	}

	// default event set for secrets
	if input.GetAllowEvents().ToDatabase() == 0 {
		e := new(types.Events)

		push := new(actions.Push)
		push.SetBranch(true)
		push.SetTag(true)

		deploy := new(actions.Deploy)
		deploy.SetCreated(true)

		e.SetPush(push)
		e.SetDeployment(deploy)

		input.SetAllowEvents(e)
	}

	if input.AllowCommand == nil {
		input.SetAllowCommand(true)
	}

	// default to not allow substitution for shared secrets
	if strings.EqualFold(input.GetType(), constants.SecretShared) && input.AllowSubstitution == nil {
		input.SetAllowSubstitution(false)
		input.SetAllowCommand(false)
	} else if input.AllowSubstitution == nil {
		input.SetAllowSubstitution(true)
	}

	// check if secret is a shared secret
	if strings.EqualFold(t, constants.SecretShared) {
		// update the team instead of repo
		input.SetTeam(n)
		input.Repo = nil
	}

	// send API call to create the secret
	s, err := secret.FromContext(c, e).Create(ctx, t, o, n, input)
	if err != nil {
		retErr := fmt.Errorf("unable to create secret %s for %s service: %w", entry, e, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// update log fields from create response
	fields = logrus.Fields{
		"secret_id":   s.GetID(),
		"secret_name": s.GetName(),
		"secret_org":  s.GetOrg(),
		"secret_repo": s.GetRepo(),
		"secret_type": s.GetType(),
	}

	// check if secret is a shared secret
	if strings.EqualFold(t, constants.SecretShared) {
		// update log fields for shared secret
		delete(fields, "secret_repo")
		fields["secret_team"] = s.GetTeam()
	}

	l.WithFields(fields).Infof("created secret %s for %s service", entry, e)

	c.JSON(http.StatusOK, s.Sanitize())
}

// validateAllowlist is a helper function for create/update secret API handlers to verify allowlist configurations.
func validateAllowlist(ctx context.Context, db database.Interface, currentAllowlist []string, s *types.Secret) error {
	switch s.GetType() {
	case constants.SecretRepo:
		if len(s.GetRepoAllowlist()) > 0 {
			return fmt.Errorf("repo secrets cannot have repository allowlists")
		}
	default:
		for _, r := range s.GetRepoAllowlist() {
			split := strings.Split(r, "/")
			if len(split) != 2 {
				return fmt.Errorf("repo %s in allowlist is in an invalid format. Must use `org`/`repo`", r)
			}

			if strings.EqualFold(s.GetType(), constants.SecretOrg) && split[0] != s.GetOrg() {
				return fmt.Errorf("repo %s is not in org %s", r, s.GetOrg())
			}

			if !slices.Contains(currentAllowlist, r) {
				_, err := db.GetRepoForOrg(ctx, split[0], split[1])
				if err != nil {
					return fmt.Errorf("repo %s in allowlist does not exist. Enable the repo in Vela and be sure to check casing", r)
				}
			}
		}
	}

	return nil
}
