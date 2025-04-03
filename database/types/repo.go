// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"encoding/base64"
	"errors"

	"github.com/lib/pq"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyRepoFullName defines the error type when a
	// Repo type has an empty FullName field provided.
	ErrEmptyRepoFullName = errors.New("empty repo full_name provided")

	// ErrEmptyRepoHash defines the error type when a
	// Repo type has an empty Hash field provided.
	ErrEmptyRepoHash = errors.New("empty repo hash provided")

	// ErrEmptyRepoName defines the error type when a
	// Repo type has an empty Name field provided.
	ErrEmptyRepoName = errors.New("empty repo name provided")

	// ErrEmptyRepoOrg defines the error type when a
	// Repo type has an empty Org field provided.
	ErrEmptyRepoOrg = errors.New("empty repo org provided")

	// ErrEmptyRepoUserID defines the error type when a
	// Repo type has an empty UserID field provided.
	ErrEmptyRepoUserID = errors.New("empty repo user_id provided")

	// ErrEmptyRepoVisibility defines the error type when a
	// Repo type has an empty Visibility field provided.
	ErrEmptyRepoVisibility = errors.New("empty repo visibility provided")

	// ErrExceededTopicsLimit defines the error type when a
	// Repo type has Topics field provided that exceeds the database limit.
	ErrExceededTopicsLimit = errors.New("exceeded topics limit")
)

// Repo is the database representation of a repo.
type (
	Repo struct {
		ID              sql.NullInt64  `sql:"id"`
		UserID          sql.NullInt64  `sql:"user_id"`
		Hash            sql.NullString `sql:"hash"`
		Org             sql.NullString `sql:"org"`
		Name            sql.NullString `sql:"name"`
		FullName        sql.NullString `sql:"full_name"`
		Link            sql.NullString `sql:"link"`
		Clone           sql.NullString `sql:"clone"`
		Branch          sql.NullString `sql:"branch"`
		Topics          pq.StringArray `sql:"topics"           gorm:"type:varchar(1020)"`
		BuildLimit      sql.NullInt32  `sql:"build_limit"`
		Timeout         sql.NullInt32  `sql:"timeout"`
		Counter         sql.NullInt64  `sql:"counter"`
		Visibility      sql.NullString `sql:"visibility"`
		Private         sql.NullBool   `sql:"private"`
		Trusted         sql.NullBool   `sql:"trusted"`
		Active          sql.NullBool   `sql:"active"`
		AllowEvents     sql.NullInt64  `sql:"allow_events"`
		PipelineType    sql.NullString `sql:"pipeline_type"`
		PreviousName    sql.NullString `sql:"previous_name"`
		ApproveBuild    sql.NullString `sql:"approve_build"`
		ApprovalTimeout sql.NullInt32  `sql:"approval_timeout"`
		InstallID       sql.NullInt64  `sql:"install_id"`

		Owner User `gorm:"foreignKey:UserID"`
	}
)

// Decrypt will manipulate the existing repo hash by
// base64 decoding that value. Then, a AES-256 cipher
// block is created from the encryption key in order to
// decrypt the base64 decoded secret value.
func (r *Repo) Decrypt(key string) error {
	// base64 decode the encrypted repo hash
	decoded, err := base64.StdEncoding.DecodeString(r.Hash.String)
	if err != nil {
		return err
	}

	// decrypt the base64 decoded repo hash
	decrypted, err := util.Decrypt(key, decoded)
	if err != nil {
		return err
	}

	// set the decrypted repo hash
	r.Hash = sql.NullString{
		String: string(decrypted),
		Valid:  true,
	}

	// decrypt owner
	// Note: In UpdateRepo() (database/repo/update.go), the incoming API repo object
	// is cast to a database repo object. The owner object isn't set in this process
	// resulting in a zero value for the owner object. A check is performed here
	// before decrypting to prevent "unable to decrypt repo..." errors.
	if r.Owner.ID.Valid {
		err = r.Owner.Decrypt(key)
		if err != nil {
			return err
		}
	}

	return nil
}

// Encrypt will manipulate the existing repo hash by
// creating a AES-256 cipher block from the encryption
// key in order to encrypt the repo hash. Then, the
// repo hash is base64 encoded for transport across
// network boundaries.
func (r *Repo) Encrypt(key string) error {
	// encrypt the repo hash
	encrypted, err := util.Encrypt(key, []byte(r.Hash.String))
	if err != nil {
		return err
	}

	// base64 encode the encrypted repo hash to make it network safe
	r.Hash = sql.NullString{
		String: base64.StdEncoding.EncodeToString(encrypted),
		Valid:  true,
	}

	return nil
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the Repo type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (r *Repo) Nullify() *Repo {
	if r == nil {
		return nil
	}

	// check if the ID field should be false
	if r.ID.Int64 == 0 {
		r.ID.Valid = false
	}

	// check if the UserID field should be false
	if r.UserID.Int64 == 0 {
		r.UserID.Valid = false
	}

	// check if the Hash field should be false
	if len(r.Hash.String) == 0 {
		r.Hash.Valid = false
	}

	// check if the Org field should be false
	if len(r.Org.String) == 0 {
		r.Org.Valid = false
	}

	// check if the Name field should be false
	if len(r.Name.String) == 0 {
		r.Name.Valid = false
	}

	// check if the FullName field should be false
	if len(r.FullName.String) == 0 {
		r.FullName.Valid = false
	}

	// check if the Link field should be false
	if len(r.Link.String) == 0 {
		r.Link.Valid = false
	}

	// check if the Clone field should be false
	if len(r.Clone.String) == 0 {
		r.Clone.Valid = false
	}

	// check if the Branch field should be false
	if len(r.Branch.String) == 0 {
		r.Branch.Valid = false
	}

	// check if the BuildLimit field should be false
	if r.BuildLimit.Int32 == 0 {
		r.BuildLimit.Valid = false
	}

	// check if the Timeout field should be false
	if r.Timeout.Int32 == 0 {
		r.Timeout.Valid = false
	}

	// check if the AllowEvents field should be false
	if r.AllowEvents.Int64 == 0 {
		r.AllowEvents.Valid = false
	}

	// check if the Visibility field should be false
	if len(r.Visibility.String) == 0 {
		r.Visibility.Valid = false
	}

	// check if the PipelineType field should be false
	if len(r.PipelineType.String) == 0 {
		r.PipelineType.Valid = false
	}

	// check if the PreviousName field should be false
	if len(r.PreviousName.String) == 0 {
		r.PreviousName.Valid = false
	}

	// check if the ApproveForkBuild field should be false
	if len(r.ApproveBuild.String) == 0 {
		r.ApproveBuild.Valid = false
	}

	// check if the ApprovalTimeout field should be false
	if r.ApprovalTimeout.Int32 == 0 {
		r.ApprovalTimeout.Valid = false
	}

	return r
}

// ToAPI converts the Repo type
// to an API Repo type.
func (r *Repo) ToAPI() *api.Repo {
	repo := new(api.Repo)

	var owner *api.User
	if r.Owner.ID.Valid {
		owner = r.Owner.ToAPI()
	} else {
		owner = new(api.User)
		owner.SetID(r.UserID.Int64)
	}

	repo.SetID(r.ID.Int64)
	repo.SetOwner(owner.Crop())
	repo.SetHash(r.Hash.String)
	repo.SetOrg(r.Org.String)
	repo.SetName(r.Name.String)
	repo.SetFullName(r.FullName.String)
	repo.SetLink(r.Link.String)
	repo.SetClone(r.Clone.String)
	repo.SetBranch(r.Branch.String)
	repo.SetTopics(r.Topics)
	repo.SetBuildLimit(r.BuildLimit.Int32)
	repo.SetTimeout(r.Timeout.Int32)
	repo.SetCounter(r.Counter.Int64)
	repo.SetVisibility(r.Visibility.String)
	repo.SetPrivate(r.Private.Bool)
	repo.SetTrusted(r.Trusted.Bool)
	repo.SetActive(r.Active.Bool)
	repo.SetAllowEvents(api.NewEventsFromMask(r.AllowEvents.Int64))
	repo.SetPipelineType(r.PipelineType.String)
	repo.SetPreviousName(r.PreviousName.String)
	repo.SetApproveBuild(r.ApproveBuild.String)
	repo.SetApprovalTimeout(r.ApprovalTimeout.Int32)
	repo.SetInstallID(r.InstallID.Int64)

	return repo
}

// Validate verifies the necessary fields for
// the Repo type are populated correctly.
func (r *Repo) Validate() error {
	// verify the UserID field is populated
	if r.UserID.Int64 <= 0 {
		return ErrEmptyRepoUserID
	}

	// verify the Hash field is populated
	if len(r.Hash.String) == 0 {
		return ErrEmptyRepoHash
	}

	// verify the Org field is populated
	if len(r.Org.String) == 0 {
		return ErrEmptyRepoOrg
	}

	// verify the Name field is populated
	if len(r.Name.String) == 0 {
		return ErrEmptyRepoName
	}

	// verify the FullName field is populated
	if len(r.FullName.String) == 0 {
		return ErrEmptyRepoFullName
	}

	// verify the Visibility field is populated
	if len(r.Visibility.String) == 0 {
		return ErrEmptyRepoVisibility
	}

	// calculate total size of favorites while sanitizing entries
	total := 0

	for i, t := range r.Topics {
		r.Topics[i] = util.Sanitize(t)
		total += len(t)
	}

	// verify the Favorites field is within the database constraints
	// len is to factor in number of comma separators included in the database field,
	// removing 1 due to the last item not having an appended comma
	if (total + len(r.Topics) - 1) > constants.TopicsMaxSize {
		return ErrExceededTopicsLimit
	}

	// ensure that all Repo string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	r.Org = sql.NullString{String: util.Sanitize(r.Org.String), Valid: r.Org.Valid}
	r.Name = sql.NullString{String: util.Sanitize(r.Name.String), Valid: r.Name.Valid}
	r.FullName = sql.NullString{String: util.Sanitize(r.FullName.String), Valid: r.FullName.Valid}
	r.Link = sql.NullString{String: util.Sanitize(r.Link.String), Valid: r.Link.Valid}
	r.Clone = sql.NullString{String: util.Sanitize(r.Clone.String), Valid: r.Clone.Valid}
	r.Branch = sql.NullString{String: util.Sanitize(r.Branch.String), Valid: r.Branch.Valid}
	r.Visibility = sql.NullString{String: util.Sanitize(r.Visibility.String), Valid: r.Visibility.Valid}
	r.PipelineType = sql.NullString{String: util.Sanitize(r.PipelineType.String), Valid: r.PipelineType.Valid}

	return nil
}

// RepoFromAPI converts the API Repo type
// to a database repo type.
func RepoFromAPI(r *api.Repo) *Repo {
	repo := &Repo{
		ID:              sql.NullInt64{Int64: r.GetID(), Valid: true},
		UserID:          sql.NullInt64{Int64: r.GetOwner().GetID(), Valid: true},
		Hash:            sql.NullString{String: r.GetHash(), Valid: true},
		Org:             sql.NullString{String: r.GetOrg(), Valid: true},
		Name:            sql.NullString{String: r.GetName(), Valid: true},
		FullName:        sql.NullString{String: r.GetFullName(), Valid: true},
		Link:            sql.NullString{String: r.GetLink(), Valid: true},
		Clone:           sql.NullString{String: r.GetClone(), Valid: true},
		Branch:          sql.NullString{String: r.GetBranch(), Valid: true},
		Topics:          pq.StringArray(r.GetTopics()),
		BuildLimit:      sql.NullInt32{Int32: r.GetBuildLimit(), Valid: true},
		Timeout:         sql.NullInt32{Int32: r.GetTimeout(), Valid: true},
		Counter:         sql.NullInt64{Int64: r.GetCounter(), Valid: true},
		Visibility:      sql.NullString{String: r.GetVisibility(), Valid: true},
		Private:         sql.NullBool{Bool: r.GetPrivate(), Valid: true},
		Trusted:         sql.NullBool{Bool: r.GetTrusted(), Valid: true},
		Active:          sql.NullBool{Bool: r.GetActive(), Valid: true},
		AllowEvents:     sql.NullInt64{Int64: r.GetAllowEvents().ToDatabase(), Valid: true},
		PipelineType:    sql.NullString{String: r.GetPipelineType(), Valid: true},
		PreviousName:    sql.NullString{String: r.GetPreviousName(), Valid: true},
		ApproveBuild:    sql.NullString{String: r.GetApproveBuild(), Valid: true},
		ApprovalTimeout: sql.NullInt32{Int32: r.GetApprovalTimeout(), Valid: true},
		InstallID:       sql.NullInt64{Int64: r.GetInstallID(), Valid: true},
	}

	return repo.Nullify()
}
