// SPDX-License-Identifier: Apache-2.0

package types

import "fmt"

// Artifact is the API representation of an artifact.
//
// swagger:model Artifact
type Artifact struct {
	ID           *int64  `json:"id,omitempty"`
	BuildID      *int64  `json:"build_id,omitempty"`
	FileName     *string `json:"file_name,omitempty"`
	ObjectPath   *string `json:"object_path,omitempty"`
	FileSize     *int64  `json:"file_size,omitempty"`
	FileType     *string `json:"file_type,omitempty"`
	PresignedURL *string `json:"presigned_url,omitempty"`
	CreatedAt    *int64  `json:"created_at,omitempty"`
}

// GetID returns the ID field.
//
// When the provided Artifact type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Artifact) GetID() int64 {
	// return zero value if Artifact type or ID field is nil
	if a == nil || a.ID == nil {
		return 0
	}

	return *a.ID
}

// GetBuildID returns the BuildID field.
//
// When the provided Artifact type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Artifact) GetBuildID() int64 {
	// return zero value if Artifact type or BuildID field is nil
	if a == nil || a.BuildID == nil {
		return 0
	}

	return *a.BuildID
}

// GetFileName returns the FileName field.
//
// When the provided Artifact type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Artifact) GetFileName() string {
	// return zero value if Artifact type or FileName field is nil
	if a == nil || a.FileName == nil {
		return ""
	}

	return *a.FileName
}

// GetObjectPath returns the FilePath field.
//
// When the provided Artifact type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Artifact) GetObjectPath() string {
	// return zero value if Artifact type or FilePath field is nil
	if a == nil || a.ObjectPath == nil {
		return ""
	}

	return *a.ObjectPath
}

// GetFileSize returns the FileSize field.
//
// When the provided Artifact type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Artifact) GetFileSize() int64 {
	// return zero value if Artifact type or FileSize field is nil
	if a == nil || a.FileSize == nil {
		return 0
	}

	return *a.FileSize
}

// GetFileType returns the FileType field.
//
// When the provided Artifact type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Artifact) GetFileType() string {
	// return zero value if Artifact type or FileType field is nil
	if a == nil || a.FileType == nil {
		return ""
	}

	return *a.FileType
}

// GetPresignedURL returns the PresignedUrl field.
//
// When the provided Artifact type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Artifact) GetPresignedURL() string {
	// return zero value if Artifact type or PresignedUrl field is nil
	if a == nil || a.PresignedURL == nil {
		return ""
	}

	return *a.PresignedURL
}

// GetCreatedAt returns the CreatedAt field.
//
// When the provided Artifact type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (a *Artifact) GetCreatedAt() int64 {
	// return zero value if Artifact type or CreatedAt field is nil
	if a == nil || a.CreatedAt == nil {
		return 0
	}

	return *a.CreatedAt
}

// SetID sets the ID field.
func (a *Artifact) SetID(v int64) {
	// return if Artifact type is nil
	if a == nil {
		return
	}
	// set the ID field
	a.ID = &v
}

// SetBuildID sets the BuildID field.
func (a *Artifact) SetBuildID(v int64) {
	// return if Artifact type is nil
	if a == nil {
		return
	}
	// set the BuildID field
	a.BuildID = &v
}

// SetFileName sets the FileName field.
func (a *Artifact) SetFileName(v string) {
	// return if Artifact type is nil
	if a == nil {
		return
	}
	// set the FileName field
	a.FileName = &v
}

// SetObjectPath sets the ObjectPath field.
func (a *Artifact) SetObjectPath(v string) {
	// return if Artifact type is nil
	if a == nil {
		return
	}
	// set the ObjectPath field
	a.ObjectPath = &v
}

// SetFileSize sets the FileSize field.
func (a *Artifact) SetFileSize(v int64) {
	// return if Artifact type is nil
	if a == nil {
		return
	}
	// set the FileSize field
	a.FileSize = &v
}

// SetFileType sets the FileType field.
func (a *Artifact) SetFileType(v string) {
	// return if Artifact type is nil
	if a == nil {
		return
	}
	// set the FileType field
	a.FileType = &v
}

// SetPresignedURL sets the PresignedUrl field.
func (a *Artifact) SetPresignedURL(v string) {
	// return if Artifact type is nil
	if a == nil {
		return
	}
	// set the PresignedUrl field
	a.PresignedURL = &v
}

// SetCreatedAt sets the CreatedAt field.
func (a *Artifact) SetCreatedAt(v int64) {
	// return if Artifact type is nil
	if a == nil {
		return
	}
	// set the CreatedAt field
	a.CreatedAt = &v
}

// String implements the Stringer interface for the Artifact type.
func (a *Artifact) String() string {
	return fmt.Sprintf("FileName: %s, ObjectPath: %s, FileType: %s", a.GetFileName(), a.GetObjectPath(), a.GetFileType())
}
