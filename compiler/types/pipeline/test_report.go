// SPDX-License-Identifier: Apache-2.0

package pipeline

// Artifacts represents the structure for test report configuration.
type (
	// ArtifactsSlice is the pipleine representation
	//of a slice of Artifacts.
	//
	// swagger:model PipelineTestReportSlice
	ArtifactsSlice []*Artifacts

	// Artifacts is the pipeline representation
	// of a test report for a pipeline.
	//
	// swagger:model PipelineTestReport
	Artifacts struct {
		Paths []string `yaml:"paths,omitempty"     json:"paths,omitempty"`
	}
)

// Empty returns true if the provided Artifacts is empty.
func (t *Artifacts) Empty() bool {
	// return true if paths field is empty
	if len(t.Paths) == 0 {
		return true
	}

	// return false if Paths are provided
	return false
}
