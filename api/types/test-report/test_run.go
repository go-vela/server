package test_report

import "time"

// TestRun is the API types representation of a test run for a pipeline.
//
// swagger:model TestRun
type TestRun struct {
	ID                      *int64                   `json:"id,omitempty"`
	TestRunSystemAttributes *TestRunSystemAttributes `json:"test_run_system_attributes,omitempty"`
	TotalTestCount          *int                     `json:"total_test_count,omitempty"`
	TotalPassingCount       *int                     `json:"total_passing_count,omitempty"`
	TotalSkipCount          *int                     `json:"total_skip_count,omitempty"`
	TotalFailureCount       *int                     `json:"total_failure_count,omitempty"`
	Passed                  *bool                    `json:"passed,omitempty"`
	Created                 *int64                   `json:"created,omitempty"`
	Started                 *int64                   `json:"started,omitempty"`
	Finished                *int64                   `json:"finished,omitempty"`
	AverageDuration         *int64                   `json:"average_duration,omitempty"`
	SlowestTestCaseDuration *int64                   `json:"slowest_test_case_duration,omitempty"`
	WallClockDuration       *int64                   `json:"wall_clock_duration,omitempty"`
}

// In Vela, int64 is used to stored Unix timestamps instead of float64 for time.Time
// created_timestamp is replaced by Created
// Original Projektor model has cumulative_duration field which is not present in this model
// Created, Started, Finished are Vela standard model fields
// AverageDuration, SlowestTestCaseDuration, WallClockDuration might be taken out and calculated on the fly
// like in the Duration method

// Duration calculates and returns the total amount of
// time the build ran for in a human-readable format.
func (b *TestRun) Duration() string {
	// check if the build doesn't have a started timestamp
	if b.GetStarted() == 0 {
		return "..."
	}

	// capture started unix timestamp from the build
	started := time.Unix(b.GetStarted(), 0)

	// check if the build doesn't have a finished timestamp
	if b.GetFinished() == 0 {
		// return the duration in a human-readable form by
		// subtracting the build started time from the
		// current time rounded to the nearest second
		return time.Since(started).Round(time.Second).String()
	}

	// capture finished unix timestamp from the build
	finished := time.Unix(b.GetFinished(), 0)

	// calculate the duration by subtracting the build
	// started time from the build finished time
	duration := finished.Sub(started)

	// return the duration in a human-readable form
	return duration.String()
}

// GetID returns the ID field.
//
// When the provided TestRun type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *TestRun) GetID() int64 {
	// return zero value if TestRun type or ID field is nil
	if b == nil || b.ID == nil {
		return 0
	}
	return *b.ID
}

// SetID sets the ID field.
func (b *TestRun) SetID(v int64) {
	if b == nil {
		return
	}
	b.ID = &v
}

// GetTestRunSystemAttributes returns the TestRunSystemAttributes field.
//
// When the provided TestRun type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *TestRun) GetTestRunSystemAttributes() TestRunSystemAttributes {
	// return zero value if TestRun type or TestRunSystemAttributes field is nil
	if b == nil || b.TestRunSystemAttributes == nil {
		return TestRunSystemAttributes{}
	}
	return *b.TestRunSystemAttributes
}

// SetTestRunSystemAttributes sets the TestRunSystemAttributes field.
func (b *TestRun) SetTestRunSystemAttributes(v TestRunSystemAttributes) {
	if b == nil {
		return
	}
	b.TestRunSystemAttributes = &v
}

// GetTotalTestCount returns the TotalTestCount field.
//
// When the provided TestRun type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *TestRun) GetTotalTestCount() int {
	// return zero value if TestRun type or TotalTestCount field is nil
	if b == nil || b.TotalTestCount == nil {
		return 0
	}
	return *b.TotalTestCount
}

// SetTotalTestCount sets the TotalTestCount field.
func (b *TestRun) SetTotalTestCount(v int) {
	if b == nil {
		return
	}
	b.TotalTestCount = &v
}

// GetTotalPassingCount returns the TotalPassingCount field.
//
// When the provided TestRun type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *TestRun) GetTotalPassingCount() int {
	// return zero value if TestRun type or TotalPassingCount field is nil
	if b == nil || b.TotalPassingCount == nil {
		return 0
	}
	return *b.TotalPassingCount
}

// SetTotalPassingCount sets the TotalPassingCount field.
func (b *TestRun) SetTotalPassingCount(v int) {
	if b == nil {
		return
	}
	b.TotalPassingCount = &v
}

// GetTotalSkipCount returns the TotalSkipCount field.
//
// When the provided TestRun type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *TestRun) GetTotalSkipCount() int {
	// return zero value if TestRun type or TotalSkipCount field is nil
	if b == nil || b.TotalSkipCount == nil {
		return 0
	}
	return *b.TotalSkipCount
}

// SetTotalSkipCount sets the TotalSkipCount field.
func (b *TestRun) SetTotalSkipCount(v int) {
	if b == nil {
		return
	}
	b.TotalSkipCount = &v
}

// GetTotalFailureCount returns the TotalFailureCount field.
//
// When the provided TestRun type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *TestRun) GetTotalFailureCount() int {
	// return zero value if TestRun type or TotalFailureCount field is nil
	if b == nil || b.TotalFailureCount == nil {
		return 0
	}
	return *b.TotalFailureCount
}

// SetTotalFailureCount sets the TotalFailureCount field.
func (b *TestRun) SetTotalFailureCount(v int) {
	if b == nil {
		return
	}
	b.TotalFailureCount = &v
}

// GetPassed returns the Passed field.
//
// When the provided TestRun type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *TestRun) GetPassed() bool {
	// return zero value if TestRun type or Passed field is nil
	if b == nil || b.Passed == nil {
		return false
	}
	return *b.Passed
}

// SetPassed sets the Passed field.
func (b *TestRun) SetPassed(v bool) {
	if b == nil {
		return
	}
	b.Passed = &v
}

// GetCreated returns the Created field.
//
// When the provided TestRun type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *TestRun) GetCreated() int64 {
	// return zero value if TestRun type or Created field is nil
	if b == nil || b.Created == nil {
		return 0
	}
	return *b.Created
}

// SetCreated sets the Created field.
func (b *TestRun) SetCreated(v int64) {
	if b == nil {
		return
	}
	b.Created = &v
}

// GetStarted returns the Started field.
//
// When the provided TestRun type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *TestRun) GetStarted() int64 {
	// return zero value if TestRun type or Started field is nil
	if b == nil || b.Started == nil {
		return 0
	}
	return *b.Started
}

// SetStarted sets the Started field.
func (b *TestRun) SetStarted(v int64) {
	if b == nil {
		return
	}
	b.Started = &v
}

// GetFinished returns the Finished field.
//
// When the provided TestRun type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *TestRun) GetFinished() int64 {
	// return zero value if TestRun type or Finished field is nil
	if b == nil || b.Finished == nil {
		return 0
	}
	return *b.Finished
}

// SetFinished sets the Finished field.
func (b *TestRun) SetFinished(v int64) {
	if b == nil {
		return
	}
	b.Finished = &v
}

// GetAverageDuration returns the AverageDuration field.
//
// When the provided TestRun type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *TestRun) GetAverageDuration() int64 {
	// return zero value if TestRun type or AverageDuration field is nil
	if b == nil || b.AverageDuration == nil {
		return 0
	}
	return *b.AverageDuration
}

// SetAverageDuration sets the AverageDuration field.
func (b *TestRun) SetAverageDuration(v int64) {
	if b == nil {
		return
	}
	b.AverageDuration = &v
}

// GetSlowestTestCaseDuration returns the SlowestTestCaseDuration field.
//
// When the provided TestRun type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *TestRun) GetSlowestTestCaseDuration() int64 {
	// return zero value if TestRun type or SlowestTestCaseDuration field is nil
	if b == nil || b.SlowestTestCaseDuration == nil {
		return 0
	}
	return *b.SlowestTestCaseDuration
}

// SetSlowestTestCaseDuration sets the SlowestTestCaseDuration field.
func (b *TestRun) SetSlowestTestCaseDuration(v int64) {
	if b == nil {
		return
	}
	b.SlowestTestCaseDuration = &v
}

// GetWallClockDuration returns the WallClockDuration field.
//
// When the provided TestRun type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *TestRun) GetWallClockDuration() int64 {
	// return zero value if TestRun type or WallClockDuration field is nil
	if b == nil || b.WallClockDuration == nil {
		return 0
	}
	return *b.WallClockDuration
}

// SetWallClockDuration sets the WallClockDuration field.
func (b *TestRun) SetWallClockDuration(v int64) {
	if b == nil {
		return
	}
	b.WallClockDuration = &v
}
