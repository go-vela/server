// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"errors"

	api "github.com/go-vela/server/api/types"
)

var (
	// ErrEmptyReportBuildID defines the error type when a
	// TestReport type has an empty BuildID field provided.
	ErrEmptyReportBuildID = errors.New("empty report build_id provided")
)

// TestReport is the database representation of a report.
type TestReport struct {
	ID      sql.NullInt64 `sql:"id"`
	BuildID sql.NullInt64 `sql:"build_id"`
	Created sql.NullInt64 `sql:"created"`

	// References to related objects
	Build *Build `gorm:"foreignKey:BuildID"`
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the TestReport type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (r *TestReport) Nullify() *TestReport {
	if r == nil {
		return nil
	}

	// check if the ID field should be false
	if r.ID.Int64 == 0 {
		r.ID.Valid = false
	}

	// check if the BuildID field should be false
	if r.BuildID.Int64 == 0 {
		r.BuildID.Valid = false
	}

	// check if the Created field should be false
	if r.Created.Int64 == 0 {
		r.Created.Valid = false
	}

	return r
}

// ToAPI converts the TestReport type
// to an API TestReport type.
func (r *TestReport) ToAPI() *api.TestReport {
	report := new(api.TestReport)
	report.SetID(r.ID.Int64)
	report.SetBuildID(r.BuildID.Int64)
	report.SetCreated(r.Created.Int64)

	// set Repo based on presence of repo data
	//var tra *api.TestReportAttachments
	//if r.Attachments.ID.Valid {
	//	tra = r.Attachments.ToAPI()
	//} else {
	//	tra = new(api.TestReportAttachments)
	//	tra.SetID(r.Attachments.ID.Int64)
	//}
	//
	//report.SetReportAttachments(tra)
	//// Convert attachments if available
	//attachment := new(api.TestReportAttachments)
	//attachment.SetID(r.Attachments.ID.Int64)
	//attachment.SetTestReportID(report.GetID())
	//attachment.SetFileName(r.Attachments.FileName.String)
	//attachment.SetObjectPath(r.Attachments.ObjectPath.String)
	//attachment.SetFileSize(r.Attachments.FileSize.Int64)
	//attachment.SetFileType(r.Attachments.FileType.String)
	//attachment.SetPresignedUrl(r.Attachments.PresignedUrl.String)
	//attachment.SetCreated(r.Attachments.Created.Int64)

	return report
}

// Validate verifies the necessary fields for
// the TestReport type are populated correctly.
func (r *TestReport) Validate() error {
	// verify the BuildID field is populated
	if r.BuildID.Int64 <= 0 {
		return ErrEmptyReportBuildID
	}

	// Also validate any attachments
	//r.Attachments.FileName = sql.NullString{String: util.Sanitize(r.Attachments.FileName.String), Valid: r.Attachments.FileName.Valid}
	//r.Attachments.ObjectPath = sql.NullString{String: util.Sanitize(r.Attachments.ObjectPath.String), Valid: r.Attachments.ObjectPath.Valid}
	//r.Attachments.FileType = sql.NullString{String: util.Sanitize(r.Attachments.FileType.String), Valid: r.Attachments.FileType.Valid}
	//r.Attachments.PresignedUrl = sql.NullString{String: util.Sanitize(r.Attachments.PresignedUrl.String), Valid: r.Attachments.PresignedUrl.Valid}

	return nil
}

// TestReportFromAPI converts the API TestReport type
// to a database report type.
func TestReportFromAPI(r *api.TestReport) *TestReport {
	report := &TestReport{
		ID:      sql.NullInt64{Int64: r.GetID(), Valid: r.GetID() > 0},
		BuildID: sql.NullInt64{Int64: r.GetBuildID(), Valid: r.GetBuildID() > 0},
		Created: sql.NullInt64{Int64: r.GetCreated(), Valid: r.GetCreated() > 0},
	}

	return report.Nullify()
}
