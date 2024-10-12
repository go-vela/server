// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"

	"github.com/go-vela/server/compiler/types/raw"
)

// Deployment is the API representation of a deployment.
//
// swagger:model Deployment
type Deployment struct {
	ID          *int64              `json:"id,omitempty"`
	Number      *int64              `json:"number,omitempty"`
	Repo        *Repo               `json:"repo,omitempty"`
	URL         *string             `json:"url,omitempty"`
	Commit      *string             `json:"commit,omitempty"`
	Ref         *string             `json:"ref,omitempty"`
	Task        *string             `json:"task,omitempty"`
	Target      *string             `json:"target,omitempty"`
	Description *string             `json:"description,omitempty"`
	Payload     *raw.StringSliceMap `json:"payload,omitempty"`
	CreatedAt   *int64              `json:"created_at,omitempty"`
	CreatedBy   *string             `json:"created_by,omitempty"`
	Builds      []*Build            `json:"builds,omitempty"`
}

// GetID returns the ID field.
//
// When the provided Deployment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Deployment) GetID() int64 {
	// return zero value if Deployment type or ID field is nil
	if d == nil || d.ID == nil {
		return 0
	}

	return *d.ID
}

// GetNumber returns the Number field.
//
// When the provided Deployment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Deployment) GetNumber() int64 {
	// return zero value if Deployment type or ID field is nil
	if d == nil || d.Number == nil {
		return 0
	}

	return *d.Number
}

// GetRepo returns the Repo field.
//
// When the provided Deployment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Deployment) GetRepo() *Repo {
	// return zero value if Deployment type or Repo field is nil
	if d == nil || d.Repo == nil {
		return new(Repo)
	}

	return d.Repo
}

// GetURL returns the URL field.
//
// When the provided Deployment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Deployment) GetURL() string {
	// return zero value if Deployment type or URL field is nil
	if d == nil || d.URL == nil {
		return ""
	}

	return *d.URL
}

// GetCommit returns the Commit field.
//
// When the provided Deployment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Deployment) GetCommit() string {
	// return zero value if Deployment type or Commit field is nil
	if d == nil || d.Commit == nil {
		return ""
	}

	return *d.Commit
}

// GetRef returns the Ref field.
//
// When the provided Deployment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Deployment) GetRef() string {
	// return zero value if Deployment type or Ref field is nil
	if d == nil || d.Ref == nil {
		return ""
	}

	return *d.Ref
}

// GetTask returns the Task field.
//
// When the provided Deployment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Deployment) GetTask() string {
	// return zero value if Deployment type or Task field is nil
	if d == nil || d.Task == nil {
		return ""
	}

	return *d.Task
}

// GetTarget returns the Target field.
//
// When the provided Deployment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Deployment) GetTarget() string {
	// return zero value if Deployment type or Target field is nil
	if d == nil || d.Target == nil {
		return ""
	}

	return *d.Target
}

// GetDescription returns the Description field.
//
// When the provided Deployment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Deployment) GetDescription() string {
	// return zero value if Deployment type or Description field is nil
	if d == nil || d.Description == nil {
		return ""
	}

	return *d.Description
}

// GetPayload returns the Payload field.
//
// When the provided Deployment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Deployment) GetPayload() map[string]string {
	// return zero value if Deployment type or Description field is nil
	if d == nil || d.Payload == nil {
		return map[string]string{}
	}

	return *d.Payload
}

// GetCreatedAt returns the CreatedAt field.
//
// When the provided Deployment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Deployment) GetCreatedAt() int64 {
	// return zero value if Deployment type or CreatedAt field is nil
	if d == nil || d.CreatedAt == nil {
		return 0
	}

	return *d.CreatedAt
}

// GetCreatedBy returns the CreatedBy field.
//
// When the provided Deployment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Deployment) GetCreatedBy() string {
	// return zero value if Deployment type or CreatedBy field is nil
	if d == nil || d.CreatedBy == nil {
		return ""
	}

	return *d.CreatedBy
}

// GetBuilds returns the Builds field.
//
// When the provided Deployment type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (d *Deployment) GetBuilds() []*Build {
	if d == nil || d.Builds == nil {
		return []*Build{}
	}

	return d.Builds
}

// SetID sets the ID field.
//
// When the provided Deployment type is nil, it
// will set nothing and immediately return.
func (d *Deployment) SetID(v int64) {
	// return if Deployment type is nil
	if d == nil {
		return
	}

	d.ID = &v
}

// SetNumber sets the Number field.
//
// When the provided Deployment type is nil, it
// will set nothing and immediately return.
func (d *Deployment) SetNumber(v int64) {
	// return if Deployment type is nil
	if d == nil {
		return
	}

	d.Number = &v
}

// SetRepo sets the Repo field.
//
// When the provided Deployment type is nil, it
// will set nothing and immediately return.
func (d *Deployment) SetRepo(v *Repo) {
	// return if Deployment type is nil
	if d == nil {
		return
	}

	d.Repo = v
}

// SetURL sets the URL field.
//
// When the provided Deployment type is nil, it
// will set nothing and immediately return.
func (d *Deployment) SetURL(v string) {
	// return if Deployment type is nil
	if d == nil {
		return
	}

	d.URL = &v
}

// SetCommit sets the Commit field.
//
// When the provided Deployment type is nil, it
// will set nothing and immediately return.
func (d *Deployment) SetCommit(v string) {
	// return if Deployment type is nil
	if d == nil {
		return
	}

	d.Commit = &v
}

// SetRef sets the Ref field.
//
// When the provided Deployment type is nil, it
// will set nothing and immediately return.
func (d *Deployment) SetRef(v string) {
	// return if Deployment type is nil
	if d == nil {
		return
	}

	d.Ref = &v
}

// SetTask sets the Task field.
//
// When the provided Deployment type is nil, it
// will set nothing and immediately return.
func (d *Deployment) SetTask(v string) {
	// return if Deployment type is nil
	if d == nil {
		return
	}

	d.Task = &v
}

// SetTarget sets the Target field.
//
// When the provided Deployment type is nil, it
// will set nothing and immediately return.
func (d *Deployment) SetTarget(v string) {
	// return if Deployment type is nil
	if d == nil {
		return
	}

	d.Target = &v
}

// SetDescription sets the Description field.
//
// When the provided Deployment type is nil, it
// will set nothing and immediately return.
func (d *Deployment) SetDescription(v string) {
	// return if Deployment type is nil
	if d == nil {
		return
	}

	d.Description = &v
}

// SetPayload sets the Payload field.
//
// When the provided Deployment type is nil, it
// will set nothing and immediately return.
func (d *Deployment) SetPayload(v raw.StringSliceMap) {
	// return if Deployment type is nil
	if d == nil {
		return
	}

	d.Payload = &v
}

// SetCreatedAt sets the CreatedAt field.
//
// When the provided Deployment type is nil, it
// will set nothing and immediately return.
func (d *Deployment) SetCreatedAt(v int64) {
	// return if Deployment type is nil
	if d == nil {
		return
	}

	d.CreatedAt = &v
}

// SetCreatedBy sets the CreatedBy field.
//
// When the provided Deployment type is nil, it
// will set nothing and immediately return.
func (d *Deployment) SetCreatedBy(v string) {
	// return if Deployment type is nil
	if d == nil {
		return
	}

	d.CreatedBy = &v
}

// SetBuilds sets the Builds field.
//
// When the provided Deployment type is nil, it
// will set nothing and immediately return.
func (d *Deployment) SetBuilds(b []*Build) {
	// return if Deployment type is nil
	if d == nil {
		return
	}

	d.Builds = b
}

// String implements the Stringer interface for the Deployment type.
func (d *Deployment) String() string {
	return fmt.Sprintf(`{
  Commit: %s,
  CreatedAt: %d,
  CreatedBy: %s,
  Description: %s,
  ID: %d,
  Number: %d,
  Ref: %s,
  Repo: %v,
  Target: %s,
  Task: %s,
  URL: %s,
  Payload: %s,
  Builds: %d,
}`,
		d.GetCommit(),
		d.GetCreatedAt(),
		d.GetCreatedBy(),
		d.GetDescription(),
		d.GetID(),
		d.GetNumber(),
		d.GetRef(),
		d.GetRepo(),
		d.GetTarget(),
		d.GetTask(),
		d.GetURL(),
		d.GetPayload(),
		len(d.GetBuilds()),
	)
}
