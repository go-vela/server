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

// Purge removes the test report configuration from the pipeline
// if it does not match the provided ruledata. If both results
// and attachments are provided, then an empty test report is returned.
//func (t *TestReport) Purge(r *RuleData) (*TestReport, error) {
//	// return an empty test report if both results and attachments are provided
//	if len(t.Results) > 0 && len(t.Attachments) > 0 {
//		return nil, fmt.Errorf("cannot have both results and attachments in the test report")
//	}
//
//	// purge results if provided
//	if len(t.Results) > 0 {
//		t.Results = ""
//	}
//
//	// purge attachments if provided
//	if len(t.Attachments) > 0 {
//		t.Attachments = ""
//	}
//
//	// return the purged test report
//	return t, nil
//}

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
