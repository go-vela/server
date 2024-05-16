// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypes_Compiler_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		compiler *Compiler
		want     *Compiler
	}{
		{
			compiler: testCompilerSettings(),
			want:     testCompilerSettings(),
		},
		{
			compiler: new(Compiler),
			want:     new(Compiler),
		},
	}

	// run tests
	for _, test := range tests {
		if !reflect.DeepEqual(test.compiler.GetCloneImage(), test.want.GetCloneImage()) {
			t.Errorf("GetCloneImage is %v, want %v", test.compiler.GetCloneImage(), test.want.GetCloneImage())
		}

		if !reflect.DeepEqual(test.compiler.GetTemplateDepth(), test.want.GetTemplateDepth()) {
			t.Errorf("GetTemplateDepth is %v, want %v", test.compiler.GetTemplateDepth(), test.want.GetTemplateDepth())
		}

		if !reflect.DeepEqual(test.compiler.GetStarlarkExecLimit(), test.want.GetStarlarkExecLimit()) {
			t.Errorf("GetStarlarkExecLimit is %v, want %v", test.compiler.GetStarlarkExecLimit(), test.want.GetStarlarkExecLimit())
		}
	}
}

func TestTypes_Compiler_Setters(t *testing.T) {
	// setup types
	var cs *Compiler

	// setup tests
	tests := []struct {
		compiler *Compiler
		want     *Compiler
	}{
		{
			compiler: testCompilerSettings(),
			want:     testCompilerSettings(),
		},
		{
			compiler: cs,
			want:     new(Compiler),
		},
	}

	// run tests
	for _, test := range tests {
		test.compiler.SetCloneImage(test.want.GetCloneImage())

		if !reflect.DeepEqual(test.compiler.GetCloneImage(), test.want.GetCloneImage()) {
			t.Errorf("SetCloneImage is %v, want %v", test.compiler.GetCloneImage(), test.want.GetCloneImage())
		}

		test.compiler.SetTemplateDepth(test.want.GetTemplateDepth())

		if !reflect.DeepEqual(test.compiler.GetTemplateDepth(), test.want.GetTemplateDepth()) {
			t.Errorf("SetTemplateDepth is %v, want %v", test.compiler.GetTemplateDepth(), test.want.GetTemplateDepth())
		}

		test.compiler.SetStarlarkExecLimit(test.want.GetStarlarkExecLimit())

		if !reflect.DeepEqual(test.compiler.GetStarlarkExecLimit(), test.want.GetStarlarkExecLimit()) {
			t.Errorf("SetStarlarkExecLimit is %v, want %v", test.compiler.GetStarlarkExecLimit(), test.want.GetStarlarkExecLimit())
		}
	}
}

func TestTypes_Compiler_String(t *testing.T) {
	// setup types
	cs := testCompilerSettings()

	want := fmt.Sprintf(`{
  CloneImage: %s,
  TemplateDepth: %d,
  StarlarkExecLimit: %d,
}`,
		cs.GetCloneImage(),
		cs.GetTemplateDepth(),
		cs.GetStarlarkExecLimit(),
	)

	// run test
	got := cs.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testCompilerSettings is a test helper function to create a Compiler
// type with all fields set to a fake value.
func testCompilerSettings() *Compiler {
	cs := new(Compiler)

	cs.SetCloneImage("target/vela-git:latest")
	cs.SetTemplateDepth(1)
	cs.SetStarlarkExecLimit(100)

	return cs
}
