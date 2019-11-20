// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/v26/github"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

// Config gets the pipeline configuration from the GitHub repo.
func (c *client) Config(u *library.User, org, name, ref string) ([]byte, error) {
	client := c.newClientToken(*u.Token)

	// set the reference for the options to capture the pipeline configuration
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	// send API call to capture the .carave.yml pipeline configuration
	data, _, resp, err := client.Repositories.GetContents(ctx, org, name, ".vela.yml", opts)
	if err != nil {
		if resp.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}

	// data is not nil if .vela.yml exists
	if data != nil {
		strData, err := data.GetContent()
		if err != nil {
			return nil, err
		}

		return []byte(strData), nil
	}

	// send API call to capture the .carave.yaml pipeline configuration
	data, _, resp, err = client.Repositories.GetContents(ctx, org, name, ".vela.yaml", opts)
	if err != nil {
		if resp.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}

	// data is not nil if .vela.yaml exists
	if data != nil {
		strData, err := data.GetContent()
		if err != nil {
			return nil, err
		}

		return []byte(strData), nil
	}

	return nil, fmt.Errorf("No valid pipeline configuration file (.vela.yml or .vela.yaml) found")
}

// Disable deactivates a repo by deleting the webhook.
func (c *client) Disable(u *library.User, org, name string) error {
	client := c.newClientToken(*u.Token)

	// send API call to capture the hooks for the repo
	hooks, _, err := client.Repositories.ListHooks(ctx, org, name, nil)
	if err != nil {
		return err
	}

	// since 0 might be a real value (though unlikely?)
	var id *int64

	// iterate through each element in the hooks
	for _, hook := range hooks {
		// skip if the hook has no ID
		if hook.ID == nil {
			continue
		}

		// cast url from hook configuration to string
		hookURL := hook.Config["url"].(string)

		// capture hook ID if the hook url matches
		if hookURL == fmt.Sprintf("%s/webhook", c.LocalHost) {
			id = hook.ID
		}
	}

	// skip if we got no hook ID
	if id == nil {
		return nil
	}

	// send API call to delete the webhook
	_, err = client.Repositories.DeleteHook(ctx, org, name, *id)

	return err
}

// Enable activates a repo by creating the webhook.
func (c *client) Enable(u *library.User, org, name string) (string, error) {
	client := c.newClientToken(*u.Token)

	// create the hook object to make the API call
	hook := &github.Hook{
		Events: []string{
			constants.EventPush,
			constants.EventPull,
			constants.EventDeploy,
		},
		Config: map[string]interface{}{
			"url":          fmt.Sprintf("%s/webhook", c.LocalHost),
			"content_type": "form",
		},
		Active: github.Bool(true),
	}

	// send API call to create the webhook
	_, resp, err := client.Repositories.CreateHook(ctx, org, name, hook)

	switch resp.StatusCode {
	case 422:
		return "", fmt.Errorf("Repo already enabled")
	case 404:
		return "", fmt.Errorf("Repo not found")
	}

	// create the URL for the repo
	url := fmt.Sprintf("%s/%s/%s", c.URL, org, name)

	return url, err
}

// Status sends the commit status for the given SHA from the GitHub repo.
func (c *client) Status(u *library.User, b *library.Build, org, name string) error {
	client := c.newClientToken(*u.Token)

	context := fmt.Sprintf("%s/%s", c.StatusContext, b.GetEvent())
	url := fmt.Sprintf("%s/%s/%s/%d", c.LocalHost, org, name, b.GetNumber())

	var state string
	var description string
	// set the state and description for the status context
	// depending on what the status of the build is
	switch b.GetStatus() {
	case constants.StatusRunning, constants.StatusPending:
		state = "pending"
		description = "the build is pending"
	case constants.StatusSuccess:
		state = "success"
		description = "the build was successful"
	case constants.StatusFailure:
		state = "failure"
		description = "the build has failed"
	case constants.StatusKilled:
		state = "failure"
		description = "the build was killed"
	default:
		state = "error"
		description = "there was an error"
	}

	// create the status object to make the API call
	status := &github.RepoStatus{
		Context:     github.String(context),
		Description: github.String(description),
		State:       github.String(state),
		TargetURL:   github.String(url),
	}

	// send API call to create the status context for the commit
	_, _, err := client.Repositories.CreateStatus(ctx, org, name, b.GetCommit(), status)

	return err
}

// create GitHub OAuth client with user's token

// ListChanges sends the list of files changed for a none pull request event.
func (c *client) ListChanges(u *library.User, r *library.Repo, sha string) ([]string, error) {
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

// ListChanges sends the list of files changed for a pull request event.
func (c *client) ListChangesPR(u *library.User, r *library.Repo, number int) ([]string, error) {
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

// ListUserRepos returns a list of all repos the user has 'admin' privileges to.
func (c *client) ListUserRepos(u *library.User) ([]*library.Repo, error) {
	// create GitHub OAuth client with user's token
	client := c.newClientToken(u.GetToken())
	r := []*github.Repository{}
	f := []*library.Repo{}

	// set the max per page for the options to capture the list of repos
	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100}, // 100 is max
	}

	// loop to capture *ALL* the repos
	for {
		// send API call to capture the user's repos
		repos, resp, err := client.Repositories.List(ctx, "", opts)
		if err != nil {
			return nil, fmt.Errorf("Repositories.List returned error: %v", err)
		}

		r = append(r, repos...)

		// break the loop if there is no more results to page through
		if resp.NextPage == 0 {
			break
		}

		opts.Page = resp.NextPage
	}

	// iterate through each repo for the user
	for _, repo := range r {
		// skip if the repo is archived or disabled
		if repo.GetArchived() || repo.GetDisabled() {
			continue
		}

		// skip if the user does not have admin access to the repo
		val, ok := repo.GetPermissions()["admin"]
		if !ok || !val {
			continue
		}

		f = append(f, toLibraryRepo(*repo))
	}

	return f, nil
}

// toLibraryRepo does a partial conversion of a github repo to a library repo
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
