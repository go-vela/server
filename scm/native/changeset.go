// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// Changeset captures the list of files changed for a commit.
func (c *client) Changeset(u *library.User, r *library.Repo, sha string) ([]string, error) {
	logrus.Tracef("Capturing commit changeset for %s/commit/%s", r.GetFullName(), sha)
	return nil, nil
}

// ChangesetPR captures the list of files changed for a pull request.
func (c *client) ChangesetPR(u *library.User, r *library.Repo, number int) ([]string, error) {
	logrus.Tracef("Capturing pull request changeset for %s/pull/%d", r.GetFullName(), number)
	return nil, nil
}
