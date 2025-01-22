package test_report

// PerformanceResults is the API types representation of performance results for a test run.
//
// swagger:model PerformanceResults
type PerformanceResults struct {
	ID                *int64   `json:"id,omitempty"`
	TestRun           *TestRun `json:"test_run,omitempty"`
	TestRunPublicID   *string  `json:"test_run_public_id,omitempty"`
	Name              *string  `json:"name,omitempty"`
	RequestCount      *int64   `json:"request_count,omitempty"`
	RequestsPerSecond *float64 `json:"requests_per_second,omitempty"`
	Average           *float64 `json:"average,omitempty"`
	Maximum           *float64 `json:"maximum,omitempty"`
	P95               *float64 `json:"p95,omitempty"`
} // GetID returns the ID field.
// When the provided PerformanceResults type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *PerformanceResults) GetID() int64 {
	// return zero value if PerformanceResults type or ID field is nil
	if p == nil || p.ID == nil {
		return 0
	}

	return *p.ID
}

// GetTestRun returns the TestRun field.
//
// When the provided PerformanceResults type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *PerformanceResults) GetTestRun() TestRun {
	// return zero value if PerformanceResults type or TestRun field is nil
	if p == nil || p.TestRun == nil {
		return TestRun{}
	}

	return *p.TestRun
}

// GetTestRunPublicID returns the TestRunPublicID field.
//
// When the provided PerformanceResults type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *PerformanceResults) GetTestRunPublicID() string {
	// return zero value if PerformanceResults type or TestRunPublicID field is nil
	if p == nil || p.TestRunPublicID == nil {
		return ""
	}

	return *p.TestRunPublicID
}

// GetName returns the Name field.
//
// When the provided PerformanceResults type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *PerformanceResults) GetName() string {
	// return zero value if PerformanceResults type or Name field is nil
	if p == nil || p.Name == nil {
		return ""
	}

	return *p.Name
}

// GetRequestCount returns the RequestCount field.
//
// When the provided PerformanceResults type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *PerformanceResults) GetRequestCount() int64 {
	// return zero value if PerformanceResults type or RequestCount field is nil
	if p == nil || p.RequestCount == nil {
		return 0
	}

	return *p.RequestCount
}

// GetRequestsPerSecond returns the RequestsPerSecond field.
//
// When the provided PerformanceResults type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *PerformanceResults) GetRequestsPerSecond() float64 {
	// return zero value if PerformanceResults type or RequestsPerSecond field is nil
	if p == nil || p.RequestsPerSecond == nil {
		return 0.0
	}

	return *p.RequestsPerSecond
}

// GetAverage returns the Average field.
//
// When the provided PerformanceResults type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *PerformanceResults) GetAverage() float64 {
	// return zero value if PerformanceResults type or Average field is nil
	if p == nil || p.Average == nil {
		return 0.0
	}

	return *p.Average
}

// GetMaximum returns the Maximum field.
//
// When the provided PerformanceResults type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *PerformanceResults) GetMaximum() float64 {
	// return zero value if PerformanceResults type or Maximum field is nil
	if p == nil || p.Maximum == nil {
		return 0.0
	}

	return *p.Maximum
}

// GetP95 returns the P95 field.
//
// When the provided PerformanceResults type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *PerformanceResults) GetP95() float64 {
	// return zero value if PerformanceResults type or P95 field is nil
	if p == nil || p.P95 == nil {
		return 0.0
	}

	return *p.P95
}
