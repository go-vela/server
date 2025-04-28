// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"errors"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyReportID defines the error type when a
	// TestReportAttachment type has an empty ReportID field provided.
	ErrEmptyReportID = errors.New("empty report_id provided")
	// ErrEmptyFileName defines the error type when a
	// TestReportAttachment type has an empty FileName field provided.
	ErrEmptyFileName = errors.New("empty file_name provided")
	// ErrEmptyObjectPath defines the error type when a
	// TestReportAttachment type has an empty ObjectPath field provided.
	ErrEmptyObjectPath = errors.New("empty object_path provided")
	// ErrEmptyFileSize defines the error type when a
	// TestReportAttachment type has an empty FileSize field provided.
	ErrEmptyFileSize = errors.New("empty file_size provided")
	// ErrEmptyFileType defines the error type when a
	// TestReportAttachment type has an empty FileType field provided.
	ErrEmptyFileType = errors.New("empty file_type provided")
	// ErrEmptyPresignedUrl defines the error type when a
	// TestReportAttachment type has an empty PresignedUrl field provided.
	ErrEmptyPresignedUrl = errors.New("empty presigned_url provided")
)

type TestReportAttachment struct {
	ID           sql.NullInt64  `sql:"id"`
	TestReportID sql.NullInt64  `sql:"test_report_id"`
	FileName     sql.NullString `sql:"file_name"`
	ObjectPath   sql.NullString `sql:"object_path"`
	FileSize     sql.NullInt64  `sql:"file_size"`
	FileType     sql.NullString `sql:"file_type"`
	PresignedUrl sql.NullString `sql:"presigned_url"`
	Created      sql.NullInt64  `sql:"created"`

	// References to related objects
	TestReport *TestReport `gorm:"foreignKey:TestReportID"`
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
// When a field within the TestReportAttachment type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (a *TestReportAttachment) Nullify() *TestReportAttachment {
	if a == nil {
		return nil
	}

	// check if the ID field should be false
	if a.ID.Int64 == 0 {
		a.ID.Valid = false
	}

	// check if the TestReportID field should be false
	if a.TestReportID.Int64 == 0 {
		a.TestReportID.Valid = false
	}

	// check if the Created field should be false
	if a.Created.Int64 == 0 {
		a.Created.Valid = false
	}

	return a
}

// ToAPI converts the TestReportAttachment type
// to the API representation of the type.
func (a *TestReportAttachment) ToAPI() *api.TestReportAttachments {
	attachment := new(api.TestReportAttachments)
	attachment.SetID(a.ID.Int64)

	// var tr *api.TestReport
	// if a.TestReport.ID.Valid {
	// 	tr = a.TestReport.ToAPI()
	// } else {
	// 	tr = new(api.TestReport)
	// tr.SetID(a.TestReportID.Int64)
	// }

	attachment.SetTestReportID(a.TestReportID.Int64)
	attachment.SetFileName(a.FileName.String)
	attachment.SetObjectPath(a.ObjectPath.String)
	attachment.SetFileSize(a.FileSize.Int64)
	attachment.SetFileType(a.FileType.String)
	attachment.SetPresignedUrl(a.PresignedUrl.String)
	attachment.SetCreated(a.Created.Int64)

	return attachment
}

// Validate ensures the TestReportAttachment type is valid
// by checking if the required fields are set.
func (a *TestReportAttachment) Validate() error {
	if !a.TestReportID.Valid {
		return ErrEmptyReportID
	}

	if !a.FileName.Valid {
		return ErrEmptyFileName
	}

	if !a.ObjectPath.Valid {
		return ErrEmptyObjectPath
	}

	if !a.FileSize.Valid {
		return ErrEmptyFileSize
	}

	if !a.FileType.Valid {
		return ErrEmptyFileType
	}

	if !a.PresignedUrl.Valid {
		return ErrEmptyPresignedUrl
	}

	// ensure that all ReportAttachment fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	a.FileName = sql.NullString{String: util.Sanitize(a.FileName.String), Valid: a.FileName.Valid}
	a.ObjectPath = sql.NullString{String: util.Sanitize(a.ObjectPath.String), Valid: a.ObjectPath.Valid}
	a.FileType = sql.NullString{String: util.Sanitize(a.FileType.String), Valid: a.FileType.Valid}
	a.PresignedUrl = sql.NullString{String: util.Sanitize(a.PresignedUrl.String), Valid: a.PresignedUrl.Valid}

	return nil
}

// TestReportAttachmentFromAPI converts the API TestReportAttachments type
// to a database report attachment type.
func TestReportAttachmentFromAPI(r *api.TestReportAttachments) *TestReportAttachment {
	attachment := &TestReportAttachment{
		ID:           sql.NullInt64{Int64: r.GetID(), Valid: r.GetID() > 0},
		TestReportID: sql.NullInt64{Int64: r.GetTestReportID(), Valid: r.GetTestReportID() > 0},
		FileName:     sql.NullString{String: r.GetFileName(), Valid: len(r.GetFileName()) > 0},
		ObjectPath:   sql.NullString{String: r.GetObjectPath(), Valid: len(r.GetObjectPath()) > 0},
		FileSize:     sql.NullInt64{Int64: r.GetFileSize(), Valid: r.GetFileSize() > 0},
		FileType:     sql.NullString{String: r.GetFileType(), Valid: len(r.GetFileType()) > 0},
		PresignedUrl: sql.NullString{String: r.GetPresignedUrl(), Valid: len(r.GetPresignedUrl()) > 0},
		Created:      sql.NullInt64{Int64: r.GetCreated(), Valid: r.GetCreated() > 0},
	}

	return attachment.Nullify()
}
