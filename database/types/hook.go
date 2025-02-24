// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"errors"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyHookNumber defines the error type when a
	// Hook type has an empty Number field provided.
	ErrEmptyHookNumber = errors.New("empty webhook number provided")

	// ErrEmptyHookRepoID defines the error type when a
	// Hook type has an empty RepoID field provided.
	ErrEmptyHookRepoID = errors.New("empty webhook repo_id provided")

	// ErrEmptyHookSourceID defines the error type when a
	// Hook type has an empty SourceID field provided.
	ErrEmptyHookSourceID = errors.New("empty webhook source_id provided")

	// ErrEmptyHookWebhookID defines the error type when a
	// Hook type has an empty WebhookID field provided.
	ErrEmptyHookWebhookID = errors.New("empty webhook webhook_id provided")
)

// Hook is the database representation of a webhook for a repo.
type Hook struct {
	ID          sql.NullInt64  `sql:"id"`
	RepoID      sql.NullInt64  `sql:"repo_id"`
	BuildID     sql.NullInt64  `sql:"build_id"`
	Number      sql.NullInt64  `sql:"number"`
	SourceID    sql.NullString `sql:"source_id"`
	Created     sql.NullInt64  `sql:"created"`
	Host        sql.NullString `sql:"host"`
	Event       sql.NullString `sql:"event"`
	EventAction sql.NullString `sql:"event_action"`
	Branch      sql.NullString `sql:"branch"`
	Error       sql.NullString `sql:"error"`
	Status      sql.NullString `sql:"status"`
	Link        sql.NullString `sql:"link"`
	WebhookID   sql.NullInt64  `sql:"webhook_id"`

	Repo  Repo  `gorm:"foreignKey:RepoID"`
	Build Build `gorm:"foreignKey:BuildID"`
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the Hook type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (h *Hook) Nullify() *Hook {
	if h == nil {
		return nil
	}

	// check if the ID field should be false
	if h.ID.Int64 == 0 {
		h.ID.Valid = false
	}

	// check if the RepoID field should be false
	if h.RepoID.Int64 == 0 {
		h.RepoID.Valid = false
	}

	// check if the BuildID field should be false
	if h.BuildID.Int64 == 0 {
		h.BuildID.Valid = false
	}

	// check if the Number field should be false
	if h.Number.Int64 == 0 {
		h.Number.Valid = false
	}

	// check if the SourceID field should be false
	if len(h.SourceID.String) == 0 {
		h.SourceID.Valid = false
	}

	// check if the Created field should be false
	if h.Created.Int64 == 0 {
		h.Created.Valid = false
	}

	// check if the Host field should be false
	if len(h.Host.String) == 0 {
		h.Host.Valid = false
	}

	// check if the Event field should be false
	if len(h.Event.String) == 0 {
		h.Event.Valid = false
	}

	// check if the EventAction field should be false
	if len(h.EventAction.String) == 0 {
		h.EventAction.Valid = false
	}

	// check if the Branch field should be false
	if len(h.Branch.String) == 0 {
		h.Branch.Valid = false
	}

	// check if the Error field should be false
	if len(h.Error.String) == 0 {
		h.Error.Valid = false
	}

	// check if the Status field should be false
	if len(h.Status.String) == 0 {
		h.Status.Valid = false
	}

	// check if the Link field should be false
	if len(h.Link.String) == 0 {
		h.Link.Valid = false
	}

	// check if the WebhookID field should be false
	if h.WebhookID.Int64 == 0 {
		h.WebhookID.Valid = false
	}

	return h
}

// ToAPI converts the Hook type
// to an API Hook type.
func (h *Hook) ToAPI() *api.Hook {
	hook := new(api.Hook)

	if h.BuildID.Int64 > 0 {
		hook.SetBuild(h.Build.ToAPI())
	}

	hook.SetID(h.ID.Int64)
	hook.SetRepo(h.Repo.ToAPI())
	hook.SetNumber(h.Number.Int64)
	hook.SetSourceID(h.SourceID.String)
	hook.SetCreated(h.Created.Int64)
	hook.SetHost(h.Host.String)
	hook.SetEvent(h.Event.String)
	hook.SetEventAction(h.EventAction.String)
	hook.SetBranch(h.Branch.String)
	hook.SetError(h.Error.String)
	hook.SetStatus(h.Status.String)
	hook.SetLink(h.Link.String)
	hook.SetWebhookID(h.WebhookID.Int64)

	return hook
}

// Validate verifies the necessary fields for
// the Hook type are populated correctly.
func (h *Hook) Validate() error {
	// verify the RepoID field is populated
	if h.RepoID.Int64 <= 0 {
		return ErrEmptyHookRepoID
	}

	// verify the Number field is populated
	if h.Number.Int64 <= 0 {
		return ErrEmptyHookNumber
	}

	// verify the SourceID field is populated
	if len(h.SourceID.String) <= 0 {
		return ErrEmptyHookSourceID
	}

	// verify the WebhookID field is populated
	if h.WebhookID.Int64 <= 0 {
		return ErrEmptyHookWebhookID
	}

	// ensure that all Hook string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	h.SourceID = sql.NullString{String: util.Sanitize(h.SourceID.String), Valid: h.SourceID.Valid}
	h.Host = sql.NullString{String: util.Sanitize(h.Host.String), Valid: h.Host.Valid}
	h.Event = sql.NullString{String: util.Sanitize(h.Event.String), Valid: h.Event.Valid}
	h.EventAction = sql.NullString{String: util.Sanitize(h.EventAction.String), Valid: h.EventAction.Valid}
	h.Branch = sql.NullString{String: util.Sanitize(h.Branch.String), Valid: h.Branch.Valid}
	h.Error = sql.NullString{String: util.Sanitize(h.Error.String), Valid: h.Error.Valid}
	h.Status = sql.NullString{String: util.Sanitize(h.Status.String), Valid: h.Status.Valid}
	h.Link = sql.NullString{String: util.Sanitize(h.Link.String), Valid: h.Link.Valid}

	return nil
}

// HookFromAPI converts the API Hook type
// to a database Hook type.
func HookFromAPI(h *api.Hook) *Hook {
	hook := &Hook{
		ID:          sql.NullInt64{Int64: h.GetID(), Valid: true},
		RepoID:      sql.NullInt64{Int64: h.GetRepo().GetID(), Valid: true},
		BuildID:     sql.NullInt64{Int64: h.GetBuild().GetID(), Valid: true},
		Number:      sql.NullInt64{Int64: int64(h.GetNumber()), Valid: true},
		SourceID:    sql.NullString{String: h.GetSourceID(), Valid: true},
		Created:     sql.NullInt64{Int64: h.GetCreated(), Valid: true},
		Host:        sql.NullString{String: h.GetHost(), Valid: true},
		Event:       sql.NullString{String: h.GetEvent(), Valid: true},
		EventAction: sql.NullString{String: h.GetEventAction(), Valid: true},
		Branch:      sql.NullString{String: h.GetBranch(), Valid: true},
		Error:       sql.NullString{String: h.GetError(), Valid: true},
		Status:      sql.NullString{String: h.GetStatus(), Valid: true},
		Link:        sql.NullString{String: h.GetLink(), Valid: true},
		WebhookID:   sql.NullInt64{Int64: h.GetWebhookID(), Valid: true},
	}

	return hook.Nullify()
}
