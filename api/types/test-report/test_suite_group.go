package test_report

// TestSuiteGroup is the API types representation of a group of test suites for a test run.
//
// swagger:model TestSuiteGroup
type TestSuiteGroup struct {
	ID         *int64   `json:"id,omitempty"`
	TestRun    *TestRun `json:"test_run,omitempty"`
	GroupName  *string  `json:"group_name,omitempty"`
	GroupLabel *string  `json:"group_label,omitempty"`
	Directory  *string  `json:"directory,omitempty"`
}

// SetID sets the ID field.
//
// When the provided TestSuiteGroup type is nil, it
// will set nothing and immediately return.
func (t *TestSuiteGroup) SetID(v int64) {
	// return if TestSuiteGroup type is nil
	if t == nil {
		return
	}

	t.ID = &v
}

// GetID returns the ID field.
//
// When the provided TestSuiteGroup type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuiteGroup) GetID() int64 {
	// return zero value if TestSuiteGroup type or ID field is nil
	if t == nil || t.ID == nil {
		return 0
	}

	return *t.ID
}

// SetTestRun sets the TestRun field.
//
// When the provided TestSuiteGroup type is nil, it
// will set nothing and immediately return.
func (t *TestSuiteGroup) SetTestRun(v TestRun) {
	// return if TestSuiteGroup type is nil
	if t == nil {
		return
	}

	t.TestRun = &v
}

// GetTestRun returns the TestRun field.
//
// When the provided TestSuiteGroup type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuiteGroup) GetTestRun() TestRun {
	// return zero value if TestSuiteGroup type or TestRun field is nil
	if t == nil || t.TestRun == nil {
		return TestRun{}
	}

	return *t.TestRun
}

// SetGroupName sets the GroupName field.
//
// When the provided TestSuiteGroup type is nil, it
// will set nothing and immediately return.
func (t *TestSuiteGroup) SetGroupName(v string) {
	// return if TestSuiteGroup type is nil
	if t == nil {
		return
	}

	t.GroupName = &v
}

// GetGroupName returns the GroupName field.
//
// When the provided TestSuiteGroup type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuiteGroup) GetGroupName() string {
	// return zero value if TestSuiteGroup type or GroupName field is nil
	if t == nil || t.GroupName == nil {
		return ""
	}

	return *t.GroupName
}

// SetGroupLabel sets the GroupLabel field.
//
// When the provided TestSuiteGroup type is nil, it
// will set nothing and immediately return.
func (t *TestSuiteGroup) SetGroupLabel(v string) {
	// return if TestSuiteGroup type is nil
	if t == nil {
		return
	}

	t.GroupLabel = &v
}

// GetGroupLabel returns the GroupLabel field.
//
// When the provided TestSuiteGroup type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuiteGroup) GetGroupLabel() string {
	// return zero value if TestSuiteGroup type or GroupLabel field is nil
	if t == nil || t.GroupLabel == nil {
		return ""
	}

	return *t.GroupLabel
}

// SetDirectory sets the Directory field.
//
// When the provided TestSuiteGroup type is nil, it
// will set nothing and immediately return.
func (t *TestSuiteGroup) SetDirectory(v string) {
	// return if TestSuiteGroup type is nil
	if t == nil {
		return
	}

	t.Directory = &v
}

// GetDirectory returns the Directory field.
//
// When the provided TestSuiteGroup type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuiteGroup) GetDirectory() string {
	// return zero value if TestSuiteGroup type or Directory field is nil
	if t == nil || t.Directory == nil {
		return ""
	}

	return *t.Directory
}
