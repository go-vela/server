// SPDX-License-Identifier: Apache-2.0

package perm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/repo"
	secMiddleware "github.com/go-vela/server/router/middleware/secret"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/secret"
	"github.com/go-vela/server/util"
)

// MustPlatformAdmin ensures the user has admin access to the platform.
func MustPlatformAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		cl := claims.Retrieve(c)

		l.Debugf("verifying user %s is a platform admin", cl.Subject)

		switch {
		case cl.IsAdmin:
			return

		default:
			if strings.EqualFold(cl.TokenType, constants.WorkerBuildTokenType) {
				l.WithFields(logrus.Fields{
					"claims_repo":  cl.Repo,
					"claims_build": cl.BuildID,
				}).Warnf("attempted access of admin endpoint with build token by %s", cl.Subject)
			}

			retErr := fmt.Errorf("user %s is not a platform admin", cl.Subject)
			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}
	}
}

// MustWorkerRegisterToken ensures the token is a registration token retrieved by a platform admin.
func MustWorkerRegisterToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		cl := claims.Retrieve(c)

		l.Debugf("verifying user %s has a registration token for worker", cl.Subject)

		switch cl.TokenType {
		case constants.WorkerRegisterTokenType:
			return
		case constants.ServerWorkerTokenType:
			if strings.EqualFold(cl.Subject, "vela-worker") {
				return
			}

			retErr := fmt.Errorf("server-worker token provided but does not match configuration")
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		default:
			retErr := fmt.Errorf("invalid token type: must provide a worker registration token")
			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}
	}
}

// MustWorkerAuthToken ensures the token is a  worker auth token.
func MustWorkerAuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		cl := claims.Retrieve(c)

		l.Debugf("verifying worker %s has a valid auth token", cl.Subject)

		// global permissions bypass
		if cl.IsAdmin {
			l.Debugf("user %s has platform admin permissions", cl.Subject)

			return
		}

		switch cl.TokenType {
		case constants.WorkerAuthTokenType, constants.WorkerRegisterTokenType:
			return
		case constants.ServerWorkerTokenType:
			if strings.EqualFold(cl.Subject, "vela-worker") {
				return
			}

			retErr := fmt.Errorf("server-worker token provided but does not match configuration")
			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		default:
			retErr := fmt.Errorf("invalid token type: must provide a worker auth token")
			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}
	}
}

// MustBuildAccess ensures the token is a build token for the appropriate build.
func MustBuildAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		cl := claims.Retrieve(c)
		b := build.Retrieve(c)

		// global permissions bypass
		if cl.IsAdmin {
			l.Debugf("user %s has platform admin permissions", cl.Subject)

			return
		}

		l.Debugf("verifying worker %s has a valid build token", cl.Subject)

		// validate token type and match build id in request with build id in token claims
		switch cl.TokenType {
		case constants.WorkerBuildTokenType:
			if b.GetID() == cl.BuildID {
				return
			}

			l.WithFields(logrus.Fields{
				"claims_repo":  cl.Repo,
				"claims_build": cl.BuildID,
			}).Warnf("build token for build %d attempted to be used for build %d by %s", cl.BuildID, b.GetID(), cl.Subject)

			fallthrough
		default:
			retErr := fmt.Errorf("invalid token: must provide matching worker build token")
			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}
	}
}

// MustIDRequestToken ensures the token is a valid ID request token for the appropriate build.
func MustIDRequestToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		cl := claims.Retrieve(c)
		b := build.Retrieve(c)

		logrus.Debugf("verifying worker %s has a valid build token", cl.Subject)

		// verify expected type
		if !strings.EqualFold(cl.TokenType, constants.IDRequestTokenType) {
			retErr := fmt.Errorf("invalid token: must provide a valid request ID token")
			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}

		// if build is not in a running state, then an ID token should not be needed
		if !strings.EqualFold(b.GetStatus(), constants.StatusRunning) {
			util.HandleError(c, http.StatusBadRequest, fmt.Errorf("invalid request"))

			return
		}

		// verify expected build id
		if b.GetID() != cl.BuildID {
			l.WithFields(logrus.Fields{
				"claims_repo":  cl.Repo,
				"claims_build": cl.BuildID,
			}).Warnf("request ID token for build %d attempted to be used for %s build %d by %s", cl.BuildID, b.GetStatus(), b.GetID(), cl.Subject)

			retErr := fmt.Errorf("invalid token")
			util.HandleError(c, http.StatusUnauthorized, retErr)
		}
	}
}

// MustAdmin ensures the user has admin access to the repo.
func MustAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		r := repo.Retrieve(c)
		u := user.Retrieve(c)
		ctx := c.Request.Context()

		l.Debugf("verifying user %s has 'admin' permissions for repo %s", u.GetName(), r.GetFullName())

		if u.GetAdmin() {
			return
		}

		// query source to determine requesters permissions for the repo using the requester's token
		perm, err := scm.FromContext(c).RepoAccess(ctx, u.GetName(), u.GetToken(), r.GetOrg(), r.GetName())
		if err != nil {
			// requester may not have permissions to use the Github API endpoint (requires read access)
			// try again using the repo owner token
			//
			// https://docs.github.com/en/rest/reference/repos#get-repository-permissions-for-a-user
			perm, err = scm.FromContext(c).RepoAccess(ctx, u.GetName(), r.GetOwner().GetToken(), r.GetOrg(), r.GetName())
			if err != nil {
				l.Errorf("unable to get user %s access level for repo %s", u.GetName(), r.GetFullName())
			}
		}

		switch perm {
		//nolint:goconst // ignore making constant
		case "admin":
			return
		default:
			retErr := fmt.Errorf("user %s does not have 'admin' permissions for the repo %s", u.GetName(), r.GetFullName())

			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}
	}
}

// MustWrite ensures the user has admin or write access to the repo.
func MustWrite() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		r := repo.Retrieve(c)
		u := user.Retrieve(c)
		ctx := c.Request.Context()

		l.Debugf("verifying user %s has 'write' permissions for repo %s", u.GetName(), r.GetFullName())

		if u.GetAdmin() {
			return
		}

		// query source to determine requesters permissions for the repo using the requester's token
		perm, err := scm.FromContext(c).RepoAccess(ctx, u.GetName(), u.GetToken(), r.GetOrg(), r.GetName())
		if err != nil {
			// requester may not have permissions to use the Github API endpoint (requires read access)
			// try again using the repo owner token
			//
			// https://docs.github.com/en/rest/reference/repos#get-repository-permissions-for-a-user
			perm, err = scm.FromContext(c).RepoAccess(ctx, u.GetName(), r.GetOwner().GetToken(), r.GetOrg(), r.GetName())
			if err != nil {
				l.Errorf("unable to get user %s access level for repo %s", u.GetName(), r.GetFullName())
			}
		}

		switch perm {
		case "admin":
			return
		case "write":
			return
		default:
			retErr := fmt.Errorf("user %s does not have 'write' permissions for the repo %s", u.GetName(), r.GetFullName())

			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}
	}
}

// MustRead ensures the user has admin, write or read access to the repo.
func MustRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		cl := claims.Retrieve(c)
		r := repo.Retrieve(c)
		u := user.Retrieve(c)
		ctx := c.Request.Context()

		// check if the repo visibility field is set to public
		if strings.EqualFold(r.GetVisibility(), constants.VisibilityPublic) {
			l.Debugf("skipping 'read' check for repo %s with %s visibility for user %s", r.GetFullName(), r.GetVisibility(), u.GetName())

			return
		}

		// return if request is from worker with build token access
		if strings.EqualFold(cl.TokenType, constants.WorkerBuildTokenType) {
			b := build.Retrieve(c)
			if cl.BuildID == b.GetID() {
				return
			}

			retErr := fmt.Errorf("subject %s does not have 'read' permissions for repo %s", cl.Subject, r.GetFullName())

			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}

		l.Debugf("verifying user %s has 'read' permissions for repo %s", u.GetName(), r.GetFullName())

		// return if user is platform admin
		if u.GetAdmin() {
			return
		}

		// query source to determine requesters permissions for the repo using the requester's token
		perm, err := scm.FromContext(c).RepoAccess(ctx, u.GetName(), u.GetToken(), r.GetOrg(), r.GetName())
		if err != nil {
			// requester may not have permissions to use the Github API endpoint (requires read access)
			// try again using the repo owner token
			//
			// https://docs.github.com/en/rest/reference/repos#get-repository-permissions-for-a-user
			perm, err = scm.FromContext(c).RepoAccess(ctx, u.GetName(), r.GetOwner().GetToken(), r.GetOrg(), r.GetName())
			if err != nil {
				l.Errorf("unable to get user %s access level for repo %s", u.GetName(), r.GetFullName())
			}
		}

		switch perm {
		case "admin":
			return
		case "write":
			return
		case "read":
			return
		default:
			retErr := fmt.Errorf("user %s does not have 'read' permissions for repo %s", u.GetName(), r.GetFullName())

			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}
	}
}

// MustSecretAdmin ensures the user has admin access to the org, repo or team.
func MustSecretAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		l := c.MustGet("logger").(*logrus.Entry)
		cl := claims.Retrieve(c)
		u := user.Retrieve(c)
		e := util.PathParameter(c, "engine")
		t := util.PathParameter(c, "type")
		o := util.PathParameter(c, "org")
		n := util.PathParameter(c, "name")
		s := strings.TrimPrefix(util.PathParameter(c, "secret"), "/")
		m := c.Request.Method
		ctx := c.Request.Context()

		// create log fields from API metadata
		fields := logrus.Fields{
			"secret_engine": e,
			"secret_org":    o,
			"secret_repo":   n,
			"secret_type":   t,
		}

		// check if secret is a shared secret
		if strings.EqualFold(t, constants.SecretShared) {
			// update log fields from API metadata
			delete(fields, "repo")
			fields["secret_team"] = n
		}

		logger := l.WithFields(fields)

		if u.GetAdmin() {
			return
		}

		entry := fmt.Sprintf("%s/%s/%s/%s", t, o, n, s)

		var (
			dbSecret *api.Secret
			err      error
		)

		// fetch secret from backend if in path
		if s != "" {
			dbSecret, err = secret.FromContext(c, e).Get(ctx, t, o, n, s)
			if err != nil {
				retErr := fmt.Errorf("unable to get secret %s from %s service: %w", entry, e, err)

				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}

			secMiddleware.ToContext(c, dbSecret)
		}

		// if caller is worker with build token, verify it has access to requested secret
		if strings.EqualFold(cl.TokenType, constants.WorkerBuildTokenType) {
			org, repo := util.SplitFullName(cl.Repo)

			// fetch repo from backend for owner credentials
			dbRepo, err := database.FromContext(c).GetRepoForOrg(ctx, org, repo)
			if err != nil {
				retErr := fmt.Errorf("unable to get repo %s/%s from the database: %w", org, repo, err)

				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}

			switch t {
			case constants.SecretShared:
				return
			case constants.SecretOrg:
				logger.Debugf("verifying subject %s has token permissions for org %s", cl.Subject, o)

				// if org SCMID on secret matches the SCMID from the requested org, return
				if dbRepo.GetOrgSCMID() == dbSecret.GetOrgSCMID() {
					return
				}

				logger.Warnf("build token for build %s/%d attempted to be used for secret %s/%s by %s", cl.Repo, cl.BuildID, o, s, cl.Subject)

				retErr := fmt.Errorf("subject %s does not have token permissions for the org %s", cl.Subject, o)

				util.HandleError(c, http.StatusUnauthorized, retErr)

				return

			case constants.SecretRepo:
				logger.Debugf("verifying subject %s has token permissions for repo %s/%s", cl.Subject, o, n)

				sourceRepo, _, err := scm.FromContext(c).GetRepo(ctx, dbRepo.GetOwner(), dbRepo)
				if err != nil {
					logger.Errorf("unable to get repo %s/%s from the SCM: %v", o, n, err)
				}

				if dbSecret.GetRepoSCMID() == sourceRepo.GetSCMID() {
					return
				}

				logger.Warnf("build token for build %s/%d attempted to be used for secret %s/%s/%s by %s", cl.Repo, cl.BuildID, o, n, s, cl.Subject)

				retErr := fmt.Errorf("subject %s does not have token permissions for the repo %s/%s", cl.Subject, o, n)

				util.HandleError(c, http.StatusUnauthorized, retErr)

				return
			}
		}

		switch t {
		case constants.SecretOrg:
			logger.Debugf("verifying user %s has 'admin' permissions for org %s", u.GetName(), o)

			_, oid, err := scm.FromContext(c).GetOrgIdentifiers(ctx, u, o)
			if err != nil {
				logger.Errorf("unable to get org %s identifiers for user %s: %v", o, u.GetName(), err)
			}

			perm, err := scm.FromContext(c).OrgAccess(ctx, u, o)
			if err != nil {
				logger.Errorf("unable to get user %s access level for org %s: %v", u.GetName(), o, err)
			}

			if !strings.EqualFold(perm, "admin") || (dbSecret.GetID() != 0 && oid != dbSecret.GetOrgSCMID()) {
				retErr := fmt.Errorf("user %s does not have 'admin' permissions for the org %s", u.GetName(), o)

				util.HandleError(c, http.StatusUnauthorized, retErr)

				return
			}
		case constants.SecretRepo:
			logger.Debugf("verifying user %s has 'admin' permissions for repo %s/%s", u.GetName(), o, n)

			sourceRepo, _, err := scm.FromContext(c).GetRepo(ctx, u, &api.Repo{Org: &o, Name: &n})
			if err != nil {
				logger.Errorf("unable to get repo %s/%s from the SCM: %v", o, n, err)
			}

			perm, err := scm.FromContext(c).RepoAccess(ctx, u.GetName(), u.GetToken(), o, n)
			if err != nil {
				logger.Errorf("unable to get user %s access level for repo %s/%s: %v", u.GetName(), o, n, err)
			}

			if !strings.EqualFold(perm, "admin") || (dbSecret.GetID() != 0 && sourceRepo.GetSCMID() != dbSecret.GetRepoSCMID()) {
				retErr := fmt.Errorf("user %s does not have 'admin' permissions for the repo %s/%s", u.GetName(), o, n)

				util.HandleError(c, http.StatusUnauthorized, retErr)

				return
			}
		case constants.SecretShared:
			if n == "*" && m == "GET" {
				// check if user is accessing shared secrets in personal org
				if strings.EqualFold(o, u.GetName()) {
					logger.WithFields(logrus.Fields{
						"secret_org": o,
					}).Debugf("skipping gathering teams for user %s with org %s", u.GetName(), o)

					return
				}

				logger.Debugf("gathering teams user %s is a member of in the org %s", u.GetName(), o)

				teams, err := scm.FromContext(c).ListUsersTeamsForOrg(ctx, u, o)
				if err != nil {
					logger.Errorf("unable to get users %s teams for org %s: %v", u.GetName(), o, err)
				}

				if len(teams) == 0 {
					retErr := fmt.Errorf("user %s is not a member of any team for the org %s", u.GetName(), o)

					util.HandleError(c, http.StatusUnauthorized, retErr)

					return
				}
			} else {
				logger.Debugf("verifying user %s has 'admin' permissions for team %s/%s", u.GetName(), o, n)

				perm, err := scm.FromContext(c).TeamAccess(ctx, u, o, n, dbSecret.GetOrgSCMID(), dbSecret.GetTeamSCMID())
				if err != nil {
					logger.Errorf("unable to get user %s access level for team %s/%s: %v", u.GetName(), o, n, err)
				}

				if !strings.EqualFold(perm, "admin") {
					retErr := fmt.Errorf("user %s does not have 'admin' permissions for the team %s/%s", u.GetName(), o, n)

					util.HandleError(c, http.StatusUnauthorized, retErr)

					return
				}
			}

		default:
			retErr := fmt.Errorf("invalid secret type: %v", t)

			util.HandleError(c, http.StatusBadRequest, retErr)

			return
		}
	}
}
