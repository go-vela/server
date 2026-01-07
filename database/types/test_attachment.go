// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"errors"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyTestReportID defines the error type when a
	// TestAttachment type has an empty TestReportID field provided.
	ErrEmptyTestReportID = errors.New("empty test_report_id provided")
	// ErrEmptyFileName defines the error type when a
	// TestAttachment type has an empty FileName field provided.
	ErrEmptyFileName = errors.New("empty file_name provided")
	// ErrEmptyObjectPath defines the error type when a
	// TestAttachment type has an empty ObjectPath field provided.
	ErrEmptyObjectPath = errors.New("empty object_path provided")
	// ErrEmptyFileSize defines the error type when a
	// TestAttachment type has an empty FileSize field provided.
	ErrEmptyFileSize = errors.New("empty file_size provided")
	// ErrEmptyFileType defines the error type when a
	// TestAttachment type has an empty FileType field provided.
	ErrEmptyFileType = errors.New("empty file_type provided")
	// ErrEmptyPresignedURL defines the error type when a
	// TestAttachment type has an empty PresignedUrl field provided.
	ErrEmptyPresignedURL = errors.New("empty presigned_url provided")
)

type TestAttachment struct {
	ID           sql.NullInt64  `sql:"id"`
	TestReportID sql.NullInt64  `sql:"test_report_id"`
	FileName     sql.NullString `sql:"file_name"`
	ObjectPath   sql.NullString `sql:"object_path"`
	FileSize     sql.NullInt64  `sql:"file_size"`
	FileType     sql.NullString `sql:"file_type"`
	PresignedURL sql.NullString `sql:"presigned_url"`
	CreatedAt    sql.NullInt64  `sql:"created_at"`

	// References to related objects
	TestReport *TestReport `gorm:"foreignKey:TestReportID"`
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
// When a field within the TestAttachment type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (ta *TestAttachment) Nullify() *TestAttachment {
	if ta == nil {
		return nil
	}

	// check if the ID field should be false
	if ta.ID.Int64 == 0 {
		ta.ID.Valid = false
	}

	// check if the TestReportID field should be false
	if ta.TestReportID.Int64 == 0 {
		ta.TestReportID.Valid = false
	}

	// check if the Created field should be false
	if ta.CreatedAt.Int64 == 0 {
		ta.CreatedAt.Valid = false
	}

	return ta
}

// ToAPI converts the TestAttachment type
// to the API representation of the type.
func (ta *TestAttachment) ToAPI() *api.TestAttachment {
	attachment := new(api.TestAttachment)
	attachment.SetID(ta.ID.Int64)

	// var tr *api.Artifacts
	// if ta.Artifacts.ID.Valid {
	// 	tr = ta.Artifacts.ToAPI()
	// } else {
	// 	tr = new(api.Artifacts)
	// tr.SetID(ta.TestReportID.Int64)
	// }

	attachment.SetTestReportID(ta.TestReportID.Int64)
	attachment.SetFileName(ta.FileName.String)
	attachment.SetObjectPath(ta.ObjectPath.String)
	attachment.SetFileSize(ta.FileSize.Int64)
	attachment.SetFileType(ta.FileType.String)
	attachment.SetPresignedURL(ta.PresignedURL.String)
	attachment.SetCreatedAt(ta.CreatedAt.Int64)

	return attachment
}

// Validate ensures the TestAttachment type is valid
// by checking if the required fields are set.
func (ta *TestAttachment) Validate() error {
	// verify the TestReportID field is populated
	if !ta.TestReportID.Valid || ta.TestReportID.Int64 <= 0 {
		return ErrEmptyTestReportID
	}

	// verify the FileName field is populated
	if !ta.FileName.Valid || len(ta.FileName.String) == 0 {
		return ErrEmptyFileName
	}

	// verify the ObjectPath field is populated
	if !ta.ObjectPath.Valid || len(ta.ObjectPath.String) == 0 {
		return ErrEmptyObjectPath
	}

	// verify the FileType field is populated
	if !ta.FileType.Valid || len(ta.FileType.String) == 0 {
		return ErrEmptyFileType
	}

	// Note: FileSize and PresignedUrl are optional during creation
	// They may be set later in the workflow

	// ensure that all TestAttachment string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	ta.FileName = sql.NullString{String: util.Sanitize(ta.FileName.String), Valid: ta.FileName.Valid}
	ta.ObjectPath = sql.NullString{String: util.Sanitize(ta.ObjectPath.String), Valid: ta.ObjectPath.Valid}
	ta.FileType = sql.NullString{String: util.Sanitize(ta.FileType.String), Valid: ta.FileType.Valid}

	// Only sanitize PresignedUrl if it's provided
	if ta.PresignedURL.Valid {
		ta.PresignedURL = sql.NullString{String: util.Sanitize(ta.PresignedURL.String), Valid: ta.PresignedURL.Valid}
	}

	return nil
}

// TestAttachmentFromAPI converts the API TestAttachment type
// to a database report attachment type.
func TestAttachmentFromAPI(ta *api.TestAttachment) *TestAttachment {
	attachment := &TestAttachment{
		ID:           sql.NullInt64{Int64: ta.GetID(), Valid: ta.GetID() > 0},
		TestReportID: sql.NullInt64{Int64: ta.GetTestReportID(), Valid: ta.GetTestReportID() > 0},
		FileName:     sql.NullString{String: ta.GetFileName(), Valid: len(ta.GetFileName()) > 0},
		ObjectPath:   sql.NullString{String: ta.GetObjectPath(), Valid: len(ta.GetObjectPath()) > 0},
		FileSize:     sql.NullInt64{Int64: ta.GetFileSize(), Valid: ta.GetFileSize() > 0},
		FileType:     sql.NullString{String: ta.GetFileType(), Valid: len(ta.GetFileType()) > 0},
		PresignedURL: sql.NullString{String: ta.GetPresignedURL(), Valid: len(ta.GetPresignedURL()) > 0},
		CreatedAt:    sql.NullInt64{Int64: ta.GetCreatedAt(), Valid: ta.GetCreatedAt() > 0},
	}

	return attachment.Nullify()
}
