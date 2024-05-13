// SPDX-License-Identifier: Apache-2.0

package settings

import "fmt"

type Compiler struct {
	CloneImage        *string `json:"clone_image,omitempty" yaml:"clone_image,omitempty"`
	TemplateDepth     *int    `json:"template_depth,omitempty" yaml:"template_depth,omitempty"`
	StarlarkExecLimit *uint64 `json:"starlark_exec_limit,omitempty" yaml:"starlark_exec_limit,omitempty"`
}

// GetCloneImage returns the CloneImage field.
//
// When the provided Compiler type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (cs *Compiler) GetCloneImage() string {
	// return zero value if Settings type or CloneImage field is nil
	if cs == nil || cs.CloneImage == nil {
		return ""
	}

	return *cs.CloneImage
}

// GetTemplateDepth returns the TemplateDepth field.
//
// When the provided Compiler type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (cs *Compiler) GetTemplateDepth() int {
	// return zero value if Settings type or TemplateDepth field is nil
	if cs == nil || cs.TemplateDepth == nil {
		return 0
	}

	return *cs.TemplateDepth
}

// GetStarlarkExecLimit returns the StarlarkExecLimit field.
//
// When the provided Compiler type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (cs *Compiler) GetStarlarkExecLimit() uint64 {
	// return zero value if Compiler type or StarlarkExecLimit field is nil
	if cs == nil || cs.StarlarkExecLimit == nil {
		return 0
	}

	return *cs.StarlarkExecLimit
}

// SetCloneImage sets the CloneImage field.
//
// When the provided Compiler type is nil, it
// will set nothing and immediately return.
func (cs *Compiler) SetCloneImage(v string) {
	// return if Compiler type is nil
	if cs == nil {
		return
	}

	cs.CloneImage = &v
}

// SetTemplateDepth sets the TemplateDepth field.
//
// When the provided Compiler type is nil, it
// will set nothing and immediately return.
func (cs *Compiler) SetTemplateDepth(v int) {
	// return if Compiler type is nil
	if cs == nil {
		return
	}

	cs.TemplateDepth = &v
}

// SetStarlarkExecLimit sets the StarlarkExecLimit field.
//
// When the provided Compiler type is nil, it
// will set nothing and immediately return.
func (cs *Compiler) SetStarlarkExecLimit(v uint64) {
	// return if Compiler type is nil
	if cs == nil {
		return
	}

	cs.StarlarkExecLimit = &v
}

// String implements the Stringer interface for the Compiler type.
func (cs *Compiler) String() string {
	return fmt.Sprintf(`{
  CloneImage: %s,
  TemplateDepth: %d,
  StarlarkExecLimit: %d,
}`,
		cs.GetCloneImage(),
		cs.GetTemplateDepth(),
		cs.GetStarlarkExecLimit(),
	)
}

// CompilerMockEmpty returns an empty Compiler type.
func CompilerMockEmpty() Compiler {
	cs := Compiler{}
	cs.SetCloneImage("")
	cs.SetTemplateDepth(0)
	cs.SetStarlarkExecLimit(0)

	return cs
}
