package types

import "fmt"

// TestReport is the API representation of a test report for a pipeline.
//
// swagger:model TestReport
type TestReport struct {
	Results     *string `json:"results,omitempty"`
	Attachments *string `json:"attachments,omitempty"`
}

// GetResults returns the Results field.
//
// When the provided TestReport type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestReport) GetResults() string {
	// return zero value if TestReport type or Results field is nil
	if t == nil || t.Results == nil {
		return ""
	}

	return *t.Results
}

// GetAttachments returns the Attachments field.
//
// When the provided TestReport type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestReport) GetAttachments() string {
	// return zero value if TestReport type or Attachments field is nil
	if t == nil || t.Attachments == nil {
		return ""
	}

	return *t.Attachments
}

// SetResults sets the Results field.
func (t *TestReport) SetResults(v string) {
	// return if TestReport type is nil
	if t == nil {
		return
	}
	// set the Results field
	t.Results = &v
}

// SetAttachments sets the Attachments field.
func (t *TestReport) SetAttachments(v string) {
	// return if TestReport type is nil
	if t == nil {
		return
	}
	// set the Attachments field
	t.Attachments = &v
}

// String implements the Stringer interface for the TestReport type.
func (t *TestReport) String() string {
	return fmt.Sprintf("Results: %s, Attachments: %s", t.GetResults(), t.GetAttachments())
}
