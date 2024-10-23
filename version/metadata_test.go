// SPDX-License-Identifier: Apache-2.0

package version

import (
	"fmt"
	"reflect"
	"testing"
)

func TestVersion_Metadata_String(t *testing.T) {
	// setup types
	m := &Metadata{
		Architecture:    "amd64",
		BuildDate:       "1970-1-1T00:00:00Z",
		Compiler:        "gc",
		GitCommit:       "abcdef123456789",
		GoVersion:       "1.19.0",
		OperatingSystem: "linux",
	}

	want := fmt.Sprintf(
		metaFormat,
		m.Architecture,
		m.BuildDate,
		m.Compiler,
		m.GitCommit,
		m.GoVersion,
		m.OperatingSystem,
	)

	// run test
	got := m.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}
