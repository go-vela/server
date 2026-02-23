// SPDX-License-Identifier: Apache-2.0

package models

type CheckRun struct {
	ID          int64  `json:"id"`
	Context     string `json:"context"`
	Repo        string `json:"repo"`
	BuildNumber int64  `json:"build_number"`
}
