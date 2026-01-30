// SPDX-License-Identifier: Apache-2.0

package pipeline

import "testing"

func TestPipeline_Artifacts_Empty(t *testing.T) {
	// setup tests
	tests := []struct {
		artifacts *Artifacts
		want      bool
	}{
		{
			artifacts: &Artifacts{Paths: []string{"foo"}},
			want:      false,
		},
		{
			artifacts: new(Artifacts),
			want:      true,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.artifacts.Empty()

		if got != test.want {
			t.Errorf("Empty is %v, want %t", got, test.want)
		}
	}
}
