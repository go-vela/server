// SPDX-License-Identifier: Apache-2.0

package types

// CleanReport is the API types representation of an admin clean report of resources in the database.
//
// swagger:model CleanReport
type CleanReport struct {
	Builds                int64 `json:"builds"`
	Executables           int64 `json:"executables"`
	PendingApprovalBuilds int64 `json:"pending_approval_builds"`
	Steps                 int64 `json:"steps"`
	Services              int64 `json:"services"`
}
