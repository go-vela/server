// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package perm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// MustPlatformAdmin ensures the user has admin access to the platform.
func MustPlatformAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := user.Retrieve(c)

		// update engine logger with API metadata
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
		logrus.WithFields(logrus.Fields{
			"user": u.GetName(),
		}).Debugf("verifying user %s is a platform admin", u.GetName())

		switch {
		case globalPerms(u):
			return

		default:
			retErr := fmt.Errorf("user %s is not a platform admin", u.GetName())
			util.HandleError(c, http.StatusUnauthorized, retErr)

			return
		}
	}
}

// MustSecretAdmin ensures the user has admin access to the org, repo or team.
func MustSecretAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := user.Retrieve(c)
		e := c.Param("engine")
		t := c.Param("type")
		o := c.Param("org")
		n := c.Param("name")
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

		if globalPerms(u) {
			return
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

		if globalPerms(u) {
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
		// nolint: goconst // ignore making constant
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

		if globalPerms(u) {
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

		logger.Debugf("verifying user %s has 'read' permissions for repo %s", u.GetName(), r.GetFullName())

		if globalPerms(u) {
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

// helper function to check if the user is a platform admin.
func globalPerms(user *library.User) bool {
	switch {
	// Agents have full access to endpoints
	case user.GetName() == "vela-worker":
		return true
	// platform admins have full access to endpoints
	case user.GetAdmin():
		return true
	}

	return false
}
