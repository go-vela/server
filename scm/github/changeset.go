// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/v73/github"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
)

// Changeset captures the list of files changed for a commit.
func (c *Client) Changeset(ctx context.Context, r *api.Repo, sha string) ([]string, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": r.GetOwner().GetName(),
	}).Tracef("capturing commit changeset for %s/commit/%s", r.GetFullName(), sha)

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, r.GetOwner().GetToken())
	s := []string{}

	// set the max per page for the options to capture the commit
	opts := github.ListOptions{PerPage: 100} // 100 is max

	// send API call to capture the commit
	commit, _, err := client.Repositories.GetCommit(ctx, r.GetOrg(), r.GetName(), sha, &opts)
	if err != nil {
		return nil, fmt.Errorf("Repositories.GetCommit returned error: %w", err)
	}

	// iterate through each file in the commit
	for _, f := range commit.Files {
		s = append(s, f.GetFilename())
	}

	return s, nil
}

// ChangesetPR captures the list of files changed for a pull request.
func (c *Client) ChangesetPR(ctx context.Context, r *api.Repo, number int) ([]string, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": r.GetOwner().GetName(),
	}).Tracef("capturing pull request changeset for %s/pull/%d", r.GetFullName(), number)

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, r.GetOwner().GetToken())
	s := []string{}
	f := []*github.CommitFile{}

	// set the max per page for the options to capture the list of repos
	opts := github.ListOptions{PerPage: 100} // 100 is max

	for {
		// send API call to capture the files from the pull request
		files, resp, err := client.PullRequests.ListFiles(ctx, r.GetOrg(), r.GetName(), number, &opts)
		if err != nil {
			return nil, fmt.Errorf("PullRequests.ListFiles returned error: %w", err)
		}

		f = append(f, files...)

		// break the loop if there is no more results to page through
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	// iterate through each file in the pull request
	for _, file := range f {
		s = append(s, file.GetFilename())
	}

	return s, nil
}
