// SPDX-License-Identifier: Apache-2.0

package pipeline

import "testing"

func TestPipeline_Git_Empty(t *testing.T) {
	// setup tests
	tests := []struct {
		git  *Git
		want bool
	}{
		{
			git: &Git{Token: Token{
				Repositories: []string{"hello-world"},
			}},
			want: false,
		},
		{
			git:  &Git{Token{Repositories: []string{}}},
			want: true,
		},
		{
			git:  new(Git),
			want: true,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.git.Empty()

		if got != test.want {
			t.Errorf("Empty is %v, want %t", got, test.want)
		}
	}
}
