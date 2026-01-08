// SPDX-License-Identifier: Apache-2.0

package pipeline

import "testing"

func TestPipeline_Artifact_Empty(t *testing.T) {
	// setup tests
	tests := []struct {
		artifact *Artifact
		want     bool
	}{
		{
			artifact: &Artifact{Paths: []string{"foo"}},
			want:     false,
		},
		{
			artifact: new(Artifact),
			want:     true,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.artifact.Empty()

		if got != test.want {
			t.Errorf("Empty is %v, want %t", got, test.want)
		}
	}
}
