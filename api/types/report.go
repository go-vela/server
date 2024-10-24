// SPDX-License-Identifier: Apache-2.0

package types

// Report represents the Vela checks report for a build.
type Report struct {
	Title            *string       `json:"title,omitempty"`
	Summary          *string       `json:"summary,omitempty"`
	Text             *string       `json:"text,omitempty"`
	AnnotationsCount *int          `json:"annotations_count,omitempty"`
	AnnotationsURL   *string       `json:"annotations_url,omitempty"`
	Annotations      []*Annotation `json:"annotations,omitempty"`
}

// Annotation represents the Vela annotation for a report.
type Annotation struct {
	Path            *string `json:"path,omitempty"`
	StartLine       *int    `json:"start_line,omitempty"`
	EndLine         *int    `json:"end_line,omitempty"`
	StartColumn     *int    `json:"start_column,omitempty"`
	EndColumn       *int    `json:"end_column,omitempty"`
	AnnotationLevel *string `json:"annotation_level,omitempty"`
	Message         *string `json:"message,omitempty"`
	Title           *string `json:"title,omitempty"`
	RawDetails      *string `json:"raw_details,omitempty"`
}

// GetTitle returns the Title field.
//
// When the provided Report type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *Report) GetTitle() string {
	// return zero value if Step type or ID field is nil
	if r == nil || r.Title == nil {
		return ""
	}

	return *r.Title
}

// GetSummary returns the Summary field.
//
// When the provided Report type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *Report) GetSummary() string {
	// return zero value if Step type or ID field is nil
	if r == nil || r.Summary == nil {
		return ""
	}

	return *r.Summary
}

// GetText returns the Text field.
//
// When the provided Report type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *Report) GetText() string {
	// return zero value if Step type or ID field is nil
	if r == nil || r.Text == nil {
		return ""
	}

	return *r.Text
}

// GetAnnotationsCount returns the AnnotationsCount field.
//
// When the provided Report type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *Report) GetAnnotationsCount() int {
	// return zero value if Step type or ID field is nil
	if r == nil || r.AnnotationsCount == nil {
		return 0
	}

	return *r.AnnotationsCount
}

// GetAnnotationsURL returns the AnnotationsURL field.
//
// When the provided Report type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *Report) GetAnnotationsURL() string {
	// return zero value if Step type or ID field is nil
	if r == nil || r.AnnotationsURL == nil {
		return ""
	}

	return *r.AnnotationsURL
}

// GetAnnotations returns the Annotations field.
//
// When the provided Report type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *Report) GetAnnotations() []*Annotation {
	// return zero value if Step type or ID field is nil
	if r == nil || r.Annotations == nil {
		return []*Annotation{}
	}

	return r.Annotations
}

// GetPath returns the Path field.
//
// When the provided Annotation type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Annotation) GetPath() string {
	// return zero value if Step type or ID field is nil
	if a == nil || a.Path == nil {
		return ""
	}

	return *a.Path
}

// GetStartLine returns the StartLine field.
//
// When the provided Annotation type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Annotation) GetStartLine() int {
	// return zero value if Step type or ID field is nil
	if a == nil || a.StartLine == nil {
		return 0
	}

	return *a.StartLine
}

// GetEndLine returns the EndLine field.
//
// When the provided Annotation type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Annotation) GetEndLine() int {
	// return zero value if Step type or ID field is nil
	if a == nil || a.EndLine == nil {
		return 0
	}

	return *a.EndLine
}

// GetStartColumn returns the StartColumn field.
//
// When the provided Annotation type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Annotation) GetStartColumn() int {
	// return zero value if Step type or ID field is nil
	if a == nil || a.StartColumn == nil {
		return 0
	}

	return *a.StartColumn
}

// GetEndColumn returns the EndColumn field.
//
// When the provided Annotation type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Annotation) GetEndColumn() int {
	// return zero value if Step type or ID field is nil
	if a == nil || a.EndColumn == nil {
		return 0
	}

	return *a.EndColumn
}

// GetAnnotationLevel returns the AnnotationLevel field.
//
// When the provided Annotation type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Annotation) GetAnnotationLevel() string {
	// return zero value if Step type or ID field is nil
	if a == nil || a.AnnotationLevel == nil {
		return ""
	}

	return *a.AnnotationLevel
}

// GetMessage returns the Message field.
//
// When the provided Annotation type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Annotation) GetMessage() string {
	// return zero value if Step type or ID field is nil
	if a == nil || a.Message == nil {
		return ""
	}

	return *a.Message
}

// GetTitle returns the Title field.
//
// When the provided Annotation type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Annotation) GetTitle() string {
	// return zero value if Step type or ID field is nil
	if a == nil || a.Title == nil {
		return ""
	}

	return *a.Title
}

// GetRawDetails returns the RawDetails field.
//
// When the provided Annotation type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Annotation) GetRawDetails() string {
	// return zero value if Step type or ID field is nil
	if a == nil || a.RawDetails == nil {
		return ""
	}

	return *a.RawDetails
}
