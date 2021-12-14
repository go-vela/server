// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"strings"

	"github.com/go-vela/types/library"
	"github.com/google/go-github/v39/github"
)

// OrgAccess captures the user's access level for an org.
func (c *client) OrgAccess(u *library.User, org string) (string, error) {
	c.Logger.Tracef("capturing %s access level to org %s", u.GetName(), org)

	// if user is accessing personal org
	if strings.EqualFold(org, *u.Name) {
		// nolint: goconst // ignore making constant
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
func (c *client) RepoAccess(u *library.User, token, org, repo string) (string, error) {
	c.Logger.Tracef("capturing %s access level to repo %s/%s", u.GetName(), org, repo)

	// create github oauth client with the given token
	client := c.newClientToken(token)

	// send API call to capture repo access level for user
	perm, _, err := client.Repositories.GetPermissionLevel(ctx, org, repo, u.GetName())
	if err != nil {
		return "", err
	}

	return perm.GetPermission(), nil
}

// TeamAccess captures the user's access level for a team.
func (c *client) TeamAccess(u *library.User, org, team string) (string, error) {
	c.Logger.Tracef("capturing %s access level to team %s/%s", u.GetName(), org, team)
	// create GitHub OAuth client with user's token
	client := c.newClientToken(u.GetToken())
	teams := []*github.Team{}

	// set the max per page for the options to capture the list of repos
	opts := github.ListOptions{PerPage: 100} // 100 is max

	for {
		// send API call to list all teams for the user
		uTeams, resp, err := client.Teams.ListUserTeams(ctx, &opts)
		if err != nil {
			return "", err
		}

		teams = append(teams, uTeams...)

		// break the loop if there is no more results to page through
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
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

// ListUsersTeamsForOrg captures the user's teams for an org.
func (c *client) ListUsersTeamsForOrg(u *library.User, org string) ([]string, error) {
	c.Logger.Tracef("capturing %s team membership for org %s", u.GetName(), org)

	// create GitHub OAuth client with user's token
	client := c.newClientToken(u.GetToken())
	teams := []*github.Team{}

	// set the max per page for the options to capture the list of repos
	opts := github.ListOptions{PerPage: 100} // 100 is max

	for {
		// send API call to list all teams for the user
		uTeams, resp, err := client.Teams.ListUserTeams(ctx, &opts)
		if err != nil {
			return []string{""}, err
		}

		teams = append(teams, uTeams...)

		// break the loop if there is no more results to page through
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	var userTeams []string

	// iterate through each element in the teams and filter teams for specified org
	for _, t := range teams {
		// skip the org if does not match the org we are checking
		if strings.EqualFold(org, t.GetOrganization().GetLogin()) {
			userTeams = append(userTeams, t.GetName())
		}
	}

	return userTeams, nil
}
