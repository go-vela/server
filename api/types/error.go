// SPDX-License-Identifier: Apache-2.0

package types

import "fmt"

// Error is the json error message from the server for a given http response.
//
// swagger:model Error
type Error struct {
	Message *string `json:"error"`
}

func (e *Error) String() string {
	return fmt.Sprintf("%+v", *e)
}
