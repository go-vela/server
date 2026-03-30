// SPDX-License-Identifier: Apache-2.0

package types

// Installation is the json error message from the server for a given http response.
//
// swagger:model Installation
type Installation struct {
	InstallID *int64  `json:"install_id"`
	Target    *string `json:"target"`
}

func (i *Installation) GetInstallID() int64 {
	if i == nil || i.InstallID == nil {
		return 0
	}
	return *i.InstallID
}

func (i *Installation) GetTarget() string {
	if i == nil || i.Target == nil {
		return ""
	}
	return *i.Target
}

func (i *Installation) SetInstallID(id int64) {
	if i == nil {
		return
	}

	i.InstallID = &id
}

func (i *Installation) SetTarget(target string) {
	if i == nil {
		return
	}

	i.Target = &target
}
