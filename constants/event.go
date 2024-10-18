// SPDX-License-Identifier: Apache-2.0

package constants

// Build and repo events.
const (
	// EventComment defines the event type for comments added to a pull request.
	EventComment = "comment"

	// EventDelete defines the event type for build and repo delete events.
	EventDelete = "delete"

	// EventDeploy defines the event type for build and repo deployment events.
	EventDeploy = "deployment"

	// EventPull defines the event type for build and repo pull_request events.
	EventPull = "pull_request"

	// EventPush defines the event type for build and repo push events.
	EventPush = "push"

	// EventRepository defines the general event type for repo management.
	EventRepository = "repository"

	// EventSchedule defines the event type for build and repo schedule events.
	EventSchedule = "schedule"

	// EventTag defines the event type for build and repo tag events.
	EventTag = "tag"

	// Alternates for common user inputs that do not match our set constants.

	// EventPullAlternate defines the alternate event type for build and repo pull_request events.
	EventPullAlternate = "pull"

	// EventDeployAlternate defines the alternate event type for build and repo deployment events.
	EventDeployAlternate = "deploy"
)
