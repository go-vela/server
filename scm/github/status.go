// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v81/github"
	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/cache/models"
	"github.com/go-vela/server/constants"
)

// Status sends the commit status for the given SHA from the GitHub repo.
func (c *Client) Status(ctx context.Context, b *api.Build, token string, checkRuns []models.CheckRun) ([]models.CheckRun, error) {
	c.Logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   b.GetRepo().GetOrg(),
		"repo":  b.GetRepo().GetName(),
	}).Tracef("setting commit status for %s/%s/%d @ %s", b.GetRepo().GetOrg(), b.GetRepo().GetName(), b.GetNumber(), b.GetCommit())

	// only report opened, synchronize, and reopened action types for pull_request events
	if strings.EqualFold(b.GetEvent(), constants.EventPull) && !strings.EqualFold(b.GetEventAction(), constants.ActionOpened) &&
		!strings.EqualFold(b.GetEventAction(), constants.ActionSynchronize) && !strings.EqualFold(b.GetEventAction(), constants.ActionReopened) {
		return nil, nil
	}

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, token)

	if b.GetRepo().GetInstallID() != 0 && b.GetEvent() != constants.EventDeploy {
		return checkRun(ctx, client, c.config.WebUIAddress, c.config.StatusContext, b, checkRuns)
	}

	err := commitStatus(ctx, client, c.config.WebUIAddress, c.config.StatusContext, b)
	if err != nil {
		return nil, err
	}

	return checkRuns, nil
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
			Description: github.Ptr(description),
			Environment: github.Ptr(b.GetDeploy()),
			State:       github.Ptr(state),
		}

		// provide "Details" link in GitHub UI if server was configured with it
		if len(addr) > 0 {
			status.LogURL = github.Ptr(url)
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
			Context:     github.Ptr(context),
			Description: github.Ptr(description),
			State:       github.Ptr(state),
		}

		// provide "Details" link in GitHub UI if server was configured with it
		if len(addr) > 0 && b.GetStatus() != constants.StatusSkipped {
			status.TargetURL = github.Ptr(url)
		}

		// send API call to create the status context for the commit
		_, _, err := ghClient.Repositories.CreateStatus(ctx, b.GetRepo().GetOrg(), b.GetRepo().GetName(), b.GetCommit(), status)
		if err != nil {
			return err
		}
	}

	return nil
}

// checkRun creates or updates a check run for the build.
func checkRun(ctx context.Context, ghClient *github.Client, addr, statusCtx string, b *api.Build, checkRuns []models.CheckRun) ([]models.CheckRun, error) {
	state, conclusion, description, url := parseCheckRunStatus(b.GetStatus(), addr, b.GetRepo().GetFullName(), b.GetNumber(), 0)

	result := checkRuns

	title := fmt.Sprintf("Vela Build #%d • %s", b.GetNumber(), description)

	summary := buildCheckRunSummary(b)
	text := buildCheckRunText(b, url)
	startedAt := buildCheckRunStartedAt(b)
	completedAt := buildCheckRunCompletedAt(b)

	if len(checkRuns) == 0 {
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
			checkOpts := github.CreateCheckRunOptions{
				Name:      context,
				HeadSHA:   b.GetCommit(),
				Status:    github.Ptr(state),
				StartedAt: startedAt,
				Output: &github.CheckRunOutput{
					Title:   github.Ptr(title),
					Summary: github.Ptr(summary),
					Text:    github.Ptr(text),
				},
			}

			if conclusion != "" {
				checkOpts.Conclusion = github.Ptr(conclusion)
				checkOpts.CompletedAt = completedAt
			}

			// provide "Details" link in GitHub UI if server was configured with it
			if len(addr) > 0 && b.GetStatus() != constants.StatusSkipped {
				checkOpts.DetailsURL = github.Ptr(url)
			}

			check, _, err := ghClient.Checks.CreateCheckRun(ctx, b.GetRepo().GetOrg(), b.GetRepo().GetName(), checkOpts)
			if err != nil {
				return nil, err
			}

			result = append(result, models.CheckRun{
				ID:          check.GetID(),
				Context:     context,
				Repo:        b.GetRepo().GetFullName(),
				BuildNumber: b.GetNumber(),
			})
		}

		return result, nil
	}

	for _, checkRun := range checkRuns {
		checkOpts := github.UpdateCheckRunOptions{
			Name:   checkRun.Context,
			Status: github.Ptr(state),
			Output: &github.CheckRunOutput{
				Title:   github.Ptr(title),
				Summary: github.Ptr(summary),
				Text:    github.Ptr(text),
			},
		}

		if conclusion != "" {
			checkOpts.Conclusion = github.Ptr(conclusion)
			checkOpts.CompletedAt = completedAt
		}

		// provide "Details" link in GitHub UI if server was configured with it
		if len(addr) > 0 && b.GetStatus() != constants.StatusSkipped {
			checkOpts.DetailsURL = github.Ptr(url)
		}

		_, _, err := ghClient.Checks.UpdateCheckRun(ctx, b.GetRepo().GetOrg(), b.GetRepo().GetName(), checkRun.ID, checkOpts)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// StepStatus sends the commit status for the given SHA from the GitHub repo.
func (c *Client) StepStatus(ctx context.Context, b *api.Build, s *api.Step, token string, checkRuns []models.CheckRun) ([]models.CheckRun, error) {
	c.Logger.WithFields(logrus.Fields{
		"step":  s.GetName(),
		"build": b.GetNumber(),
		"org":   b.GetRepo().GetOrg(),
		"repo":  b.GetRepo().GetName(),
	}).Tracef("setting commit status for %s/%s/%d @ %s", b.GetRepo().GetOrg(), b.GetRepo().GetName(), b.GetNumber(), b.GetCommit())

	// no commit statuses on deployments
	if strings.EqualFold(b.GetEvent(), constants.EventDeploy) {
		return nil, nil
	}

	// create GitHub OAuth client with user's token
	client := c.newOAuthTokenClient(ctx, token)

	if b.GetRepo().GetInstallID() != 0 {
		return stepCheckRun(ctx, client, c.config.WebUIAddress, c.config.StatusContext, b, s, checkRuns)
	}

	err := stepCommitStatus(ctx, client, c.config.WebUIAddress, c.config.StatusContext, b, s)
	if err != nil {
		return nil, err
	}

	return checkRuns, nil
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
			Context:     github.Ptr(context),
			Description: github.Ptr(description),
			State:       github.Ptr(state),
		}

		// provide "Details" link in GitHub UI if server was configured with it
		if len(addr) > 0 && b.GetStatus() != constants.StatusSkipped {
			status.TargetURL = github.Ptr(url)
		}

		// send API call to create the status context for the commit
		_, _, err := ghClient.Repositories.CreateStatus(ctx, b.GetRepo().GetOrg(), b.GetRepo().GetName(), b.GetCommit(), status)
		if err != nil {
			return err
		}
	}

	return nil
}

func stepCheckRun(ctx context.Context, ghClient *github.Client, addr, statusCtx string, b *api.Build, s *api.Step, checkRuns []models.CheckRun) ([]models.CheckRun, error) {
	state, conclusion, description, url := parseCheckRunStatus(s.GetStatus(), addr, b.GetRepo().GetFullName(), b.GetNumber(), s.GetNumber())

	result := checkRuns

	if len(checkRuns) == 0 {
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
			checkOpts := github.CreateCheckRunOptions{
				Name:       context,
				HeadSHA:    b.GetCommit(),
				ExternalID: github.Ptr(strconv.FormatInt(b.GetID(), 10)),
				Status:     github.Ptr(state),
				Output: &github.CheckRunOutput{
					Title:   github.Ptr(description),
					Summary: github.Ptr(description),
				},
			}

			if conclusion != "" {
				checkOpts.Conclusion = github.Ptr(conclusion)
			}

			// provide "Details" link in GitHub UI if server was configured with it
			if len(addr) > 0 && b.GetStatus() != constants.StatusSkipped {
				checkOpts.DetailsURL = github.Ptr(url)
			}

			check, _, err := ghClient.Checks.CreateCheckRun(ctx, b.GetRepo().GetOrg(), b.GetRepo().GetName(), checkOpts)
			if err != nil {
				return nil, err
			}

			result = append(result, models.CheckRun{
				ID:          check.GetID(),
				Context:     context,
				Repo:        b.GetRepo().GetFullName(),
				BuildNumber: b.GetNumber(),
			})
		}

		return result, nil
	}

	for _, checkRun := range checkRuns {
		checkOpts := github.UpdateCheckRunOptions{
			Name:   checkRun.Context,
			Status: github.Ptr(state),
			Output: &github.CheckRunOutput{
				Title:   github.Ptr(description),
				Summary: github.Ptr(fmt.Sprintf("%s#%d", b.GetLink(), s.GetNumber())),
			},
		}

		if conclusion != "" {
			checkOpts.Conclusion = github.Ptr(conclusion)
		}

		// provide "Details" link in GitHub UI if server was configured with it
		if len(addr) > 0 && b.GetStatus() != constants.StatusSkipped {
			checkOpts.DetailsURL = github.Ptr(url)
		}

		_, _, err := ghClient.Checks.UpdateCheckRun(ctx, b.GetRepo().GetOrg(), b.GetRepo().GetName(), checkRun.ID, checkOpts)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// parseCommitStatus is a helper function to determine the url, state, and description for a commit status.
func parseCommitStatus(status, addr, repo string, buildNumber int64, stepNumber int32) (state string, description string, url string) {
	url = fmt.Sprintf("%s/%s/%d", addr, repo, buildNumber)
	target := "build"

	if stepNumber != 0 {
		url = fmt.Sprintf("%s#%d", url, stepNumber)
		target = "step"
	}

	switch status {
	case constants.StatusRunning, constants.StatusPending:
		state = "pending"
		description = fmt.Sprintf("the %s is %s", target, status)
	case constants.StatusPendingApproval:
		state = "pending"
		description = fmt.Sprintf("the %s needs approval from repo admin to run", target)
	case constants.StatusSuccess:
		state = "success"
		description = fmt.Sprintf("the %s was successful", target)
	case constants.StatusFailure:
		state = "failure"
		description = fmt.Sprintf("the %s has failed", target)
	case constants.StatusCanceled:
		state = "failure"
		description = fmt.Sprintf("the %s was canceled", target)
	case constants.StatusKilled:
		state = "failure"
		description = fmt.Sprintf("the %s was killed", target)
	case constants.StatusSkipped:
		state = "success"
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

	return
}

// parseCheckRunStatus is a helper function to determine the url, state, description, and conclusion for a check run.
func parseCheckRunStatus(status, addr, repo string, buildNumber int64, stepNumber int32) (state string, conclusion string, description string, url string) {
	url = fmt.Sprintf("%s/%s/%d", addr, repo, buildNumber)
	target := "build"

	if stepNumber != 0 {
		url = fmt.Sprintf("%s#%d", url, stepNumber)
		target = "step"
	}

	// set the state and description for the status context
	// depending on what the status of the build is
	switch status {
	case constants.StatusRunning:
		state = "in_progress"
		description = fmt.Sprintf("the %s is %s", target, status)
	case constants.StatusPending:
		state = "queued"
		description = fmt.Sprintf("the %s is %s", target, status)
	case constants.StatusPendingApproval:
		state = "queued"
		description = fmt.Sprintf("the %s needs approval from repo admin to run", target)
	case constants.StatusSuccess:
		state = "completed"
		conclusion = "success"
		description = fmt.Sprintf("the %s was successful", target)
	case constants.StatusFailure:
		state = "completed"
		conclusion = "failure"
		description = fmt.Sprintf("the %s has failed", target)
	case constants.StatusCanceled:
		state = "completed"
		conclusion = "cancelled"
		description = fmt.Sprintf("the %s was canceled", target)
	case constants.StatusKilled:
		state = "completed"
		conclusion = "cancelled"
		description = fmt.Sprintf("the %s was killed", target)
	case constants.StatusSkipped:
		state = "completed"
		conclusion = "skipped"
		description = fmt.Sprintf("the %s was skipped as no steps/stages found", target)
	default:
		state = "completed"
		conclusion = "failure"

		// if there is no build, then this status update is from a failed compilation
		if buildNumber == 0 && stepNumber == 0 {
			description = "error compiling pipeline - check audit for more information"
			url = fmt.Sprintf("%s/%s/hooks", addr, repo)
		} else {
			description = "there was an error"
		}
	}

	return
}

// buildCheckRunSummary creates the summary section of the check run output.
func buildCheckRunSummary(b *api.Build) string {
	return fmt.Sprintf(
		"### Build\n- Repo: `%s`\n- Build: `%d`\n- Event: `%s`\n-Action: `%s`\n- Branch: `%s`\n- Status: `%s`\n\n### Commit\n- SHA: `%s`\n- Sender: `%s`",
		b.GetRepo().GetFullName(),
		b.GetNumber(),
		b.GetEvent(),
		b.GetEventAction(),
		b.GetBranch(),
		b.GetStatus(),
		b.GetCommit(),
		b.GetSender(),
	)
}

// buildCheckRunText creates the text section of the check run output.
func buildCheckRunText(b *api.Build, detailsURL string) string {
	return fmt.Sprintf(
		"### Links\n- [Open build in Vela](%s)\n\n### Context\n- Ref: `%s`\n- Message: %s",
		detailsURL,
		b.GetRef(),
		b.GetMessage(),
	)
}

func buildCheckRunStartedAt(b *api.Build) *github.Timestamp {
	if b.GetStarted() > 0 {
		return &github.Timestamp{Time: time.Unix(b.GetStarted(), 0).UTC()}
	}

	return &github.Timestamp{Time: time.Now().UTC()}
}

func buildCheckRunCompletedAt(b *api.Build) *github.Timestamp {
	if b.GetFinished() > 0 {
		return &github.Timestamp{Time: time.Unix(b.GetFinished(), 0).UTC()}
	}

	return &github.Timestamp{Time: time.Now().UTC()}
}
