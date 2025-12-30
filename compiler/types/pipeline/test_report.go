// SPDX-License-Identifier: Apache-2.0

package pipeline

// TestReport represents the structure for test report configuration.
type (
	// TestReportSlice is the pipleine representation
	//of a slice of TestReport.
	//
	// swagger:model PipelineTestReportSlice
	TestReportSlice []*TestReport

	// TestReport is the pipeline representation
	// of a test report for a pipeline.
	//
	// swagger:model PipelineTestReport
	TestReport struct {
		Results     []string `yaml:"results,omitempty"     json:"results,omitempty"`
		Attachments []string `yaml:"attachments,omitempty" json:"attachments,omitempty"`
	}
)

// Empty returns true if the provided test report is empty.
func (t *TestReport) Empty() bool {
	// return true if every test report field is empty
	if len(t.Results) == 0 &&
		len(t.Attachments) == 0 {
		return true
	}

	// return false if any of the test report fields are provided
	return false
}
