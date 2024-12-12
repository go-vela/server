package test_report

// CodeCoverageStats is the API types representation of code coverage stats for a pipeline.
//
// swagger:model CodeCoverageStats
type CodeCoverageStats struct {
	ID               *int64           `json:"id,omitempty"`
	CodeCoverageRun  *CodeCoverageRun `json:"code_coverage_run,omitempty"`
	Scope            *string          `json:"scope,omitempty"`
	StatementCovered *int             `json:"statement_covered,omitempty"`
	StatementMissed  *int             `json:"statement_missed,omitempty"`
	LineCovered      *int             `json:"line_covered,omitempty"`
	LineMissed       *int             `json:"line_missed,omitempty"`
	BranchCovered    *int             `json:"branch_covered,omitempty"`
	BranchMissed     *int             `json:"branch_missed,omitempty"`
}

// GetID returns the ID field.
//
// When the provided CodeCoverageStats type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageStats) GetID() int64 {
	if c == nil || c.ID == nil {
		return 0
	}
	return *c.ID
}

// SetID sets the ID field.
//
// When the provided CodeCoverageStats type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageStats) SetID(v int64) {
	if c == nil {
		return
	}
	c.ID = &v
}

// GetCodeCoverageRun returns the CodeCoverageRun field.
//
// When the provided CodeCoverageStats type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageStats) GetCodeCoverageRun() *CodeCoverageRun {
	if c == nil || c.CodeCoverageRun == nil {
		return new(CodeCoverageRun)
	}
	return c.CodeCoverageRun
}

// SetCodeCoverageRun sets the CodeCoverageRun field.
//
// When the provided CodeCoverageStats type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageStats) SetCodeCoverageRun(v *CodeCoverageRun) {
	if c == nil {
		return
	}
	c.CodeCoverageRun = v
}

// GetScope returns the Scope field.
//
// When the provided CodeCoverageStats type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageStats) GetScope() string {
	if c == nil || c.Scope == nil {
		return ""
	}
	return *c.Scope
}

// SetScope sets the Scope field.
//
// When the provided CodeCoverageStats type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageStats) SetScope(v string) {
	if c == nil {
		return
	}
	c.Scope = &v
}

// GetStatementCovered returns the StatementCovered field.
//
// When the provided CodeCoverageStats type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageStats) GetStatementCovered() int {
	if c == nil || c.StatementCovered == nil {
		return 0
	}
	return *c.StatementCovered
}

// SetStatementCovered sets the StatementCovered field.
//
// When the provided CodeCoverageStats type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageStats) SetStatementCovered(v int) {
	if c == nil {
		return
	}
	c.StatementCovered = &v
}

// GetStatementMissed returns the StatementMissed field.
//
// When the provided CodeCoverageStats type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageStats) GetStatementMissed() int {
	if c == nil || c.StatementMissed == nil {
		return 0
	}
	return *c.StatementMissed
}

// SetStatementMissed sets the StatementMissed field.
//
// When the provided CodeCoverageStats type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageStats) SetStatementMissed(v int) {
	if c == nil {
		return
	}
	c.StatementMissed = &v
}

// GetLineCovered returns the LineCovered field.
//
// When the provided CodeCoverageStats type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageStats) GetLineCovered() int {
	if c == nil || c.LineCovered == nil {
		return 0
	}
	return *c.LineCovered
}

// SetLineCovered sets the LineCovered field.
//
// When the provided CodeCoverageStats type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageStats) SetLineCovered(v int) {
	if c == nil {
		return
	}
	c.LineCovered = &v
}

// GetLineMissed returns the LineMissed field.
//
// When the provided CodeCoverageStats type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageStats) GetLineMissed() int {
	if c == nil || c.LineMissed == nil {
		return 0
	}
	return *c.LineMissed
}

// SetLineMissed sets the LineMissed field.
//
// When the provided CodeCoverageStats type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageStats) SetLineMissed(v int) {
	if c == nil {
		return
	}
	c.LineMissed = &v
}

// GetBranchCovered returns the BranchCovered field.
//
// When the provided CodeCoverageStats type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageStats) GetBranchCovered() int {
	if c == nil || c.BranchCovered == nil {
		return 0
	}
	return *c.BranchCovered
}

// SetBranchCovered sets the BranchCovered field.
//
// When the provided CodeCoverageStats type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageStats) SetBranchCovered(v int) {
	if c == nil {
		return
	}
	c.BranchCovered = &v
}

// GetBranchMissed returns the BranchMissed field.
//
// When the provided CodeCoverageStats type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (c *CodeCoverageStats) GetBranchMissed() int {
	if c == nil || c.BranchMissed == nil {
		return 0
	}
	return *c.BranchMissed
}

// SetBranchMissed sets the BranchMissed field.
//
// When the provided CodeCoverageStats type is nil, it
// will set nothing and immediately return.
func (c *CodeCoverageStats) SetBranchMissed(v int) {
	if c == nil {
		return
	}
	c.BranchMissed = &v
}
