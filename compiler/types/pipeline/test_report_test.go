// SPDX-License-Identifier: Apache-2.0

package pipeline

import "testing"

func TestPipeline_TestReport_Empty(t *testing.T) {
	// setup tests
	tests := []struct {
		report *Artifacts
		want   bool
	}{
		{
			report: &Artifacts{Paths: []string{"foo"}},
			want:   false,
		},
		{
			report: new(Artifacts),
			want:   true,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.report.Empty()

		if got != test.want {
			t.Errorf("Empty is %v, want %t", got, test.want)
		}
	}
}
