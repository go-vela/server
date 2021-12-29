// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jenkins-x/go-scm/scm"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// ConfigBackoff is a wrapper for Config that will retry five times if the function
// fails to retrieve the yaml/yml file.
// nolint: lll // ignore long line length due to input arguments
func (c *client) ConfigBackoff(u *library.User, r *library.Repo, ref string) (data []byte, err error) {
	// number of times to retry
	retryLimit := 5

	for i := 0; i < retryLimit; i++ {
		// attempt to fetch the config
		data, err = c.Config(u, r, ref)

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
func (c *client) Config(u *library.User, r *library.Repo, ref string) ([]byte, error) {
	logrus.Tracef("Capturing configuration file for %s/%s/commit/%s", r.GetOrg(), r.GetName(), ref)

	// create GitHub OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return nil, err
	}

	files := []string{".vela.yml", ".vela.yaml"}

	if strings.EqualFold(r.GetPipelineType(), constants.PipelineTypeStarlark) {
		files = append(files, ".vela.star", ".vela.py")
	}

	for _, file := range files {
		// send API call to capture the .vela.yml pipeline configuration
		content, resp, err := client.Contents.Find(ctx, r.GetFullName(), file, ref)
		if err != nil {
			if resp.Status != http.StatusNotFound {
				return nil, err
			}
		}

		// data is not nil if .vela.yml exists
		if content != nil {
			return content.Data, nil
		}
	}

	return nil, fmt.Errorf("no valid pipeline configuration file (%s) found", strings.Join(files, ","))
}

// Disable deactivates a repo by deleting the webhook.
func (c *client) Disable(u *library.User, org, name string) error {
	logrus.Tracef("Deleting repository webhook for %s/%s", org, name)

	// create GitHub OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return err
	}

	opts := scm.ListOptions{
		Size: 100,
	}

	// send API call to capture the hooks for the repo
	hooks, _, err := client.Repositories.ListHooks(ctx, fmt.Sprintf("%s/%s", org, name), opts)
	if err != nil {
		return err
	}

	// accounting for situations in which multiple hooks have been
	// associated with this vela instance, which causes some
	// disable, repair, enable operations to act in undesirable ways
	var ids []string

	// iterate through each element in the hooks
	for _, hook := range hooks {
		// skip if the hook has no ID
		if strings.EqualFold(hook.ID, "0") {
			continue
		}

		// capture hook ID if the hook url matches
		if strings.EqualFold(hook.Target, fmt.Sprintf("%s/webhook", c.config.ServerWebhookAddress)) {
			ids = append(ids, hook.ID)
		}
	}

	// skip if we have no hook IDs
	if len(ids) == 0 {
		return nil
	}

	// go through all found hook IDs and delete them
	for _, id := range ids {
		// send API call to delete the webhook
		_, err = client.Repositories.DeleteHook(ctx, fmt.Sprintf("%s/%s", org, name), id)
	}

	return err
}

// Enable activates a repo by creating the webhook.
func (c *client) Enable(u *library.User, org, name, secret string) (string, error) {
	logrus.Tracef("Creating repository webhook for %s/%s", org, name)

	// create GitHub OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return "", err
	}

	// create the hook object to make the API call
	hook := &scm.HookInput{
		Events: scm.HookEvents{
			Push:         true,
			PullRequest:  true,
			Deployment:   true,
			IssueComment: true,
		},
		Target: fmt.Sprintf("%s/webhook", c.config.ServerWebhookAddress),
		Secret: secret,
	}

	// send API call to create the webhook
	_, resp, err := client.Repositories.CreateHook(ctx, fmt.Sprintf("%s/%s", org, name), hook)

	switch resp.Status {
	case http.StatusUnprocessableEntity:
		return "", fmt.Errorf("repo already enabled")
	case http.StatusNotFound:
		return "", fmt.Errorf("repo not found")
	}

	// create the URL for the repo
	url := fmt.Sprintf("%s/%s/%s", c.config.Address, org, name)

	return url, err
}

// Status sends the commit status for the given SHA from the GitHub repo.
func (c *client) Status(u *library.User, b *library.Build, org, name string) error {
	logrus.Tracef("Setting commit status for %s/%s/%d @ %s", org, name, b.GetNumber(), b.GetCommit())

	// create GitHub OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return err
	}

	var (
		state       scm.State
		description string
	)

	// set the state and description for the status context
	// depending on what the status of the build is
	switch b.GetStatus() {
	case constants.StatusRunning, constants.StatusPending:
		state = scm.StatePending
		description = fmt.Sprintf("the build is %s", b.GetStatus())
	case constants.StatusSuccess:
		state = scm.StateSuccess
		description = "the build was successful"
	case constants.StatusFailure:
		state = scm.StateFailure
		description = "the build has failed"
	case constants.StatusCanceled:
		state = scm.StateFailure
		description = "the build was canceled"
	case constants.StatusKilled:
		state = scm.StateFailure
		description = "the build was killed"
	case constants.StatusSkipped:
		state = scm.StateSuccess
		description = "build was skipped as no steps/stages found"
	default:
		state = scm.StateError
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
		status := &scm.DeploymentStatusInput{
			Description: description,
			Environment: b.GetDeploy(),
			State:       state.String(),
		}

		// provide "Details" link in GitHub UI if server was configured with it
		if len(c.config.WebUIAddress) > 0 {
			status.LogLink = fmt.Sprintf("%s/%s/%s/%d", c.config.WebUIAddress, org, name, b.GetNumber())
		}

		_, _, err = client.Deployments.CreateStatus(
			ctx,
			fmt.Sprintf("%s/%s", org, name),
			strconv.Itoa(number),
			status,
		)

		return err
	}

	// create the status object to make the API call
	status := &scm.StatusInput{
		Label: fmt.Sprintf("%s/%s", c.config.StatusContext, b.GetEvent()),
		Desc:  description,
		State: state,
	}

	// provide "Details" link in GitHub UI if server was configured with it
	if len(c.config.WebUIAddress) > 0 && b.GetStatus() != constants.StatusSkipped {
		status.Target = fmt.Sprintf("%s/%s/%s/%d", c.config.WebUIAddress, org, name, b.GetNumber())
	}

	// send API call to create the status context for the commit
	_, _, err = client.Repositories.CreateStatus(
		ctx,
		fmt.Sprintf("%s/%s", org, name),
		b.GetCommit(),
		status,
	)

	return err
}

// GetRepo gets repo information from Github.
func (c *client) GetRepo(u *library.User, r *library.Repo) (*library.Repo, error) {
	logrus.Tracef("Retrieving repository information for %s", r.GetFullName())

	// create GitHub OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return nil, err
	}

	// send an API call to get the repo info
	repo, _, err := client.Repositories.Find(ctx, r.GetFullName())
	if err != nil {
		return nil, err
	}

	return toLibraryRepo(*repo), nil
}

// ListUserRepos returns a list of all repos the user has access to.
func (c *client) ListUserRepos(u *library.User) ([]*library.Repo, error) {
	logrus.Tracef("Listing source repositories for %s", u.GetName())

	// create GitHub OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return nil, err
	}

	r := []*scm.Repository{}
	f := []*library.Repo{}

	// set the max per page for the options to capture the list of repos
	opts := scm.ListOptions{
		Size: 100, // 100 is max
	}

	// loop to capture *ALL* the repos
	for {
		// send API call to capture the user's repos
		repos, resp, err := client.Repositories.List(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("unable to list user repos: %v", err)
		}

		r = append(r, repos...)

		// break the loop if there is no more results to page through
		if resp.Page.Next == 0 {
			break
		}

		opts.Page = resp.Page.Next
	}

	// iterate through each repo for the user
	for _, repo := range r {
		if repo.Archived {
			continue
		}

		f = append(f, toLibraryRepo(*repo))
	}

	return f, nil
}

// GetPullRequest defines a function that retrieves
// a pull request for a repo.
// nolint:lll // function signature is lengthy
func (c *client) GetPullRequest(u *library.User, r *library.Repo, number int) (string, string, string, string, error) {
	logrus.Tracef("Listing source repositories for %s", u.GetName())

	// create GitHub OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return "", "", "", "", err
	}

	pull, _, err := client.PullRequests.Find(ctx, r.GetFullName(), number)
	if err != nil {
		return "", "", "", "", err
	}

	commit := pull.Head.Sha
	branch := pull.Base.Ref
	baseref := pull.Base.Ref
	headref := pull.Head.Ref

	return commit, branch, baseref, headref, nil
}

// GetHTMLURL retrieves the html_url from repository contents from the GitHub repo.
func (c *client) GetHTMLURL(u *library.User, org, repo, name, ref string) (string, error) {
	logrus.Tracef("Capturing html_url for %s/%s/%s@%s", org, repo, name, ref)

	// create GitHub OAuth client with user's token
	client, err := c.newClientToken(*u.Token)
	if err != nil {
		return "", err
	}

	// send API call to capture the repository contents for org/repo/name at the ref provided
	content, _, err := client.Contents.Find(ctx, fmt.Sprintf("%s/%s", org, repo), name, ref)
	if err != nil {
		return "", err
	}

	// data is not nil if the file exists
	if content != nil {
		// TODO: The library doesn't expose the URL data.
		// We might need to massage this path or make a feature request.
		URL := content.Path
		if err != nil {
			return "", err
		}

		return URL, nil
	}

	return "", fmt.Errorf("no valid repository contents found")
}

// toLibraryRepo does a partial conversion of a github repo to a library repo.
func toLibraryRepo(scm scm.Repository) *library.Repo {
	return &library.Repo{
		Org:      &scm.Namespace,
		Name:     &scm.Name,
		FullName: &scm.FullName,
		Link:     &scm.Link,
		Clone:    &scm.Clone,
		Branch:   &scm.Branch,
		Private:  &scm.Private,
	}
}
