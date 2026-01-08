// SPDX-License-Identifier: Apache-2.0

package pipeline

// ArtifactSlice is the pipeline representation
// of a slice of artifacts.
//
// swagger:model PipelineArtifactSlice
type ArtifactSlice []*Artifacts

// Artifact is the pipeline representation
// of artifacts for a pipeline.
//
// swagger:model PipelineArtifact
type Artifacts struct {
	Paths []string `yaml:"paths,omitempty" json:"paths,omitempty"`
}

// Empty returns true if the provided Artifact is empty.
func (a *Artifacts) Empty() bool {
	// return true if paths field is empty
	if len(a.Paths) == 0 {
		return true
	}

	// return false if Paths are provided
	return false
}
