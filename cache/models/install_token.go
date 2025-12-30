// SPDX-License-Identifier: Apache-2.0

package models

type InstallToken struct {
	Token        string            `json:"token"`
	InstallID    int64             `json:"install_id"`
	Repositories []string          `json:"repositories"`
	Permissions  map[string]string `json:"permissions"`
	Expiration   int64             `json:"expiration"`
	Timeout      int32             `json:"timeout"`
}
