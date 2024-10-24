// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
)

func TestResp_String(t *testing.T) {
	// setup types
	str := "foo"
	e := &Error{
		Message: &str,
	}
	want := fmt.Sprintf("%+v", *e)

	// run test
	got := e.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}
