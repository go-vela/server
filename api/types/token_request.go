// SPDX-License-Identifier: Apache-2.0

package types

// TokenRequest is the API representation of an install token request from worker.
//
// swagger:model TokenRequest
type TokenRequest struct {
	Repositories []string          `json:"repositories,omitempty"`
	Permissions  map[string]string `json:"permissions,omitempty"`
}
