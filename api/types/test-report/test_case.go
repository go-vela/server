package test_report

import "time"

// TestCases is the API types representation of a test case.
//
// swagger:model TestCase

type TestCase struct {
	ID           *int64     `json:"id,omitempty"`
	TestSuite    *TestSuite `json:"test_suite,omitempty"`
	Idx          *int64     `json:"idx,omitempty"`
	Name         *string    `json:"name,omitempty"`
	PackageName  *string    `json:"package_name,omitempty"`
	ClassName    *string    `json:"class_name,omitempty"`
	Created      *int64     `json:"created,omitempty"`
	Started      *int64     `json:"started,omitempty"`
	Finished     *int64     `json:"finished,omitempty"`
	Passed       *bool      `json:"passed,omitempty"`
	Skipped      *bool      `json:"skipped,omitempty"`
	SystemOut    *string    `json:"system_out,omitempty"`
	SystemErr    *string    `json:"system_err,omitempty"`
	HasSystemOut *bool      `json:"has_system_out,omitempty"`
	HasSystemErr *bool      `json:"has_system_err,omitempty"`
}

// Duration calculates and returns the total amount of
// time the test case ran for in a human-readable format.
func (t *TestCase) Duration() string {
	// check if the test case doesn't have a started timestamp
	if t.GetStarted() == 0 {
		return "..."
	}

	// capture started unix timestamp from the test case
	started := time.Unix(t.GetStarted(), 0)

	// check if the test case doesn't have a finished timestamp
	if t.GetFinished() == 0 {
		// return the duration in a human-readable form by
		// subtracting the test case started time from the
		// current time rounded to the nearest second
		return time.Since(started).Round(time.Second).String()
	}

	// capture finished unix timestamp from the test case
	finished := time.Unix(t.GetFinished(), 0)

	// calculate the duration by subtracting the test case
	// started time from the test case finished time
	duration := finished.Sub(started)

	// return the duration in a human-readable form
	return duration.String()
}

// SetID sets the ID field.
//
// When the provided TestCase type is nil, it
// will set nothing and immediately return.
func (t *TestCase) SetID(v int64) {
	// return if TestCase type is nil
	if t == nil {
		return
	}

	t.ID = &v
}

// GetID returns the ID field.
//
// When the provided TestCase type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestCase) GetID() int64 {
	// return zero value if TestCase type or ID field is nil
	if t == nil || t.ID == nil {
		return 0
	}

	return *t.ID
}

// SetIdx sets the Idx field.
//
// When the provided TestCase type is nil, it
// will set nothing and immediately return.
func (t *TestCase) SetIdx(v int64) {
	// return if TestCase type is nil
	if t == nil {
		return
	}

	t.Idx = &v
}

// GetIdx returns the Idx field.
//
// When the provided TestCase type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestCase) GetIdx() int64 {
	// return zero value if TestCase type or Idx field is nil
	if t == nil || t.Idx == nil {
		return 0
	}

	return *t.Idx
}

// SetName sets the Name field.
//
// When the provided TestCase type is nil, it
// will set nothing and immediately return.
func (t *TestCase) SetName(v string) {
	// return if TestCase type is nil
	if t == nil {
		return
	}

	t.Name = &v
}

// GetName returns the Name field.
//
// When the provided TestCase type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestCase) GetName() string {
	// return zero value if TestCase type or Name field is nil
	if t == nil || t.Name == nil {
		return ""
	}

	return *t.Name
}

// SetPackageName sets the PackageName field.
//
// When the provided TestCase type is nil, it
// will set nothing and immediately return.
func (t *TestCase) SetPackageName(v string) {
	// return if TestCase type is nil
	if t == nil {
		return
	}

	t.PackageName = &v
}

// GetPackageName returns the PackageName field.
//
// When the provided TestCase type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestCase) GetPackageName() string {
	// return zero value if TestCase type or PackageName field is nil
	if t == nil || t.PackageName == nil {
		return ""
	}

	return *t.PackageName
}

// SetClassName sets the ClassName field.
//
// When the provided TestCase type is nil, it
// will set nothing and immediately return.
func (t *TestCase) SetClassName(v string) {
	// return if TestCase type is nil
	if t == nil {
		return
	}

	t.ClassName = &v
}

// GetClassName returns the ClassName field.
//
// When the provided TestCase type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestCase) GetClassName() string {
	// return zero value if TestCase type or ClassName field is nil
	if t == nil || t.ClassName == nil {
		return ""
	}

	return *t.ClassName
}

// SetCreated sets the Created field.
//
// When the provided TestCase type is nil, it
// will set nothing and immediately return.
func (t *TestCase) SetCreated(v int64) {
	// return if TestCase type is nil
	if t == nil {
		return
	}

	t.Created = &v
}

// GetCreated returns the Created field.
//
// When the provided TestCase type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestCase) GetCreated() int64 {
	// return zero value if TestCase type or Created field is nil
	if t == nil || t.Created == nil {
		return 0
	}

	return *t.Created
}

// SetStarted sets the Started field.
//
// When the provided TestCase type is nil, it
// will set nothing and immediately return.
func (t *TestCase) SetStarted(v int64) {
	// return if TestCase type is nil
	if t == nil {
		return
	}

	t.Started = &v
}

// GetStarted returns the Started field.
//
// When the provided TestCase type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestCase) GetStarted() int64 {
	// return zero value if TestCase type or Started field is nil
	if t == nil || t.Started == nil {
		return 0
	}

	return *t.Started
}

// SetFinished sets the Finished field.
//
// When the provided TestCase type is nil, it
// will set nothing and immediately return.
func (t *TestCase) SetFinished(v int64) {
	// return if TestCase type is nil
	if t == nil {
		return
	}

	t.Finished = &v
}

// GetFinished returns the Finished field.
//
// When the provided TestCase type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestCase) GetFinished() int64 {
	// return zero value if TestCase type or Finished field is nil
	if t == nil || t.Finished == nil {
		return 0
	}

	return *t.Finished
}

// SetPassed sets the Passed field.
//
// When the provided TestCase type is nil, it
// will set nothing and immediately return.
func (t *TestCase) SetPassed(v bool) {
	// return if TestCase type is nil
	if t == nil {
		return
	}

	t.Passed = &v
}

// GetPassed returns the Passed field.
//
// When the provided TestCase type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestCase) GetPassed() bool {
	// return zero value if TestCase type or Passed field is nil
	if t == nil || t.Passed == nil {
		return false
	}

	return *t.Passed
}

// SetSkipped sets the Skipped field.
//
// When the provided TestCase type is nil, it
// will set nothing and immediately return.
func (t *TestCase) SetSkipped(v bool) {
	// return if TestCase type is nil
	if t == nil {
		return
	}

	t.Skipped = &v
}

// GetSkipped returns the Skipped field.
//
// When the provided TestCase type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestCase) GetSkipped() bool {
	// return zero value if TestCase type or Skipped field is nil
	if t == nil || t.Skipped == nil {
		return false
	}

	return *t.Skipped
}

// SetSystemOut sets the SystemOut field.
//
// When the provided TestCase type is nil, it
// will set nothing and immediately return.
func (t *TestCase) SetSystemOut(v string) {
	// return if TestCase type is nil
	if t == nil {
		return
	}

	t.SystemOut = &v
}

// GetSystemOut returns the SystemOut field.
//
// When the provided TestCase type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestCase) GetSystemOut() string {
	// return zero value if TestCase type or SystemOut field is nil
	if t == nil || t.SystemOut == nil {
		return ""
	}

	return *t.SystemOut
}

// SetSystemErr sets the SystemErr field.
//
// When the provided TestCase type is nil, it
// will set nothing and immediately return.
func (t *TestCase) SetSystemErr(v string) {
	// return if TestCase type is nil
	if t == nil {
		return
	}

	t.SystemErr = &v
}

// GetSystemErr returns the SystemErr field.
//
// When the provided TestCase type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestCase) GetSystemErr() string {
	// return zero value if TestCase type or SystemErr field is nil
	if t == nil || t.SystemErr == nil {
		return ""
	}

	return *t.SystemErr
}

// SetHasSystemOut sets the HasSystemOut field.
//
// When the provided TestCase type is nil, it
// will set nothing and immediately return.
func (t *TestCase) SetHasSystemOut(v bool) {
	// return if TestCase type is nil
	if t == nil {
		return
	}

	t.HasSystemOut = &v
}

// GetHasSystemOut returns the HasSystemOut field.
//
// When the provided TestCase type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestCase) GetHasSystemOut() bool {
	// return zero value if TestCase type or HasSystemOut field is nil
	if t == nil || t.HasSystemOut == nil {
		return false
	}

	return *t.HasSystemOut
}

// SetHasSystemErr sets the HasSystemErr field.
//
// When the provided TestCase type is nil, it
// will set nothing and immediately return.
func (t *TestCase) SetHasSystemErr(v bool) {
	// return if TestCase type is nil
	if t == nil {
		return
	}

	t.HasSystemErr = &v
}

// GetHasSystemErr returns the HasSystemErr field.
//
// When the provided TestCase type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestCase) GetHasSystemErr() bool {
	// return zero value if TestCase type or HasSystemErr field is nil
	if t == nil || t.HasSystemErr == nil {
		return false
	}

	return *t.HasSystemErr
}
