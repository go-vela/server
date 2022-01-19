// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/google/go-github/v42/github"
)

// ProcessWebhook parses the webhook from a repo.
func (c *client) ProcessWebhook(request *http.Request) (*types.Webhook, error) {
	c.Logger.Tracef("processing GitHub webhook")

	h := new(library.Hook)
	h.SetNumber(1)
	h.SetSourceID(request.Header.Get("X-GitHub-Delivery"))
	h.SetCreated(time.Now().UTC().Unix())
	h.SetHost("github.com")
	h.SetEvent(request.Header.Get("X-GitHub-Event"))
	h.SetStatus(constants.StatusSuccess)

	if len(request.Header.Get("X-GitHub-Enterprise-Host")) > 0 {
		h.SetHost(request.Header.Get("X-GitHub-Enterprise-Host"))
	}

	payload, err := github.ValidatePayload(request, nil)
	if err != nil {
		return &types.Webhook{Hook: h}, nil
	}

	// parse the payload from the webhook
	event, err := github.ParseWebHook(github.WebHookType(request), payload)
	if err != nil {
		return &types.Webhook{Hook: h}, nil
	}

	// process the event from the webhook
	switch event := event.(type) {
	case *github.PushEvent:
		return c.processPushEvent(h, event)
	case *github.PullRequestEvent:
		return c.processPREvent(h, event)
	case *github.DeploymentEvent:
		return c.processDeploymentEvent(h, event)
	case *github.IssueCommentEvent:
		return c.processIssueCommentEvent(h, event)
	}

	return &types.Webhook{Hook: h}, nil
}

// VerifyWebhook verifies the webhook from a repo.
func (c *client) VerifyWebhook(request *http.Request, r *library.Repo) error {
	c.Logger.WithFields(logrus.Fields{
		"org":  r.GetOrg(),
		"repo": r.GetName(),
	}).Tracef("verifying GitHub webhook for %s", r.GetFullName())

	_, err := github.ValidatePayload(request, []byte(r.GetHash()))
	if err != nil {
		return err
	}

	return nil
}

// processPushEvent is a helper function to process the push event.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) processPushEvent(h *library.Hook, payload *github.PushEvent) (*types.Webhook, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  payload.GetRepo().GetOwner().GetLogin(),
		"repo": payload.GetRepo().GetName(),
	}).Tracef("processing push GitHub webhook for %s", payload.GetRepo().GetFullName())

	repo := payload.GetRepo()

	// convert payload to library repo
	r := new(library.Repo)
	r.SetOrg(repo.GetOwner().GetLogin())
	r.SetName(repo.GetName())
	r.SetFullName(repo.GetFullName())
	r.SetLink(repo.GetHTMLURL())
	r.SetClone(repo.GetCloneURL())
	r.SetBranch(repo.GetDefaultBranch())
	r.SetPrivate(repo.GetPrivate())

	// convert payload to library build
	b := new(library.Build)
	b.SetEvent(constants.EventPush)
	b.SetClone(repo.GetCloneURL())
	b.SetSource(payload.GetHeadCommit().GetURL())
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventPush, repo.GetHTMLURL()))
	b.SetMessage(payload.GetHeadCommit().GetMessage())
	b.SetCommit(payload.GetHeadCommit().GetID())
	b.SetSender(payload.GetSender().GetLogin())
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

	return &types.Webhook{
		Comment: "",
		Hook:    h,
		Repo:    r,
		Build:   b,
	}, nil
}

// processPREvent is a helper function to process the pull_request event.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) processPREvent(h *library.Hook, payload *github.PullRequestEvent) (*types.Webhook, error) {
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
		return &types.Webhook{Hook: h}, nil
	}

	// skip if the pull request action is not opened or synchronize
	if !strings.EqualFold(payload.GetAction(), "opened") &&
		!strings.EqualFold(payload.GetAction(), "synchronize") {
		return &types.Webhook{Hook: h}, nil
	}

	// capture the repo from the payload
	repo := payload.GetRepo()

	// convert payload to library repo
	r := new(library.Repo)
	r.SetOrg(repo.GetOwner().GetLogin())
	r.SetName(repo.GetName())
	r.SetFullName(repo.GetFullName())
	r.SetLink(repo.GetHTMLURL())
	r.SetClone(repo.GetCloneURL())
	r.SetBranch(repo.GetDefaultBranch())
	r.SetPrivate(repo.GetPrivate())

	// convert payload to library build
	b := new(library.Build)
	b.SetEvent(constants.EventPull)
	b.SetClone(repo.GetCloneURL())
	b.SetSource(payload.GetPullRequest().GetHTMLURL())
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventPull, repo.GetHTMLURL()))
	b.SetMessage(payload.GetPullRequest().GetTitle())
	b.SetCommit(payload.GetPullRequest().GetHead().GetSHA())
	b.SetSender(payload.GetSender().GetLogin())
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
	}

	// ensure the build email is set
	if len(b.GetEmail()) == 0 {
		b.SetEmail(payload.GetPullRequest().GetHead().GetUser().GetEmail())
	}

	return &types.Webhook{
		Comment:  "",
		PRNumber: payload.GetNumber(),
		Hook:     h,
		Repo:     r,
		Build:    b,
	}, nil
}

// processDeploymentEvent is a helper function to process the deployment event.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) processDeploymentEvent(h *library.Hook, payload *github.DeploymentEvent) (*types.Webhook, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  payload.GetRepo().GetOwner().GetLogin(),
		"repo": payload.GetRepo().GetName(),
	}).Tracef("processing deployment GitHub webhook for %s", payload.GetRepo().GetFullName())

	// capture the repo from the payload
	repo := payload.GetRepo()

	// convert payload to library repo
	r := new(library.Repo)
	r.SetOrg(repo.GetOwner().GetLogin())
	r.SetName(repo.GetName())
	r.SetFullName(repo.GetFullName())
	r.SetLink(repo.GetHTMLURL())
	r.SetClone(repo.GetCloneURL())
	r.SetBranch(repo.GetDefaultBranch())
	r.SetPrivate(repo.GetPrivate())

	// convert payload to library build
	b := new(library.Build)
	b.SetEvent(constants.EventDeploy)
	b.SetClone(repo.GetCloneURL())
	b.SetDeploy(payload.GetDeployment().GetEnvironment())
	b.SetSource(payload.GetDeployment().GetURL())
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventDeploy, repo.GetHTMLURL()))
	b.SetMessage(payload.GetDeployment().GetDescription())
	b.SetCommit(payload.GetDeployment().GetSHA())
	b.SetSender(payload.GetSender().GetLogin())
	b.SetAuthor(payload.GetDeployment().GetCreator().GetLogin())
	b.SetEmail(payload.GetDeployment().GetCreator().GetEmail())
	b.SetBranch(payload.GetDeployment().GetRef())
	b.SetRef(payload.GetDeployment().GetRef())

	// check if payload is provided within request
	//
	// use a length of 2 because the payload will
	// never be nil even if no payload is provided.
	//
	// sending an API request to GitHub with no
	// payload provided yields a default of `{}`.
	//
	// nolint: gomnd // ignore magic number
	if len(payload.GetDeployment().Payload) > 2 {
		deployPayload := make(map[string]string)
		// unmarshal the payload into the expected map[string]string format
		err := json.Unmarshal(payload.GetDeployment().Payload, &deployPayload)
		if err != nil {
			return &types.Webhook{}, err
		}

		// check if the map is empty
		if len(deployPayload) != 0 {
			// set the payload info on the build
			b.SetDeployPayload(deployPayload)
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
	h.SetLink(
		fmt.Sprintf("https://%s/%s/settings/hooks", h.GetHost(), r.GetFullName()),
	)

	return &types.Webhook{
		Comment: "",
		Hook:    h,
		Repo:    r,
		Build:   b,
	}, nil
}

// processIssueCommentEvent is a helper function to process the issue comment event.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) processIssueCommentEvent(h *library.Hook, payload *github.IssueCommentEvent) (*types.Webhook, error) {
	c.Logger.WithFields(logrus.Fields{
		"org":  payload.GetRepo().GetOwner().GetLogin(),
		"repo": payload.GetRepo().GetName(),
	}).Tracef("processing issue_comment GitHub webhook for %s", payload.GetRepo().GetFullName())

	// update the hook object
	h.SetEvent(constants.EventComment)
	h.SetLink(
		fmt.Sprintf("https://%s/%s/settings/hooks", h.GetHost(), payload.GetRepo().GetFullName()),
	)

	// skip if the comment action is deleted
	if strings.EqualFold(payload.GetAction(), "deleted") {
		// return &types.Webhook{Hook: h}, nil
		return &types.Webhook{
			Comment: payload.GetComment().GetBody(),
			Hook:    h,
		}, nil
	}

	// capture the repo from the payload
	repo := payload.GetRepo()

	// convert payload to library repo
	r := new(library.Repo)
	r.SetOrg(repo.GetOwner().GetLogin())
	r.SetName(repo.GetName())
	r.SetFullName(repo.GetFullName())
	r.SetLink(repo.GetHTMLURL())
	r.SetClone(repo.GetCloneURL())
	r.SetBranch(repo.GetDefaultBranch())
	r.SetPrivate(repo.GetPrivate())

	// convert payload to library build
	b := new(library.Build)
	b.SetEvent(constants.EventComment)
	b.SetClone(repo.GetCloneURL())
	b.SetSource(payload.Issue.GetHTMLURL())
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventComment, repo.GetHTMLURL()))
	b.SetMessage(payload.Issue.GetTitle())
	b.SetSender(payload.GetSender().GetLogin())
	b.SetAuthor(payload.GetIssue().GetUser().GetLogin())
	b.SetEmail(payload.GetIssue().GetUser().GetEmail())
	// treat as non-pull-request comment by default and
	// set ref to default branch for the repo
	b.SetRef(fmt.Sprintf("refs/heads/%s", r.GetBranch()))

	pr := 0
	// override ref and pull request number if this is
	// a comment on a pull request
	if payload.GetIssue().IsPullRequest() {
		b.SetRef(fmt.Sprintf("refs/pull/%d/head", payload.GetIssue().GetNumber()))
		pr = payload.GetIssue().GetNumber()
	}

	return &types.Webhook{
		Comment:  payload.GetComment().GetBody(),
		PRNumber: pr,
		Hook:     h,
		Repo:     r,
		Build:    b,
	}, nil
}
