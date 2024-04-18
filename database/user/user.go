// SPDX-License-Identifier: Apache-2.0

package user

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types/constants"
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

type (
	// config represents the settings required to create the engine that implements the UserInterface interface.
	config struct {
		// specifies the encryption key to use for the User engine
		EncryptionKey string
		// specifies to skip creating tables and indexes for the User engine
		SkipCreation bool
	}

	// engine represents the user functionality that implements the UserInterface interface.
	engine struct {
		// engine configuration settings used in user functions
		config *config

		ctx context.Context

		// gorm.io/gorm database client used in user functions
		//
		// https://pkg.go.dev/gorm.io/gorm#DB
		client *gorm.DB

		// sirupsen/logrus logger used in user functions
		//
		// https://pkg.go.dev/github.com/sirupsen/logrus#Entry
		logger *logrus.Entry
	}

	// User is the database representation of a user.
	User struct {
		ID           sql.NullInt64  `sql:"id"`
		Name         sql.NullString `sql:"name"`
		RefreshToken sql.NullString `sql:"refresh_token"`
		Token        sql.NullString `sql:"token"`
		Favorites    pq.StringArray `sql:"favorites" gorm:"type:varchar(5000)"`
		Active       sql.NullBool   `sql:"active"`
		Admin        sql.NullBool   `sql:"admin"`
		Dashboards   pq.StringArray `sql:"dashboards" gorm:"type:varchar(5000)"`
	}
)

// New creates and returns a Vela service for integrating with users in the database.
//
//nolint:revive // ignore returning unexported engine
func New(opts ...EngineOpt) (*engine, error) {
	// create new User engine
	e := new(engine)

	// create new fields
	e.client = new(gorm.DB)
	e.config = new(config)
	e.logger = new(logrus.Entry)

	// apply all provided configuration options
	for _, opt := range opts {
		err := opt(e)
		if err != nil {
			return nil, err
		}
	}

	// check if we should skip creating user database objects
	if e.config.SkipCreation {
		e.logger.Warning("skipping creation of users table and indexes in the database")

		return e, nil
	}

	// create the users table
	err := e.CreateUserTable(e.ctx, e.client.Config.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("unable to create %s table: %w", constants.TableUser, err)
	}

	// create the indexes for the users table
	err = e.CreateUserIndexes(e.ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to create indexes for %s table: %w", constants.TableUser, err)
	}

	return e, nil
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

	// calculate totalFavorites size of favorites
	totalFavorites := 0
	for _, f := range u.Favorites {
		totalFavorites += len(f)
	}

	// verify the Favorites field is within the database constraints
	// len is to factor in number of comma separators included in the database field,
	// removing 1 due to the last item not having an appended comma
	if (totalFavorites + len(u.Favorites) - 1) > constants.FavoritesMaxSize {
		return ErrExceededFavoritesLimit
	}

	// calculate totalDashboards size of dashboards
	totalDashboards := 0
	for _, d := range u.Dashboards {
		totalDashboards += len(d)
	}

	// verify the Dashboards field is within the database constraints
	// len is to factor in number of comma separators included in the database field,
	// removing 1 due to the last item not having an appended comma
	if (totalDashboards + len(u.Dashboards) - 1) > constants.FavoritesMaxSize {
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

// FromAPI converts the API User type
// to a database User type.
func FromAPI(u *api.User) *User {
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
