// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/go-github/v81/github"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/cache"
	"github.com/go-vela/server/compiler/types/yaml"
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
func (c *Client) GetNetrcPassword(ctx context.Context, db database.Interface, tknCache cache.Service, b *api.Build, g yaml.Git) (string, int64, error) {
	r := b.GetRepo()
	u := b.GetRepo().GetOwner()

	l := c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	})

	l.Tracef("getting netrc password for %s/%s", r.GetOrg(), r.GetName())

	// no GitHub App configured, use legacy oauth token
	if c.AppClient == nil {
		return u.GetToken(), 0, nil
	}

	var err error

	// repos that the token has access to
	// providing no repos, nil, or empty slice will default the token permissions to the list
	// of repos added to the installation
	repos := g.Repositories

	// enforce max number of repos allowed for token
	//
	// this prevents a large number of access checks for the repo owner
	if len(repos) > constants.GitTokenRepoLimit {
		return u.GetToken(), 0, fmt.Errorf("number of repositories specified (%d) exceeds the maximum allowed (%d)", len(repos), constants.GitTokenRepoLimit)
	}

	// ensure build repo is included in list
	if !slices.Contains(repos, r.GetName()) {
		repos = append(repos, r.GetName())
	}

	// permissions that are applied to the token for every repo provided
	// providing no permissions, nil, or empty map will default to the permissions
	// of the GitHub App installation
	//
	// the Vela compiler follows a least-privileged-defaults model where
	// the list contains only the triggering repo, unless provided in the git yaml block
	//
	// the below map is the default
	permissions := map[string]string{
		"contents": constants.PermissionRead,
		"checks":   constants.PermissionWrite,
		"statuses": constants.PermissionWrite,
	}

	if b.GetEvent() == constants.EventDeploy {
		permissions["deployments"] = constants.PermissionWrite
	}

	if len(g.Permissions) > 0 {
		permissions = g.Permissions

		normalizePermissions(permissions)
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
			return u.GetToken(), 0, fmt.Errorf("repository owner does not have adequate permissions to request install token for repository %s/%s", r.GetOrg(), repo)
		}
	}

	id := r.GetInstallID()

	// if the source scm repo has an install ID but the Vela db record does not
	// then use the source repo to create an installation token
	if id == 0 {
		// list all installations (a.k.a. orgs) where the GitHub App is installed
		installations, _, err := c.AppClient.Apps.ListInstallations(ctx, &github.ListOptions{})
		if err != nil {
			l.Tracef("unable to list github app installations: %s", err.Error())

			return u.GetToken(), 0, err
		}

		// iterate through the list of installations
		for _, install := range installations {
			// find the installation that matches the org for the repo
			if strings.EqualFold(install.GetAccount().GetLogin(), r.GetOrg()) {
				if install.GetRepositorySelection() == constants.AppInstallRepositoriesSelectionSelected {
					installationCanReadRepo, err := c.installationCanReadRepo(ctx, r, install)
					if err != nil {
						l.Tracef("unable to check if installation for org %s can read repo %s: %s", install.GetAccount().GetLogin(), r.GetFullName(), err.Error())

						return u.GetToken(), 0, nil
					}

					if !installationCanReadRepo {
						l.Tracef("installation for org %s exists but does not have access to repo %s", install.GetAccount().GetLogin(), r.GetFullName())

						return u.GetToken(), 0, nil
					}
				}

				id = install.GetID()
			}
		}
	}

	// the app might not be installed therefore we retain backwards compatibility via the user oauth token
	// https://docs.github.com/en/apps/creating-github-apps/authenticating-with-a-github-app/authenticating-as-a-github-app-installation
	// the optional list of repos and permissions are driven by yaml
	installToken, err := c.NewAppInstallationToken(ctx, id, repos, permissions)
	if err != nil {
		// return the legacy token along with no error for backwards compatibility
		// todo: return an error based based on app installation requirements
		l.Tracef("unable to create github app installation token for repos %v with permissions %v: %v", repos, permissions, err)

		return u.GetToken(), 0, nil
	}

	if installToken != nil && len(installToken.Token) != 0 {
		l.Tracef("using github app installation token for %s/%s", r.GetOrg(), r.GetName())

		// (optional) sync the install ID with the repo
		if db != nil && r.GetInstallID() != id {
			r.SetInstallID(id)

			_, err = db.UpdateRepo(ctx, r)
			if err != nil {
				c.Logger.Tracef("unable to update repo with install ID %d: %v", id, err)
			}
		}

		if tknCache != nil {
			err = tknCache.StoreInstallToken(ctx, installToken, r.GetTimeout())
			if err != nil {
				l.Tracef("unable to store installation token in cache: %v", err)

				return "", 0, fmt.Errorf("unable to store installation token in cache: %w", err)
			}
		}

		return installToken.Token, installToken.Expiration, nil
	}

	l.Tracef("using user oauth token for %s/%s", r.GetOrg(), r.GetName())

	return u.GetToken(), 0, nil
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

	installationCanReadRepo, err := c.installationCanReadRepo(ctx, r, installation)
	if err != nil {
		return r, err
	}

	if installationCanReadRepo {
		r.SetInstallID(installation.GetID())
	}

	return r, nil
}

// normalizePermissions ensures minimum required permissions are set for installation token generation.
func normalizePermissions(perms map[string]string) {
	if _, ok := perms["contents"]; !ok {
		perms["contents"] = "read"
	}

	if _, ok := perms["checks"]; !ok {
		perms["checks"] = "write"
	}

	if _, ok := perms["statuses"]; !ok {
		perms["statuses"] = "write"
	}
}
