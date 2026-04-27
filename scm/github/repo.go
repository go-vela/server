// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/go-github/v84/github"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

// ConfigBackoff is a wrapper for Config that will retry five times if the function
// fails to retrieve the yaml/yml file.
func (c *Client) ConfigBackoff(ctx context.Context, r *api.Repo, ref, token string) (data []byte, err error) {
	// number of times to retry
	retryLimit := 5

	for i := range retryLimit {
		logrus.Debugf("fetching config file - Attempt %d", i+1)
		// attempt to fetch the config
		data, err = c.Config(ctx, r, ref, token)

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
func (c *Client) Config(ctx context.Context, r *api.Repo, ref, token string) ([]byte, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("capturing configuration file for %s/commit/%s", r.GetFullName(), ref)

	// create GitHub OAuth client with user's token
	client := c.newTokenClient(ctx, token)

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
			if resp == nil || resp.StatusCode != http.StatusNotFound {
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
	client := c.newUserOAuthTokenClient(ctx, u)

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
func (c *Client) Enable(ctx context.Context, u *api.User, r *api.Repo) (*api.Hook, string, error) {
	return c.CreateWebhook(ctx, u, r)
}

// CreateWebhook creates a repo's webhook.
func (c *Client) CreateWebhook(ctx context.Context, u *api.User, r *api.Repo) (*api.Hook, string, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Tracef("creating repository webhook for %s/%s", r.GetOrg(), r.GetName())

	// create GitHub OAuth client with user's token
	client := c.newUserOAuthTokenClient(ctx, u)

	// create the hook object to make the API call
	hook := &github.Hook{
		Events: webhookConfigEvents(r),
		Config: &github.HookConfig{
			URL:         new(c.config.ServerWebhookAddress),
			ContentType: new("form"),
			Secret:      new(r.GetHash()),
		},
		Active: new(true),
	}

	// send API call to create the webhook
	hookInfo, resp, err := client.Repositories.CreateHook(ctx, r.GetOrg(), r.GetName(), hook)

	// create the first hook for the repo and record its ID from GitHub
	webhook := new(api.Hook)
	webhook.SetWebhookID(hookInfo.GetID())
	webhook.SetSourceID(r.GetName() + "-" + eventInitialize)
	webhook.SetCreated(hookInfo.GetCreatedAt().Unix())
	webhook.SetEvent(eventInitialize)
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
	client := c.newUserOAuthTokenClient(ctx, u)

	// create the hook object to make the API call
	hook := &github.Hook{
		Events: webhookConfigEvents(r),
		Config: &github.HookConfig{
			URL:         new(c.config.ServerWebhookAddress),
			ContentType: new("form"),
			Secret:      new(r.GetHash()),
		},
		Active: new(true),
	}

	// send API call to update the webhook
	_, resp, err := client.Repositories.EditHook(ctx, r.GetOrg(), r.GetName(), hookID, hook)

	// track if webhook exists in GitHub; a missing webhook
	// indicates the webhook has been manually deleted from GitHub
	return resp.StatusCode != http.StatusNotFound, err
}

// webhookConfigEvents returns a list of events to subscribe to for the webhook based on the repo's allowed events.
func webhookConfigEvents(r *api.Repo) []string {
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

	// subscribe to merge_group event if repo has merge group events
	if len(r.GetMergeQueueEvents()) > 0 {
		events = append(events, eventMergeGroup)
	}

	return events
}

// GetRepo gets repo information from Github.
func (c *Client) GetRepo(ctx context.Context, u *api.User, r *api.Repo) (*api.Repo, int, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": u.GetName(),
	}).Tracef("retrieving repository information for %s", r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newUserOAuthTokenClient(ctx, u)

	// send an API call to get the repo info
	repo, resp, err := client.Repositories.Get(ctx, r.GetOrg(), r.GetName())
	if err != nil {
		var code int
		if resp != nil {
			code = resp.StatusCode
		} else {
			code = http.StatusInternalServerError
		}

		return nil, code, err
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
	client := c.newUserOAuthTokenClient(ctx, u)

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
	client := c.newUserOAuthTokenClient(ctx, u)

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
		Org:         gr.GetOwner().Login,
		Name:        gr.Name,
		FullName:    gr.FullName,
		Link:        gr.HTMLURL,
		Clone:       gr.CloneURL,
		Branch:      gr.DefaultBranch,
		Topics:      &gr.Topics,
		Private:     gr.Private,
		Visibility:  &visibility,
		CustomProps: &gr.CustomProperties,
	}
}

// GetPullRequest defines a function that retrieves
// a pull request for a repo.
func (c *Client) GetPullRequest(ctx context.Context, r *api.Repo, number int, token string) (string, string, string, string, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": r.GetOwner().GetName(),
	}).Tracef("retrieving pull request %d for repo %s", number, r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newTokenClient(ctx, token)

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
	client := c.newUserOAuthTokenClient(ctx, u)

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
func (c *Client) GetBranch(ctx context.Context, r *api.Repo, branch, token string) (string, string, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
		"user": r.GetOwner().GetName(),
	}).Tracef("retrieving branch %s for repo %s", branch, r.GetFullName())

	// create GitHub OAuth client with user's token
	client := c.newTokenClient(ctx, token)

	maxRedirects := 3

	data, _, err := client.Repositories.GetBranch(ctx, r.GetOrg(), r.GetName(), branch, maxRedirects)
	if err != nil {
		return "", "", err
	}

	return data.GetName(), data.GetCommit().GetSHA(), nil
}

// GetNetrcPassword returns a scoped access token for cloning repos during a build.
// It takes the repo owner's user access token and scopes it down to the requested
// repositories and permissions using the GitHub App's client credentials.
func (c *Client) GetNetrcPassword(ctx context.Context, b *api.Build, repos []string, perms map[string]string) (string, error) {
	r := b.GetRepo()

	l := c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	})

	l.Tracef("getting netrc password for %s/%s", r.GetOrg(), r.GetName())

	// ensure build repo is included in list
	if !slices.Contains(repos, r.GetName()) {
		repos = append(repos, r.GetName())
	}

	// permissions that are applied to the scoped token
	//
	// the Vela compiler follows a least-privileged-defaults model where
	// the list contains only the triggering repo, unless provided in the git yaml block
	//
	// the below map is the default
	permissions := map[string]string{
		"contents": constants.PermissionRead,
	}

	if len(perms) > 0 {
		permissions = perms

		if _, ok := perms["contents"]; !ok {
			permissions["contents"] = "read"
		}
	}

	// create a scoped access token from the repo owner's user token
	scopedToken, err := c.CreateScopedAccessToken(ctx, r, repos, permissions)
	if err != nil {
		return "", fmt.Errorf("unable to create scoped access token for repos %v with permissions %v: %w", repos, permissions, err)
	}

	if len(scopedToken.Token) == 0 {
		return "", fmt.Errorf("scoped access token is empty for repo %s/%s", r.GetOrg(), r.GetName())
	}

	l.Tracef("using scoped access token for %s/%s", r.GetOrg(), r.GetName())

	return scopedToken.Token, nil
}

// SyncRepoWithInstallation ensures the repo is synchronized with the scm installation, if it exists.
func (c *Client) SyncRepoWithInstallation(ctx context.Context, r *api.Repo) (*api.Repo, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("syncing app installation for repo %s/%s", r.GetOrg(), r.GetName())

	// no GitHub App configured, skip
	if c.AppClient == nil {
		return r, nil
	}

	installations, _, err := c.AppClient.Apps.ListInstallations(ctx, &github.ListOptions{})
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

	installationCanReadRepo, err := c.installationCanReadRepo(ctx, r.GetOrg(), r.GetName(), installation)
	if err != nil {
		return r, err
	}

	if installationCanReadRepo {
		r.SetInstallID(installation.GetID())
	}

	return r, nil
}

// GeneratePermissionToken generates a token for checking permissions for an installation.
func (c *Client) GeneratePermissionToken(ctx context.Context, installID int64) (string, error) {
	if c.AppClient != nil && installID != 0 {
		opts := &github.InstallationTokenOptions{
			Permissions: &github.InstallationPermissions{
				Metadata:     new("read"),
				Contents:     new("read"),
				PullRequests: new("read"),
			},
		}
		// create installation token for the repo
		t, _, err := c.AppClient.Apps.CreateInstallationToken(ctx, installID, opts)
		if err != nil {
			return "", err
		}

		return t.GetToken(), nil
	}

	return "", fmt.Errorf("github app client not configured or install ID is 0")
}
