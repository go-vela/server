// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jenkins-x/go-scm/scm"

	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// ProcessWebhook parses the webhook from a repo.
func (c *client) ProcessWebhook(request *http.Request) (*types.Webhook, error) {
	logrus.Tracef("Processing SCM webhook")

	// create SCM OAuth client with user's token
	client, err := c.newClientToken("")
	if err != nil {
		return nil, err
	}

	h := extractHook(request, c.config.Kind)

	webhook, err := client.Webhooks.Parse(request, nil)
	if err != nil {
		return &types.Webhook{Hook: h}, nil
	}

	h.SetEvent(fmt.Sprintf("%v", webhook.Kind()))

	switch v := webhook.(type) {
	case *scm.PushHook:
		return processPushEvent(h, v)
	case *scm.PullRequestHook:
		return processPREvent(h, v)
	case *scm.DeployHook:
		return processDeploymentEvent(h, v)
	case *scm.PullRequestCommentHook:
		return processIssueCommentEvent(h, v)
	}

	return nil, nil
}

// VerifyWebhook verifies the webhook from a repo.
func (c *client) VerifyWebhook(request *http.Request, r *library.Repo) error {
	logrus.Tracef("Verifying SCM webhook for %s", r.GetFullName())
	// no-op
	return nil
}

func extractHook(request *http.Request, kind string) *library.Hook {
	h := new(library.Hook)
	h.SetNumber(1)
	h.SetCreated(time.Now().UTC().Unix())
	//TODO: What do here...
	// h.SetHost("github.com")
	h.SetStatus(constants.StatusSuccess)

	switch kind {
	case "bitbucket", "bitbucketcloud":
		h.SetSourceID(request.Header.Get("X-Hook-UUID"))
	case "github":
		h.SetSourceID(request.Header.Get("X-GitHub-Delivery"))

		if len(request.Header.Get("X-GitHub-Enterprise-Host")) > 0 {
			h.SetHost(request.Header.Get("X-GitHub-Enterprise-Host"))
		}
	}

	return h
}

// processPushEvent is a helper function to process the push event.
func processPushEvent(h *library.Hook, payload *scm.PushHook) (*types.Webhook, error) {
	logrus.Tracef("processing push SCM webhook for %s", payload.Repo.FullName)

	// convert payload to library repo
	r := toRepo(&payload.Repo)

	// convert payload to library build
	b := new(library.Build)
	b.SetEvent(constants.EventPush)
	b.SetClone(payload.Repo.Clone)
	b.SetSource(payload.Commit.Link)
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventPush, payload.Repo.Link))
	b.SetMessage(payload.Commit.Message)
	b.SetCommit(payload.Commit.Sha)
	b.SetSender(payload.Sender.Login)
	b.SetAuthor(payload.Commit.Author.Login)
	b.SetEmail(payload.Commit.Author.Email)
	b.SetBranch(strings.TrimPrefix(payload.Ref, "refs/heads/"))
	b.SetRef(payload.Ref)
	b.SetBaseRef(payload.BaseRef)

	// update the hook object
	h.SetBranch(b.GetBranch())
	h.SetEvent(constants.EventPush)
	h.SetLink(
		fmt.Sprintf("https://%s/%s/settings/hooks", h.GetHost(), r.GetFullName()),
	)

	// ensure the build author is set
	if len(b.GetAuthor()) == 0 {
		b.SetAuthor(payload.Commit.Committer.Name)
	}

	// ensure the build email is set
	if len(b.GetEmail()) == 0 {
		b.SetEmail(payload.Commit.Committer.Email)
	}

	// handle when push event is a tag
	if strings.HasPrefix(b.GetRef(), "refs/tags/") {
		// set the proper event for the hook
		h.SetEvent(constants.EventTag)
		// set the proper event for the build
		b.SetEvent(constants.EventTag)

		// set the proper branch from the base ref
		if strings.HasPrefix(payload.BaseRef, "refs/heads/") {
			b.SetBranch(strings.TrimPrefix(payload.BaseRef, "refs/heads/"))
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
func processPREvent(h *library.Hook, payload *scm.PullRequestHook) (*types.Webhook, error) {
	logrus.Tracef("processing pull_request SCM webhook for %s", payload.Repo.FullName)

	// update the hook object
	h.SetBranch(payload.PullRequest.Base.Ref)
	h.SetEvent(constants.EventPull)
	h.SetLink(
		fmt.Sprintf("https://%s/%s/settings/hooks", h.GetHost(), payload.Repo.FullName),
	)

	// if the pull request state isn't open we ignore it
	if payload.PullRequest.State != "open" {
		return &types.Webhook{Hook: h}, nil
	}

	// skip if the pull request action is not opened or synchronize
	if !strings.EqualFold(payload.Action.String(), "opened") &&
		!strings.EqualFold(payload.Action.String(), "synchronize") {
		return &types.Webhook{Hook: h}, nil
	}

	// convert payload to library repo
	r := toRepo(&payload.Repo)

	// convert payload to library build
	b := new(library.Build)
	b.SetEvent(constants.EventPull)
	b.SetClone(payload.Repo.Clone)
	b.SetSource(payload.PullRequest.Link)
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventPull, payload.Repo.Link))
	b.SetMessage(payload.PullRequest.Title)
	b.SetCommit(payload.PullRequest.Head.Sha)
	b.SetSender(payload.Sender.Login)
	b.SetAuthor(payload.PullRequest.Author.Login)
	b.SetEmail(payload.PullRequest.Author.Email)
	b.SetBranch(payload.PullRequest.Base.Ref)
	b.SetRef(fmt.Sprintf("refs/pull/%d/head", payload.PullRequest.Number))
	b.SetBaseRef(payload.PullRequest.Base.Ref)
	b.SetHeadRef(payload.PullRequest.Head.Ref)

	// ensure the build reference is set
	if payload.PullRequest.Merged {
		b.SetRef(fmt.Sprintf("refs/pull/%d/merge", payload.PullRequest.Number))
	}

	return &types.Webhook{
		Comment:  "",
		PRNumber: payload.PullRequest.Number,
		Hook:     h,
		Repo:     r,
		Build:    b,
	}, nil
}

// processDeploymentEvent is a helper function to process the deployment event.
func processDeploymentEvent(h *library.Hook, payload *scm.DeployHook) (*types.Webhook, error) {
	logrus.Tracef("processing deployment SCM webhook for %s", payload.Repo.FullName)

	// convert payload to library repo
	r := toRepo(&payload.Repo)

	// convert payload to library build
	//
	// Note: deployment only has "Author" and not sender/creator info for user
	b := new(library.Build)
	b.SetEvent(constants.EventDeploy)
	b.SetClone(payload.Repo.Link)
	b.SetDeploy(payload.Deployment.Environment)
	b.SetSource(payload.Deployment.Link)
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventDeploy, payload.Repo.Link))
	b.SetMessage(payload.Deployment.Description)
	b.SetCommit(payload.Deployment.Sha)
	b.SetSender(payload.Deployment.Author.Login)
	b.SetAuthor(payload.Deployment.Author.Login)
	b.SetEmail(payload.Deployment.Author.Email)
	b.SetBranch(payload.Deployment.Ref)
	b.SetRef(payload.Deployment.Ref)

	// check if payload is provided within request
	if payload.Deployment.Payload != nil {
		// set the payload info on the build
		b.SetDeployPayload(toMap(payload.Deployment.Payload))
	}

	// handle when the ref is a sha or short sha
	if strings.HasPrefix(b.GetCommit(), b.GetRef()) || b.GetCommit() == b.GetRef() {
		// set the proper branch for the build
		b.SetBranch(r.GetBranch())
		// set the proper ref for the build
		b.SetRef(fmt.Sprintf("refs/heads/%s", b.GetBranch()))
	}

	// update the hook object
	h.SetBranch(b.GetBranch())
	h.SetEvent(constants.EventDeploy)
	h.SetLink(
		fmt.Sprintf("https://%s/%s/settings/hooks", h.GetHost(), r.GetFullName()),
	)

	// handle when the ref is a branch
	if !strings.HasPrefix(b.GetRef(), "refs/") {
		// set the proper ref for the build
		b.SetRef(fmt.Sprintf("refs/heads/%s", b.GetBranch()))
	}

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
func processIssueCommentEvent(h *library.Hook, payload *scm.PullRequestCommentHook) (*types.Webhook, error) {
	logrus.Tracef("processing issue_comment SCM webhook for %s", payload.Repo.FullName)

	// convert payload to library repo
	r := toRepo(&payload.Repo)

	// update the hook object
	h.SetEvent(constants.EventComment)
	h.SetLink(
		fmt.Sprintf("https://%s/%s/settings/hooks", h.GetHost(), payload.Repo.FullName),
	)

	// skip if the comment action is deleted
	if strings.EqualFold(payload.Action.String(), "deleted") {
		// return &types.Webhook{Hook: h}, nil
		return &types.Webhook{
			Comment: payload.Comment.Body,
			Hook:    h,
		}, nil
	}

	// convert payload to library build
	b := new(library.Build)
	b.SetEvent(constants.EventComment)
	b.SetClone(payload.Repo.Clone)
	b.SetSource(payload.Comment.Link)
	b.SetTitle(fmt.Sprintf("%s received from %s", constants.EventComment, payload.Repo.Link))
	b.SetMessage(payload.PullRequest.Title)
	b.SetSender(payload.Sender.Login)
	b.SetAuthor(payload.Comment.Author.Login)
	b.SetEmail(payload.Comment.Author.Email)
	// treat as non-pull-request comment by default and
	// set ref to default branch for the repo
	b.SetRef(fmt.Sprintf("refs/pull/%d/head", payload.PullRequest.Number))

	return &types.Webhook{
		Comment:  payload.Comment.Body,
		PRNumber: payload.PullRequest.Number,
		Hook:     h,
		Repo:     r,
		Build:    b,
	}, nil
}

// helper function to conver scm repo into library
func toRepo(repo *scm.Repository) *library.Repo {
	// convert payload to library repo
	r := new(library.Repo)
	r.SetOrg(repo.Namespace)
	r.SetName(repo.Name)
	r.SetFullName(repo.FullName)
	r.SetLink(repo.Link)
	r.SetClone(repo.Clone)
	r.SetBranch(repo.Branch)
	r.SetPrivate(repo.Private)

	return r
}

// helper function to convert interface deployment payload into a map of strings
func toMap(src interface{}) map[string]string {
	set, ok := src.(map[string]interface{})
	if !ok {
		return nil
	}
	dst := map[string]string{}
	for k, v := range set {
		dst[k] = fmt.Sprint(v)
	}
	return dst
}
