// SPDX-License-Identifier: Apache-2.0

package types

import (
	"testing"
)

func TestTypes_Sanitize(t *testing.T) {
	// setup tests
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{
			name:  "percent",
			value: `%`,
			want:  `%`,
		},
		{
			name:  "quoted",
			value: `"hello"`,
			want:  `"hello"`,
		},
		{
			name:  "email",
			value: `OctoKitty@github.com`,
			want:  `OctoKitty@github.com`,
		},
		{
			name:  "url",
			value: `https://github.com/go-vela`,
			want:  `https://github.com/go-vela`,
		},
		{
			name:  "encoded",
			value: `+ added foo %25 + updated bar %22 +`,
			want:  `+ added foo %25 + updated bar %22 +`,
		},
		{
			name:  "html with headers",
			value: `Merge pull request #1 from me/patch-1\n\n<h1>hello</h1> is now <h2>hello</h2>`,
			want:  `Merge pull request #1 from me/patch-1\n\nhello is now hello`,
		},
		{
			name:  "html with email",
			value: `Co-authored-by: OctoKitty <OctoKitty@github.com>`,
			want:  `Co-authored-by: OctoKitty `,
		},
		{
			name:  "html with href",
			value: `<a onblur="alert(secret)" href="http://www.google.com">Google</a>`,
			want:  `Google`,
		},
		{
			name:  "local cross-site script",
			value: `<script>alert('XSS')</script>`,
			want:  ``,
		},
		{
			name:  "remote cross-site script",
			value: `<SCRIPT/XSS SRC="http://ha.ckers.org/xss.js"></SCRIPT>`,
			want:  ``,
		},
		{
			name:  "embedded cross-site script",
			value: `%3cDIV%20STYLE%3d%22width%3a%20expression(alert('XSS'))%3b%22%3e`,
			want:  ``,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := sanitize(test.value)

			if got != test.want {
				t.Errorf("sanitize is %v, want %v", got, test.want)
			}
		})
	}
}
