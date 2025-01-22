package test_report

// ResultsProcessingFailure is the API types representation of a failure processing test results.
//
// swagger:model ResultsProcessingFailure
type ResultsProcessingFailure struct {
	ID             *int64  `json:"id,omitempty"`
	Body           *string `json:"body,omitempty"`
	Created        *int64  `json:"created,omitempty"`
	FailureMessage *string `json:"failure_message,omitempty"`
	FailureType    *string `json:"failure_type,omitempty"`
	BodyType       *string `json:"body_type,omitempty"`
}

// SetID sets the ID field.
//
// When the provided ResultsProcessingFailure type is nil, it
// will set nothing and immediately return.
func (r *ResultsProcessingFailure) SetID(v int64) {
	// return if ResultsProcessingFailure type is nil
	if r == nil {
		return
	}

	r.ID = &v
}

// GetID returns the ID field.
//
// When the provided ResultsProcessingFailure type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *ResultsProcessingFailure) GetID() int64 {
	// return zero value if ResultsProcessingFailure type or ID field is nil
	if r == nil || r.ID == nil {
		return 0
	}

	return *r.ID
}

// SetBody sets the Body field.
//
// When the provided ResultsProcessingFailure type is nil, it
// will set nothing and immediately return.
func (r *ResultsProcessingFailure) SetBody(v string) {
	// return if ResultsProcessingFailure type is nil
	if r == nil {
		return
	}

	r.Body = &v
}

// GetBody returns the Body field.
//
// When the provided ResultsProcessingFailure type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *ResultsProcessingFailure) GetBody() string {
	// return zero value if ResultsProcessingFailure type or Body field is nil
	if r == nil || r.Body == nil {
		return ""
	}

	return *r.Body
}

// SetCreated sets the Created field.
//
// When the provided ResultsProcessingFailure type is nil, it
// will set nothing and immediately return.
func (r *ResultsProcessingFailure) SetCreated(v int64) {
	// return if ResultsProcessingFailure type is nil
	if r == nil {
		return
	}

	r.Created = &v
}

// GetCreated returns the Created field.
//
// When the provided ResultsProcessingFailure type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *ResultsProcessingFailure) GetCreated() int64 {
	// return zero value if ResultsProcessingFailure type or Created field is nil
	if r == nil || r.Created == nil {
		return 0
	}

	return *r.Created
}

// SetFailureMessage sets the FailureMessage field.
//
// When the provided ResultsProcessingFailure type is nil, it
// will set nothing and immediately return.
func (r *ResultsProcessingFailure) SetFailureMessage(v string) {
	// return if ResultsProcessingFailure type is nil
	if r == nil {
		return
	}

	r.FailureMessage = &v
}

// GetFailureMessage returns the FailureMessage field.
//
// When the provided ResultsProcessingFailure type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *ResultsProcessingFailure) GetFailureMessage() string {
	// return zero value if ResultsProcessingFailure type or FailureMessage field is nil
	if r == nil || r.FailureMessage == nil {
		return ""
	}

	return *r.FailureMessage
}

// SetFailureType sets the FailureType field.
//
// When the provided ResultsProcessingFailure type is nil, it
// will set nothing and immediately return.
func (r *ResultsProcessingFailure) SetFailureType(v string) {
	// return if ResultsProcessingFailure type is nil
	if r == nil {
		return
	}

	r.FailureType = &v
}

// GetFailureType returns the FailureType field.
//
// When the provided ResultsProcessingFailure type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *ResultsProcessingFailure) GetFailureType() string {
	// return zero value if ResultsProcessingFailure type or FailureType field is nil
	if r == nil || r.FailureType == nil {
		return ""
	}

	return *r.FailureType
}

// SetBodyType sets the BodyType field.
//
// When the provided ResultsProcessingFailure type is nil, it
// will set nothing and immediately return.
func (r *ResultsProcessingFailure) SetBodyType(v string) {
	// return if ResultsProcessingFailure type is nil
	if r == nil {
		return
	}

	r.BodyType = &v
}

// GetBodyType returns the BodyType field.
//
// When the provided ResultsProcessingFailure type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *ResultsProcessingFailure) GetBodyType() string {
	// return zero value if ResultsProcessingFailure type or BodyType field is nil
	if r == nil || r.BodyType == nil {
		return ""
	}

	return *r.BodyType
}
