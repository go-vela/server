// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
)

// Hook is the API representation of a webhook.
//
// swagger:model Webhook
type Hook struct {
	ID          *int64  `json:"id,omitempty"`
	Repo        *Repo   `json:"repo,omitempty"`
	Build       *Build  `json:"build,omitempty"`
	Number      *int64  `json:"number,omitempty"`
	SourceID    *string `json:"source_id,omitempty"`
	Created     *int64  `json:"created,omitempty"`
	Host        *string `json:"host,omitempty"`
	Event       *string `json:"event,omitempty"`
	EventAction *string `json:"event_action,omitempty"`
	Branch      *string `json:"branch,omitempty"`
	Error       *string `json:"error,omitempty"`
	Status      *string `json:"status,omitempty"`
	Link        *string `json:"link,omitempty"`
	WebhookID   *int64  `json:"webhook_id,omitempty"`
}

// GetID returns the ID field.
//
// When the provided Hook type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (h *Hook) GetID() int64 {
	// return zero value if Hook type or ID field is nil
	if h == nil || h.ID == nil {
		return 0
	}

	return *h.ID
}

// GetRepo returns the Repo field.
//
// When the provided Hook type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (h *Hook) GetRepo() *Repo {
	// return zero value if Hook type or Repo field is nil
	if h == nil || h.Repo == nil {
		return new(Repo)
	}

	return h.Repo
}

// GetBuild returns the Build field.
//
// When the provided Hook type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (h *Hook) GetBuild() *Build {
	// return zero value if Hook type or Build field is nil
	if h == nil || h.Build == nil {
		return new(Build)
	}

	return h.Build
}

// GetNumber returns the Number field.
//
// When the provided Hook type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (h *Hook) GetNumber() int64 {
	// return zero value if Hook type or BuildID field is nil
	if h == nil || h.Number == nil {
		return 0
	}

	return *h.Number
}

// GetSourceID returns the SourceID field.
//
// When the provided Hook type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (h *Hook) GetSourceID() string {
	// return zero value if Hook type or SourceID field is nil
	if h == nil || h.SourceID == nil {
		return ""
	}

	return *h.SourceID
}

// GetCreated returns the Created field.
//
// When the provided Hook type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (h *Hook) GetCreated() int64 {
	// return zero value if Hook type or Created field is nil
	if h == nil || h.Created == nil {
		return 0
	}

	return *h.Created
}

// GetHost returns the Host field.
//
// When the provided Hook type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (h *Hook) GetHost() string {
	// return zero value if Hook type or Host field is nil
	if h == nil || h.Host == nil {
		return ""
	}

	return *h.Host
}

// GetEvent returns the Event field.
//
// When the provided Hook type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (h *Hook) GetEvent() string {
	// return zero value if Hook type or Event field is nil
	if h == nil || h.Event == nil {
		return ""
	}

	return *h.Event
}

// GetEventAction returns the EventAction field.
//
// When the provided Hook type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (h *Hook) GetEventAction() string {
	// return zero value if Hook type or EventAction field is nil
	if h == nil || h.EventAction == nil {
		return ""
	}

	return *h.EventAction
}

// GetBranch returns the Branch field.
//
// When the provided Hook type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (h *Hook) GetBranch() string {
	// return zero value if Hook type or Branch field is nil
	if h == nil || h.Branch == nil {
		return ""
	}

	return *h.Branch
}

// GetError returns the Error field.
//
// When the provided Hook type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (h *Hook) GetError() string {
	// return zero value if Hook type or Error field is nil
	if h == nil || h.Error == nil {
		return ""
	}

	return *h.Error
}

// GetStatus returns the Status field.
//
// When the provided Hook type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (h *Hook) GetStatus() string {
	// return zero value if Hook type or Status field is nil
	if h == nil || h.Status == nil {
		return ""
	}

	return *h.Status
}

// GetLink returns the Link field.
//
// When the provided Hook type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (h *Hook) GetLink() string {
	// return zero value if Hook type or Link field is nil
	if h == nil || h.Link == nil {
		return ""
	}

	return *h.Link
}

// GetWebhookID returns the WebhookID field.
//
// When the provided Hook type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (h *Hook) GetWebhookID() int64 {
	// return zero value if Hook type or WebhookID field is nil
	if h == nil || h.WebhookID == nil {
		return 0
	}

	return *h.WebhookID
}

// SetID sets the ID field.
//
// When the provided Hook type is nil, it
// will set nothing and immediately return.
func (h *Hook) SetID(v int64) {
	// return if Hook type is nil
	if h == nil {
		return
	}

	h.ID = &v
}

// SetRepo sets the Repo field.
//
// When the provided Hook type is nil, it
// will set nothing and immediately return.
func (h *Hook) SetRepo(v *Repo) {
	// return if Hook type is nil
	if h == nil {
		return
	}

	h.Repo = v
}

// SetBuild sets the Build field.
//
// When the provided Hook type is nil, it
// will set nothing and immediately return.
func (h *Hook) SetBuild(v *Build) {
	// return if Hook type is nil
	if h == nil {
		return
	}

	h.Build = v
}

// SetNumber sets the Number field.
//
// When the provided Hook type is nil, it
// will set nothing and immediately return.
func (h *Hook) SetNumber(v int64) {
	// return if Hook type is nil
	if h == nil {
		return
	}

	h.Number = &v
}

// SetSourceID sets the SourceID field.
//
// When the provided Hook type is nil, it
// will set nothing and immediately return.
func (h *Hook) SetSourceID(v string) {
	// return if Hook type is nil
	if h == nil {
		return
	}

	h.SourceID = &v
}

// SetCreated sets the Created field.
//
// When the provided Hook type is nil, it
// will set nothing and immediately return.
func (h *Hook) SetCreated(v int64) {
	// return if Hook type is nil
	if h == nil {
		return
	}

	h.Created = &v
}

// SetHost sets the Host field.
//
// When the provided Hook type is nil, it
// will set nothing and immediately return.
func (h *Hook) SetHost(v string) {
	// return if Hook type is nil
	if h == nil {
		return
	}

	h.Host = &v
}

// SetEvent sets the Event field.
//
// When the provided Hook type is nil, it
// will set nothing and immediately return.
func (h *Hook) SetEvent(v string) {
	// return if Hook type is nil
	if h == nil {
		return
	}

	h.Event = &v
}

// SetEventAction sets the EventAction field.
//
// When the provided Hook type is nil, it
// will set nothing and immediately return.
func (h *Hook) SetEventAction(v string) {
	// return if Hook type is nil
	if h == nil {
		return
	}

	h.EventAction = &v
}

// SetBranch sets the Branch field.
//
// When the provided Hook type is nil, it
// will set nothing and immediately return.
func (h *Hook) SetBranch(v string) {
	// return if Hook type is nil
	if h == nil {
		return
	}

	h.Branch = &v
}

// SetError sets the Error field.
//
// When the provided Hook type is nil, it
// will set nothing and immediately return.
func (h *Hook) SetError(v string) {
	// return if Hook type is nil
	if h == nil {
		return
	}

	h.Error = &v
}

// SetStatus sets the Status field.
//
// When the provided Hook type is nil, it
// will set nothing and immediately return.
func (h *Hook) SetStatus(v string) {
	// return if Hook type is nil
	if h == nil {
		return
	}

	h.Status = &v
}

// SetLink sets the Link field.
//
// When the provided Hook type is nil, it
// will set nothing and immediately return.
func (h *Hook) SetLink(v string) {
	// return if Hook type is nil
	if h == nil {
		return
	}

	h.Link = &v
}

// SetWebhookID sets the WebhookID field.
//
// When the provided Hook type is nil, it
// will set nothing and immediately return.
func (h *Hook) SetWebhookID(v int64) {
	// return if Hook type is nil
	if h == nil {
		return
	}

	h.WebhookID = &v
}

// String implements the Stringer interface for the Hook type.
func (h *Hook) String() string {
	return fmt.Sprintf(`{
  Branch: %s,
  Build: %v,
  Created: %d,
  Error: %s,
  Event: %s,
  EventAction: %s,
  Host: %s,
  ID: %d,
  Link: %s,
  Number: %d,
  Repo: %v,
  SourceID: %s,
  Status: %s,
  WebhookID: %d,
}`,
		h.GetBranch(),
		h.GetBuild(),
		h.GetCreated(),
		h.GetError(),
		h.GetEvent(),
		h.GetEventAction(),
		h.GetHost(),
		h.GetID(),
		h.GetLink(),
		h.GetNumber(),
		h.GetRepo(),
		h.GetSourceID(),
		h.GetStatus(),
		h.GetWebhookID(),
	)
}
