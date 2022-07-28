// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"fmt"

	"github.com/go-vela/types/library"
	"github.com/google/go-github/v44/github"
)

// CreateIssue is function that creates an issue in the configured feedback repository
// using a provided description
func (c *client) CreateIssue(u *library.User, owner, repo, title, desc, page string) error {
	body := fmt.Sprintf(
		"## Description\n %s\n\nSubmitted By: @%s\nPage: %s", desc, u.GetName(), page)
	labels := &[]string{"UI_FEEDBACK"}
	issue := &github.IssueRequest{
		Title:  github.String(title),
		Body:   github.String(body),
		Labels: labels,
	}

	// create GitHub OAuth client with user's token
	client := c.newClientToken(*u.Token)

	_, _, err := client.Issues.Create(ctx, owner, repo, issue)
	if err != nil {
		return err
	}

	return nil
}
