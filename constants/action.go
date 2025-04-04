// SPDX-License-Identifier: Apache-2.0

package constants

// Build and repo events.
const (
	// ActionOpened defines the action for opening pull requests.
	ActionOpened = "opened"

	// ActionCreated defines the action for creating deployments or issue comments.
	ActionCreated = "created"

	// ActionEdited defines the action for the editing of pull requests or issue comments.
	ActionEdited = "edited"

	// ActionRenamed defines the action for renaming a repository.
	ActionRenamed = "renamed"

	// ActionReopened defines the action for re-opening a pull request (or issue).
	ActionReopened = "reopened"

	// ActionSynchronize defines the action for the synchronizing of pull requests.
	ActionSynchronize = "synchronize"

	// ActionLabeled defines the action for the labeling of pull requests.
	ActionLabeled = "labeled"

	// ActionUnlabeled defines the action for the unlabeling of pull requests.
	ActionUnlabeled = "unlabeled"

	// ActionMerged defines the action for the merging of a pull request.
	ActionMerged = "merged"

	// ActionClosed defines the action for closing pull requests (unmerged).
	ActionClosed = "closed"

	// ActionTransferred defines the action for transferring repository ownership.
	ActionTransferred = "transferred"

	// ActionBranch defines the action for deleting a branch.
	ActionBranch = "branch"

	// ActionTag defines the action for deleting a tag.
	ActionTag = "tag"

	// ActionRun defines the action for running a schedule.
	ActionRun = "run"
)
