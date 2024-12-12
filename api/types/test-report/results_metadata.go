package test_report

// ResultsMetadata is the API types representation of metadata for a test run.
//
// swagger:model ResultsMetadata
type ResultsMetadata struct {
	ID      *int64   `json:"id,omitempty"`
	TestRun *TestRun `json:"test_run,omitempty"`
	CI      *bool    `json:"ci,omitempty"`
	Group   *string  `json:"group,omitempty"`
}

// SetID sets the ID field.
//
// When the provided ResultsMetadata type is nil, it
// will set nothing and immediately return.
func (r *ResultsMetadata) SetID(v int64) {
	// return if ResultsMetadata type is nil
	if r == nil {
		return
	}

	r.ID = &v
}

// GetID returns the ID field.
//
// When the provided ResultsMetadata type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *ResultsMetadata) GetID() int64 {
	// return zero value if ResultsMetadata type or ID field is nil
	if r == nil || r.ID == nil {
		return 0
	}

	return *r.ID
}

// SetTestRun sets the TestRun field.
//
// When the provided ResultsMetadata type is nil, it
// will set nothing and immediately return.
func (r *ResultsMetadata) SetTestRun(v TestRun) {
	// return if ResultsMetadata type is nil
	if r == nil {
		return
	}

	r.TestRun = &v
}

// GetTestRun returns the TestRun field.
//
// When the provided ResultsMetadata type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *ResultsMetadata) GetTestRun() TestRun {
	// return zero value if ResultsMetadata type or TestRun field is nil
	if r == nil || r.TestRun == nil {
		return TestRun{}
	}

	return *r.TestRun
}

// SetCI sets the CI field.
//
// When the provided ResultsMetadata type is nil, it
// will set nothing and immediately return.
func (r *ResultsMetadata) SetCI(v bool) {
	// return if ResultsMetadata type is nil
	if r == nil {
		return
	}

	r.CI = &v
}

// GetCI returns the CI field.
//
// When the provided ResultsMetadata type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *ResultsMetadata) GetCI() bool {
	// return zero value if ResultsMetadata type or CI field is nil
	if r == nil || r.CI == nil {
		return false
	}

	return *r.CI
}

// SetGroup sets the Group field.
//
// When the provided ResultsMetadata type is nil, it
// will set nothing and immediately return.
func (r *ResultsMetadata) SetGroup(v string) {
	// return if ResultsMetadata type is nil
	if r == nil {
		return
	}

	r.Group = &v
}

// GetGroup returns the Group field.
//
// When the provided ResultsMetadata type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *ResultsMetadata) GetGroup() string {
	// return zero value if ResultsMetadata type or Group field is nil
	if r == nil || r.Group == nil {
		return ""
	}

	return *r.Group
}
