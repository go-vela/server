package test_report

// ResultsProcessing is the API types representation of processing results for a test run.
//
// swagger:model ResultsProcessing
type ResultsProcessing struct {
	ID           *int64  `json:"id,omitempty"`
	Status       *string `json:"status,omitempty"`
	ErrorMessage *string `json:"error_message,omitempty"`
	Created      *int64  `json:"created,omitempty"`
	//Created replaced created_timestamp in Projektor model
	// Following the same pattern as other fields in Vela project
	// PublicID is replaced by ID for consistency across all models
}

// SetID sets the PublicID field.
//
// When the provided ResultsProcessing type is nil, it
// will set nothing and immediately return.
func (r *ResultsProcessing) SetID(v int64) {
	// return if ResultsProcessing type is nil
	if r == nil {
		return
	}

	r.ID = &v
}

// GetPublicID returns the PublicID field.
//
// When the provided ResultsProcessing type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *ResultsProcessing) GetPublicID() int64 {
	// return zero value if ResultsProcessing type or PublicID field is nil
	if r == nil || r.ID == nil {
		return 0
	}

	return *r.ID
}

// SetStatus sets the Status field.
//
// When the provided ResultsProcessing type is nil, it
// will set nothing and immediately return.
func (r *ResultsProcessing) SetStatus(v string) {
	// return if ResultsProcessing type is nil
	if r == nil {
		return
	}

	r.Status = &v
}

// GetStatus returns the Status field.
//
// When the provided ResultsProcessing type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *ResultsProcessing) GetStatus() string {
	// return zero value if ResultsProcessing type or Status field is nil
	if r == nil || r.Status == nil {
		return ""
	}

	return *r.Status
}

// SetErrorMessage sets the ErrorMessage field.
//
// When the provided ResultsProcessing type is nil, it
// will set nothing and immediately return.
func (r *ResultsProcessing) SetErrorMessage(v string) {
	// return if ResultsProcessing type is nil
	if r == nil {
		return
	}

	r.ErrorMessage = &v
}

// GetErrorMessage returns the ErrorMessage field.
//
// When the provided ResultsProcessing type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *ResultsProcessing) GetErrorMessage() string {
	// return zero value if ResultsProcessing type or ErrorMessage field is nil
	if r == nil || r.ErrorMessage == nil {
		return ""
	}

	return *r.ErrorMessage
}

// SetCreated sets the Created field.
//
// When the provided ResultsProcessing type is nil, it
// will set nothing and immediately return.
func (r *ResultsProcessing) SetCreated(v int64) {
	// return if ResultsProcessing type is nil
	if r == nil {
		return
	}

	r.Created = &v
}

// GetCreated returns the Created field.
//
// When the provided ResultsProcessing type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *ResultsProcessing) GetCreated() int64 {
	// return zero value if ResultsProcessing type or Created field is nil
	if r == nil || r.Created == nil {
		return 0
	}

	return *r.Created
}
