// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"

	"github.com/go-vela/types/library"
	"github.com/jenkins-x/go-scm/scm"

	"github.com/sirupsen/logrus"
)

// Changeset captures the list of files changed for a commit.
func (c *client) Changeset(u *library.User, r *library.Repo, sha string) ([]string, error) {
	logrus.Tracef("Capturing commit changeset for %s/commit/%s", r.GetFullName(), sha)

	// create GitHub OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return nil, err
	}

	// send API call to capture the commit
	changes, _, err := client.Git.ListChanges(ctx, r.GetFullName(), sha, scm.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("Repositories.ListChanges returned error: %v", err)
	}

	s := []string{}

	// iterate through each file in the commit
	for _, change := range changes {
		s = append(s, change.Path)
	}

	return s, nil
}

// ChangesetPR captures the list of files changed for a pull request.
func (c *client) ChangesetPR(u *library.User, r *library.Repo, number int) ([]string, error) {
	logrus.Tracef("Capturing pull request changeset for %s/pull/%d", r.GetFullName(), number)

	// create GitHub OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return nil, err
	}

	changes := []*scm.Change{}
	opts := scm.ListOptions{
		Size: 100,
	}

	for {
		// send API call to capture the files from the pull request
		PRChanges, resp, err := client.PullRequests.ListChanges(ctx, r.GetFullName(), number, opts)
		if err != nil {
			return nil, fmt.Errorf("PullRequests.ListChanges returned error: %v", err)
		}

		changes = append(changes, PRChanges...)

		// break the loop if there is no more results to page through
		if resp.Page.Next == 0 {
			break
		}

		opts.Page = resp.Page.Next
	}

	s := []string{}

	// iterate through each file in the commit
	for _, change := range changes {
		s = append(s, change.Path)
	}

	return s, nil
}
