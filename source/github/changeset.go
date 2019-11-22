// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"fmt"

	"github.com/go-vela/types/library"
)

// Changeset sends the list of files changed for a none pull request event.
func (c *client) Changeset(u *library.User, r *library.Repo, sha string) ([]string, error) {
	// create GitHub OAuth client with user's token
	client := c.newClientToken(u.GetToken())
	s := []string{}

	// send API call to get the commit
	commit, _, err := client.Repositories.GetCommit(ctx, r.GetOrg(), r.GetName(), sha)
	if err != nil {
		return nil, fmt.Errorf("Repositories.GetCommit returned error: %v", err)
	}

	// iterate through each file in the commit
	for _, f := range commit.Files {
		s = append(s, f.GetFilename())
	}

	return s, nil
}

// ChangesetPR sends the list of files changed for a pull request event.
func (c *client) ChangesetPR(u *library.User, r *library.Repo, number int) ([]string, error) {
	s := []string{}
	client := c.newClientToken(u.GetToken())

	// send API call to get the files from the pull request
	files, _, err := client.PullRequests.ListFiles(ctx, r.GetOrg(), r.GetName(), number, nil)
	if err != nil {
		return nil, fmt.Errorf("PullRequests.ListFiles returned error: %v", err)
	}

	// iterate through each file in the pull request
	for _, f := range files {
		s = append(s, f.GetFilename())
	}

	return s, nil
}
