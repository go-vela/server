// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"strings"

	"github.com/google/go-github/v70/github"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
)

// OrgAccess captures the user's access level for an org.
func (c *client) OrgAccess(ctx context.Context, u *api.User, org string) (string, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  org,
		"user": u.GetName(),
	}).Tracef("capturing %s access level to org %s", u.GetName(), org)

	// check if user is accessing personal org
	if strings.EqualFold(org, u.GetName()) {
		c.Logger.WithFields(logrus.Fields{
			"org":  org,
			"user": u.GetName(),
		}).Debugf("skipping access level check for user %s with org %s", u.GetName(), org)

		//nolint:goconst // ignore making constant
		return "admin", nil
	}

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, *u.Token)

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
func (c *client) RepoAccess(ctx context.Context, name, token, org, repo string) (string, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  org,
		"repo": repo,
		"user": name,
	}).Tracef("capturing %s access level to repo %s/%s", name, org, repo)

	// check if user is accessing repo in personal org
	if strings.EqualFold(org, name) {
		c.Logger.WithFields(logrus.Fields{
			"org":  org,
			"repo": repo,
			"user": name,
		}).Debugf("skipping access level check for user %s with repo %s/%s", name, org, repo)

		return "admin", nil
	}

	// create github oauth client with the given token
	client := c.newOAuthTokenClient(ctx, token)

	// send API call to capture repo access level for user
	perm, _, err := client.Repositories.GetPermissionLevel(ctx, org, repo, name)
	if err != nil {
		return "", err
	}

	return perm.GetPermission(), nil
}

// TeamAccess captures the user's access level for a team.
func (c *client) TeamAccess(ctx context.Context, u *api.User, org, team string) (string, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  org,
		"team": team,
		"user": u.GetName(),
	}).Tracef("capturing %s access level to team %s/%s", u.GetName(), org, team)

	// check if user is accessing team in personal org
	if strings.EqualFold(org, u.GetName()) {
		c.Logger.WithFields(logrus.Fields{
			"org":  org,
			"team": team,
			"user": u.GetName(),
		}).Debugf("skipping access level check for user %s with team %s/%s", u.GetName(), org, team)

		return "admin", nil
	}

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, u.GetToken())
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
func (c *client) ListUsersTeamsForOrg(ctx context.Context, u *api.User, org string) ([]string, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  org,
		"user": u.GetName(),
	}).Tracef("capturing %s team membership for org %s", u.GetName(), org)

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, u.GetToken())
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

// RepoContributor lists all contributors from a repository and checks if the sender is one of the contributors.
func (c *client) RepoContributor(ctx context.Context, owner *api.User, sender, org, repo string) (bool, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  org,
		"repo": repo,
		"user": sender,
	}).Tracef("capturing %s contributor status for repo %s/%s", sender, org, repo)

	// create GitHub OAuth client with repo owner's token
	client := c.newOAuthTokenClient(ctx, owner.GetToken())

	// set the max per page for the options to capture the list of repos
	opts := github.ListContributorsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100, // 100 is max
		},
	}

	for {
		// send API call to list all contributors for repository
		contributors, resp, err := client.Repositories.ListContributors(ctx, org, repo, &opts)
		if err != nil {
			return false, err
		}

		// match login to sender to see if they are a contributor
		//
		// check this as we page through the results to spare API
		for _, contributor := range contributors {
			if strings.EqualFold(contributor.GetLogin(), sender) {
				return true, nil
			}
		}

		// break the loop if there is no more results to page through
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	return false, nil
}
