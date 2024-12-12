package test_report

import "time"

// TestSuite is the API types representation of a test suite for a test run.
//
// swagger:model TestSuite
type TestSuite struct {
	ID             *int64          `json:"id,omitempty"`
	TestRun        *TestRun        `json:"test_run,omitempty"`
	TestSuiteGroup *TestSuiteGroup `json:"test_suite_group,omitempty"`
	Idx            *int64          `json:"idx,omitempty"`
	PackageName    *string         `json:"package_name,omitempty"`
	ClassName      *string         `json:"class_name,omitempty"`
	TestCount      *int64          `json:"test_count,omitempty"`
	PassingCount   *int64          `json:"passing_count,omitempty"`
	SkippedCount   *int64          `json:"skipped_count,omitempty"`
	FailureCount   *int64          `json:"failure_count,omitempty"`
	StartTs        *int64          `json:"start_ts,omitempty"`
	Hostname       *string         `json:"hostname,omitempty"`
	Started        *int64          `json:"started,omitempty"`
	Finished       *int64          `json:"finished,omitempty"`
	SystemOut      *string         `json:"system_out,omitempty"`
	SystemErr      *string         `json:"system_err,omitempty"`
	HasSystemOut   *bool           `json:"has_system_out,omitempty"`
	HasSystemErr   *bool           `json:"has_system_err,omitempty"`
	FileName       *string         `json:"file_name,omitempty"`
}

// Duration calculates and returns the total amount of
// time the test suite ran for in a human-readable format.
func (t *TestSuite) Duration() string {
	// check if the test suite doesn't have a started timestamp
	if t.GetStarted() == 0 {
		return "..."
	}

	// capture started unix timestamp from the test suite
	started := time.Unix(t.GetStarted(), 0)

	// check if the test suite doesn't have a finished timestamp
	if t.GetFinished() == 0 {
		// return the duration in a human-readable form by
		// subtracting the test suite started time from the
		// current time rounded to the nearest second
		return time.Since(started).Round(time.Second).String()
	}

	// capture finished unix timestamp from the test suite
	finished := time.Unix(t.GetFinished(), 0)

	// calculate the duration by subtracting the test suite
	// started time from the test suite finished time
	duration := finished.Sub(started)

	// return the duration in a human-readable form
	return duration.String()
}

// SetID sets the ID field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetID(v int64) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.ID = &v
}

// GetID returns the ID field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetID() int64 {
	// return zero value if TestSuite type or ID field is nil
	if t == nil || t.ID == nil {
		return 0
	}

	return *t.ID
}

// SetTestRun sets the TestRun field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetTestRun(v TestRun) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.TestRun = &v
}

// GetTestRun returns the TestRun field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetTestRun() TestRun {
	// return zero value if TestSuite type or TestRun field is nil
	if t == nil || t.TestRun == nil {
		return TestRun{}
	}

	return *t.TestRun
}

// SetTestSuiteGroup sets the TestSuiteGroup field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetTestSuiteGroup(v TestSuiteGroup) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.TestSuiteGroup = &v
}

// GetTestSuiteGroup returns the TestSuiteGroup field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetTestSuiteGroup() TestSuiteGroup {
	// return zero value if TestSuite type or TestSuiteGroup field is nil
	if t == nil || t.TestSuiteGroup == nil {
		return TestSuiteGroup{}
	}

	return *t.TestSuiteGroup
}

// SetIdx sets the Idx field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetIdx(v int64) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.Idx = &v
}

// GetIdx returns the Idx field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetIdx() int64 {
	// return zero value if TestSuite type or Idx field is nil
	if t == nil || t.Idx == nil {
		return 0
	}

	return *t.Idx
}

// SetPackageName sets the PackageName field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetPackageName(v string) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.PackageName = &v
}

// GetPackageName returns the PackageName field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetPackageName() string {
	// return zero value if TestSuite type or PackageName field is nil
	if t == nil || t.PackageName == nil {
		return ""
	}

	return *t.PackageName
}

// SetClassName sets the ClassName field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetClassName(v string) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.ClassName = &v
}

// GetClassName returns the ClassName field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetClassName() string {
	// return zero value if TestSuite type or ClassName field is nil
	if t == nil || t.ClassName == nil {
		return ""
	}

	return *t.ClassName
}

// SetTestCount sets the TestCount field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetTestCount(v int64) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.TestCount = &v
}

// GetTestCount returns the TestCount field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetTestCount() int64 {
	// return zero value if TestSuite type or TestCount field is nil
	if t == nil || t.TestCount == nil {
		return 0
	}

	return *t.TestCount
}

// SetPassingCount sets the PassingCount field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetPassingCount(v int64) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.PassingCount = &v
}

// GetPassingCount returns the PassingCount field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetPassingCount() int64 {
	// return zero value if TestSuite type or PassingCount field is nil
	if t == nil || t.PassingCount == nil {
		return 0
	}

	return *t.PassingCount
}

// SetSkippedCount sets the SkippedCount field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetSkippedCount(v int64) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.SkippedCount = &v
}

// GetSkippedCount returns the SkippedCount field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetSkippedCount() int64 {
	// return zero value if TestSuite type or SkippedCount field is nil
	if t == nil || t.SkippedCount == nil {
		return 0
	}

	return *t.SkippedCount
}

// SetFailureCount sets the FailureCount field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetFailureCount(v int64) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.FailureCount = &v
}

// GetFailureCount returns the FailureCount field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetFailureCount() int64 {
	// return zero value if TestSuite type or FailureCount field is nil
	if t == nil || t.FailureCount == nil {
		return 0
	}

	return *t.FailureCount
}

// SetStartTs sets the StartTs field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetStartTs(v int64) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.StartTs = &v
}

// GetStartTs returns the StartTs field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetStartTs() int64 {
	// return zero value if TestSuite type or StartTs field is nil
	if t == nil || t.StartTs == nil {
		return 0
	}

	return *t.StartTs
}

// SetHostname sets the Hostname field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetHostname(v string) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.Hostname = &v
}

// GetHostname returns the Hostname field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetHostname() string {
	// return zero value if TestSuite type or Hostname field is nil
	if t == nil || t.Hostname == nil {
		return ""
	}

	return *t.Hostname
}

// SetStarted sets the Started field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetStarted(v int64) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.Started = &v
}

// GetStarted returns the Started field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetStarted() int64 {
	// return zero value if TestSuite type or Started field is nil
	if t == nil || t.Started == nil {
		return 0
	}

	return *t.Started
}

// SetFinished sets the Finished field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetFinished(v int64) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.Finished = &v
}

// GetFinished returns the Finished field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetFinished() int64 {
	// return zero value if TestSuite type or Finished field is nil
	if t == nil || t.Finished == nil {
		return 0
	}

	return *t.Finished
}

// SetSystemOut sets the SystemOut field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetSystemOut(v string) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.SystemOut = &v
}

// GetSystemOut returns the SystemOut field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetSystemOut() string {
	// return zero value if TestSuite type or SystemOut field is nil
	if t == nil || t.SystemOut == nil {
		return ""
	}

	return *t.SystemOut
}

// SetSystemErr sets the SystemErr field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetSystemErr(v string) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.SystemErr = &v
}

// GetSystemErr returns the SystemErr field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetSystemErr() string {
	// return zero value if TestSuite type or SystemErr field is nil
	if t == nil || t.SystemErr == nil {
		return ""
	}

	return *t.SystemErr
}

// SetHasSystemOut sets the HasSystemOut field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetHasSystemOut(v bool) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.HasSystemOut = &v
}

// GetHasSystemOut returns the HasSystemOut field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetHasSystemOut() bool {
	// return zero value if TestSuite type or HasSystemOut field is nil
	if t == nil || t.HasSystemOut == nil {
		return false
	}

	return *t.HasSystemOut
}

// SetHasSystemErr sets the HasSystemErr field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetHasSystemErr(v bool) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.HasSystemErr = &v
}

// GetHasSystemErr returns the HasSystemErr field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetHasSystemErr() bool {
	// return zero value if TestSuite type or HasSystemErr field is nil
	if t == nil || t.HasSystemErr == nil {
		return false
	}

	return *t.HasSystemErr
}

// SetFileName sets the FileName field.
//
// When the provided TestSuite type is nil, it
// will set nothing and immediately return.
func (t *TestSuite) SetFileName(v string) {
	// return if TestSuite type is nil
	if t == nil {
		return
	}

	t.FileName = &v
}

// GetFileName returns the FileName field.
//
// When the provided TestSuite type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestSuite) GetFileName() string {
	// return zero value if TestSuite type or FileName field is nil
	if t == nil || t.FileName == nil {
		return ""
	}

	return *t.FileName
}
