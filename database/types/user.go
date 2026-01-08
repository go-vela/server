// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"regexp"

	"github.com/lib/pq"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/util"
)

var (
	// userRegex defines the regex pattern for validating
	// the Name field for the User type.
	userRegex = regexp.MustCompile("^[a-zA-Z0-9_-]{0,38}$")

	// ErrEmptyUserName defines the error type when a
	// User type has an empty Name field provided.
	ErrEmptyUserName = errors.New("empty user name provided")

	// ErrEmptyUserRefreshToken defines the error type when a
	// User type has an empty RefreshToken field provided.
	ErrEmptyUserRefreshToken = errors.New("empty user refresh token provided")

	// ErrEmptyUserToken defines the error type when a
	// User type has an empty Token field provided.
	ErrEmptyUserToken = errors.New("empty user token provided")

	// ErrInvalidUserName defines the error type when a
	// User type has an invalid Name field provided.
	ErrInvalidUserName = errors.New("invalid user name provided")

	// ErrExceededFavoritesLimit defines the error type when a
	// User type has Favorites field provided that exceeds the database limit.
	ErrExceededFavoritesLimit = errors.New("exceeded favorites limit")

	// ErrExceededDashboardsLimit defines the error type when a
	// User type has Dashboards field provided that exceeds the database limit.
	ErrExceededDashboardsLimit = errors.New("exceeded dashboards limit")
)

// User is the database representation of a user.
type User struct {
	ID           sql.NullInt64  `sql:"id"`
	Name         sql.NullString `sql:"name"`
	RefreshToken sql.NullString `sql:"refresh_token"`
	Token        sql.NullString `sql:"token"`
	Favorites    pq.StringArray `sql:"favorites"     gorm:"type:varchar(5000)"`
	Active       sql.NullBool   `sql:"active"`
	Admin        sql.NullBool   `sql:"admin"`
	Dashboards   pq.StringArray `sql:"dashboards"    gorm:"type:varchar(5000)"`
}

// Decrypt will manipulate the existing user tokens by
// base64 decoding them. Then, a AES-256 cipher
// block is created from the encryption key in order to
// decrypt the base64 decoded user tokens.
func (u *User) Decrypt(key string) error {
	// base64 decode the encrypted user token
	decoded, err := base64.StdEncoding.DecodeString(u.Token.String)
	if err != nil {
		return err
	}

	// decrypt the base64 decoded user token
	decrypted, err := util.Decrypt(key, decoded)
	if err != nil {
		return err
	}

	// set the decrypted user token
	u.Token = sql.NullString{
		String: string(decrypted),
		Valid:  true,
	}

	// base64 decode the encrypted user refresh token
	decoded, err = base64.StdEncoding.DecodeString(u.RefreshToken.String)
	if err != nil {
		return err
	}

	// decrypt the base64 decoded user refresh token
	decrypted, err = util.Decrypt(key, decoded)
	if err != nil {
		return err
	}

	// set the decrypted user refresh token
	u.RefreshToken = sql.NullString{
		String: string(decrypted),
		Valid:  true,
	}

	return nil
}

// Encrypt will manipulate the existing user tokens by
// creating a AES-256 cipher block from the encryption
// key in order to encrypt the user tokens. Then, the
// user tokens are base64 encoded for transport across
// network boundaries.
func (u *User) Encrypt(key string) error {
	// encrypt the user token
	encrypted, err := util.Encrypt(key, []byte(u.Token.String))
	if err != nil {
		return err
	}

	// base64 encode the encrypted user token to make it network safe
	u.Token = sql.NullString{
		String: base64.StdEncoding.EncodeToString(encrypted),
		Valid:  true,
	}

	// encrypt the user refresh token
	encrypted, err = util.Encrypt(key, []byte(u.RefreshToken.String))
	if err != nil {
		return err
	}

	// base64 encode the encrypted user refresh token to make it network safe
	u.RefreshToken = sql.NullString{
		String: base64.StdEncoding.EncodeToString(encrypted),
		Valid:  true,
	}

	return nil
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the User type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (u *User) Nullify() *User {
	if u == nil {
		return nil
	}

	// check if the ID field should be false
	if u.ID.Int64 == 0 {
		u.ID.Valid = false
	}

	// check if the Name field should be false
	if len(u.Name.String) == 0 {
		u.Name.Valid = false
	}

	// check if the RefreshToken field should be false
	if len(u.RefreshToken.String) == 0 {
		u.RefreshToken.Valid = false
	}

	// check if the Token field should be false
	if len(u.Token.String) == 0 {
		u.Token.Valid = false
	}

	return u
}

// ToAPI converts the User type
// to an API User type.
func (u *User) ToAPI() *api.User {
	user := new(api.User)

	user.SetID(u.ID.Int64)
	user.SetName(u.Name.String)
	user.SetRefreshToken(u.RefreshToken.String)
	user.SetToken(u.Token.String)
	user.SetActive(u.Active.Bool)
	user.SetAdmin(u.Admin.Bool)
	user.SetFavorites(u.Favorites)
	user.SetDashboards(u.Dashboards)

	return user
}

// Validate verifies the necessary fields for
// the User type are populated correctly.
func (u *User) Validate() error {
	// verify the Name field is populated
	if len(u.Name.String) == 0 {
		return ErrEmptyUserName
	}

	// verify the Token field is populated
	if len(u.Token.String) == 0 {
		return ErrEmptyUserToken
	}

	// verify the Name field is valid
	if !userRegex.MatchString(u.Name.String) {
		return ErrInvalidUserName
	}

	// calculate total size of favorites as it would be stored in PostgreSQL
	// PostgreSQL stores arrays as: {value1,value2,value3}
	// GitHub repo names (owner/repo) only contain alphanumeric, hyphens, underscores, periods, and slashes
	// so they never require quoting or escaping in PostgreSQL arrays
	total := 2 // account for opening and closing curly braces: { }

	for i, f := range u.Favorites {
		total += len(f)

		// add comma separator (except for last item)
		if i < len(u.Favorites)-1 {
			total++
		}
	}

	// verify the Favorites field is within the database constraints
	if total > constants.FavoritesMaxSize {
		return ErrExceededFavoritesLimit
	}

	// validate number of dashboards
	if len(u.Dashboards) > constants.UserDashboardLimit {
		return ErrExceededDashboardsLimit
	}

	// ensure that all User string fields
	// that can be returned as JSON are sanitized
	// to avoid unsafe HTML content
	u.Name = sql.NullString{String: util.Sanitize(u.Name.String), Valid: u.Name.Valid}

	// ensure that all Favorites are sanitized
	// to avoid unsafe HTML content
	for i, v := range u.Favorites {
		u.Favorites[i] = util.Sanitize(v)
	}

	return nil
}

// UserFromAPI converts the API User type
// to a database User type.
func UserFromAPI(u *api.User) *User {
	user := &User{
		ID:           sql.NullInt64{Int64: u.GetID(), Valid: true},
		Name:         sql.NullString{String: u.GetName(), Valid: true},
		RefreshToken: sql.NullString{String: u.GetRefreshToken(), Valid: true},
		Token:        sql.NullString{String: u.GetToken(), Valid: true},
		Active:       sql.NullBool{Bool: u.GetActive(), Valid: true},
		Admin:        sql.NullBool{Bool: u.GetAdmin(), Valid: true},
		Favorites:    pq.StringArray(u.GetFavorites()),
		Dashboards:   pq.StringArray(u.GetDashboards()),
	}

	return user.Nullify()
}
