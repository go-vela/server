// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"fmt"

	"github.com/go-vela/types/library"
	"github.com/google/go-github/v37/github"

	"github.com/sirupsen/logrus"
)

// Changeset captures the list of files changed for a commit.
func (c *client) Changeset(u *library.User, r *library.Repo, sha string) ([]string, error) {
	logrus.Tracef("Capturing commit changeset for %s/commit/%s", r.GetFullName(), sha)

	// create GitHub OAuth client with user's token
	client := c.newClientToken(u.GetToken())
	s := []string{}

	// send API call to capture the commit
	commit, _, err := client.Repositories.GetCommit(ctx, r.GetOrg(), r.GetName(), sha)
	if err != nil {
		// nolint: golint // ignore capitalized error message
		return nil, fmt.Errorf("Repositories.GetCommit returned error: %v", err)
	}

	// iterate through each file in the commit
	for _, f := range commit.Files {
		s = append(s, f.GetFilename())
	}

	return s, nil
}

// ChangesetPR captures the list of files changed for a pull request.
func (c *client) ChangesetPR(u *library.User, r *library.Repo, number int) ([]string, error) {
	logrus.Tracef("Capturing pull request changeset for %s/pull/%d", r.GetFullName(), number)

	// create GitHub OAuth client with user's token
	client := c.newClientToken(u.GetToken())
	s := []string{}
	f := []*github.CommitFile{}

	// set the max per page for the options to capture the list of repos
	opts := github.ListOptions{PerPage: 100} // 100 is max

	for {
		// send API call to capture the files from the pull request
		files, resp, err := client.PullRequests.ListFiles(ctx, r.GetOrg(), r.GetName(), number, &opts)
		if err != nil {
			// nolint: golint // ignore capitalized error message
			return nil, fmt.Errorf("PullRequests.ListFiles returned error: %v", err)
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
