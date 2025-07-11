// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v73/github"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/yaml/yaml"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
)

// ConfigBackoff is a wrapper for Config that will retry five times if the function
// fails to retrieve the yaml/yml file.
func (c *Client) ConfigBackoff(ctx context.Context, u *api.User, r *api.Repo, ref string) (data []byte, err error) {
	// number of times to retry
	retryLimit := 5

	for i := 0; i < retryLimit; i++ {
		logrus.Debugf("fetching config file - Attempt %d", i+1)
		// attempt to fetch the config
		data, err = c.Config(ctx, u, r, ref)

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
func (c *Client) Config(ctx context.Context, u *api.User, r *api.Repo, ref string) ([]byte, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Tracef("capturing configuration file for %s/commit/%s", r.GetFullName(), ref)

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, *u.Token)

	// default pipeline file names
	files := []string{".vela.yml", ".vela.yaml"}

	// starlark support - prefer .star/.py, use default as fallback
	if strings.EqualFold(r.GetPipelineType(), constants.PipelineTypeStarlark) {
		files = append([]string{".vela.star", ".vela.py"}, files...)
	}

	// set the reference for the options to capture the pipeline configuration
	opts := &github.RepositoryContentGetOptions{
		Ref: ref,
	}

	for _, file := range files {
		// send API call to capture the .vela.yml pipeline configuration
		data, _, resp, err := client.Repositories.GetContents(ctx, r.GetOrg(), r.GetName(), file, opts)
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
	}

	return nil, fmt.Errorf("no valid pipeline configuration file (%s) found", strings.Join(files, ","))
}

// Disable deactivates a repo by deleting the webhook.
func (c *Client) Disable(ctx context.Context, u *api.User, org, name string) error {
	return c.DestroyWebhook(ctx, u, org, name)
}

// DestroyWebhook deletes a repo's webhook.
func (c *Client) DestroyWebhook(ctx context.Context, u *api.User, org, name string) error {
	c.Logger.WithFields(logrus.Fields{
		"org":  org,
		"repo": name,
		"user": u.GetName(),
	}).Tracef("deleting repository webhooks for %s/%s", org, name)

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, *u.Token)

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

		// capture hook ID if the hook url matches
		if strings.EqualFold(hook.GetConfig().GetURL(), c.config.ServerWebhookAddress) {
			ids = append(ids, hook.GetID())
		}
	}

	// skip if we have no hook IDs
	if len(ids) == 0 {
		c.Logger.WithFields(logrus.Fields{
			"org":  org,
			"repo": name,
			"user": u.GetName(),
		}).Warnf("no repository webhooks matching %s found for %s/%s", c.config.ServerWebhookAddress, org, name)

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
func (c *Client) Enable(ctx context.Context, u *api.User, r *api.Repo, h *api.Hook) (*api.Hook, string, error) {
	return c.CreateWebhook(ctx, u, r, h)
}

// CreateWebhook creates a repo's webhook.
func (c *Client) CreateWebhook(ctx context.Context, u *api.User, r *api.Repo, h *api.Hook) (*api.Hook, string, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Tracef("creating repository webhook for %s/%s", r.GetOrg(), r.GetName())

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, *u.Token)

	// always listen to repository events in case of repo name change
	events := []string{eventRepository, constants.EventCustomProperties}

	// subscribe to comment event if any comment action is allowed
	if r.GetAllowEvents().GetComment().GetCreated() ||
		r.GetAllowEvents().GetComment().GetEdited() {
		events = append(events, eventIssueComment)
	}

	// subscribe to deployment event if allowed
	if r.GetAllowEvents().GetDeployment().GetCreated() {
		events = append(events, eventDeployment)
	}

	// subscribe to pull_request event if any PR action is allowed
	if r.GetAllowEvents().GetPullRequest().GetOpened() ||
		r.GetAllowEvents().GetPullRequest().GetEdited() ||
		r.GetAllowEvents().GetPullRequest().GetSynchronize() {
		events = append(events, eventPullRequest)
	}

	// subscribe to push event if branch push or tag is allowed
	if r.GetAllowEvents().GetPush().GetBranch() ||
		r.GetAllowEvents().GetPush().GetTag() {
		events = append(events, eventPush)
	}

	// create the hook object to make the API call
	hook := &github.Hook{
		Events: events,
		Config: &github.HookConfig{
			URL:         github.Ptr(c.config.ServerWebhookAddress),
			ContentType: github.Ptr("form"),
			Secret:      github.Ptr(r.GetHash()),
		},
		Active: github.Ptr(true),
	}

	// send API call to create the webhook
	hookInfo, resp, err := client.Repositories.CreateHook(ctx, r.GetOrg(), r.GetName(), hook)

	// create the first hook for the repo and record its ID from GitHub
	webhook := new(api.Hook)
	webhook.SetWebhookID(hookInfo.GetID())
	webhook.SetSourceID(r.GetName() + "-" + eventInitialize)
	webhook.SetCreated(hookInfo.GetCreatedAt().Unix())
	webhook.SetEvent(eventInitialize)
	webhook.SetNumber(h.GetNumber() + 1)
	webhook.SetStatus(constants.StatusSuccess)

	switch resp.StatusCode {
	case http.StatusUnprocessableEntity:
		return nil, "", fmt.Errorf("repo already enabled")
	case http.StatusNotFound:
		return nil, "", fmt.Errorf("repo not found")
	}

	// create the URL for the repo
	url := fmt.Sprintf("%s/%s/%s", c.config.Address, r.GetOrg(), r.GetName())

	return webhook, url, err
}

// Update edits a repo webhook.
func (c *Client) Update(ctx context.Context, u *api.User, r *api.Repo, hookID int64) (bool, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Tracef("updating repository webhook for %s/%s", r.GetOrg(), r.GetName())

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, *u.Token)

	// always listen to repository events in case of repo name change
	events := []string{eventRepository, constants.EventCustomProperties}

	// subscribe to comment event if any comment action is allowed
	if r.GetAllowEvents().GetComment().GetCreated() ||
		r.GetAllowEvents().GetComment().GetEdited() {
		events = append(events, eventIssueComment)
	}

	// subscribe to deployment event if allowed
	if r.GetAllowEvents().GetDeployment().GetCreated() {
		events = append(events, eventDeployment)
	}

	// subscribe to pull_request event if any PR action is allowed
	if r.GetAllowEvents().GetPullRequest().GetOpened() ||
		r.GetAllowEvents().GetPullRequest().GetEdited() ||
		r.GetAllowEvents().GetPullRequest().GetSynchronize() {
		events = append(events, eventPullRequest)
	}

	// subscribe to push event if branch push or tag is allowed
	if r.GetAllowEvents().GetPush().GetBranch() ||
		r.GetAllowEvents().GetPush().GetTag() {
		events = append(events, eventPush)
	}

	// create the hook object to make the API call
	hook := &github.Hook{
		Events: events,
		Config: &github.HookConfig{
			URL:         github.Ptr(c.config.ServerWebhookAddress),
			ContentType: github.Ptr("form"),
			Secret:      github.Ptr(r.GetHash()),
		},
		Active: github.Ptr(true),
	}

	// send API call to update the webhook
	_, resp, err := client.Repositories.EditHook(ctx, r.GetOrg(), r.GetName(), hookID, hook)

	// track if webhook exists in GitHub; a missing webhook
	// indicates the webhook has been manually deleted from GitHub
	return resp.StatusCode != http.StatusNotFound, err
}

// Status sends the commit status for the given SHA from the GitHub repo.
func (c *Client) Status(ctx context.Context, u *api.User, b *api.Build, org, name string) error {
	c.Logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   org,
		"repo":  name,
		"user":  u.GetName(),
	}).Tracef("setting commit status for %s/%s/%d @ %s", org, name, b.GetNumber(), b.GetCommit())

	// only report opened, synchronize, and reopened action types for pull_request events
	if strings.EqualFold(b.GetEvent(), constants.EventPull) && !strings.EqualFold(b.GetEventAction(), constants.ActionOpened) &&
		!strings.EqualFold(b.GetEventAction(), constants.ActionSynchronize) && !strings.EqualFold(b.GetEventAction(), constants.ActionReopened) {
		return nil
	}

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, *u.Token)

	context := fmt.Sprintf("%s/%s", c.config.StatusContext, b.GetEvent())
	url := fmt.Sprintf("%s/%s/%s/%d", c.config.WebUIAddress, org, name, b.GetNumber())

	var (
		state       string
		description string
	)

	// set the state and description for the status context
	// depending on what the status of the build is
	switch b.GetStatus() {
	case constants.StatusRunning, constants.StatusPending:
		//nolint:goconst // ignore making constant
		state = "pending"
		description = fmt.Sprintf("the build is %s", b.GetStatus())
	case constants.StatusPendingApproval:
		state = "pending"
		description = "build needs approval from repo admin to run"
	case constants.StatusSuccess:
		//nolint:goconst // ignore making constant
		state = "success"
		description = "the build was successful"
	case constants.StatusFailure:
		//nolint:goconst // ignore making constant
		state = "failure"
		description = "the build has failed"
	case constants.StatusCanceled:
		state = "failure"
		description = "the build was canceled"
	case constants.StatusKilled:
		state = "failure"
		description = "the build was killed"
	case constants.StatusSkipped:
		state = "success"
		description = "build was skipped as no steps/stages found"
	default:
		state = "error"

		// if there is no build, then this status update is from a failed compilation
		if b.GetID() == 0 {
			description = "error compiling pipeline - check audit for more information"
			url = fmt.Sprintf("%s/%s/%s/hooks", c.config.WebUIAddress, org, name)
		} else {
			description = "there was an error"
		}
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
			Description: github.Ptr(description),
			Environment: github.Ptr(b.GetDeploy()),
			State:       github.Ptr(state),
		}

		// provide "Details" link in GitHub UI if server was configured with it
		if len(c.config.WebUIAddress) > 0 {
			status.LogURL = github.Ptr(url)
		}

		_, _, err = client.Repositories.CreateDeploymentStatus(ctx, org, name, int64(number), status)

		return err
	}

	// create the status object to make the API call
	status := &github.RepoStatus{
		Context:     github.Ptr(context),
		Description: github.Ptr(description),
		State:       github.Ptr(state),
	}

	// provide "Details" link in GitHub UI if server was configured with it
	if len(c.config.WebUIAddress) > 0 && b.GetStatus() != constants.StatusSkipped {
		status.TargetURL = github.Ptr(url)
	}

	// send API call to create the status context for the commit
	_, _, err := client.Repositories.CreateStatus(ctx, org, name, b.GetCommit(), status)

	return err
}

// StepStatus sends the commit status for the given SHA to the GitHub repo with the step as the context.
func (c *Client) StepStatus(ctx context.Context, u *api.User, b *api.Build, s *api.Step, org, name string) error {
	c.Logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   org,
		"repo":  name,
		"user":  u.GetName(),
	}).Tracef("setting commit status for %s/%s/%d @ %s", org, name, b.GetNumber(), b.GetCommit())

	// no commit statuses on deployments
	if strings.EqualFold(b.GetEvent(), constants.EventDeploy) {
		return nil
	}

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, *u.Token)

	context := fmt.Sprintf("%s/%s/%s", c.config.StatusContext, b.GetEvent(), s.GetReportAs())
	url := fmt.Sprintf("%s/%s/%s/%d#%d", c.config.WebUIAddress, org, name, b.GetNumber(), s.GetNumber())

	var (
		state       string
		description string
	)

	// set the state and description for the status context
	// depending on what the status of the step is
	switch s.GetStatus() {
	case constants.StatusRunning, constants.StatusPending, constants.StatusPendingApproval:
		state = "pending"
		description = fmt.Sprintf("the step is %s", s.GetStatus())
	case constants.StatusSuccess:
		state = "success"
		description = "the step was successful"
	case constants.StatusFailure:
		state = "failure"
		description = "the step has failed"
	case constants.StatusCanceled:
		state = "failure"
		description = "the step was canceled"
	case constants.StatusKilled:
		state = "failure"
		description = "the step was killed"
	case constants.StatusSkipped:
		state = "failure"
		description = "step was skipped or never ran"
	default:
		state = "error"
		description = "there was an error"
	}

	// create the status object to make the API call
	status := &github.RepoStatus{
		Context:     github.Ptr(context),
		Description: github.Ptr(description),
		State:       github.Ptr(state),
	}

	// provide "Details" link in GitHub UI if server was configured with it
	if len(c.config.WebUIAddress) > 0 && b.GetStatus() != constants.StatusSkipped {
		status.TargetURL = github.Ptr(url)
	}

	// send API call to create the status context for the commit
	_, _, err := client.Repositories.CreateStatus(ctx, org, name, b.GetCommit(), status)

	return err
}

// GetRepo gets repo information from Github.
func (c *Client) GetRepo(ctx context.Context, u *api.User, r *api.Repo) (*api.Repo, int, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Tracef("retrieving repository information for %s", r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, u.GetToken())

	// send an API call to get the repo info
	repo, resp, err := client.Repositories.Get(ctx, r.GetOrg(), r.GetName())
	if err != nil {
		return nil, resp.StatusCode, err
	}

	return toAPIRepo(*repo), resp.StatusCode, nil
}

// GetOrgAndRepoName returns the name of the org and the repository in the SCM.
func (c *Client) GetOrgAndRepoName(ctx context.Context, u *api.User, o string, r string) (string, string, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  o,
		"repo": r,
		"user": u.GetName(),
	}).Tracef("retrieving repository information for %s/%s", o, r)

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, u.GetToken())

	// send an API call to get the repo info
	repo, _, err := client.Repositories.Get(ctx, o, r)
	if err != nil {
		return "", "", err
	}

	return repo.GetOwner().GetLogin(), repo.GetName(), nil
}

// ListUserRepos returns a list of all repos the user has access to.
func (c *Client) ListUserRepos(ctx context.Context, u *api.User) ([]string, error) {
	c.Logger.WithFields(logrus.Fields{
		"user": u.GetName(),
	}).Tracef("listing source repositories for %s", u.GetName())

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, u.GetToken())

	r := []*github.Repository{}
	f := []string{}

	// set the max per page for the options to capture the list of repos
	opts := &github.RepositoryListByAuthenticatedUserOptions{
		ListOptions: github.ListOptions{PerPage: 100}, // 100 is max
	}

	// loop to capture *ALL* the repos
	for {
		// send API call to capture the user's repos
		repos, resp, err := client.Repositories.ListByAuthenticatedUser(ctx, opts)
		if err != nil {
			return nil, fmt.Errorf("unable to list user repos: %w", err)
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
		// skip if the repo is void
		// fixes an issue with GitHub replication being out of sync
		if repo == nil {
			continue
		}

		// skip if the repo is archived or disabled
		if repo.GetArchived() || repo.GetDisabled() {
			continue
		}

		f = append(f, repo.GetFullName())
	}

	return f, nil
}

// toAPIRepo does a partial conversion of a github repo to a API repo.
func toAPIRepo(gr github.Repository) *api.Repo {
	var visibility string

	// setting the visbility to match the SCM visbility
	switch gr.GetVisibility() {
	// if gh resp does not have visibility field, use private
	case "":
		if gr.GetPrivate() {
			visibility = constants.VisibilityPrivate
		} else {
			visibility = constants.VisibilityPublic
		}
	case "private":
		visibility = constants.VisibilityPrivate
	default:
		visibility = constants.VisibilityPublic
	}

	return &api.Repo{
		Org:        gr.GetOwner().Login,
		Name:       gr.Name,
		FullName:   gr.FullName,
		Link:       gr.HTMLURL,
		Clone:      gr.CloneURL,
		Branch:     gr.DefaultBranch,
		Topics:     &gr.Topics,
		Private:    gr.Private,
		Visibility: &visibility,
	}
}

// GetPullRequest defines a function that retrieves
// a pull request for a repo.
func (c *Client) GetPullRequest(ctx context.Context, r *api.Repo, number int) (string, string, string, string, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": r.GetOwner().GetName(),
	}).Tracef("retrieving pull request %d for repo %s", number, r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, r.GetOwner().GetToken())

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
func (c *Client) GetHTMLURL(ctx context.Context, u *api.User, org, repo, name, ref string) (string, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  org,
		"repo": repo,
		"user": u.GetName(),
	}).Tracef("capturing html_url for %s/%s/%s@%s", org, repo, name, ref)

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, *u.Token)

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

		return URL, nil
	}

	return "", fmt.Errorf("no valid repository contents found")
}

// GetBranch defines a function that retrieves a branch for a repo.
func (c *Client) GetBranch(ctx context.Context, r *api.Repo, branch string) (string, string, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": r.GetOwner().GetName(),
	}).Tracef("retrieving branch %s for repo %s", branch, r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, r.GetOwner().GetToken())

	maxRedirects := 3

	data, _, err := client.Repositories.GetBranch(ctx, r.GetOrg(), r.GetName(), branch, maxRedirects)
	if err != nil {
		return "", "", err
	}

	return data.GetName(), data.GetCommit().GetSHA(), nil
}

// GetNetrcPassword returns a clone token using the repo's github app installation if it exists.
// If not, it defaults to the user OAuth token.
func (c *Client) GetNetrcPassword(ctx context.Context, db database.Interface, r *api.Repo, u *api.User, g yaml.Git) (string, error) {
	l := c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	})

	l.Tracef("getting netrc password for %s/%s", r.GetOrg(), r.GetName())

	// no GitHub App configured, use legacy oauth token
	if c.AppsTransport == nil {
		return u.GetToken(), nil
	}

	var err error

	// repos that the token has access to
	// providing no repos, nil, or empty slice will default the token permissions to the list
	// of repos added to the installation
	repos := g.Repositories

	// use triggering repo as a restrictive default
	if len(repos) == 0 {
		repos = []string{r.GetName()}
	}

	// permissions that are applied to the token for every repo provided
	// providing no permissions, nil, or empty map will default to the permissions
	// of the GitHub App installation
	//
	// the Vela compiler follows a least-privileged-defaults model where
	// the list contains only the triggering repo, unless provided in the git yaml block
	//
	// the default is contents:read and checks:write
	ghPermissions := &github.InstallationPermissions{
		Contents: github.Ptr(AppInstallPermissionRead),
		Checks:   github.Ptr(AppInstallPermissionWrite),
	}

	permissions := g.Permissions
	if permissions == nil {
		permissions = map[string]string{}
	}

	for resource, perm := range permissions {
		ghPermissions, err = ApplyInstallationPermissions(resource, perm, ghPermissions)
		if err != nil {
			return u.GetToken(), err
		}
	}

	// verify repo owner has `write` access to listed repositories before provisioning install token
	//
	// this prevents an app installed across the org from bypassing restrictions
	for _, repo := range g.Repositories {
		if repo == r.GetName() {
			continue
		}

		access, err := c.RepoAccess(ctx, u.GetName(), u.GetToken(), r.GetOrg(), repo)
		if err != nil || (access != constants.PermissionAdmin && access != constants.PermissionWrite) {
			return u.GetToken(), fmt.Errorf("repository owner does not have adequate permissions to request install token for repository %s/%s", r.GetOrg(), repo)
		}
	}

	// the app might not be installed therefore we retain backwards compatibility via the user oauth token
	// https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app-installation
	// the optional list of repos and permissions are driven by yaml
	installToken, installID, err := c.newGithubAppInstallationRepoToken(ctx, r, repos, ghPermissions)
	if err != nil {
		// return the legacy token along with no error for backwards compatibility
		// todo: return an error based based on app installation requirements
		l.Tracef("unable to create github app installation token for repos %v with permissions %v: %v", repos, permissions, err)

		return u.GetToken(), nil
	}

	if installToken != nil && len(installToken.GetToken()) != 0 {
		l.Tracef("using github app installation token for %s/%s", r.GetOrg(), r.GetName())

		// (optional) sync the install ID with the repo
		if db != nil && r.GetInstallID() != installID {
			r.SetInstallID(installID)

			_, err = db.UpdateRepo(ctx, r)
			if err != nil {
				c.Logger.Tracef("unable to update repo with install ID %d: %v", installID, err)
			}
		}

		return installToken.GetToken(), nil
	}

	l.Tracef("using user oauth token for %s/%s", r.GetOrg(), r.GetName())

	return u.GetToken(), nil
}

// SyncRepoWithInstallation ensures the repo is synchronized with the scm installation, if it exists.
func (c *Client) SyncRepoWithInstallation(ctx context.Context, r *api.Repo) (*api.Repo, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("syncing app installation for repo %s/%s", r.GetOrg(), r.GetName())

	// no GitHub App configured, skip
	if c.AppsTransport == nil {
		return r, nil
	}

	client, err := c.newGithubAppClient()
	if err != nil {
		return r, err
	}

	installations, _, err := client.Apps.ListInstallations(ctx, &github.ListOptions{})
	if err != nil {
		return r, err
	}

	var installation *github.Installation

	for _, install := range installations {
		if strings.EqualFold(install.GetAccount().GetLogin(), r.GetOrg()) {
			installation = install
		}
	}

	if installation == nil {
		return r, nil
	}

	installationCanReadRepo, err := c.installationCanReadRepo(ctx, r, installation)
	if err != nil {
		return r, err
	}

	if installationCanReadRepo {
		r.SetInstallID(installation.GetID())
	}

	return r, nil
}
