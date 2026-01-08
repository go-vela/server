// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"errors"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyBuildID defines the error type when a
	// Artifact type has an empty BuildID field provided.
	ErrEmptyBuildID = errors.New("empty build_id provided")
	// ErrEmptyFileName defines the error type when a
	// Artifact type has an empty FileName field provided.
	ErrEmptyFileName = errors.New("empty file_name provided")
	// ErrEmptyObjectPath defines the error type when a
	// Artifact type has an empty ObjectPath field provided.
	ErrEmptyObjectPath = errors.New("empty object_path provided")
	// ErrEmptyFileSize defines the error type when a
	// Artifact type has an empty FileSize field provided.
	ErrEmptyFileSize = errors.New("empty file_size provided")
	// ErrEmptyFileType defines the error type when a
	// Artifact type has an empty FileType field provided.
	ErrEmptyFileType = errors.New("empty file_type provided")
	// ErrEmptyPresignedURL defines the error type when a
	// Artifact type has an empty PresignedUrl field provided.
	ErrEmptyPresignedURL = errors.New("empty presigned_url provided")
)

type Artifact struct {
	ID           sql.NullInt64  `sql:"id"`
	BuildID      sql.NullInt64  `sql:"build_id"`
	FileName     sql.NullString `sql:"file_name"`
	ObjectPath   sql.NullString `sql:"object_path"`
	FileSize     sql.NullInt64  `sql:"file_size"`
	FileType     sql.NullString `sql:"file_type"`
	PresignedURL sql.NullString `sql:"presigned_url"`
	CreatedAt    sql.NullInt64  `sql:"created_at"`
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
// When a field within the Artifact type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (a *Artifact) Nullify() *Artifact {
	if a == nil {
		return nil
	}

	// check if the ID field should be false
	if a.ID.Int64 == 0 {
		a.ID.Valid = false
	}

	// check if the BuildID field should be false
	if a.BuildID.Int64 == 0 {
		a.BuildID.Valid = false
	}

	// check if the Created field should be false
	if a.CreatedAt.Int64 == 0 {
		a.CreatedAt.Valid = false
	}

	return a
}

// ToAPI converts the Artifact type
// to the API representation of the type.
func (a *Artifact) ToAPI() *api.Artifact {
	artifact := new(api.Artifact)
	artifact.SetID(a.ID.Int64)

	// var tr *api.Artifacts
	// if a.Artifacts.ID.Valid {
	// 	tr = a.Artifacts.ToAPI()
	// } else {
	// 	tr = new(api.Artifacts)
	// tr.SetID(a.BuildID.Int64)
	// }

	artifact.SetBuildID(a.BuildID.Int64)
	artifact.SetFileName(a.FileName.String)
	artifact.SetObjectPath(a.ObjectPath.String)
	artifact.SetFileSize(a.FileSize.Int64)
	artifact.SetFileType(a.FileType.String)
	artifact.SetPresignedURL(a.PresignedURL.String)
	artifact.SetCreatedAt(a.CreatedAt.Int64)

	return artifact
}

// Validate ensures the Artifact type is valid
// by checking if the required fields are set.
func (a *Artifact) Validate() error {
	// verify the BuildID field is populated
	if !a.BuildID.Valid || a.BuildID.Int64 <= 0 {
		return ErrEmptyBuildID
	}

	// verify the FileName field is populated
	if !a.FileName.Valid || len(a.FileName.String) == 0 {
		return ErrEmptyFileName
	}

	// verify the ObjectPath field is populated
	if !a.ObjectPath.Valid || len(a.ObjectPath.String) == 0 {
		return ErrEmptyObjectPath
	}

	// verify the FileType field is populated
	if !a.FileType.Valid || len(a.FileType.String) == 0 {
		return ErrEmptyFileType
	}

	// Note: FileSize and PresignedUrl are optional during creation
	// They may be set later in the workflow

	// ensure that all Artifact string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	a.FileName = sql.NullString{String: util.Sanitize(a.FileName.String), Valid: a.FileName.Valid}
	a.ObjectPath = sql.NullString{String: util.Sanitize(a.ObjectPath.String), Valid: a.ObjectPath.Valid}
	a.FileType = sql.NullString{String: util.Sanitize(a.FileType.String), Valid: a.FileType.Valid}

	// Only sanitize PresignedUrl if it's provided
	if a.PresignedURL.Valid {
		a.PresignedURL = sql.NullString{String: util.Sanitize(a.PresignedURL.String), Valid: a.PresignedURL.Valid}
	}

	return nil
}

// ArtifactFromAPI converts the API Artifact type
// to a database report artifact type.
func ArtifactFromAPI(a *api.Artifact) *Artifact {
	artifact := &Artifact{
		ID:           sql.NullInt64{Int64: a.GetID(), Valid: a.GetID() > 0},
		BuildID:      sql.NullInt64{Int64: a.GetBuildID(), Valid: a.GetBuildID() > 0},
		FileName:     sql.NullString{String: a.GetFileName(), Valid: len(a.GetFileName()) > 0},
		ObjectPath:   sql.NullString{String: a.GetObjectPath(), Valid: len(a.GetObjectPath()) > 0},
		FileSize:     sql.NullInt64{Int64: a.GetFileSize(), Valid: a.GetFileSize() > 0},
		FileType:     sql.NullString{String: a.GetFileType(), Valid: len(a.GetFileType()) > 0},
		PresignedURL: sql.NullString{String: a.GetPresignedURL(), Valid: len(a.GetPresignedURL()) > 0},
		CreatedAt:    sql.NullInt64{Int64: a.GetCreatedAt(), Valid: a.GetCreatedAt() > 0},
	}

	return artifact.Nullify()
}
