package test_report

// TestFailure is the API types representation of a test for a pipeline.
//
// swagger:model TestFailure
type TestFailure struct {
	ID             *int64  `json:"id,omitempty"`
	TestCaseID     *int64  `json:"test_case_id,omitempty"`
	FailureMessage *string `json:"failure_message,omitempty"`
	FailureType    *string `json:"failure_type,omitempty"`
	FailureText    *string `json:"failure_text,omitempty"`
}

// GetID returns the ID field.
//
// When the provided TestFailure type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestFailure) GetID() int64 {
	if t == nil || t.ID == nil {
		return 0
	}
	return *t.ID
}

// GetTestCaseID returns the TestCaseID field.
//
// When the provided TestFailure type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestFailure) GetTestCaseID() int64 {
	if t == nil || t.TestCaseID == nil {
		return 0
	}
	return *t.TestCaseID
}

// GetFailureMessage returns the FailureMessage field.
//
// When the provided TestFailure type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestFailure) GetFailureMessage() string {
	if t == nil || t.FailureMessage == nil {
		return ""
	}
	return *t.FailureMessage
}

// GetFailureType returns the FailureType field.
//
// When the provided TestFailure type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestFailure) GetFailureType() string {
	if t == nil || t.FailureType == nil {
		return ""
	}
	return *t.FailureType
}

// GetFailureText returns the FailureText field.
//
// When the provided TestFailure type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestFailure) GetFailureText() string {
	if t == nil || t.FailureText == nil {
		return ""
	}
	return *t.FailureText
}

// SetID sets the ID field.
//
// When the provided TestFailure type is nil, it
// will set nothing and immediately return.
func (t *TestFailure) SetID(v int64) {
	// return if TestFailure type is nil
	if t == nil {
		return
	}

	t.ID = &v
}

// SetTestCaseID sets the TestCaseID field.
//
// When the provided TestFailure type is nil, it
// will set nothing and immediately return.
func (t *TestFailure) SetTestCaseID(v int64) {
	// return if TestFailure type is nil
	if t == nil {
		return
	}

	t.TestCaseID = &v
}

// SetFailureMessage sets the FailureMessage field.
//
// When the provided TestFailure type is nil, it
// will set nothing and immediately return.
func (t *TestFailure) SetFailureMessage(v string) {
	// return if TestFailure type is nil
	if t == nil {
		return
	}

	t.FailureMessage = &v
}

// SetFailureType sets the FailureType field.
//
// When the provided TestFailure type is nil, it
// will set nothing and immediately return.
func (t *TestFailure) SetFailureType(v string) {
	// return if TestFailure type is nil
	if t == nil {
		return
	}

	t.FailureType = &v
}

// SetFailureText sets the FailureText field.
//
// When the provided TestFailure type is nil, it
// will set nothing and immediately return.
func (t *TestFailure) SetFailureText(v string) {
	// return if TestFailure type is nil
	if t == nil {
		return
	}

	t.FailureText = &v
}
