// SPDX-License-Identifier: Apache-2.0

package constants

// Build and step statuses.
const (
	// StatusError defines the status type for build and step error statuses.
	StatusError = "error"

	// StatusFailure defines the status type for build and step failure statuses.
	StatusFailure = "failure"

	// StatusKilled defines the status type for build and step killed statuses.
	StatusKilled = "killed"

	// StatusCanceled defines the status type for build and step canceled statuses.
	StatusCanceled = "canceled"

	// StatusPending defines the status type for build and step pending statuses.
	StatusPending = "pending"

	// StatusPendingApproval defines the status type for a build waiting to be approved to run.
	StatusPendingApproval = "pending approval"

	// StatusRunning defines the status type for build and step running statuses.
	StatusRunning = "running"

	// StatusSuccess defines the status type for build and step success statuses.
	StatusSuccess = "success"

	// StatusSkipped defines the status type for build and step skipped statuses.
	StatusSkipped = "skipped"
)
