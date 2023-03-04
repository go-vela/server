// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package perm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/claims"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"
)

// MustPlatformAdmin ensures the user has admin access to the platform.
func MustPlatformAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		cl := claims.Retrieve(c)

		// update engine logger with API metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
		logrus.WithFields(logrus.Fields{
			"user": cl.Subject,
		}).Debugf("verifying user %s is a platform admin", cl.Subject)

		switch {
		case cl.IsAdmin:
			return

		default:
			if strings.EqualFold(cl.TokenType, constants.WorkerBuildTokenType) {
				logrus.WithFields(logrus.Fields{
					"user":  cl.Subject,
					"repo":  cl.Repo,
					"build": cl.BuildID,
				}).Warnf("attempted access of admin endpoint with build token from %s", cl.Subject)
			}

			retErr := fmt.Errorf("user %s is not a platform admin", cl.Subject)
			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}
	}
}

// MustWorker ensures the request is coming from an agent.
func MustWorker() gin.HandlerFunc {
	return func(c *gin.Context) {
		cl := claims.Retrieve(c)

		// global permissions bypass
		if cl.IsAdmin {
			logrus.WithFields(logrus.Fields{
				"user": cl.Subject,
			}).Debugf("user %s has platform admin permissions", cl.Subject)

			return
		}

		// update engine logger with API metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
		logrus.WithFields(logrus.Fields{
			"subject": cl.Subject,
		}).Debugf("verifying user %s is a worker", cl.Subject)

		// validate claims as worker
		switch {
		case strings.EqualFold(cl.Subject, "vela-worker") && strings.EqualFold(cl.TokenType, constants.ServerWorkerTokenType):
			return

		default:
			retErr := fmt.Errorf("user %s is not a worker", cl.Subject)
			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}
	}
}

// MustBuildAccess ensures the token is a build token for the appropriate build.
func MustBuildAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		cl := claims.Retrieve(c)
		b := build.Retrieve(c)

		// global permissions bypass
		if cl.IsAdmin {
			logrus.WithFields(logrus.Fields{
				"user": cl.Subject,
			}).Debugf("user %s has platform admin permissions", cl.Subject)

			return
		}

		// update engine logger with API metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
		logrus.WithFields(logrus.Fields{
			"worker": cl.Subject,
		}).Debugf("verifying worker %s has a valid build token", cl.Subject)

		// validate token type and match build id in request with build id in token claims
		switch cl.TokenType {
		case constants.WorkerBuildTokenType:
			if b.GetID() == cl.BuildID {
				return
			}

			logrus.WithFields(logrus.Fields{
				"user":  cl.Subject,
				"repo":  cl.Repo,
				"build": cl.BuildID,
			}).Warnf("build token for build %d attempted to be used for build %d by %s", cl.BuildID, b.GetID(), cl.Subject)

			fallthrough
		default:
			retErr := fmt.Errorf("invalid token: must provide matching worker build token")
			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}
	}
}

// MustSecretAdmin ensures the user has admin access to the org, repo or team.
func MustSecretAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		cl := claims.Retrieve(c)
		u := user.Retrieve(c)
		e := util.PathParameter(c, "engine")
		t := util.PathParameter(c, "type")
		o := util.PathParameter(c, "org")
		n := util.PathParameter(c, "name")
		s := util.PathParameter(c, "secret")
		m := c.Request.Method

		// create log fields from API metadata
		fields := logrus.Fields{
			"engine": e,
			"org":    o,
			"repo":   n,
			"type":   t,
			"user":   u.GetName(),
		}

		// check if secret is a shared secret
		if strings.EqualFold(t, constants.SecretShared) {
			// update log fields from API metadata
			fields = logrus.Fields{
				"engine": e,
				"org":    o,
				"team":   n,
				"type":   t,
				"user":   u.GetName(),
			}
		}

		// update engine logger with API metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
		logger := logrus.WithFields(fields)

		if u.GetAdmin() {
			return
		}

		// if caller is worker with build token, verify it has access to requested secret
		if strings.EqualFold(cl.TokenType, constants.WorkerBuildTokenType) {
			// split repo full name into org and repo
			repoSlice := strings.Split(cl.Repo, "/")
			if len(repoSlice) != 2 {
				logger.Errorf("unable to parse repo claim in build token")
			}

			org := repoSlice[0]
			repo := repoSlice[1]

			switch t {
			case constants.SecretShared:
				return
			case constants.SecretOrg:
				logger.Debugf("verifying subject %s has token permissions for org %s", cl.Subject, o)

				if strings.EqualFold(org, o) {
					return
				}

				logger.Warnf("build token for build %s/%d attempted to be used for secret %s/%s by %s", cl.Repo, cl.BuildID, o, s, cl.Subject)

				retErr := fmt.Errorf("subject %s does not have token permissions for the org %s", cl.Subject, o)

				util.HandleError(c, http.StatusUnauthorized, retErr)

				return

			case constants.SecretRepo:
				logger.Debugf("verifying subject %s has token permissions for repo %s/%s", cl.Subject, o, n)

				if strings.EqualFold(org, o) && strings.EqualFold(repo, n) {
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

			perm, err := scm.FromContext(c).OrgAccess(u, o)
			if err != nil {
				logger.Errorf("unable to get user %s access level for org %s: %v", u.GetName(), o, err)
			}

			if !strings.EqualFold(perm, "admin") {
				retErr := fmt.Errorf("user %s does not have 'admin' permissions for the org %s", u.GetName(), o)

				util.HandleError(c, http.StatusUnauthorized, retErr)

				return
			}
		case constants.SecretRepo:
			logger.Debugf("verifying user %s has 'admin' permissions for repo %s/%s", u.GetName(), o, n)

			perm, err := scm.FromContext(c).RepoAccess(u, u.GetToken(), o, n)
			if err != nil {
				logger.Errorf("unable to get user %s access level for repo %s/%s: %v", u.GetName(), o, n, err)
			}

			if !strings.EqualFold(perm, "admin") {
				retErr := fmt.Errorf("user %s does not have 'admin' permissions for the repo %s/%s", u.GetName(), o, n)

				util.HandleError(c, http.StatusUnauthorized, retErr)

				return
			}
		case constants.SecretShared:
			if n == "*" && m == "GET" {
				// check if user is accessing shared secrets in personal org
				if strings.EqualFold(o, u.GetName()) {
					logger.WithFields(logrus.Fields{
						"org":  o,
						"user": u.GetName(),
					}).Warnf("skipping gathering teams for user %s with org %s", u.GetName(), o)
					return
				}

				logger.Debugf("gathering teams user %s is a member of in the org %s", u.GetName(), o)

				teams, err := scm.FromContext(c).ListUsersTeamsForOrg(u, o)
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

				perm, err := scm.FromContext(c).TeamAccess(u, o, n)
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

// MustAdmin ensures the user has admin access to the repo.
func MustAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		o := org.Retrieve(c)
		r := repo.Retrieve(c)
		u := user.Retrieve(c)

		// update engine logger with API metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
		logger := logrus.WithFields(logrus.Fields{
			"org":  o,
			"repo": r.GetName(),
			"user": u.GetName(),
		})

		logger.Debugf("verifying user %s has 'admin' permissions for repo %s", u.GetName(), r.GetFullName())

		if u.GetAdmin() {
			return
		}

		// query source to determine requesters permissions for the repo using the requester's token
		perm, err := scm.FromContext(c).RepoAccess(u, u.GetToken(), r.GetOrg(), r.GetName())
		if err != nil {
			// requester may not have permissions to use the Github API endpoint (requires read access)
			// try again using the repo owner token
			//
			// https://docs.github.com/en/rest/reference/repos#get-repository-permissions-for-a-user
			ro, err := database.FromContext(c).GetUser(r.GetUserID())
			if err != nil {
				retErr := fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err)

				util.HandleError(c, http.StatusBadRequest, retErr)

				return
			}

			perm, err = scm.FromContext(c).RepoAccess(u, ro.GetToken(), r.GetOrg(), r.GetName())
			if err != nil {
				logger.Errorf("unable to get user %s access level for repo %s", u.GetName(), r.GetFullName())
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
		o := org.Retrieve(c)
		r := repo.Retrieve(c)
		u := user.Retrieve(c)

		// update engine logger with API metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
		logger := logrus.WithFields(logrus.Fields{
			"org":  o,
			"repo": r.GetName(),
			"user": u.GetName(),
		})

		logger.Debugf("verifying user %s has 'write' permissions for repo %s", u.GetName(), r.GetFullName())

		if u.GetAdmin() {
			return
		}

		// query source to determine requesters permissions for the repo using the requester's token
		perm, err := scm.FromContext(c).RepoAccess(u, u.GetToken(), r.GetOrg(), r.GetName())
		if err != nil {
			// requester may not have permissions to use the Github API endpoint (requires read access)
			// try again using the repo owner token
			//
			// https://docs.github.com/en/rest/reference/repos#get-repository-permissions-for-a-user
			ro, err := database.FromContext(c).GetUser(r.GetUserID())
			if err != nil {
				retErr := fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err)

				util.HandleError(c, http.StatusBadRequest, retErr)

				return
			}

			perm, err = scm.FromContext(c).RepoAccess(u, ro.GetToken(), r.GetOrg(), r.GetName())
			if err != nil {
				logger.Errorf("unable to get user %s access level for repo %s", u.GetName(), r.GetFullName())
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
		cl := claims.Retrieve(c)
		o := org.Retrieve(c)
		r := repo.Retrieve(c)
		u := user.Retrieve(c)

		// update engine logger with API metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
		logger := logrus.WithFields(logrus.Fields{
			"org":  o,
			"repo": r.GetName(),
			"user": u.GetName(),
		})

		// check if the repo visibility field is set to public
		if strings.EqualFold(r.GetVisibility(), constants.VisibilityPublic) {
			logger.Debugf("skipping 'read' check for repo %s with %s visibility for user %s", r.GetFullName(), r.GetVisibility(), u.GetName())

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

		logger.Debugf("verifying user %s has 'read' permissions for repo %s", u.GetName(), r.GetFullName())

		// return if user is platform admin
		if u.GetAdmin() {
			return
		}

		// query source to determine requesters permissions for the repo using the requester's token
		perm, err := scm.FromContext(c).RepoAccess(u, u.GetToken(), r.GetOrg(), r.GetName())
		if err != nil {
			// requester may not have permissions to use the Github API endpoint (requires read access)
			// try again using the repo owner token
			//
			// https://docs.github.com/en/rest/reference/repos#get-repository-permissions-for-a-user
			ro, err := database.FromContext(c).GetUser(r.GetUserID())
			if err != nil {
				retErr := fmt.Errorf("unable to get owner for %s: %w", r.GetFullName(), err)

				util.HandleError(c, http.StatusBadRequest, retErr)

				return
			}

			perm, err = scm.FromContext(c).RepoAccess(u, ro.GetToken(), r.GetOrg(), r.GetName())
			if err != nil {
				logger.Errorf("unable to get user %s access level for repo %s", u.GetName(), r.GetFullName())
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
