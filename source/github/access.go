// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"strings"

	"github.com/go-vela/types/library"
)

// OrgAccess captures the user's access level for an org.
func (c *client) OrgAccess(u *library.User, org string) (string, error) {
	// if user is accessing personal org
	if strings.EqualFold(org, *u.Name) {
		return "admin", nil
	}

	// create GitHub OAuth client with user's token
	client := c.newClientToken(*u.Token)

	// send API call to capture org access level for user
	membership, _, err := client.Organizations.GetOrgMembership(ctx, *u.Name, org)
	if err != nil {
		return "", err
	}

	// return their access level if they are an active user
	if membership.GetState() == "active" {
		return membership.GetRole(), nil
	}

	return "", nil
}

// RepoAccess captures the user's access level for a repo.
func (c *client) RepoAccess(u *library.User, org, repo string) (string, error) {
	// create github oauth client with user's token
	client := c.newClientToken(*u.Token)

	// send API call to capture repo access level for user
	perm, _, err := client.Repositories.GetPermissionLevel(ctx, org, repo, u.GetName())
	if err != nil {
		return "", err
	}

	return perm.GetPermission(), nil
}

// TeamAccess captures the user's access level for a team.
func (c *client) TeamAccess(u *library.User, org, team string) (string, error) {
	// create GitHub OAuth client with user's token
	client := c.newClientToken(u.GetToken())

	// send API call to list all teams for the user
	teams, _, err := client.Teams.ListUserTeams(ctx, nil)
	if err != nil {
		return "", err
	}

	// iterate through each element in the teams
	for _, t := range teams {
		// skip the team if does not match the team we are checking
		if !strings.EqualFold(team, t.GetName()) {
			continue
		}

		// skip the org if does not match the org we are checking
		if !strings.EqualFold(org, t.GetOrganization().GetLogin()) {
			continue
		}

		// return admin access if the user is a part of that team
		return "admin", nil
	}

	return "", nil
}
