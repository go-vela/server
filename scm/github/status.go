// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/go-github/v84/github"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
)

const (
	// Below are github status constants.
	StatePending = "pending"
	StateSuccess = "success"
	StateFailure = "failure"
	StateError   = "error"
	//nolint:misspell // GitHub uses cancelled
	StateCancelled  = "cancelled"
	StateSkipped    = "skipped"
	StateQueued     = "queued"
	StateInProgress = "in_progress"
	StateCompleted  = "completed"
)

// GenerateStatusToken generates a token for setting commit status on the SCM provider.
func (c *Client) GenerateStatusToken(ctx context.Context, b *api.Build) string {
	tknRepo := []string{b.GetRepo().GetName()}
	tknPerms := map[string]string{"statuses": constants.PermissionWrite}

	if b.GetEvent() == constants.EventDeploy {
		tknPerms = map[string]string{"deployments": constants.PermissionWrite}
	}

	tkn, err := c.NewAppInstallationToken(ctx, b.GetRepo().GetInstallID(), tknRepo, tknPerms)
	if err != nil {
		c.Logger.Errorf("unable to generate status token for build %d: %v", b.GetNumber(), err)
		return ""
	}

	return tkn
}

// Status sends the commit status for the given SHA from the GitHub repo.
func (c *Client) Status(ctx context.Context, b *api.Build, token string) error {
	c.Logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   b.GetRepo().GetOrg(),
		"repo":  b.GetRepo().GetName(),
	}).Tracef("setting commit status for %s/%s/%d @ %s", b.GetRepo().GetOrg(), b.GetRepo().GetName(), b.GetNumber(), b.GetCommit())

	// only report opened, synchronize, and reopened action types for pull_request events
	if strings.EqualFold(b.GetEvent(), constants.EventPull) && !strings.EqualFold(b.GetEventAction(), constants.ActionOpened) &&
		!strings.EqualFold(b.GetEventAction(), constants.ActionSynchronize) && !strings.EqualFold(b.GetEventAction(), constants.ActionReopened) {
		return nil
	}

	// create token client
	client := c.newTokenClient(ctx, token)

	err := commitStatus(ctx, client, c.config.WebUIAddress, c.config.StatusContext, b)
	if err != nil {
		return err
	}

	return nil
}

// commitStatus sends a commit status update to GitHub for non-GitHub app installed repos or deployment events.
func commitStatus(ctx context.Context, ghClient *github.Client, addr, statusCtx string, b *api.Build) error {
	state, description, url := parseCommitStatus(b.GetStatus(), addr, b.GetRepo().GetFullName(), b.GetNumber(), 0)

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
			Description: new(description),
			Environment: new(b.GetDeploy()),
			State:       new(state),
		}

		// provide "Details" link in GitHub UI if server was configured with it
		if len(addr) > 0 {
			status.LogURL = new(url)
		}

		_, _, err = ghClient.Repositories.CreateDeploymentStatus(ctx, b.GetRepo().GetOrg(), b.GetRepo().GetName(), int64(number), status)

		return err
	}

	var contexts []string

	if b.GetEvent() == constants.EventMergeGroup {
		for _, e := range b.GetRepo().GetMergeQueueEvents() {
			context := fmt.Sprintf("%s/%s", statusCtx, e)
			contexts = append(contexts, context)
		}
	} else {
		contexts = append(contexts, fmt.Sprintf("%s/%s", statusCtx, b.GetEvent()))
	}

	for _, context := range contexts {
		// create the status object to make the API call
		status := github.RepoStatus{
			Context:     new(context),
			Description: new(description),
			State:       new(state),
		}

		// provide "Details" link in GitHub UI if server was configured with it
		if len(addr) > 0 && b.GetStatus() != constants.StatusSkipped {
			status.TargetURL = new(url)
		}

		// send API call to create the status context for the commit
		_, _, err := ghClient.Repositories.CreateStatus(ctx, b.GetRepo().GetOrg(), b.GetRepo().GetName(), b.GetCommit(), status)
		if err != nil {
			return err
		}
	}

	return nil
}

// StepStatus sends the commit status for the given SHA from the GitHub repo.
func (c *Client) StepStatus(ctx context.Context, b *api.Build, s *api.Step, token string) error {
	c.Logger.WithFields(logrus.Fields{
		"step":  s.GetName(),
		"build": b.GetNumber(),
		"org":   b.GetRepo().GetOrg(),
		"repo":  b.GetRepo().GetName(),
	}).Tracef("setting commit status for %s/%s/%d @ %s", b.GetRepo().GetOrg(), b.GetRepo().GetName(), b.GetNumber(), b.GetCommit())

	// no commit statuses on deployments
	if strings.EqualFold(b.GetEvent(), constants.EventDeploy) {
		return nil
	}

	// create GitHub OAuth client with user's token
	client := c.newTokenClient(ctx, token)

	err := stepCommitStatus(ctx, client, c.config.WebUIAddress, c.config.StatusContext, b, s)
	if err != nil {
		return err
	}

	return nil
}

func stepCommitStatus(ctx context.Context, ghClient *github.Client, addr, statusCtx string, b *api.Build, s *api.Step) error {
	state, description, url := parseCommitStatus(s.GetStatus(), addr, b.GetRepo().GetFullName(), b.GetNumber(), s.GetNumber())

	var contexts []string

	if b.GetEvent() == constants.EventMergeGroup {
		for _, e := range b.GetRepo().GetMergeQueueEvents() {
			context := fmt.Sprintf("%s/%s/%s", statusCtx, e, s.GetReportAs())
			contexts = append(contexts, context)
		}
	} else {
		contexts = append(contexts, fmt.Sprintf("%s/%s/%s", statusCtx, b.GetEvent(), s.GetReportAs()))
	}

	for _, context := range contexts {
		// create the status object to make the API call
		status := github.RepoStatus{
			Context:     new(context),
			Description: new(description),
			State:       new(state),
		}

		// provide "Details" link in GitHub UI if server was configured with it
		if len(addr) > 0 && b.GetStatus() != constants.StatusSkipped {
			status.TargetURL = new(url)
		}

		// send API call to create the status context for the commit
		_, _, err := ghClient.Repositories.CreateStatus(ctx, b.GetRepo().GetOrg(), b.GetRepo().GetName(), b.GetCommit(), status)
		if err != nil {
			return err
		}
	}

	return nil
}

// parseCommitStatus is a helper function to determine the url, state, and description for a commit status.
func parseCommitStatus(status, addr, repo string, buildNumber int64, stepNumber int32) (string, string, string) {
	var (
		url         = fmt.Sprintf("%s/%s/%d", addr, repo, buildNumber)
		target      = "build"
		state       string
		description string
	)

	if stepNumber != 0 {
		url = fmt.Sprintf("%s#%d", url, stepNumber)
		target = "step"
	}

	switch status {
	case constants.StatusRunning, constants.StatusPending:
		state = StatePending
		description = fmt.Sprintf("the %s is %s", target, status)
	case constants.StatusPendingApproval:
		state = StatePending
		description = fmt.Sprintf("the %s needs approval from repo admin to run", target)
	case constants.StatusSuccess:
		state = StateSuccess
		description = fmt.Sprintf("the %s was successful", target)
	case constants.StatusFailure:
		state = StateFailure
		description = fmt.Sprintf("the %s has failed", target)
	case constants.StatusCanceled:
		state = StateFailure
		description = fmt.Sprintf("the %s was canceled", target)
	case constants.StatusKilled:
		state = StateFailure
		description = fmt.Sprintf("the %s was killed", target)
	case constants.StatusSkipped:
		state = StateSuccess
		description = fmt.Sprintf("the %s was skipped as no steps/stages found", target)
	default:
		state = "error"

		// if there is no build, then this status update is from a failed compilation
		if buildNumber == 0 && stepNumber == 0 {
			description = "error compiling pipeline - check audit for more information"
			url = fmt.Sprintf("%s/%s/hooks", addr, repo)
		} else {
			description = "there was an error"
		}
	}

	return state, description, url
}
