// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"errors"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyBuildNumber defines the error type when a
	// Build type has an empty Number field provided.
	ErrEmptyBuildNumber = errors.New("empty build number provided")

	// ErrEmptyBuildRepoID defines the error type when a
	// Build type has an empty `RepoID` field provided.
	ErrEmptyBuildRepoID = errors.New("empty build repo_id provided")
)

const (
	// Maximum title field length.
	maxTitleLength = 1000
	// Maximum message field length.
	maxMessageLength = 2000
	// Maximum error field length.
	maxErrorLength = 1000
)

// Build is the database representation of a build for a pipeline.
type Build struct {
	ID            sql.NullInt64      `sql:"id"`
	RepoID        sql.NullInt64      `sql:"repo_id"`
	PipelineID    sql.NullInt64      `sql:"pipeline_id"`
	Number        sql.NullInt64      `sql:"number"`
	Parent        sql.NullInt64      `sql:"parent"`
	Event         sql.NullString     `sql:"event"`
	EventAction   sql.NullString     `sql:"event_action"`
	Status        sql.NullString     `sql:"status"`
	Error         sql.NullString     `sql:"error"`
	Enqueued      sql.NullInt64      `sql:"enqueued"`
	Created       sql.NullInt64      `sql:"created"`
	Started       sql.NullInt64      `sql:"started"`
	Finished      sql.NullInt64      `sql:"finished"`
	Deploy        sql.NullString     `sql:"deploy"`
	DeployNumber  sql.NullInt64      `sql:"deploy_number"`
	DeployPayload raw.StringSliceMap `sql:"deploy_payload" gorm:"type:varchar(2000)"`
	Clone         sql.NullString     `sql:"clone"`
	Source        sql.NullString     `sql:"source"`
	Title         sql.NullString     `sql:"title"`
	Message       sql.NullString     `sql:"message"`
	Commit        sql.NullString     `sql:"commit"`
	Sender        sql.NullString     `sql:"sender"`
	SenderSCMID   sql.NullString     `sql:"sender_scm_id"`
	Author        sql.NullString     `sql:"author"`
	Email         sql.NullString     `sql:"email"`
	Link          sql.NullString     `sql:"link"`
	Branch        sql.NullString     `sql:"branch"`
	Ref           sql.NullString     `sql:"ref"`
	BaseRef       sql.NullString     `sql:"base_ref"`
	HeadRef       sql.NullString     `sql:"head_ref"`
	Host          sql.NullString     `sql:"host"`
	Runtime       sql.NullString     `sql:"runtime"`
	Distribution  sql.NullString     `sql:"distribution"`
	ApprovedAt    sql.NullInt64      `sql:"approved_at"`
	ApprovedBy    sql.NullString     `sql:"approved_by"`

	Repo Repo `gorm:"foreignKey:RepoID"`
}

// Crop prepares the Build type for inserting into the database by
// trimming values that may exceed the database column limit.
func (b *Build) Crop() *Build {
	// trim the Title field to 1000 characters
	if len(b.Title.String) > maxTitleLength {
		b.Title = sql.NullString{String: b.Title.String[:maxTitleLength], Valid: true}
	}

	// trim the Message field to 2000 characters
	if len(b.Message.String) > maxMessageLength {
		b.Message = sql.NullString{String: b.Message.String[:maxMessageLength], Valid: true}
	}

	// trim the Error field to 1000 characters
	if len(b.Error.String) > maxErrorLength {
		b.Error = sql.NullString{String: b.Error.String[:maxErrorLength], Valid: true}
	}

	return b
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the Build type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
//
//nolint:gocyclo,funlen // ignore cyclomatic complexity due to number of fields
func (b *Build) Nullify() *Build {
	if b == nil {
		return nil
	}

	// check if the ID field should be false
	if b.ID.Int64 == 0 {
		b.ID.Valid = false
	}

	// check if the RepoID field should be false
	if b.RepoID.Int64 == 0 {
		b.RepoID.Valid = false
	}

	// check if the PipelineID field should be false
	if b.PipelineID.Int64 == 0 {
		b.PipelineID.Valid = false
	}

	// check if the Number field should be false
	if b.Number.Int64 == 0 {
		b.Number.Valid = false
	}

	// check if the Parent field should be false
	if b.Parent.Int64 == 0 {
		b.Parent.Valid = false
	}

	// check if the Event field should be false
	if len(b.Event.String) == 0 {
		b.Event.Valid = false
	}

	// check if the EventAction field should be false
	if len(b.EventAction.String) == 0 {
		b.EventAction.Valid = false
	}

	// check if the Status field should be false
	if len(b.Status.String) == 0 {
		b.Status.Valid = false
	}

	// check if the Error field should be false
	if len(b.Error.String) == 0 {
		b.Error.Valid = false
	}

	// check if the Enqueued field should be false
	if b.Enqueued.Int64 == 0 {
		b.Enqueued.Valid = false
	}

	// check if the Created field should be false
	if b.Created.Int64 == 0 {
		b.Created.Valid = false
	}

	// check if the Started field should be false
	if b.Started.Int64 == 0 {
		b.Started.Valid = false
	}

	// check if the Finished field should be false
	if b.Finished.Int64 == 0 {
		b.Finished.Valid = false
	}

	// check if the Deploy field should be false
	if len(b.Deploy.String) == 0 {
		b.Deploy.Valid = false
	}

	// check if the DeployNumber field should be false
	if b.DeployNumber.Int64 == 0 {
		b.DeployNumber.Valid = false
	}

	// check if the Clone field should be false
	if len(b.Clone.String) == 0 {
		b.Clone.Valid = false
	}

	// check if the Source field should be false
	if len(b.Source.String) == 0 {
		b.Source.Valid = false
	}

	// check if the Title field should be false
	if len(b.Title.String) == 0 {
		b.Title.Valid = false
	}

	// check if the Message field should be false
	if len(b.Message.String) == 0 {
		b.Message.Valid = false
	}

	// check if the Commit field should be false
	if len(b.Commit.String) == 0 {
		b.Commit.Valid = false
	}

	// check if the Sender field should be false
	if len(b.Sender.String) == 0 {
		b.Sender.Valid = false
	}

	// check if the SenderSCMID field should be false
	if len(b.SenderSCMID.String) == 0 {
		b.SenderSCMID.Valid = false
	}

	// check if the Author field should be false
	if len(b.Author.String) == 0 {
		b.Author.Valid = false
	}

	// check if the Email field should be false
	if len(b.Email.String) == 0 {
		b.Email.Valid = false
	}

	// check if the Link field should be false
	if len(b.Link.String) == 0 {
		b.Link.Valid = false
	}

	// check if the Branch field should be false
	if len(b.Branch.String) == 0 {
		b.Branch.Valid = false
	}

	// check if the Ref field should be false
	if len(b.Ref.String) == 0 {
		b.Ref.Valid = false
	}

	// check if the BaseRef field should be false
	if len(b.BaseRef.String) == 0 {
		b.BaseRef.Valid = false
	}

	// check if the HeadRef field should be false
	if len(b.HeadRef.String) == 0 {
		b.HeadRef.Valid = false
	}

	// check if the Host field should be false
	if len(b.Host.String) == 0 {
		b.Host.Valid = false
	}

	// check if the Runtime field should be false
	if len(b.Runtime.String) == 0 {
		b.Runtime.Valid = false
	}

	// check if the Distribution field should be false
	if len(b.Distribution.String) == 0 {
		b.Distribution.Valid = false
	}

	// check if the ApprovedAt field should be false
	if b.ApprovedAt.Int64 == 0 {
		b.ApprovedAt.Valid = false
	}

	// check if the ApprovedBy field should be false
	if len(b.ApprovedBy.String) == 0 {
		b.ApprovedBy.Valid = false
	}

	return b
}

// ToAPI converts the Build type
// to an API Build type.
func (b *Build) ToAPI() *api.Build {
	build := new(api.Build)

	// set Repo based on presence of repo data
	var repo *api.Repo
	if b.Repo.ID.Valid {
		repo = b.Repo.ToAPI()
	} else {
		repo = new(api.Repo)
		repo.SetID(b.RepoID.Int64)
	}

	build.SetID(b.ID.Int64)
	build.SetRepo(repo)
	build.SetPipelineID(b.PipelineID.Int64)
	build.SetNumber(int(b.Number.Int64))
	build.SetParent(int(b.Parent.Int64))
	build.SetEvent(b.Event.String)
	build.SetEventAction(b.EventAction.String)
	build.SetStatus(b.Status.String)
	build.SetError(b.Error.String)
	build.SetEnqueued(b.Enqueued.Int64)
	build.SetCreated(b.Created.Int64)
	build.SetStarted(b.Started.Int64)
	build.SetFinished(b.Finished.Int64)
	build.SetDeploy(b.Deploy.String)
	build.SetDeployNumber(b.DeployNumber.Int64)
	build.SetDeployPayload(b.DeployPayload)
	build.SetClone(b.Clone.String)
	build.SetSource(b.Source.String)
	build.SetTitle(b.Title.String)
	build.SetMessage(b.Message.String)
	build.SetCommit(b.Commit.String)
	build.SetSender(b.Sender.String)
	build.SetSenderSCMID(b.SenderSCMID.String)
	build.SetAuthor(b.Author.String)
	build.SetEmail(b.Email.String)
	build.SetLink(b.Link.String)
	build.SetBranch(b.Branch.String)
	build.SetRef(b.Ref.String)
	build.SetBaseRef(b.BaseRef.String)
	build.SetHeadRef(b.HeadRef.String)
	build.SetHost(b.Host.String)
	build.SetRuntime(b.Runtime.String)
	build.SetDistribution(b.Distribution.String)
	build.SetApprovedAt(b.ApprovedAt.Int64)
	build.SetApprovedBy(b.ApprovedBy.String)

	return build
}

// Validate verifies the necessary fields for
// the Build type are populated correctly.
func (b *Build) Validate() error {
	// verify the RepoID field is populated
	if b.RepoID.Int64 <= 0 {
		return ErrEmptyBuildRepoID
	}

	// verify the Number field is populated
	if b.Number.Int64 <= 0 {
		return ErrEmptyBuildNumber
	}

	// ensure that all Build string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	b.Event = sql.NullString{String: util.Sanitize(b.Event.String), Valid: b.Event.Valid}
	b.EventAction = sql.NullString{String: util.Sanitize(b.EventAction.String), Valid: b.EventAction.Valid}
	b.Status = sql.NullString{String: util.Sanitize(b.Status.String), Valid: b.Status.Valid}
	b.Error = sql.NullString{String: util.Sanitize(b.Error.String), Valid: b.Error.Valid}
	b.Deploy = sql.NullString{String: util.Sanitize(b.Deploy.String), Valid: b.Deploy.Valid}
	b.Clone = sql.NullString{String: util.Sanitize(b.Clone.String), Valid: b.Clone.Valid}
	b.Source = sql.NullString{String: util.Sanitize(b.Source.String), Valid: b.Source.Valid}
	b.Title = sql.NullString{String: util.Sanitize(b.Title.String), Valid: b.Title.Valid}
	b.Message = sql.NullString{String: util.Sanitize(b.Message.String), Valid: b.Message.Valid}
	b.Commit = sql.NullString{String: util.Sanitize(b.Commit.String), Valid: b.Commit.Valid}
	b.Sender = sql.NullString{String: util.Sanitize(b.Sender.String), Valid: b.Sender.Valid}
	b.SenderSCMID = sql.NullString{String: util.Sanitize(b.SenderSCMID.String), Valid: b.SenderSCMID.Valid}
	b.Author = sql.NullString{String: util.Sanitize(b.Author.String), Valid: b.Author.Valid}
	b.Email = sql.NullString{String: util.Sanitize(b.Email.String), Valid: b.Email.Valid}
	b.Link = sql.NullString{String: util.Sanitize(b.Link.String), Valid: b.Link.Valid}
	b.Branch = sql.NullString{String: util.Sanitize(b.Branch.String), Valid: b.Branch.Valid}
	b.Ref = sql.NullString{String: util.Sanitize(b.Ref.String), Valid: b.Ref.Valid}
	b.BaseRef = sql.NullString{String: util.Sanitize(b.BaseRef.String), Valid: b.BaseRef.Valid}
	b.HeadRef = sql.NullString{String: util.Sanitize(b.HeadRef.String), Valid: b.HeadRef.Valid}
	b.Host = sql.NullString{String: util.Sanitize(b.Host.String), Valid: b.Host.Valid}
	b.Runtime = sql.NullString{String: util.Sanitize(b.Runtime.String), Valid: b.Runtime.Valid}
	b.Distribution = sql.NullString{String: util.Sanitize(b.Distribution.String), Valid: b.Distribution.Valid}
	b.ApprovedBy = sql.NullString{String: util.Sanitize(b.ApprovedBy.String), Valid: b.ApprovedBy.Valid}

	return nil
}

// BuildFromAPI converts the API Build type
// to a database build type.
func BuildFromAPI(b *api.Build) *Build {
	build := &Build{
		ID:            sql.NullInt64{Int64: b.GetID(), Valid: true},
		RepoID:        sql.NullInt64{Int64: b.GetRepo().GetID(), Valid: true},
		PipelineID:    sql.NullInt64{Int64: b.GetPipelineID(), Valid: true},
		Number:        sql.NullInt64{Int64: int64(b.GetNumber()), Valid: true},
		Parent:        sql.NullInt64{Int64: int64(b.GetParent()), Valid: true},
		Event:         sql.NullString{String: b.GetEvent(), Valid: true},
		EventAction:   sql.NullString{String: b.GetEventAction(), Valid: true},
		Status:        sql.NullString{String: b.GetStatus(), Valid: true},
		Error:         sql.NullString{String: b.GetError(), Valid: true},
		Enqueued:      sql.NullInt64{Int64: b.GetEnqueued(), Valid: true},
		Created:       sql.NullInt64{Int64: b.GetCreated(), Valid: true},
		Started:       sql.NullInt64{Int64: b.GetStarted(), Valid: true},
		Finished:      sql.NullInt64{Int64: b.GetFinished(), Valid: true},
		Deploy:        sql.NullString{String: b.GetDeploy(), Valid: true},
		DeployNumber:  sql.NullInt64{Int64: b.GetDeployNumber(), Valid: true},
		DeployPayload: b.GetDeployPayload(),
		Clone:         sql.NullString{String: b.GetClone(), Valid: true},
		Source:        sql.NullString{String: b.GetSource(), Valid: true},
		Title:         sql.NullString{String: b.GetTitle(), Valid: true},
		Message:       sql.NullString{String: b.GetMessage(), Valid: true},
		Commit:        sql.NullString{String: b.GetCommit(), Valid: true},
		Sender:        sql.NullString{String: b.GetSender(), Valid: true},
		SenderSCMID:   sql.NullString{String: b.GetSenderSCMID(), Valid: true},
		Author:        sql.NullString{String: b.GetAuthor(), Valid: true},
		Email:         sql.NullString{String: b.GetEmail(), Valid: true},
		Link:          sql.NullString{String: b.GetLink(), Valid: true},
		Branch:        sql.NullString{String: b.GetBranch(), Valid: true},
		Ref:           sql.NullString{String: b.GetRef(), Valid: true},
		BaseRef:       sql.NullString{String: b.GetBaseRef(), Valid: true},
		HeadRef:       sql.NullString{String: b.GetHeadRef(), Valid: true},
		Host:          sql.NullString{String: b.GetHost(), Valid: true},
		Runtime:       sql.NullString{String: b.GetRuntime(), Valid: true},
		Distribution:  sql.NullString{String: b.GetDistribution(), Valid: true},
		ApprovedAt:    sql.NullInt64{Int64: b.GetApprovedAt(), Valid: true},
		ApprovedBy:    sql.NullString{String: b.GetApprovedBy(), Valid: true},
	}

	return build.Nullify()
}
