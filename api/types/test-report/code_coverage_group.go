package test_report

// CodeCoverageGroup is the API types representation of a group of code coverage.
//
// swagger:model CodeCoverageGroup
type CodeCoverageGroup struct {
	ID                *int64             `json:"id,omitempty"`
	CodeCoverageRun   *CodeCoverageRun   `json:"code_coverage_run,omitempty"`
	Name              *string            `json:"name,omitempty"`
	CodeCoverageStats *CodeCoverageStats `json:"code_coverage_stats,omitempty"`
}

// GetID returns the ID field.
//
// When the provided CodeCoverageGroup type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageGroup) GetID() int64 {
	if c == nil || c.ID == nil {
		return 0
	}
	return *c.ID
}

// SetID sets the ID field.
//
// When the provided CodeCoverageGroup type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageGroup) SetID(v int64) {
	if c == nil {
		return
	}
	c.ID = &v
}

// GetCodeCoverageRun returns the CodeCoverageRun field.
//
// When the provided CodeCoverageGroup type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageGroup) GetCodeCoverageRun() *CodeCoverageRun {
	if c == nil || c.CodeCoverageRun == nil {
		return new(CodeCoverageRun)
	}
	return c.CodeCoverageRun
}

// SetCodeCoverageRun sets the CodeCoverageRun field.
//
// When the provided CodeCoverageGroup type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageGroup) SetCodeCoverageRun(v *CodeCoverageRun) {
	if c == nil {
		return
	}
	c.CodeCoverageRun = v
}

// GetName returns the Name field.
//
// When the provided CodeCoverageGroup type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageGroup) GetName() string {
	if c == nil || c.Name == nil {
		return ""
	}
	return *c.Name
}

// SetName sets the Name field.
//
// When the provided CodeCoverageGroup type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageGroup) SetName(v string) {
	if c == nil {
		return
	}
	c.Name = &v
}

// GetCodeCoverageStats returns the CodeCoverageStats field.
//
// When the provided CodeCoverageGroup type is nil, or the field within
// the type is nil, it returns a new CodeCoverageStats instance.
func (c *CodeCoverageGroup) GetCodeCoverageStats() *CodeCoverageStats {
	if c == nil || c.CodeCoverageStats == nil {
		return new(CodeCoverageStats)
	}
	return c.CodeCoverageStats
}

// SetCodeCoverageStats sets the CodeCoverageStats field.
//
// When the provided CodeCoverageGroup type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageGroup) SetCodeCoverageStats(v *CodeCoverageStats) {
	if c == nil {
		return
	}
	c.CodeCoverageStats = v
}
