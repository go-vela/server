package test_report

// TestRunAttachment is the API types representation of an attachment for a test run.
//
// swagger:model TestRunAttachment
type TestRunAttachment struct {
	ID         *int64  `json:"id,omitempty"`
	FileName   *string `json:"file_name,omitempty"`
	ObjectName *string `json:"object_name,omitempty"`
	FileSize   *int64  `json:"file_size,omitempty"`
}
