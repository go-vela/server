// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypes_Template_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		template *Template
		want     *Template
	}{
		{
			template: testTemplate(),
			want:     testTemplate(),
		},
		{
			template: new(Template),
			want:     new(Template),
		},
	}

	// run tests
	for _, test := range tests {
		if test.template.GetLink() != test.want.GetLink() {
			t.Errorf("GetLink is %v, want %v", test.template.GetLink(), test.want.GetLink())
		}

		if test.template.GetName() != test.want.GetName() {
			t.Errorf("GetName is %v, want %v", test.template.GetName(), test.want.GetName())
		}

		if test.template.GetSource() != test.want.GetSource() {
			t.Errorf("GetSource is %v, want %v", test.template.GetSource(), test.want.GetSource())
		}

		if test.template.GetType() != test.want.GetType() {
			t.Errorf("GetType is %v, want %v", test.template.GetType(), test.want.GetType())
		}
	}
}

func TestTypes_Template_Setters(t *testing.T) {
	// setup types
	var tmpl *Template

	// setup tests
	tests := []struct {
		template *Template
		want     *Template
	}{
		{
			template: testTemplate(),
			want:     testTemplate(),
		},
		{
			template: tmpl,
			want:     new(Template),
		},
	}

	// run tests
	for _, test := range tests {
		test.template.SetLink(test.want.GetLink())
		test.template.SetName(test.want.GetName())
		test.template.SetSource(test.want.GetSource())
		test.template.SetType(test.want.GetType())

		if test.template.GetLink() != test.want.GetLink() {
			t.Errorf("SetLink is %v, want %v", test.template.GetLink(), test.want.GetLink())
		}

		if test.template.GetName() != test.want.GetName() {
			t.Errorf("SetName is %v, want %v", test.template.GetName(), test.want.GetName())
		}

		if test.template.GetSource() != test.want.GetSource() {
			t.Errorf("SetSource is %v, want %v", test.template.GetSource(), test.want.GetSource())
		}

		if test.template.GetType() != test.want.GetType() {
			t.Errorf("SetType is %v, want %v", test.template.GetType(), test.want.GetType())
		}
	}
}

func TestTypes_Template_String(t *testing.T) {
	// setup types
	tmpl := testTemplate()

	want := fmt.Sprintf(`{
  Link: %s,
  Name: %s,
  Source: %s,
  Type: %s,
}`,
		tmpl.GetLink(),
		tmpl.GetName(),
		tmpl.GetSource(),
		tmpl.GetType(),
	)

	// run test
	got := tmpl.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testTemplate is a test helper function to create a Template
// type with all fields set to a fake value.
func testTemplate() *Template {
	t := new(Template)

	t.SetLink("https://github.com/github/octocat/blob/branch/template.yml")
	t.SetName("template")
	t.SetSource("github.com/github/octocat/template.yml@branch")
	t.SetType("github")

	return t
}
