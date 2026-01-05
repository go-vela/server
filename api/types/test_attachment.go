// SPDX-License-Identifier: Apache-2.0

package types

import "fmt"

// TestAttachment is the API representation of a test attachment.
//
// swagger:model TestAttachment
type TestAttachment struct {
	ID           *int64  `json:"id,omitempty"`
	TestReportID *int64  `json:"test_report_id,omitempty"`
	FileName     *string `json:"file_name,omitempty"`
	ObjectPath   *string `json:"object_path,omitempty"`
	FileSize     *int64  `json:"file_size,omitempty"`
	FileType     *string `json:"file_type,omitempty"`
	PresignedURL *string `json:"presigned_url,omitempty"`
	CreatedAt    *int64  `json:"created_at,omitempty"`
}

// GetID returns the ID field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ta *TestAttachment) GetID() int64 {
	// return zero value if TestAttachment type or ID field is nil
	if ta == nil || ta.ID == nil {
		return 0
	}

	return *ta.ID
}

// GetTestReportID returns the TestReportID field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ta *TestAttachment) GetTestReportID() int64 {
	// return zero value if TestAttachment type or TestReportID field is nil
	if ta == nil || ta.TestReportID == nil {
		return 0
	}

	return *ta.TestReportID
}

// GetFileName returns the FileName field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ta *TestAttachment) GetFileName() string {
	// return zero value if TestAttachment type or FileName field is nil
	if ta == nil || ta.FileName == nil {
		return ""
	}

	return *ta.FileName
}

// GetObjectPath returns the FilePath field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ta *TestAttachment) GetObjectPath() string {
	// return zero value if TestAttachment type or FilePath field is nil
	if ta == nil || ta.ObjectPath == nil {
		return ""
	}

	return *ta.ObjectPath
}

// GetFileSize returns the FileSize field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ta *TestAttachment) GetFileSize() int64 {
	// return zero value if TestAttachment type or FileSize field is nil
	if ta == nil || ta.FileSize == nil {
		return 0
	}

	return *ta.FileSize
}

// GetFileType returns the FileType field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ta *TestAttachment) GetFileType() string {
	// return zero value if TestAttachment type or FileType field is nil
	if ta == nil || ta.FileType == nil {
		return ""
	}

	return *ta.FileType
}

// GetPresignedURL returns the PresignedUrl field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ta *TestAttachment) GetPresignedURL() string {
	// return zero value if TestAttachment type or PresignedUrl field is nil
	if ta == nil || ta.PresignedURL == nil {
		return ""
	}

	return *ta.PresignedURL
}

// GetCreatedAt returns the CreatedAt field.
//
// When the provided TestAttachment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ta *TestAttachment) GetCreatedAt() int64 {
	// return zero value if TestAttachment type or CreatedAt field is nil
	if ta == nil || ta.CreatedAt == nil {
		return 0
	}

	return *ta.CreatedAt
}

// SetID sets the ID field.
func (ta *TestAttachment) SetID(v int64) {
	// return if TestAttachment type is nil
	if ta == nil {
		return
	}
	// set the ID field
	ta.ID = &v
}

// SetTestReportID sets the TestReportID field.
func (ta *TestAttachment) SetTestReportID(v int64) {
	// return if TestAttachment type is nil
	if ta == nil {
		return
	}
	// set the TestReportID field
	ta.TestReportID = &v
}

// SetFileName sets the FileName field.
func (ta *TestAttachment) SetFileName(v string) {
	// return if TestAttachment type is nil
	if ta == nil {
		return
	}
	// set the FileName field
	ta.FileName = &v
}

// SetObjectPath sets the ObjectPath field.
func (ta *TestAttachment) SetObjectPath(v string) {
	// return if TestAttachment type is nil
	if ta == nil {
		return
	}
	// set the ObjectPath field
	ta.ObjectPath = &v
}

// SetFileSize sets the FileSize field.
func (ta *TestAttachment) SetFileSize(v int64) {
	// return if TestAttachment type is nil
	if ta == nil {
		return
	}
	// set the FileSize field
	ta.FileSize = &v
}

// SetFileType sets the FileType field.
func (ta *TestAttachment) SetFileType(v string) {
	// return if TestAttachment type is nil
	if ta == nil {
		return
	}
	// set the FileType field
	ta.FileType = &v
}

// SetPresignedURL sets the PresignedUrl field.
func (ta *TestAttachment) SetPresignedURL(v string) {
	// return if TestAttachment type is nil
	if ta == nil {
		return
	}
	// set the PresignedUrl field
	ta.PresignedURL = &v
}

// SetCreatedAt sets the CreatedAt field.
func (ta *TestAttachment) SetCreatedAt(v int64) {
	// return if TestAttachment type is nil
	if ta == nil {
		return
	}
	// set the CreatedAt field
	ta.CreatedAt = &v
}

// String implements the Stringer interface for the TestAttachment type.
func (ta *TestAttachment) String() string {
	return fmt.Sprintf("FileName: %s, ObjectPath: %s, FileType: %s", ta.GetFileName(), ta.GetObjectPath(), ta.GetFileType())
}
