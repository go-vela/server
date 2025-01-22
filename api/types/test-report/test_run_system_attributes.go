package test_report

// TestRunSystemAttributes is the API types representation of system attributes for a test run.
//
// swagger:model TestRunSystemAttributes
type TestRunSystemAttributes struct {
	ID     *int64 `json:"id,omitempty"`
	Pinned *bool  `json:"pinned,omitempty"`
}

// GetID returns the ID field.
//
// When the provided TestRunSystemAttributes type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *TestRunSystemAttributes) GetID() int64 {
	if b == nil || b.ID == nil {
		return 0
	}
	return *b.ID
}

// SetID sets the ID field.
func (b *TestRunSystemAttributes) SetID(v int64) {
	if b == nil {
		return
	}
	b.ID = &v
}

// GetPinned returns the Pinned field.
//
// When the provided TestRunSystemAttributes type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (b *TestRunSystemAttributes) GetPinned() bool {
	if b == nil || b.Pinned == nil {
		return false
	}
	return *b.Pinned
}

// SetPinned sets the Pinned field.
func (b *TestRunSystemAttributes) SetPinned(v bool) {
	if b == nil {
		return
	}
	b.Pinned = &v
}
