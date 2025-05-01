// SPDX-License-Identifier: Apache-2.0

package types

import "fmt"

// TestReport is the API representation of a test report for a pipeline.
//
// swagger:model TestReport
type TestReport struct {
	ID        *int64 `json:"id,omitempty"`
	BuildID   *int64 `json:"build_id,omitempty"`
	CreatedAt *int64 `json:"created_at,omitempty"`
}

type TestAttachment struct {
	ID           *int64  `json:"id,omitempty"`
	TestReportID *int64  `json:"test_report_id,omitempty"`
	FileName     *string `json:"file_name,omitempty"`
	ObjectPath   *string `json:"object_path,omitempty"`
	FileSize     *int64  `json:"file_size,omitempty"`
	FileType     *string `json:"file_type,omitempty"`
	PresignedUrl *string `json:"presigned_url,omitempty"`
	CreatedAt    *int64  `json:"created_at,omitempty"`
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

// GetCreatedAt returns the CreatedAt field.
//
// When the provided TestReport type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (t *TestReport) GetCreatedAt() int64 {
	// return zero value if TestReport type or CreatedAt field is nil
	if t == nil || t.CreatedAt == nil {
		return 0
	}

	return *t.CreatedAt
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

// SetCreatedAt sets the CreatedAt field.
func (t *TestReport) SetCreatedAt(v int64) {
	// return if TestReport type is nil
	if t == nil {
		return
	}
	// set the CreatedAt field
	t.CreatedAt = &v
}

// GetID returns the ID field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestAttachment) GetID() int64 {
	// return zero value if TestAttachment type or ID field is nil
	if r == nil || r.ID == nil {
		return 0
	}

	return *r.ID
}

// GetTestReportID returns the TestReportID field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestAttachment) GetTestReportID() int64 {
	// return zero value if TestAttachment type or TestReportID field is nil
	if r == nil || r.TestReportID == nil {
		return 0
	}

	return *r.TestReportID
}

// GetFileName returns the FileName field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestAttachment) GetFileName() string {
	// return zero value if TestAttachment type or FileName field is nil
	if r == nil || r.FileName == nil {
		return ""
	}

	return *r.FileName
}

// GetObjectPath returns the FilePath field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestAttachment) GetObjectPath() string {
	// return zero value if TestAttachment type or FilePath field is nil
	if r == nil || r.ObjectPath == nil {
		return ""
	}

	return *r.ObjectPath
}

// GetFileSize returns the FileSize field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestAttachment) GetFileSize() int64 {
	// return zero value if TestAttachment type or FileSize field is nil
	if r == nil || r.FileSize == nil {
		return 0
	}

	return *r.FileSize
}

// GetFileType returns the FileType field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestAttachment) GetFileType() string {
	// return zero value if TestAttachment type or FileType field is nil
	if r == nil || r.FileType == nil {
		return ""
	}

	return *r.FileType
}

// GetPresignedUrl returns the PresignedUrl field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestAttachment) GetPresignedUrl() string {
	// return zero value if TestAttachment type or PresignedUrl field is nil
	if r == nil || r.PresignedUrl == nil {
		return ""
	}

	return *r.PresignedUrl
}

// GetCreatedAt returns the CreatedAt field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (r *TestAttachment) GetCreatedAt() int64 {
	// return zero value if TestAttachment type or CreatedAt field is nil
	if r == nil || r.CreatedAt == nil {
		return 0
	}

	return *r.CreatedAt
}

// SetID sets the ID field.
func (r *TestAttachment) SetID(v int64) {
	// return if TestAttachment type is nil
	if r == nil {
		return
	}
	// set the ID field
	r.ID = &v
}

// SetTestReportID sets the TestReportID field.
func (r *TestAttachment) SetTestReportID(v int64) {
	// return if TestAttachment type is nil
	if r == nil {
		return
	}
	// set the TestReportID field
	r.TestReportID = &v
}

// SetFileName sets the FileName field.
func (r *TestAttachment) SetFileName(v string) {
	// return if TestAttachment type is nil
	if r == nil {
		return
	}
	// set the FileName field
	r.FileName = &v
}

// SetObjectPath sets the ObjectPath field.
func (r *TestAttachment) SetObjectPath(v string) {
	// return if TestAttachment type is nil
	if r == nil {
		return
	}
	// set the ObjectPath field
	r.ObjectPath = &v
}

// SetFileSize sets the FileSize field.
func (r *TestAttachment) SetFileSize(v int64) {
	// return if TestAttachment type is nil
	if r == nil {
		return
	}
	// set the FileSize field
	r.FileSize = &v
}

// SetFileType sets the FileType field.
func (r *TestAttachment) SetFileType(v string) {
	// return if TestAttachment type is nil
	if r == nil {
		return
	}
	// set the FileType field
	r.FileType = &v
}

// SetPresignedUrl sets the PresignedUrl field.
func (r *TestAttachment) SetPresignedUrl(v string) {
	// return if TestAttachment type is nil
	if r == nil {
		return
	}
	// set the PresignedUrl field
	r.PresignedUrl = &v
}

// SetCreatedAt sets the CreatedAt field.
func (r *TestAttachment) SetCreatedAt(v int64) {
	// return if TestAttachment type is nil
	if r == nil {
		return
	}
	// set the CreatedAt field
	r.CreatedAt = &v
}

// String implements the Stringer interface for the TestAttachment type.
func (r *TestAttachment) String() string {
	return fmt.Sprintf("FileName: %s, ObjectPath: %s, FileType: %s", r.GetFileName(), r.GetObjectPath(), r.GetFileType())
}
