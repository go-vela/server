// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyDeploymentNumber defines the error type when a
	// Deployment type has an empty Number field provided.
	ErrEmptyDeploymentNumber = errors.New("empty deployment number provided")

	// ErrEmptyDeploymentRepoID defines the error type when a
	// Deployment type has an empty RepoID field provided.
	ErrEmptyDeploymentRepoID = errors.New("empty deployment repo_id provided")
)

// Deployment is the database representation of a deployment for a repo.
type Deployment struct {
	ID          sql.NullInt64      `sql:"id"`
	Number      sql.NullInt64      `sql:"number"`
	RepoID      sql.NullInt64      `sql:"repo_id"`
	URL         sql.NullString     `sql:"url"`
	Commit      sql.NullString     `sql:"commit"`
	Ref         sql.NullString     `sql:"ref"`
	Task        sql.NullString     `sql:"task"`
	Target      sql.NullString     `sql:"target"`
	Description sql.NullString     `sql:"description"`
	Payload     raw.StringSliceMap `sql:"payload"`
	CreatedAt   sql.NullInt64      `sql:"created_at"`
	CreatedBy   sql.NullString     `sql:"created_by"`
	Builds      pq.StringArray     `sql:"builds"      gorm:"type:varchar(50)"`

	Repo Repo `gorm:"foreignKey:RepoID"`
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the Deployment type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (d *Deployment) Nullify() *Deployment {
	if d == nil {
		return nil
	}

	// check if the ID field should be false
	if d.ID.Int64 == 0 {
		d.ID.Valid = false
	}

	// check if the Number field should be false
	if d.Number.Int64 == 0 {
		d.Number.Valid = false
	}

	// check if the RepoID field should be false
	if d.RepoID.Int64 == 0 {
		d.RepoID.Valid = false
	}

	// check if the URL field should be false
	if len(d.URL.String) == 0 {
		d.URL.Valid = false
	}

	// check if the Commit field should be false
	if len(d.Commit.String) == 0 {
		d.Commit.Valid = false
	}

	// check if the Ref field should be false
	if len(d.Ref.String) == 0 {
		d.Ref.Valid = false
	}

	// check if the Task field should be false
	if len(d.Task.String) == 0 {
		d.Task.Valid = false
	}

	// check if the Target field should be false
	if len(d.Target.String) == 0 {
		d.Target.Valid = false
	}

	// check if the Description field should be false
	if len(d.Description.String) == 0 {
		d.Description.Valid = false
	}

	// check if the CreatedAt field should be false
	if d.CreatedAt.Int64 == 0 {
		d.CreatedAt.Valid = false
	}

	// check if the CreatedBy field should be false
	if len(d.CreatedBy.String) == 0 {
		d.CreatedBy.Valid = false
	}

	return d
}

// ToAPI converts the Deployment type
// to the API Deployment type.
func (d *Deployment) ToAPI(builds []*api.Build) *api.Deployment {
	deployment := new(api.Deployment)

	deployment.SetID(d.ID.Int64)
	deployment.SetNumber(d.Number.Int64)
	deployment.SetRepo(d.Repo.ToAPI())
	deployment.SetURL(d.URL.String)
	deployment.SetCommit(d.Commit.String)
	deployment.SetRef(d.Ref.String)
	deployment.SetTask(d.Task.String)
	deployment.SetTarget(d.Target.String)
	deployment.SetDescription(d.Description.String)
	deployment.SetPayload(d.Payload)
	deployment.SetCreatedAt(d.CreatedAt.Int64)
	deployment.SetCreatedBy(d.CreatedBy.String)

	if len(builds) > 0 {
		deployment.SetBuilds(builds)
	}

	return deployment
}

// Validate verifies the necessary fields for
// the Deployment type are populated correctly.
func (d *Deployment) Validate() error {
	// verify the RepoID field is populated
	if d.RepoID.Int64 <= 0 {
		return ErrEmptyDeploymentRepoID
	}

	// verify the Number field is populated
	if d.Number.Int64 <= 0 {
		return ErrEmptyDeploymentNumber
	}

	// ensure that all Deployment string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	d.Commit = sql.NullString{String: util.Sanitize(d.Commit.String), Valid: d.Commit.Valid}
	d.Ref = sql.NullString{String: util.Sanitize(d.Ref.String), Valid: d.Ref.Valid}
	d.Task = sql.NullString{String: util.Sanitize(d.Task.String), Valid: d.Task.Valid}
	d.Target = sql.NullString{String: util.Sanitize(d.Target.String), Valid: d.Target.Valid}
	d.Description = sql.NullString{String: util.Sanitize(d.Description.String), Valid: d.Description.Valid}

	// calculate total size of builds
	total := 0
	for _, b := range d.Builds {
		total += len(b)
	}

	// verify the Builds field is within the database constraints and evict if not
	// len is to factor in number of comma separators included in the database field,
	// removing 1 due to the last item not having an appended comma
	if diff := (total + len(d.Builds) - 1) - constants.DeployBuildsMaxSize; diff > 0 {
		for diff > 0 {
			diff = diff - (len(d.Builds[0]) + 1)
			d.Builds = d.Builds[1:]
		}
	}

	return nil
}

// DeploymentFromAPI converts the API Deployment type
// to a database Deployment type.
func DeploymentFromAPI(d *api.Deployment) *Deployment {
	buildIDs := []string{}
	for _, build := range d.GetBuilds() {
		buildIDs = append(buildIDs, fmt.Sprint(build.GetID()))
	}

	deployment := &Deployment{
		ID:          sql.NullInt64{Int64: d.GetID(), Valid: true},
		Number:      sql.NullInt64{Int64: d.GetNumber(), Valid: true},
		RepoID:      sql.NullInt64{Int64: d.GetRepo().GetID(), Valid: true},
		URL:         sql.NullString{String: d.GetURL(), Valid: true},
		Commit:      sql.NullString{String: d.GetCommit(), Valid: true},
		Ref:         sql.NullString{String: d.GetRef(), Valid: true},
		Task:        sql.NullString{String: d.GetTask(), Valid: true},
		Target:      sql.NullString{String: d.GetTarget(), Valid: true},
		Description: sql.NullString{String: d.GetDescription(), Valid: true},
		Payload:     d.GetPayload(),
		CreatedAt:   sql.NullInt64{Int64: d.GetCreatedAt(), Valid: true},
		CreatedBy:   sql.NullString{String: d.GetCreatedBy(), Valid: true},
		Builds:      buildIDs,
	}

	return deployment.Nullify()
}
