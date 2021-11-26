// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"net/http"

	"github.com/google/go-github/v39/github"

	"github.com/go-vela/types"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// ProcessWebhook parses the webhook from a repo.
func (c *client) ProcessWebhook(request *http.Request) (*types.Webhook, error) {
	logrus.Tracef("Processing GitHub webhook")
	return nil, nil
}

// VerifyWebhook verifies the webhook from a repo.
func (c *client) VerifyWebhook(request *http.Request, r *library.Repo) error {
	logrus.Tracef("Verifying GitHub webhook for %s", r.GetFullName())
	return nil
}

// processPushEvent is a helper function to process the push event.
func processPushEvent(h *library.Hook, payload *github.PushEvent) (*types.Webhook, error) {
	logrus.Tracef("processing push GitHub webhook for %s", payload.GetRepo().GetFullName())
	return nil, nil
}

// processPREvent is a helper function to process the pull_request event.
func processPREvent(h *library.Hook, payload *github.PullRequestEvent) (*types.Webhook, error) {
	logrus.Tracef("processing pull_request GitHub webhook for %s", payload.GetRepo().GetFullName())
	return nil, nil
}

// processDeploymentEvent is a helper function to process the deployment event.
//
// nolint: lll // ignore long line length due to variable names
func processDeploymentEvent(h *library.Hook, payload *github.DeploymentEvent) (*types.Webhook, error) {
	logrus.Tracef("processing deployment GitHub webhook for %s", payload.GetRepo().GetFullName())
	return nil, nil
}

// processIssueCommentEvent is a helper function to process the issue comment event.
//
// nolint: lll // ignore long line length due to variable names
func processIssueCommentEvent(h *library.Hook, payload *github.IssueCommentEvent) (*types.Webhook, error) {
	logrus.Tracef("processing issue_comment GitHub webhook for %s", payload.GetRepo().GetFullName())
	return nil, nil
}
