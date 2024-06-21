// SPDX-License-Identifier: Apache-2.0

package image

import (
	"strings"
	"testing"
)

func TestImage_ParseWithError(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		failure bool
		image   string
		want    string
	}{
		{
			name:    "image only",
			failure: false,
			image:   "golang",
			want:    "docker.io/library/golang:latest",
		},
		{
			name:    "image and tag",
			failure: false,
			image:   "golang:latest",
			want:    "docker.io/library/golang:latest",
		},
		{
			name:    "image and tag",
			failure: false,
			image:   "golang:1.14",
			want:    "docker.io/library/golang:1.14",
		},
		{
			name:    "fails with bad image",
			failure: true,
			image:   "!@#$%^&*()",
			want:    "!@#$%^&*()",
		},
		{
			name:    "fails with image sha",
			failure: true,
			image:   "1a3f5e7d9c1b3a5f7e9d1c3b5a7f9e1d3c5b7a9f1e3d5d7c9b1a3f5e7d9c1b3a",
			want:    "sha256:1a3f5e7d9c1b3a5f7e9d1c3b5a7f9e1d3c5b7a9f1e3d5d7c9b1a3f5e7d9c1b3a",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ParseWithError(test.image)

			if test.failure {
				if err == nil {
					t.Errorf("ParseWithError should have returned err")
				}

				if !strings.EqualFold(got, test.want) {
					t.Errorf("ParseWithError is %s want %s", got, test.want)
				}

				return // continue to next test
			}

			if err != nil {
				t.Errorf("ParseWithError returned err: %v", err)
			}

			if !strings.EqualFold(got, test.want) {
				t.Errorf("ParseWithError is %s want %s", got, test.want)
			}
		})
	}
}
