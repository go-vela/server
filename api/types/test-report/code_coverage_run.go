package test_report

// CodeCoverageRun is the API types representation of a code coverage run for a pipeline.
//
// swagger:model CodeCoverageRun
type CodeCoverageRun struct {
	ID              *int64  `json:"id,omitempty"`
	TestRunPublicID *string `json:"test_run_public_id,omitempty"`
}

// GetID returns the ID field.
//
// When the provided CodeCoverageRun type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageRun) GetID() int64 {
	if c == nil || c.ID == nil {
		return 0
	}
	return *c.ID
}

// GetTestRunPublicID returns the TestRunPublicID field.
//
// When the provided CodeCoverageRun type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageRun) GetTestRunPublicID() string {
	if c == nil || c.TestRunPublicID == nil {
		return ""
	}
	return *c.TestRunPublicID
}

// SetID sets the ID field.
//
// When the provided CodeCoverageRun type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageRun) SetID(v int64) {
	// return if CodeCoverageRun type is nil
	if c == nil {
		return
	}

	c.ID = &v
}

// SetTestRunPublicID sets the TestRunPublicID field.
//
// When the provided CodeCoverageRun type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageRun) SetTestRunPublicID(v string) {
	// return if CodeCoverageRun type is nil
	if c == nil {
		return
	}

	c.TestRunPublicID = &v
}
