package types

import "fmt"

// TestReport is the API representation of a test report for a pipeline.
//
// swagger:model TestReport
type TestReport struct {
	ID      *int64 `json:"id,omitempty"`
	BuildID *int64 `json:"build_id,omitempty"`
	Created *int64 `json:"created,omitempty"`
}

type TestReportAttachments struct {
	ID           *int64  `json:"id,omitempty"`
	TestReportID *int64  `json:"test_report_id,omitempty"`
	Filename     *string `json:"filename,omitempty"`
	FilePath     *string `json:"file_path,omitempty"`
	FileSize     *int64  `json:"file_size,omitempty"`
	FileType     *string `json:"file_type,omitempty"`
	PresignedUrl *string `json:"presigned_url,omitempty"`
	Created      *int64  `json:"created,omitempty"`
}

// GetID returns the ID field.
//
// When the provided TestReport type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestReport) GetID() int64 {
	// return zero value if TestReport type or ID field is nil
	if t == nil || t.ID == nil {
		return 0
	}

	return *t.ID
}

// GetBuildID returns the BuildID field.
//
// When the provided TestReport type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestReport) GetBuildID() int64 {
	// return zero value if TestReport type or BuildID field is nil
	if t == nil || t.BuildID == nil {
		return 0
	}

	return *t.BuildID
}

// GetCreated returns the Created field.
//
// When the provided TestReport type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestReport) GetCreated() int64 {
	// return zero value if TestReport type or Created field is nil
	if t == nil || t.Created == nil {
		return 0
	}

	return *t.Created
}

// SetID sets the ID field.
func (t *TestReport) SetID(v int64) {
	// return if TestReport type is nil
	if t == nil {
		return
	}
	// set the ID field
	t.ID = &v
}

// SetBuildID sets the BuildID field.
func (t *TestReport) SetBuildID(v int64) {
	// return if TestReport type is nil
	if t == nil {
		return
	}
	// set the BuildID field
	t.BuildID = &v
}

// SetCreated sets the Created field.
func (t *TestReport) SetCreated(v int64) {
	// return if TestReport type is nil
	if t == nil {
		return
	}
	// set the Created field
	t.Created = &v
}

// GetID returns the ID field.
//
// When the provided TestReportAttachments type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestReportAttachments) GetID() int64 {
	// return zero value if TestReportAttachments type or ID field is nil
	if r == nil || r.ID == nil {
		return 0
	}

	return *r.ID
}

// GetTestReportID returns the TestReportID field.
//
// When the provided TestReportAttachments type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestReportAttachments) GetTestReportID() int64 {
	// return zero value if TestReportAttachments type or TestReportID field is nil
	if r == nil || r.TestReportID == nil {
		return 0
	}

	return *r.TestReportID
}

// GetFilename returns the Filename field.
//
// When the provided TestReportAttachments type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestReportAttachments) GetFilename() string {
	// return zero value if TestReportAttachments type or Filename field is nil
	if r == nil || r.Filename == nil {
		return ""
	}

	return *r.Filename
}

// GetFilePath returns the FilePath field.
//
// When the provided TestReportAttachments type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestReportAttachments) GetFilePath() string {
	// return zero value if TestReportAttachments type or FilePath field is nil
	if r == nil || r.FilePath == nil {
		return ""
	}

	return *r.FilePath
}

// GetFileSize returns the FileSize field.
//
// When the provided TestReportAttachments type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestReportAttachments) GetFileSize() int64 {
	// return zero value if TestReportAttachments type or FileSize field is nil
	if r == nil || r.FileSize == nil {
		return 0
	}

	return *r.FileSize
}

// GetFileType returns the FileType field.
//
// When the provided TestReportAttachments type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestReportAttachments) GetFileType() string {
	// return zero value if TestReportAttachments type or FileType field is nil
	if r == nil || r.FileType == nil {
		return ""
	}

	return *r.FileType
}

// GetPresignedUrl returns the PresignedUrl field.
//
// When the provided TestReportAttachments type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestReportAttachments) GetPresignedUrl() string {
	// return zero value if TestReportAttachments type or PresignedUrl field is nil
	if r == nil || r.PresignedUrl == nil {
		return ""
	}

	return *r.PresignedUrl
}

// GetCreated returns the Created field.
//
// When the provided TestReportAttachments type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestReportAttachments) GetCreated() int64 {
	// return zero value if TestReportAttachments type or Created field is nil
	if r == nil || r.Created == nil {
		return 0
	}

	return *r.Created
}

// SetID sets the ID field.
func (r *TestReportAttachments) SetID(v int64) {
	// return if TestReportAttachments type is nil
	if r == nil {
		return
	}
	// set the ID field
	r.ID = &v
}

// SetTestReportID sets the TestReportID field.
func (r *TestReportAttachments) SetTestReportID(v int64) {
	// return if TestReportAttachments type is nil
	if r == nil {
		return
	}
	// set the TestReportID field
	r.TestReportID = &v
}

// SetFilename sets the Filename field.
func (r *TestReportAttachments) SetFilename(v string) {
	// return if TestReportAttachments type is nil
	if r == nil {
		return
	}
	// set the Filename field
	r.Filename = &v
}

// SetFilePath sets the FilePath field.
func (r *TestReportAttachments) SetFilePath(v string) {
	// return if TestReportAttachments type is nil
	if r == nil {
		return
	}
	// set the FilePath field
	r.FilePath = &v
}

// SetFileSize sets the FileSize field.
func (r *TestReportAttachments) SetFileSize(v int64) {
	// return if TestReportAttachments type is nil
	if r == nil {
		return
	}
	// set the FileSize field
	r.FileSize = &v
}

// SetFileType sets the FileType field.
func (r *TestReportAttachments) SetFileType(v string) {
	// return if TestReportAttachments type is nil
	if r == nil {
		return
	}
	// set the FileType field
	r.FileType = &v
}

// SetPresignedUrl sets the PresignedUrl field.
func (r *TestReportAttachments) SetPresignedUrl(v string) {
	// return if TestReportAttachments type is nil
	if r == nil {
		return
	}
	// set the PresignedUrl field
	r.PresignedUrl = &v
}

// SetCreated sets the Created field.
func (r *TestReportAttachments) SetCreated(v int64) {
	// return if TestReportAttachments type is nil
	if r == nil {
		return
	}
	// set the Created field
	r.Created = &v
}

// String implements the Stringer interface for the TestReportAttachments type.
func (r *TestReportAttachments) String() string {
	return fmt.Sprintf("Filename: %s, FilePath: %s, FileType: %s", r.GetFilename(), r.GetFilePath(), r.GetFileType())
}
