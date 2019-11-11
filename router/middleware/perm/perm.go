// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package perm

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/source"
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

		logrus.Debugf("Verifying user %s is a platform admin", u.GetName())
		switch {
		case globalPerms(u):
			return

		default:
			retErr := fmt.Errorf("User %s is not a platform admin", u.GetName())
			util.HandleError(c, http.StatusUnauthorized, retErr)
			return
		}
	}
}

// MustSecretAdmin ensures the user has admin access to the org, repo or team.
func MustSecretAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		u := user.Retrieve(c)

		if globalPerms(u) {
			return
		}

		t := c.Param("type")
		o := c.Param("org")
		n := c.Param("name")

		switch t {
		case constants.SecretOrg:
			logrus.Debugf("Verifying user %s has 'admin' permissions for org %s", u.GetName(), o)
			perm, err := source.FromContext(c).OrgAccess(u, o)
			if err != nil {
				logrus.Errorf("Unable to get user %s access level for org %s: %v", u.GetName(), o, err)
			}

			if !strings.EqualFold(perm, "admin") {
				retErr := fmt.Errorf("User %s does not have 'admin' permissions for the org %s", u.GetName(), o)
				util.HandleError(c, http.StatusUnauthorized, retErr)
				return
			}
		case constants.SecretRepo:
			logrus.Debugf("Verifying user %s has 'admin' permissions for repo %s/%s", u.GetName(), o, n)
			perm, err := source.FromContext(c).RepoAccess(u, o, n)
			if err != nil {
				logrus.Errorf("Unable to get user %s access level for repo %s/%s: %v", u.GetName(), o, n, err)
			}

			if !strings.EqualFold(perm, "admin") {
				retErr := fmt.Errorf("User %s does not have 'admin' permissions for the repo %s/%s", u.GetName(), o, n)
				util.HandleError(c, http.StatusUnauthorized, retErr)
				return
			}
		case constants.SecretShared:
			logrus.Debugf("Verifying user %s has 'admin' permissions for team %s/%s", u.GetName(), o, n)
			perm, err := source.FromContext(c).TeamAccess(u, o, n)
			if err != nil {
				logrus.Errorf("Unable to get user %s access level for team %s/%s: %v", u.GetName(), o, n, err)
			}

			if !strings.EqualFold(perm, "admin") {
				retErr := fmt.Errorf("User %s does not have 'admin' permissions for the team %s/%s", u.GetName(), o, n)
				util.HandleError(c, http.StatusUnauthorized, retErr)
				return
			}
		default:
			retErr := fmt.Errorf("Invalid secret type: %v", t)
			util.HandleError(c, http.StatusBadRequest, retErr)
			return
		}
	}
}

// MustAdmin ensures the user has admin access to the repo.
func MustAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := repo.Retrieve(c)
		u := user.Retrieve(c)

		logrus.Debugf("Verifying user %s has 'admin' permissions for repo %s", u.GetName(), r.GetFullName())
		if globalPerms(u) {
			return
		}

		perm, err := source.FromContext(c).RepoAccess(u, r.GetOrg(), r.GetName())
		if err != nil {
			logrus.Errorf("Unable to get user %s access level for repo %s", u.GetName(), r.GetFullName())
		}

		switch {
		case perm == "admin":
			return

		default:
			retErr := fmt.Errorf("User %s does not have 'admin' permissions for the repo %s", u.GetName(), r.GetFullName())
			util.HandleError(c, http.StatusUnauthorized, retErr)
			return
		}
	}
}

// MustWrite ensures the user has admin or write access to the repo.
func MustWrite() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := repo.Retrieve(c)
		u := user.Retrieve(c)

		logrus.Debugf("Verifying user %s has 'write' permissions for repo %s", u.GetName(), r.GetFullName())
		if globalPerms(u) {
			return
		}

		perm, err := source.FromContext(c).RepoAccess(u, r.GetOrg(), r.GetName())
		if err != nil {
			logrus.Errorf("Unable to get user %s access level for repo %s", u.GetName(), r.GetFullName())
		}

		switch {
		case perm == "admin":
			return

		case perm == "write":
			return

		default:
			retErr := fmt.Errorf("User %s does not have 'write' permissions for the repo %s", u.GetName(), r.GetFullName())
			util.HandleError(c, http.StatusUnauthorized, retErr)
			return
		}
	}
}

// MustRead ensures the user has admin, write or read access to the repo.
func MustRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := repo.Retrieve(c)
		u := user.Retrieve(c)

		logrus.Debugf("Verifying user %s has 'read' permissions for repo %s", u.GetName(), r.GetFullName())
		if globalPerms(u) {
			return
		}

		perm, err := source.FromContext(c).RepoAccess(u, r.GetOrg(), r.GetName())
		if err != nil {
			logrus.Errorf("Unable to get user %s access level for repo %s", u.GetName(), r.GetFullName())
		}

		switch {
		case perm == "admin":
			return

		case perm == "write":
			return

		case perm == "read":
			return

		default:
			retErr := fmt.Errorf("User %s does not have 'read' permissions for repo %s", u.GetName(), r.GetFullName())
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
