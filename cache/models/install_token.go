// SPDX-License-Identifier: Apache-2.0

package models

type InstallToken struct {
	Token        string            `json:"token"`
	Repositories []string          `json:"repositories"`
	Permissions  map[string]string `json:"permissions"`
	Expiration   int64             `json:"expiration"`
}
