package test_report

// CodeCoverageFile is the API types representation of a code coverage file for a pipeline.
//
// swagger:model CodeCoverageFile
type CodeCoverageFile struct {
	ID                *int64             `json:"id,omitempty"`
	CodeCoverageRun   *CodeCoverageRun   `json:"code_coverage_run,omitempty"`
	CodeCoverageGroup *CodeCoverageGroup `json:"code_coverage_group,omitempty"`
	CodeCoverageStats *CodeCoverageStats `json:"code_coverage_stats,omitempty"`
	DirectoryName     *string            `json:"directory_name,omitempty"`
	FileName          *string            `json:"file_name,omitempty"`
	MissedLines       *int               `json:"missed_lines,omitempty"`
	PartialLines      *int               `json:"partial_lines,omitempty"`
	FilePath          *string            `json:"file_path,omitempty"`
}

// GetID returns the ID field.
//
// When the provided CodeCoverageFile type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageFile) GetID() int64 {
	if c == nil || c.ID == nil {
		return 0
	}
	return *c.ID
}

// SetID sets the ID field.
//
// When the provided CodeCoverageFile type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageFile) SetID(v int64) {
	if c == nil {
		return
	}
	c.ID = &v
}

// GetCodeCoverageRun returns the CodeCoverageRun field.
//
// When the provided CodeCoverageFile type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageFile) GetCodeCoverageRun() *CodeCoverageRun {
	if c == nil || c.CodeCoverageRun == nil {
		return new(CodeCoverageRun)
	}
	return c.CodeCoverageRun
}

// SetCodeCoverageRun sets the CodeCoverageRun field.
//
// When the provided CodeCoverageFile type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageFile) SetCodeCoverageRun(v *CodeCoverageRun) {
	if c == nil {
		return
	}
	c.CodeCoverageRun = v
}

// GetCodeCoverageGroup returns the CodeCoverageGroup field.
//
// When the provided CodeCoverageFile type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageFile) GetCodeCoverageGroup() *CodeCoverageGroup {
	if c == nil || c.CodeCoverageGroup == nil {
		return new(CodeCoverageGroup)
	}
	return c.CodeCoverageGroup
}

// SetCodeCoverageGroup sets the CodeCoverageGroup field.
//
// When the provided CodeCoverageFile type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageFile) SetCodeCoverageGroup(v *CodeCoverageGroup) {
	if c == nil {
		return
	}
	c.CodeCoverageGroup = v
}

// GetCodeCoverageStats returns the CodeCoverageStats field.
//
// When the provided CodeCoverageFile type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageFile) GetCodeCoverageStats() *CodeCoverageStats {
	if c == nil || c.CodeCoverageStats == nil {
		return new(CodeCoverageStats)
	}
	return c.CodeCoverageStats
}

// SetCodeCoverageStats sets the CodeCoverageStats field.
//
// When the provided CodeCoverageFile type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageFile) SetCodeCoverageStats(v *CodeCoverageStats) {
	if c == nil {
		return
	}
	c.CodeCoverageStats = v
}

// GetDirectoryName returns the DirectoryName field.
//
// When the provided CodeCoverageFile type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageFile) GetDirectoryName() string {
	if c == nil || c.DirectoryName == nil {
		return ""
	}
	return *c.DirectoryName
}

// SetDirectoryName sets the DirectoryName field.
//
// When the provided CodeCoverageFile type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageFile) SetDirectoryName(v string) {
	if c == nil {
		return
	}
	c.DirectoryName = &v
}

// GetFileName returns the FileName field.
//
// When the provided CodeCoverageFile type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageFile) GetFileName() string {
	if c == nil || c.FileName == nil {
		return ""
	}
	return *c.FileName
}

// SetFileName sets the FileName field.
//
// When the provided CodeCoverageFile type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageFile) SetFileName(v string) {
	if c == nil {
		return
	}
	c.FileName = &v
}

// GetMissedLines returns the MissedLines field.
//
// When the provided CodeCoverageFile type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageFile) GetMissedLines() int {
	if c == nil || c.MissedLines == nil {
		return 0
	}
	return *c.MissedLines
}

// SetMissedLines sets the MissedLines field.
//
// When the provided CodeCoverageFile type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageFile) SetMissedLines(v int) {
	if c == nil {
		return
	}
	c.MissedLines = &v
}

// GetPartialLines returns the PartialLines field.
//
// When the provided CodeCoverageFile type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageFile) GetPartialLines() int {
	if c == nil || c.PartialLines == nil {
		return 0
	}
	return *c.PartialLines
}

// SetPartialLines sets the PartialLines field.
//
// When the provided CodeCoverageFile type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageFile) SetPartialLines(v int) {
	if c == nil {
		return
	}
	c.PartialLines = &v
}

// GetFilePath returns the FilePath field.
//
// When the provided CodeCoverageFile type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageFile) GetFilePath() string {
	if c == nil || c.FilePath == nil {
		return ""
	}
	return *c.FilePath
}

// SetFilePath sets the FilePath field.
//
// When the provided CodeCoverageFile type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageFile) SetFilePath(v string) {
	if c == nil {
		return
	}
	c.FilePath = &v
}
