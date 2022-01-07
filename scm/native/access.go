// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"
	"strings"

	"github.com/go-vela/types/library"
	"github.com/jenkins-x/go-scm/scm"

	"github.com/sirupsen/logrus"
)

// OrgAccess captures the user's access level for an org.
func (c *client) OrgAccess(u *library.User, org string) (string, error) {
	logrus.Tracef("Capturing %s access level to org %s", u.GetName(), org)

	// create SCM OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return "", err
	}

	opts := scm.ListOptions{}

	for {
		// send API call to capture org access level for user
		memberships, resp, err := client.Organizations.ListMemberships(ctx, opts)
		if err != nil {
			return "", err
		}

		for _, membership := range memberships {
			// return their access level if they are an active user
			if strings.EqualFold(membership.State, "active") &&
				strings.EqualFold(membership.OrganizationName, org) {
				return membership.Role, nil
			}
		}

		// break the loop if there is no more results to page through
		if resp.Page.Next == 0 {
			break
		}

		opts.Page = resp.Page.Next
	}

	return "", nil
}

// RepoAccess captures the user's access level for a repo.
func (c *client) RepoAccess(u *library.User, token, org, repo string) (string, error) {
	logrus.Tracef("Capturing %s access level to repo %s/%s", u.GetName(), org, repo)

	// create SCM OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return "", err
	}

	perm, _, err := client.Repositories.FindUserPermission(ctx,
		fmt.Sprintf("%s/%s", org, repo), u.GetName())
	if err != nil {
		return "", err
	}

	return perm, nil
}

// TeamAccess captures the user's access level for a team.
func (c *client) TeamAccess(u *library.User, org, team string) (string, error) {
	logrus.Tracef("Capturing %s access level to team %s/%s", u.GetName(), org, team)

	// create SCM OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return "", err
	}

	opts := scm.ListOptions{}
	teams := []*scm.Team{}

	for {
		// send API call to capture org access level for user
		uTeams, resp, err := client.Organizations.ListTeams(ctx, org, opts)
		if err != nil || resp == nil {
			return "", err
		}

		teams = append(teams, uTeams...)

		// break the loop if there is no more results to page through
		if resp.Page.Next == 0 {
			break
		}

		opts.Page = resp.Page.Next
	}

	for _, t := range teams {
		if strings.EqualFold(team, t.Name) {
			members, _, err := client.Organizations.ListTeamMembers(ctx, t.ID, "all", scm.ListOptions{})
			if err != nil {
				return "", err
			}

			for _, member := range members {
				if strings.EqualFold(u.GetName(), member.Login) {
					if member.IsAdmin {
						// nolint:goconst // do not make this a constant
						return "admin", nil
					}
					return "write", nil
				}
			}
		}
	}

	return "", nil
}

// ListUsersTeamsForOrg captures the user's teams for an org.
func (c *client) ListUsersTeamsForOrg(u *library.User, org string) ([]string, error) {
	logrus.Tracef("Capturing %s team membership for org %s", u.GetName(), org)

	// create SCM OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return nil, err
	}

	opts := scm.ListOptions{}
	teams := []*scm.Team{}

	for {
		// send API call to capture org access level for user
		uTeams, resp, err := client.Organizations.ListTeams(ctx, org, opts)
		if err != nil || resp == nil {
			return nil, err
		}

		teams = append(teams, uTeams...)

		// break the loop if there is no more results to page through
		if resp.Page.Next == 0 {
			break
		}

		opts.Page = resp.Page.Next
	}

	uTeams := []string{}

	// filter out the names of the orgs
	for _, t := range teams {
		uTeams = append(uTeams, t.Name)
	}

	return uTeams, nil
}
