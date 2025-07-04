// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v73/github"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/internal"
)

// ProcessWebhook parses the webhook from a repo.
//
//nolint:nilerr // ignore webhook returning nil
func (c *Client) ProcessWebhook(ctx context.Context, request *http.Request) (*internal.Webhook, error) {
	c.Logger.Tracef("processing GitHub webhook")

	// create our own record of the hook and populate its fields
	h := new(api.Hook)
	h.SetNumber(1)
	h.SetSourceID(request.Header.Get("X-GitHub-Delivery"))

	hookID, err := strconv.Atoi(request.Header.Get("X-GitHub-Hook-ID"))
	if err != nil {
		return nil, fmt.Errorf("unable to convert hook id to int64: %w", err)
	}

	h.SetWebhookID(int64(hookID))
	h.SetCreated(time.Now().UTC().Unix())
	h.SetHost("github.com")
	h.SetEvent(request.Header.Get("X-GitHub-Event"))
	h.SetStatus(constants.StatusSuccess)

	if len(request.Header.Get("X-GitHub-Enterprise-Host")) > 0 {
		h.SetHost(request.Header.Get("X-GitHub-Enterprise-Host"))
	}

	// get content type
	contentType, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	payload, err := github.ValidatePayloadFromBody(contentType, request.Body, "", nil)
	if err != nil {
		return &internal.Webhook{Hook: h}, nil
	}

	// parse the payload from the webhook
	event, err := github.ParseWebHook(github.WebHookType(request), payload)
	if err != nil {
		return &internal.Webhook{Hook: h}, nil
	}

	// process the event from the webhook
	switch event := event.(type) {
	case *github.PushEvent:
		return c.processPushEvent(ctx, h, event)
	case *github.PullRequestEvent:
		return c.processPREvent(h, event)
	case *github.DeploymentEvent:
		return c.processDeploymentEvent(h, event)
	case *github.IssueCommentEvent:
		return c.processIssueCommentEvent(h, event)
	case *github.RepositoryEvent:
		return c.processRepositoryEvent(h, event)
	case *github.InstallationEvent:
		return c.processInstallationEvent(ctx, h, event)
	case *github.InstallationRepositoriesEvent:
		return c.processInstallationRepositoriesEvent(ctx, h, event)
	case *github.CustomPropertyValuesEvent:
		return c.processCustomPropertiesEvent(h, event)
	}

	return &internal.Webhook{Hook: h}, nil
}

// VerifyWebhook verifies the webhook from a repo.
func (c *Client) VerifyWebhook(_ context.Context, request *http.Request, secret []byte) error {
	_, err := github.ValidatePayload(request, secret)
	if err != nil {
		return err
	}

	return nil
}

// RedeliverWebhook redelivers webhooks from GitHub.
func (c *Client) RedeliverWebhook(ctx context.Context, u *api.User, h *api.Hook) error {
	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, u.GetToken())

	// capture the delivery ID of the hook using GitHub API
	deliveryID, err := c.getDeliveryID(ctx, client, h)
	if err != nil {
		return err
	}

	// redeliver the webhook
	_, _, err = client.Repositories.RedeliverHookDelivery(
		ctx,
		h.GetRepo().GetOrg(),
		h.GetRepo().GetName(),
		h.GetWebhookID(), deliveryID,
	)

	if err != nil {
		var acceptedError *github.AcceptedError
		// Persist if the status received is a 202 Accepted. This
		// means the job was added to the queue for GitHub.
		if errors.As(err, &acceptedError) {
			return nil
		}

		return err
	}

	return nil
}

// processPushEvent is a helper function to process the push event.
func (c *Client) processPushEvent(_ context.Context, h *api.Hook, payload *github.PushEvent) (*internal.Webhook, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  payload.GetRepo().GetOwner().GetLogin(),
		"repo": payload.GetRepo().GetName(),
	}).Tracef("processing push GitHub webhook for %s", payload.GetRepo().GetFullName())

	repo := payload.GetRepo()
	if repo == nil {
		return &internal.Webhook{Hook: h}, nil
	}

	// convert payload to API repo
	r := new(api.Repo)
	r.SetOrg(repo.GetOwner().GetLogin())
	r.SetName(repo.GetName())
	r.SetFullName(repo.GetFullName())
	r.SetLink(repo.GetHTMLURL())
	r.SetClone(repo.GetCloneURL())
	r.SetBranch(repo.GetDefaultBranch())
	r.SetPrivate(repo.GetPrivate())
	r.SetTopics(repo.Topics)
	r.SetCustomProps(repo.CustomProperties)

	// convert payload to API build
	b := new(api.Build)
	b.SetEvent(constants.EventPush)
	b.SetClone(repo.GetCloneURL())
	b.SetSource(payload.GetHeadCommit().GetURL())
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventPush, repo.GetHTMLURL()))
	b.SetMessage(payload.GetHeadCommit().GetMessage())
	b.SetCommit(payload.GetHeadCommit().GetID())
	b.SetSender(payload.GetSender().GetLogin())
	b.SetSenderSCMID(fmt.Sprint(payload.GetSender().GetID()))
	b.SetAuthor(payload.GetHeadCommit().GetAuthor().GetLogin())
	b.SetEmail(payload.GetHeadCommit().GetAuthor().GetEmail())
	b.SetBranch(strings.TrimPrefix(payload.GetRef(), "refs/heads/"))
	b.SetRef(payload.GetRef())
	b.SetBaseRef(payload.GetBaseRef())

	// update the hook object
	h.SetBranch(b.GetBranch())
	h.SetEvent(constants.EventPush)
	h.SetLink(
		fmt.Sprintf("https://%s/%s/settings/hooks", h.GetHost(), r.GetFullName()),
	)

	// ensure the build author is set
	if len(b.GetAuthor()) == 0 {
		b.SetAuthor(payload.GetHeadCommit().GetCommitter().GetName())
	}

	// ensure the build sender is set
	if len(b.GetSender()) == 0 {
		b.SetSender(payload.GetPusher().GetName())
	}

	// ensure the build email is set
	if len(b.GetEmail()) == 0 {
		b.SetEmail(payload.GetHeadCommit().GetCommitter().GetEmail())
	}

	// handle when push event is a tag
	if strings.HasPrefix(b.GetRef(), "refs/tags/") {
		// set the proper event for the hook
		h.SetEvent(constants.EventTag)
		// set the proper event for the build
		b.SetEvent(constants.EventTag)

		// set the proper branch from the base ref
		if strings.HasPrefix(payload.GetBaseRef(), "refs/heads/") {
			b.SetBranch(strings.TrimPrefix(payload.GetBaseRef(), "refs/heads/"))
		}
	}

	// handle when push event is a delete
	if strings.EqualFold(b.GetCommit(), "") {
		b.SetCommit(payload.GetBefore())
		b.SetRef(payload.GetBefore())
		b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventDelete, repo.GetHTMLURL()))
		b.SetAuthor(payload.GetSender().GetLogin())
		b.SetSource(fmt.Sprintf("%s/commit/%s", payload.GetRepo().GetHTMLURL(), payload.GetBefore()))
		b.SetEmail(payload.GetPusher().GetEmail())

		// set the proper event for the hook
		h.SetEvent(constants.EventDelete)
		// set the proper event for the build
		b.SetEvent(constants.EventDelete)

		if strings.HasPrefix(payload.GetRef(), "refs/tags/") {
			b.SetBranch(strings.TrimPrefix(payload.GetRef(), "refs/tags/"))
			// set the proper action for the build
			b.SetEventAction(constants.ActionTag)
			// set the proper message for the build
			b.SetMessage(fmt.Sprintf("%s %s deleted", strings.TrimPrefix(payload.GetRef(), "refs/tags/"), constants.ActionTag))
		} else {
			// set the proper action for the build
			b.SetEventAction(constants.ActionBranch)
			// set the proper message for the build
			b.SetMessage(fmt.Sprintf("%s %s deleted", strings.TrimPrefix(payload.GetRef(), "refs/heads/"), constants.ActionBranch))
		}
	}

	return &internal.Webhook{
		Hook:  h,
		Repo:  r,
		Build: b,
	}, nil
}

// processPREvent is a helper function to process the pull_request event.
func (c *Client) processPREvent(h *api.Hook, payload *github.PullRequestEvent) (*internal.Webhook, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  payload.GetRepo().GetOwner().GetLogin(),
		"repo": payload.GetRepo().GetName(),
	}).Tracef("processing pull_request GitHub webhook for %s", payload.GetRepo().GetFullName())

	// update the hook object
	h.SetBranch(payload.GetPullRequest().GetBase().GetRef())
	h.SetEvent(constants.EventPull)
	h.SetLink(
		fmt.Sprintf("https://%s/%s/settings/hooks", h.GetHost(), payload.GetRepo().GetFullName()),
	)

	// if the pull request state isn't open we ignore it
	if payload.GetPullRequest().GetState() != "open" {
		return &internal.Webhook{Hook: h}, nil
	}

	// skip if the pull request action is not opened, synchronize, reopened, edited, labeled, or unlabeled
	if !strings.EqualFold(payload.GetAction(), "opened") &&
		!strings.EqualFold(payload.GetAction(), "synchronize") &&
		!strings.EqualFold(payload.GetAction(), "reopened") &&
		!strings.EqualFold(payload.GetAction(), "edited") &&
		!strings.EqualFold(payload.GetAction(), "labeled") &&
		!strings.EqualFold(payload.GetAction(), "unlabeled") {
		return &internal.Webhook{Hook: h}, nil
	}

	// capture the repo from the payload
	repo := payload.GetRepo()
	if repo == nil {
		return &internal.Webhook{Hook: h}, nil
	}

	// convert payload to API repo
	r := new(api.Repo)
	r.SetOrg(repo.GetOwner().GetLogin())
	r.SetName(repo.GetName())
	r.SetFullName(repo.GetFullName())
	r.SetLink(repo.GetHTMLURL())
	r.SetClone(repo.GetCloneURL())
	r.SetBranch(repo.GetDefaultBranch())
	r.SetPrivate(repo.GetPrivate())
	r.SetTopics(repo.Topics)
	r.SetCustomProps(repo.CustomProperties)

	// convert payload to api build
	b := new(api.Build)
	b.SetEvent(constants.EventPull)
	b.SetEventAction(payload.GetAction())
	b.SetClone(repo.GetCloneURL())
	b.SetSource(payload.GetPullRequest().GetHTMLURL())
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventPull, repo.GetHTMLURL()))
	b.SetMessage(payload.GetPullRequest().GetTitle())
	b.SetCommit(payload.GetPullRequest().GetHead().GetSHA())
	b.SetSender(payload.GetSender().GetLogin())
	b.SetSenderSCMID(fmt.Sprint(payload.GetSender().GetID()))
	b.SetAuthor(payload.GetPullRequest().GetUser().GetLogin())
	b.SetEmail(payload.GetPullRequest().GetUser().GetEmail())
	b.SetBranch(payload.GetPullRequest().GetBase().GetRef())
	b.SetRef(fmt.Sprintf("refs/pull/%d/head", payload.GetNumber()))
	b.SetBaseRef(payload.GetPullRequest().GetBase().GetRef())
	b.SetHeadRef(payload.GetPullRequest().GetHead().GetRef())

	// ensure the build reference is set
	if payload.GetPullRequest().GetMerged() {
		b.SetRef(fmt.Sprintf("refs/pull/%d/merge", payload.GetNumber()))
	}

	// ensure the build author is set
	if len(b.GetAuthor()) == 0 {
		b.SetAuthor(payload.GetPullRequest().GetHead().GetUser().GetLogin())
	}

	// ensure the build sender is set
	if len(b.GetSender()) == 0 {
		b.SetSender(payload.GetPullRequest().GetUser().GetLogin())
		b.SetSenderSCMID(fmt.Sprint(payload.GetPullRequest().GetUser().GetID()))
	}

	// ensure the build email is set
	if len(b.GetEmail()) == 0 {
		b.SetEmail(payload.GetPullRequest().GetHead().GetUser().GetEmail())
	}

	var prLabels []string
	if strings.EqualFold(payload.GetAction(), "labeled") ||
		strings.EqualFold(payload.GetAction(), "unlabeled") {
		prLabels = append(prLabels, payload.GetLabel().GetName())
	} else {
		labels := payload.GetPullRequest().Labels
		for _, label := range labels {
			prLabels = append(prLabels, label.GetName())
		}
	}

	// determine if pull request head is a fork and does not match the repo name of base

	b.SetFork(payload.GetPullRequest().GetHead().GetRepo().GetFork() &&
		!strings.EqualFold(payload.GetPullRequest().GetBase().GetRepo().GetFullName(), payload.GetPullRequest().GetHead().GetRepo().GetFullName()))

	return &internal.Webhook{
		PullRequest: internal.PullRequest{
			Number: int64(payload.GetNumber()),
			Labels: prLabels,
		},
		Hook:  h,
		Repo:  r,
		Build: b,
	}, nil
}

// processDeploymentEvent is a helper function to process the deployment event.
func (c *Client) processDeploymentEvent(h *api.Hook, payload *github.DeploymentEvent) (*internal.Webhook, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  payload.GetRepo().GetOwner().GetLogin(),
		"repo": payload.GetRepo().GetName(),
	}).Tracef("processing deployment GitHub webhook for %s", payload.GetRepo().GetFullName())

	// capture the repo from the payload
	repo := payload.GetRepo()
	if repo == nil {
		return &internal.Webhook{Hook: h}, nil
	}

	// convert payload to API repo
	r := new(api.Repo)
	r.SetOrg(repo.GetOwner().GetLogin())
	r.SetName(repo.GetName())
	r.SetFullName(repo.GetFullName())
	r.SetLink(repo.GetHTMLURL())
	r.SetClone(repo.GetCloneURL())
	r.SetBranch(repo.GetDefaultBranch())
	r.SetPrivate(repo.GetPrivate())
	r.SetTopics(repo.Topics)
	r.SetCustomProps(repo.CustomProperties)

	// convert payload to api build
	b := new(api.Build)
	b.SetEvent(constants.EventDeploy)
	b.SetEventAction(constants.ActionCreated)
	b.SetClone(repo.GetCloneURL())
	b.SetDeploy(payload.GetDeployment().GetEnvironment())
	b.SetDeployNumber(payload.GetDeployment().GetID())
	b.SetSource(payload.GetDeployment().GetURL())
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventDeploy, repo.GetHTMLURL()))
	b.SetMessage(payload.GetDeployment().GetDescription())
	b.SetCommit(payload.GetDeployment().GetSHA())
	b.SetSender(payload.GetSender().GetLogin())
	b.SetSenderSCMID(fmt.Sprint(payload.GetSender().GetID()))

	b.SetAuthor(payload.GetDeployment().GetCreator().GetLogin())
	b.SetEmail(payload.GetDeployment().GetCreator().GetEmail())
	b.SetBranch(payload.GetDeployment().GetRef())
	b.SetRef(payload.GetDeployment().GetRef())

	d := new(api.Deployment)

	d.SetNumber(payload.GetDeployment().GetID())
	d.SetURL(payload.GetDeployment().GetURL())
	d.SetCommit(payload.GetDeployment().GetSHA())
	d.SetRef(b.GetRef())
	d.SetTask(payload.GetDeployment().GetTask())
	d.SetTarget(payload.GetDeployment().GetEnvironment())
	d.SetDescription(payload.GetDeployment().GetDescription())
	d.SetCreatedAt(time.Now().Unix())
	d.SetCreatedBy(payload.GetDeployment().GetCreator().GetLogin())

	// check if payload is provided within request
	//
	// use a length of 2 because the payload will
	// never be nil even if no payload is provided.
	//
	// sending an API request to GitHub with no
	// payload provided yields a default of `{}`.
	if len(payload.GetDeployment().Payload) > 2 {
		deployPayload := make(map[string]string)
		// unmarshal the payload into the expected map[string]string format
		err := json.Unmarshal(payload.GetDeployment().Payload, &deployPayload)
		if err != nil {
			return &internal.Webhook{}, err
		}

		// check if the map is empty
		if len(deployPayload) != 0 {
			// set the payload info on the build
			b.SetDeployPayload(deployPayload)
			// set payload info on the deployment
			d.SetPayload(deployPayload)
		}
	}

	// handle when the ref is a sha or short sha
	if strings.HasPrefix(b.GetCommit(), b.GetRef()) || b.GetCommit() == b.GetRef() {
		// set the proper branch for the build
		b.SetBranch(r.GetBranch())
		// set the proper ref for the build
		b.SetRef(fmt.Sprintf("refs/heads/%s", b.GetBranch()))
	}

	// handle when the ref is a branch
	if !strings.HasPrefix(b.GetRef(), "refs/") {
		// set the proper ref for the build
		b.SetRef(fmt.Sprintf("refs/heads/%s", b.GetBranch()))
	}

	// update the hook object
	h.SetBranch(b.GetBranch())
	h.SetEvent(constants.EventDeploy)
	h.SetEventAction(constants.ActionCreated)
	h.SetLink(
		fmt.Sprintf("https://%s/%s/settings/hooks", h.GetHost(), r.GetFullName()),
	)

	return &internal.Webhook{
		Hook:       h,
		Repo:       r,
		Build:      b,
		Deployment: d,
	}, nil
}

// processIssueCommentEvent is a helper function to process the issue comment event.
func (c *Client) processIssueCommentEvent(h *api.Hook, payload *github.IssueCommentEvent) (*internal.Webhook, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  payload.GetRepo().GetOwner().GetLogin(),
		"repo": payload.GetRepo().GetName(),
	}).Tracef("processing issue_comment GitHub webhook for %s", payload.GetRepo().GetFullName())

	// update the hook object
	h.SetEvent(constants.EventComment)
	h.SetLink(
		fmt.Sprintf("https://%s/%s/settings/hooks", h.GetHost(), payload.GetRepo().GetFullName()),
	)

	// skip if the comment action is deleted or not part of a pull request
	if strings.EqualFold(payload.GetAction(), "deleted") || !payload.GetIssue().IsPullRequest() {
		// return &internal.Webhook{Hook: h}, nil
		return &internal.Webhook{
			Hook: h,
		}, nil
	}

	// capture the repo from the payload
	repo := payload.GetRepo()
	if repo == nil {
		return &internal.Webhook{Hook: h}, nil
	}

	// convert payload to API repo
	r := new(api.Repo)
	r.SetOrg(repo.GetOwner().GetLogin())
	r.SetName(repo.GetName())
	r.SetFullName(repo.GetFullName())
	r.SetLink(repo.GetHTMLURL())
	r.SetClone(repo.GetCloneURL())
	r.SetBranch(repo.GetDefaultBranch())
	r.SetPrivate(repo.GetPrivate())
	r.SetTopics(repo.Topics)
	r.SetCustomProps(repo.CustomProperties)

	// convert payload to API build
	b := new(api.Build)
	b.SetEvent(constants.EventComment)
	b.SetEventAction(payload.GetAction())
	b.SetClone(repo.GetCloneURL())
	b.SetSource(payload.Issue.GetHTMLURL())
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventComment, repo.GetHTMLURL()))
	b.SetMessage(payload.Issue.GetTitle())
	b.SetSender(payload.GetSender().GetLogin())
	b.SetSenderSCMID(fmt.Sprint(payload.GetSender().GetID()))
	b.SetAuthor(payload.GetIssue().GetUser().GetLogin())
	b.SetEmail(payload.GetIssue().GetUser().GetEmail())
	b.SetRef(fmt.Sprintf("refs/pull/%d/head", payload.GetIssue().GetNumber()))

	return &internal.Webhook{
		PullRequest: internal.PullRequest{
			Comment: payload.GetComment().GetBody(),
			Number:  int64(payload.GetIssue().GetNumber()),
		},
		Hook:  h,
		Repo:  r,
		Build: b,
	}, nil
}

// processRepositoryEvent is a helper function to process the repository event.
func (c *Client) processRepositoryEvent(h *api.Hook, payload *github.RepositoryEvent) (*internal.Webhook, error) {
	logrus.Tracef("processing repository event GitHub webhook for %s", payload.GetRepo().GetFullName())

	repo := payload.GetRepo()
	if repo == nil {
		return &internal.Webhook{Hook: h}, nil
	}

	// convert payload to API repo
	r := new(api.Repo)
	r.SetOrg(repo.GetOwner().GetLogin())
	r.SetName(repo.GetName())
	r.SetFullName(repo.GetFullName())
	r.SetLink(repo.GetHTMLURL())
	r.SetClone(repo.GetCloneURL())
	r.SetBranch(repo.GetDefaultBranch())
	r.SetPrivate(repo.GetPrivate())
	r.SetActive(!repo.GetArchived())
	r.SetTopics(repo.Topics)
	r.SetCustomProps(repo.CustomProperties)

	h.SetEvent(constants.EventRepository)
	h.SetEventAction(payload.GetAction())
	h.SetBranch(r.GetBranch())
	h.SetLink(
		fmt.Sprintf("https://%s/%s/settings/hooks", h.GetHost(), r.GetFullName()),
	)

	return &internal.Webhook{
		Hook: h,
		Repo: r,
	}, nil
}

// processCustomPropertiesEvent is a helper function to process the custom properties event.
func (c *Client) processCustomPropertiesEvent(h *api.Hook, payload *github.CustomPropertyValuesEvent) (*internal.Webhook, error) {
	logrus.Tracef("processing repository event GitHub webhook for %s", payload.GetRepo().GetFullName())

	repo := payload.GetRepo()
	if repo == nil {
		return &internal.Webhook{Hook: h}, nil
	}

	// convert payload to API repo
	r := new(api.Repo)
	r.SetOrg(repo.GetOwner().GetLogin())
	r.SetName(repo.GetName())
	r.SetFullName(repo.GetFullName())
	r.SetLink(repo.GetHTMLURL())
	r.SetClone(repo.GetCloneURL())
	r.SetBranch(repo.GetDefaultBranch())
	r.SetPrivate(repo.GetPrivate())
	r.SetActive(!repo.GetArchived())
	r.SetTopics(repo.Topics)
	r.SetCustomProps(repo.CustomProperties)

	h.SetEvent(constants.EventCustomProperties)
	h.SetEventAction(payload.GetAction())
	h.SetLink(
		fmt.Sprintf("https://%s/%s/settings/hooks", h.GetHost(), r.GetFullName()),
	)

	return &internal.Webhook{
		Hook: h,
		Repo: r,
	}, nil
}

// processInstallationEvent is a helper function to process the installation event.
func (c *Client) processInstallationEvent(_ context.Context, h *api.Hook, payload *github.InstallationEvent) (*internal.Webhook, error) {
	h.SetEvent(constants.EventInstallation)
	h.SetEventAction(payload.GetAction())

	install := new(internal.Installation)

	install.Action = payload.GetAction()
	install.ID = payload.GetInstallation().GetID()
	install.Org = payload.GetInstallation().GetAccount().GetLogin()

	switch payload.GetAction() {
	case constants.AppInstallCreated:
		for _, repo := range payload.Repositories {
			install.RepositoriesAdded = append(install.RepositoriesAdded, repo.GetName())
		}
	case constants.AppInstallDeleted:
		for _, repo := range payload.Repositories {
			install.RepositoriesRemoved = append(install.RepositoriesRemoved, repo.GetName())
		}
	}

	return &internal.Webhook{
		Hook:         h,
		Installation: install,
	}, nil
}

// processInstallationRepositoriesEvent is a helper function to process the installation repositories event.
func (c *Client) processInstallationRepositoriesEvent(_ context.Context, h *api.Hook, payload *github.InstallationRepositoriesEvent) (*internal.Webhook, error) {
	h.SetEvent(constants.EventInstallationRepositories)
	h.SetEventAction(payload.GetAction())

	install := new(internal.Installation)

	install.Action = payload.GetAction()
	install.ID = payload.GetInstallation().GetID()
	install.Org = payload.GetInstallation().GetAccount().GetLogin()

	for _, repo := range payload.RepositoriesAdded {
		install.RepositoriesAdded = append(install.RepositoriesAdded, repo.GetName())
	}

	for _, repo := range payload.RepositoriesRemoved {
		install.RepositoriesRemoved = append(install.RepositoriesRemoved, repo.GetName())
	}

	return &internal.Webhook{
		Hook:         h,
		Installation: install,
	}, nil
}

// getDeliveryID gets the last 100 webhook deliveries for a repo and
// finds the matching delivery id with the source id in the hook.
func (c *Client) getDeliveryID(ctx context.Context, ghClient *github.Client, h *api.Hook) (int64, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  h.GetRepo().GetOrg(),
		"repo": h.GetRepo().GetName(),
	}).Tracef("searching for delivery id for hook: %s", h.GetSourceID())

	// set per page to 100 to retrieve last 100 hook summaries
	opt := &github.ListCursorOptions{PerPage: 100}

	// send API call to capture delivery summaries that contain Delivery ID value
	deliveries, resp, err := ghClient.Repositories.ListHookDeliveries(
		ctx,
		h.GetRepo().GetOrg(),
		h.GetRepo().GetName(),
		h.GetWebhookID(),
		opt,
	)

	// version check: if GitHub API is older than version 3.2, this call will not work
	if resp.StatusCode == 415 {
		err = fmt.Errorf("requires GitHub version 3.2 or later")
		return 0, err
	}

	if err != nil {
		return 0, err
	}

	// cycle through delivery summaries and match Source ID/GUID. Capture Delivery ID
	for _, delivery := range deliveries {
		if delivery.GetGUID() == h.GetSourceID() {
			return delivery.GetID(), nil
		}
	}

	// if not found, webhook was not recent enough for GitHub
	err = fmt.Errorf("webhook no longer available to be redelivered")

	return 0, err
}
