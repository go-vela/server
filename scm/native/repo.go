// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"github.com/google/go-github/v39/github"

	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// ConfigBackoff is a wrapper for Config that will retry five times if the function
// fails to retrieve the yaml/yml file.
// nolint: lll // ignore long line length due to input arguments
func (c *client) ConfigBackoff(u *library.User, r *library.Repo, ref string) (data []byte, err error) {
	return
}

// Config gets the pipeline configuration from the GitHub repo.
func (c *client) Config(u *library.User, r *library.Repo, ref string) ([]byte, error) {
	logrus.Tracef("Capturing configuration file for %s/%s/commit/%s", r.GetOrg(), r.GetName(), ref)
	return nil, nil
}

// Disable deactivates a repo by deleting the webhook.
func (c *client) Disable(u *library.User, org, name string) error {
	logrus.Tracef("Deleting repository webhook for %s/%s", org, name)
	return nil
}

// Enable activates a repo by creating the webhook.
func (c *client) Enable(u *library.User, org, name, secret string) (string, error) {
	logrus.Tracef("Creating repository webhook for %s/%s", org, name)
	return "", nil
}

// Status sends the commit status for the given SHA from the GitHub repo.
func (c *client) Status(u *library.User, b *library.Build, org, name string) error {
	logrus.Tracef("Setting commit status for %s/%s/%d @ %s", org, name, b.GetNumber(), b.GetCommit())
	return nil
}

// GetRepo gets repo information from Github.
func (c *client) GetRepo(u *library.User, r *library.Repo) (*library.Repo, error) {
	logrus.Tracef("Retrieving repository information for %s", r.GetFullName())
	return nil, nil
}

// ListUserRepos returns a list of all repos the user has access to.
func (c *client) ListUserRepos(u *library.User) ([]*library.Repo, error) {
	logrus.Tracef("Listing source repositories for %s", u.GetName())
	return nil, nil
}

// toLibraryRepo does a partial conversion of a github repo to a library repo.
func toLibraryRepo(gr github.Repository) *library.Repo {
	return &library.Repo{
		Org:      gr.GetOwner().Login,
		Name:     gr.Name,
		FullName: gr.FullName,
		Link:     gr.HTMLURL,
		Clone:    gr.CloneURL,
		Branch:   gr.DefaultBranch,
		Private:  gr.Private,
	}
}

// GetPullRequest defines a function that retrieves
// a pull request for a repo.
// nolint:lll // function signature is lengthy
func (c *client) GetPullRequest(u *library.User, r *library.Repo, number int) (string, string, string, string, error) {
	logrus.Tracef("Listing source repositories for %s", u.GetName())
	return "", "", "", "", nil
}

// GetHTMLURL retrieves the html_url from repository contents from the GitHub repo.
func (c *client) GetHTMLURL(u *library.User, org, repo, name, ref string) (string, error) {
	logrus.Tracef("Capturing html_url for %s/%s/%s@%s", org, repo, name, ref)
	return "", nil
}
