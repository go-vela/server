// SPDX-License-Identifier: Apache-2.0

package constants

// Log options constants.
const (
	// LogHead defines the option for preserving beginning of logs
	// when the step has produces over the limit of log bytes.
	LogHead = "head"

	// LogTail defines the option for preserving end of logs
	// when the step has produces over the limit of log bytes.
	LogTail = "tail"
)
