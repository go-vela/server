// SPDX-License-Identifier: Apache-2.0

package settings

import "fmt"

// ImageRestriction represents a container image pattern that is either
// blocked or warned about when used in a pipeline.
type ImageRestriction struct {
	Image  *string `json:"image,omitempty"  yaml:"image,omitempty"`
	Reason *string `json:"reason,omitempty" yaml:"reason,omitempty"`
}

// GetImage returns the Image field.
//
// When the provided ImageRestriction type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ir *ImageRestriction) GetImage() string {
	if ir == nil || ir.Image == nil {
		return ""
	}

	return *ir.Image
}

// GetReason returns the Reason field.
//
// When the provided ImageRestriction type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ir *ImageRestriction) GetReason() string {
	if ir == nil || ir.Reason == nil {
		return ""
	}

	return *ir.Reason
}

// SetImage sets the Image field.
//
// When the provided ImageRestriction type is nil, it
// will set nothing and immediately return.
func (ir *ImageRestriction) SetImage(v string) {
	if ir == nil {
		return
	}

	ir.Image = &v
}

// SetReason sets the Reason field.
//
// When the provided ImageRestriction type is nil, it
// will set nothing and immediately return.
func (ir *ImageRestriction) SetReason(v string) {
	if ir == nil {
		return
	}

	ir.Reason = &v
}

// String implements the Stringer interface for the ImageRestriction type.
func (ir *ImageRestriction) String() string {
	return fmt.Sprintf(`{
  Image: %s,
  Reason: %s,
}`,
		ir.GetImage(),
		ir.GetReason(),
	)
}

type Compiler struct {
	CloneImage        *string             `json:"clone_image,omitempty"         yaml:"clone_image,omitempty"`
	TemplateDepth     *int                `json:"template_depth,omitempty"      yaml:"template_depth,omitempty"`
	StarlarkExecLimit *int64              `json:"starlark_exec_limit,omitempty" yaml:"starlark_exec_limit,omitempty"`
	BlockedImages     *[]ImageRestriction `json:"blocked_images,omitempty"      yaml:"blocked_images,omitempty"`
	WarnImages        *[]ImageRestriction `json:"warn_images,omitempty"         yaml:"warn_images,omitempty"`
}

// GetBlockedImages returns the BlockedImages field.
//
// When the provided Compiler type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (cs *Compiler) GetBlockedImages() []ImageRestriction {
	if cs == nil || cs.BlockedImages == nil {
		return []ImageRestriction{}
	}

	return *cs.BlockedImages
}

// GetWarnImages returns the WarnImages field.
//
// When the provided Compiler type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (cs *Compiler) GetWarnImages() []ImageRestriction {
	if cs == nil || cs.WarnImages == nil {
		return []ImageRestriction{}
	}

	return *cs.WarnImages
}

// SetBlockedImages sets the BlockedImages field.
//
// When the provided Compiler type is nil, it
// will set nothing and immediately return.
func (cs *Compiler) SetBlockedImages(v []ImageRestriction) {
	if cs == nil {
		return
	}

	cs.BlockedImages = &v
}

// SetWarnImages sets the WarnImages field.
//
// When the provided Compiler type is nil, it
// will set nothing and immediately return.
func (cs *Compiler) SetWarnImages(v []ImageRestriction) {
	if cs == nil {
		return
	}

	cs.WarnImages = &v
}

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
func (cs *Compiler) GetStarlarkExecLimit() int64 {
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
func (cs *Compiler) SetStarlarkExecLimit(v int64) {
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
  BlockedImages: %v,
  WarnImages: %v,
}`,
		cs.GetCloneImage(),
		cs.GetTemplateDepth(),
		cs.GetStarlarkExecLimit(),
		cs.GetBlockedImages(),
		cs.GetWarnImages(),
	)
}

// CompilerMockEmpty returns an empty Compiler type.
func CompilerMockEmpty() Compiler {
	cs := Compiler{}
	cs.SetCloneImage("")
	cs.SetTemplateDepth(0)
	cs.SetStarlarkExecLimit(0)
	cs.SetBlockedImages(nil)
	cs.SetWarnImages(nil)

	return cs
}
