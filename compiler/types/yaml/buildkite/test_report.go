package buildkite

import (
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/yaml/yaml"
)

// TestReport represents the structure for test report configuration.
type TestReport struct {
	Results     []string `yaml:"results,omitempty" json:"results,omitempty"`
	Attachments []string `yaml:"attachments,omitempty" json:"attachments,omitempty"`
}

// ToPipeline converts the TestReport type
// to a pipeline TestReport type.
func (t *TestReport) ToPipeline() *pipeline.TestReport {
	return &pipeline.TestReport{
		Results:     t.Results,
		Attachments: t.Attachments,
	}
}

// UnmarshalYAML implements the Unmarshaler interface for the TestReport type.
func (t *TestReport) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// test report we try unmarshalling to
	testReport := new(struct {
		Results     []string `yaml:"results,omitempty" json:"results,omitempty"`
		Attachments []string `yaml:"attachments,omitempty" json:"attachments,omitempty"`
	})

	// attempt to unmarshal test report type
	err := unmarshal(testReport)
	if err != nil {
		return err
	}

	// set the results field
	t.Results = testReport.Results
	// set the attachments field
	t.Attachments = testReport.Attachments

	return nil
}

func (t *TestReport) ToYAML() yaml.TestReport {
	return yaml.TestReport{
		Results:     t.Results,
		Attachments: t.Attachments,
	}
}
