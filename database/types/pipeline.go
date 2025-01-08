// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyPipelineCommit defines the error type when a
	// Pipeline type has an empty Commit field provided.
	ErrEmptyPipelineCommit = errors.New("empty pipeline commit provided")

	// ErrEmptyPipelineRef defines the error type when a
	// Pipeline type has an empty Ref field provided.
	ErrEmptyPipelineRef = errors.New("empty pipeline ref provided")

	// ErrEmptyPipelineRepoID defines the error type when a
	// Pipeline type has an empty RepoID field provided.
	ErrEmptyPipelineRepoID = errors.New("empty pipeline repo_id provided")

	// ErrEmptyPipelineType defines the error type when a
	// Pipeline type has an empty Type field provided.
	ErrEmptyPipelineType = errors.New("empty pipeline type provided")

	// ErrEmptyPipelineVersion defines the error type when a
	// Pipeline type has an empty Version field provided.
	ErrEmptyPipelineVersion = errors.New("empty pipeline version provided")

	// ErrExceededWarningsLimit defines the error type when a
	// Pipeline warnings field has too many total characters.
	ErrExceededWarningsLimit = errors.New("exceeded character limit for pipeline warnings")
)

// Pipeline is the database representation of a pipeline.
type Pipeline struct {
	ID              sql.NullInt64  `sql:"id"`
	RepoID          sql.NullInt64  `sql:"repo_id"`
	Commit          sql.NullString `sql:"commit"`
	Flavor          sql.NullString `sql:"flavor"`
	Platform        sql.NullString `sql:"platform"`
	Ref             sql.NullString `sql:"ref"`
	Type            sql.NullString `sql:"type"`
	Version         sql.NullString `sql:"version"`
	ExternalSecrets sql.NullBool   `sql:"external_secrets"`
	InternalSecrets sql.NullBool   `sql:"internal_secrets"`
	Services        sql.NullBool   `sql:"services"`
	Stages          sql.NullBool   `sql:"stages"`
	Steps           sql.NullBool   `sql:"steps"`
	Templates       sql.NullBool   `sql:"templates"`
	Warnings        pq.StringArray `sql:"warnings"         gorm:"type:varchar(5000)"`
	Data            []byte         `sql:"data"`

	Repo Repo `gorm:"foreignKey:RepoID"`
}

// Compress will manipulate the existing data for the
// pipeline by compressing that data. This produces
// a significantly smaller amount of data that is
// stored in the system.
func (p *Pipeline) Compress(level int) error {
	// compress the database pipeline data
	data, err := util.Compress(level, p.Data)
	if err != nil {
		return err
	}

	// overwrite database pipeline data with compressed pipeline data
	p.Data = data

	return nil
}

// Decompress will manipulate the existing data for the
// pipeline by decompressing that data. This allows us
// to have a significantly smaller amount of data that
// is stored in the system.
func (p *Pipeline) Decompress() error {
	// decompress the database pipeline data
	data, err := util.Decompress(p.Data)
	if err != nil {
		return err
	}

	// overwrite compressed pipeline data with decompressed pipeline data
	p.Data = data

	return nil
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the Pipeline type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (p *Pipeline) Nullify() *Pipeline {
	if p == nil {
		return nil
	}

	// check if the ID field should be false
	if p.ID.Int64 == 0 {
		p.ID.Valid = false
	}

	// check if the RepoID field should be false
	if p.RepoID.Int64 == 0 {
		p.RepoID.Valid = false
	}

	// check if the Commit field should be false
	if len(p.Commit.String) == 0 {
		p.Commit.Valid = false
	}

	// check if the Flavor field should be false
	if len(p.Flavor.String) == 0 {
		p.Flavor.Valid = false
	}

	// check if the Platform field should be false
	if len(p.Platform.String) == 0 {
		p.Platform.Valid = false
	}

	// check if the Ref field should be false
	if len(p.Ref.String) == 0 {
		p.Ref.Valid = false
	}

	// check if the Type field should be false
	if len(p.Type.String) == 0 {
		p.Type.Valid = false
	}

	// check if the Version field should be false
	if len(p.Version.String) == 0 {
		p.Version.Valid = false
	}

	return p
}

// ToAPI converts the Pipeline type
// to a API Pipeline type.
func (p *Pipeline) ToAPI() *api.Pipeline {
	pipeline := new(api.Pipeline)

	pipeline.SetID(p.ID.Int64)
	pipeline.SetRepo(p.Repo.ToAPI())
	pipeline.SetCommit(p.Commit.String)
	pipeline.SetFlavor(p.Flavor.String)
	pipeline.SetPlatform(p.Platform.String)
	pipeline.SetRef(p.Ref.String)
	pipeline.SetType(p.Type.String)
	pipeline.SetVersion(p.Version.String)
	pipeline.SetExternalSecrets(p.ExternalSecrets.Bool)
	pipeline.SetInternalSecrets(p.InternalSecrets.Bool)
	pipeline.SetServices(p.Services.Bool)
	pipeline.SetStages(p.Stages.Bool)
	pipeline.SetSteps(p.Steps.Bool)
	pipeline.SetTemplates(p.Templates.Bool)
	pipeline.SetWarnings(p.Warnings)
	pipeline.SetData(p.Data)

	return pipeline
}

// Validate verifies the necessary fields for
// the Pipeline type are populated correctly.
func (p *Pipeline) Validate() error {
	// verify the Commit field is populated
	if len(p.Commit.String) == 0 {
		return ErrEmptyPipelineCommit
	}

	// verify the Ref field is populated
	if len(p.Ref.String) == 0 {
		return ErrEmptyPipelineRef
	}

	// verify the RepoID field is populated
	if p.RepoID.Int64 <= 0 {
		return ErrEmptyPipelineRepoID
	}

	// verify the Type field is populated
	if len(p.Type.String) == 0 {
		return ErrEmptyPipelineType
	}

	// verify the Version field is populated
	if len(p.Version.String) == 0 {
		return ErrEmptyPipelineVersion
	}

	// calculate total size of warnings
	total := 0
	for _, w := range p.Warnings {
		total += len(w)
	}

	// verify the Warnings field is within the database constraints
	// len is to factor in number of comma separators included in the database field,
	// removing 1 due to the last item not having an appended comma
	if (total + len(p.Warnings) - 1) > constants.PipelineWarningsMaxSize {
		return ErrExceededWarningsLimit
	}

	// ensure that all Pipeline string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	p.Commit = sql.NullString{String: util.Sanitize(p.Commit.String), Valid: p.Commit.Valid}
	p.Flavor = sql.NullString{String: util.Sanitize(p.Flavor.String), Valid: p.Flavor.Valid}
	p.Platform = sql.NullString{String: util.Sanitize(p.Platform.String), Valid: p.Platform.Valid}
	p.Ref = sql.NullString{String: util.Sanitize(p.Ref.String), Valid: p.Ref.Valid}
	p.Type = sql.NullString{String: util.Sanitize(p.Type.String), Valid: p.Type.Valid}
	p.Version = sql.NullString{String: util.Sanitize(p.Version.String), Valid: p.Version.Valid}

	return nil
}

// PipelineFromAPI converts the API Pipeline type
// to a database Pipeline type.
func PipelineFromAPI(p *api.Pipeline) *Pipeline {
	pipeline := &Pipeline{
		ID:              sql.NullInt64{Int64: p.GetID(), Valid: true},
		RepoID:          sql.NullInt64{Int64: p.GetRepo().GetID(), Valid: true},
		Commit:          sql.NullString{String: p.GetCommit(), Valid: true},
		Flavor:          sql.NullString{String: p.GetFlavor(), Valid: true},
		Platform:        sql.NullString{String: p.GetPlatform(), Valid: true},
		Ref:             sql.NullString{String: p.GetRef(), Valid: true},
		Type:            sql.NullString{String: p.GetType(), Valid: true},
		Version:         sql.NullString{String: p.GetVersion(), Valid: true},
		ExternalSecrets: sql.NullBool{Bool: p.GetExternalSecrets(), Valid: true},
		InternalSecrets: sql.NullBool{Bool: p.GetInternalSecrets(), Valid: true},
		Services:        sql.NullBool{Bool: p.GetServices(), Valid: true},
		Stages:          sql.NullBool{Bool: p.GetStages(), Valid: true},
		Steps:           sql.NullBool{Bool: p.GetSteps(), Valid: true},
		Templates:       sql.NullBool{Bool: p.GetTemplates(), Valid: true},
		Warnings:        pq.StringArray(p.GetWarnings()),
		Data:            p.GetData(),
	}

	return pipeline.Nullify()
}
