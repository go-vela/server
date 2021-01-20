// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v29/github"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// ConfigBackoff is a wrapper for Config that will retry five times if the function
// fails to retrieve the yaml/yml file.
func (c *client) ConfigBackoff(u *library.User, org, name, ref string) (data []byte, err error) {
	// number of times to retry
	retryLimit := 5

	for i := 0; i < retryLimit; i++ {
		// attempt to fetch the config
		data, err = c.Config(u, org, name, ref)

		// return err if the last attempt returns error
		if err != nil && i == retryLimit-1 {
			return
		}

		// if data is valid break the retry loop
		if data != nil {
			break
		}

		// sleep in between retries
		sleep := time.Duration(i+1) * time.Second
		time.Sleep(sleep)
	}

	return
}

// Config gets the pipeline configuration from the GitHub repo.
func (c *client) Config(u *library.User, org, name, ref string) ([]byte, error) {
	logrus.Tracef("Capturing configuration file for %s/%s/commit/%s", org, name, ref)

	// create GitHub OAuth client with user's token
	client := c.newClientToken(*u.Token)

	// set the reference for the options to capture the pipeline configuration
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	// send API call to capture the .vela.yml pipeline configuration
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

	// send API call to capture the .vela.yaml pipeline configuration
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

	return nil, fmt.Errorf("no valid pipeline configuration file (.vela.yml or .vela.yaml) found")
}

// Disable deactivates a repo by deleting the webhook.
func (c *client) Disable(u *library.User, org, name string) error {
	logrus.Tracef("Deleting repository webhook for %s/%s", org, name)

	// create GitHub OAuth client with user's token
	client := c.newClientToken(*u.Token)

	// send API call to capture the hooks for the repo
	hooks, _, err := client.Repositories.ListHooks(ctx, org, name, nil)
	if err != nil {
		return err
	}

	// accounting for situations in which multiple hooks have been
	// associated with this vela instance, which causes some
	// disable, repair, enable operations to act in undesirable ways
	var ids []int64

	// iterate through each element in the hooks
	for _, hook := range hooks {
		// skip if the hook has no ID
		if hook.GetID() == 0 {
			continue
		}

		// cast url from hook configuration to string
		hookURL := hook.Config["url"].(string)

		// capture hook ID if the hook url matches
		if hookURL == fmt.Sprintf("%s/webhook", c.LocalHost) {
			ids = append(ids, hook.GetID())
		}
	}

	// skip if we have no hook IDs
	if len(ids) == 0 {
		return nil
	}

	// go through all found hook IDs and delete them
	for _, id := range ids {
		// send API call to delete the webhook
		_, err = client.Repositories.DeleteHook(ctx, org, name, id)
	}

	return err
}

// Enable activates a repo by creating the webhook.
func (c *client) Enable(u *library.User, org, name, secret string) (string, error) {
	logrus.Tracef("Creating repository webhook for %s/%s", org, name)

	// create GitHub OAuth client with user's token
	client := c.newClientToken(*u.Token)

	// create the hook object to make the API call
	hook := &github.Hook{
		Events: []string{
			eventPush,
			eventPullRequest,
			eventDeployment,
			eventIssueComment,
		},
		Config: map[string]interface{}{
			"url":          fmt.Sprintf("%s/webhook", c.LocalHost),
			"content_type": "form",
			"secret":       secret,
		},
		Active: github.Bool(true),
	}

	// send API call to create the webhook
	_, resp, err := client.Repositories.CreateHook(ctx, org, name, hook)

	switch resp.StatusCode {
	case http.StatusUnprocessableEntity:
		return "", fmt.Errorf("repo already enabled")
	case http.StatusNotFound:
		return "", fmt.Errorf("repo not found")
	}

	// create the URL for the repo
	url := fmt.Sprintf("%s/%s/%s", c.URL, org, name)

	return url, err
}

// Status sends the commit status for the given SHA from the GitHub repo.
func (c *client) Status(u *library.User, b *library.Build, org, name string) error {
	logrus.Tracef("Setting commit status for %s/%s/%d @ %s", org, name, b.GetNumber(), b.GetCommit())

	// create GitHub OAuth client with user's token
	client := c.newClientToken(*u.Token)

	context := fmt.Sprintf("%s/%s", c.StatusContext, b.GetEvent())
	url := fmt.Sprintf("%s/%s/%s/%d", c.WebUIHost, org, name, b.GetNumber())

	var (
		state       string
		description string
	)

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

	// check if the build event is deployment
	if strings.EqualFold(b.GetEvent(), constants.EventDeploy) {
		// parse out deployment number from build source URL
		//
		// pattern: <org>/<repo>/deployments/<deployment_id>
		var parts []string
		if strings.Contains(b.GetSource(), "/deployments/") {
			parts = strings.Split(b.GetSource(), "/deployments/")
		}

		// capture number by converting from string
		number, err := strconv.Atoi(parts[1])
		if err != nil {
			// capture number by scanning from string
			_, err := fmt.Sscanf(b.GetSource(), "%s/%d", nil, &number)
			if err != nil {
				return err
			}
		}

		// create the status object to make the API call
		status := &github.DeploymentStatusRequest{
			Description: github.String(description),
			Environment: github.String(b.GetDeploy()),
			State:       github.String(state),
		}

		// provide "Details" link in GitHub UI if server was configured with it
		if len(c.WebUIHost) > 0 {
			status.LogURL = github.String(url)
		}

		_, _, err = client.Repositories.CreateDeploymentStatus(ctx, org, name, int64(number), status)

		return err
	}

	// create the status object to make the API call
	status := &github.RepoStatus{
		Context:     github.String(context),
		Description: github.String(description),
		State:       github.String(state),
	}

	// provide "Details" link in GitHub UI if server was configured with it
	if len(c.WebUIHost) > 0 {
		status.TargetURL = github.String(url)
	}

	// send API call to create the status context for the commit
	_, _, err := client.Repositories.CreateStatus(ctx, org, name, b.GetCommit(), status)

	return err
}

// GetRepo gets repo information from Github.
func (c *client) GetRepo(u *library.User, r *library.Repo) (*library.Repo, error) {
	logrus.Tracef("Retrieving repository information for %s", r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newClientToken(u.GetToken())

	// send an API call to get the repo info
	repo, _, err := client.Repositories.Get(ctx, r.GetOrg(), r.GetName())
	if err != nil {
		return nil, err
	}

	return toLibraryRepo(*repo), nil
}

// ListUserRepos returns a list of all repos the user has access to.
func (c *client) ListUserRepos(u *library.User) ([]*library.Repo, error) {
	logrus.Tracef("Listing source repositories for %s", u.GetName())

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
			return nil, fmt.Errorf("unable to list user repos: %v", err)
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

// GetPullRequest defines a function that retrieves
// a pull request for a repo.
// nolint:lll // function signature is lengthy
func (c *client) GetPullRequest(u *library.User, r *library.Repo, number int) (string, string, string, string, error) {
	logrus.Tracef("Listing source repositories for %s", u.GetName())

	// create GitHub OAuth client with user's token
	client := c.newClientToken(u.GetToken())

	pull, _, err := client.PullRequests.Get(ctx, r.GetOrg(), r.GetName(), number)
	if err != nil {
		return "", "", "", "", err
	}

	commit := pull.GetHead().GetSHA()
	branch := pull.GetBase().GetRef()
	baseref := pull.GetBase().GetRef()
	headref := pull.GetHead().GetRef()

	return commit, branch, baseref, headref, nil
}

// GetHTMLURL retrieves the html_url from repository contents from the GitHub repo.
func (c *client) GetHTMLURL(u *library.User, org, repo, name, ref string) (string, error) {
	logrus.Tracef("Capturing html_url for %s/%s/%s@%s", org, repo, name, ref)

	// create GitHub OAuth client with user's token
	client := c.newClientToken(*u.Token)

	// set the reference for the options to capture the repository contents
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	// send API call to capture the repository contents for org/repo/name at the ref provided
	data, _, _, err := client.Repositories.GetContents(ctx, org, repo, name, opts)
	if err != nil {
		return "", err
	}

	// data is not nil if the file exists
	if data != nil {
		URL := data.GetHTMLURL()
		if err != nil {
			return "", err
		}

		return URL, nil
	}

	return "", fmt.Errorf("no valid repository contents found")
}
