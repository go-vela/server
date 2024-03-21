// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
)

// Settings is the API representation of platform settings.
//
// swagger:model Settings
type Settings struct {
	ID     *int64  `json:"id,omitempty"`
	FooNum *int64  `json:"foo_num,omitempty"`
	FooStr *string `json:"foo_str,omitempty"`
}

// String implements the Stringer interface for the Settings type.
func (w *Settings) String() string {
	return fmt.Sprintf(`{
  ID: %d,
}`,
		1, // w.GetID(),
	)
}
