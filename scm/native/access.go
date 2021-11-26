// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// OrgAccess captures the user's access level for an org.
func (c *client) OrgAccess(u *library.User, org string) (string, error) {
	logrus.Tracef("Capturing %s access level to org %s", u.GetName(), org)
	return "", nil
}

// RepoAccess captures the user's access level for a repo.
func (c *client) RepoAccess(u *library.User, token, org, repo string) (string, error) {
	logrus.Tracef("Capturing %s access level to repo %s/%s", u.GetName(), org, repo)
	return "", nil
}

// TeamAccess captures the user's access level for a team.
func (c *client) TeamAccess(u *library.User, org, team string) (string, error) {
	logrus.Tracef("Capturing %s access level to team %s/%s", u.GetName(), org, team)
	return "", nil
}

// ListUsersTeamsForOrg captures the user's teams for an org.
func (c *client) ListUsersTeamsForOrg(u *library.User, org string) ([]string, error) {
	logrus.Tracef("Capturing %s team membership for org %s", u.GetName(), org)
	return nil, nil
}
