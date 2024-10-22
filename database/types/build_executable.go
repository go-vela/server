// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"encoding/base64"
	"errors"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/util"
)

var (
	// ErrEmptyBuildExecutableBuildID defines the error type when a
	// BuildExecutable type has an empty BuildID field provided.
	ErrEmptyBuildExecutableBuildID = errors.New("empty build_executable build_id provided")
)

// BuildExecutable is the database representation of a BuildExecutable.
type BuildExecutable struct {
	ID      sql.NullInt64 `sql:"id"`
	BuildID sql.NullInt64 `sql:"build_id"`
	Data    []byte        `sql:"data"`
}

// Compress will manipulate the existing data for the
// BuildExecutable by compressing that data. This produces
// a significantly smaller amount of data that is
// stored in the system.
func (b *BuildExecutable) Compress(level int) error {
	// compress the database BuildExecutable data
	data, err := util.Compress(level, b.Data)
	if err != nil {
		return err
	}

	// overwrite database BuildExecutable data with compressed BuildExecutable data
	b.Data = data

	return nil
}

// Decompress will manipulate the existing data for the
// BuildExecutable by decompressing that data. This allows us
// to have a significantly smaller amount of data that
// is stored in the system.
func (b *BuildExecutable) Decompress() error {
	// decompress the database BuildExecutable data
	data, err := util.Decompress(b.Data)
	if err != nil {
		return err
	}

	// overwrite compressed BuildExecutable data with decompressed BuildExecutable data
	b.Data = data

	return nil
}

// Decrypt will manipulate the existing executable data by
// base64 decoding that value. Then, a AES-256 cipher
// block is created from the encryption key in order to
// decrypt the base64 decoded secret value.
func (b *BuildExecutable) Decrypt(key string) error {
	dst := make([]byte, base64.StdEncoding.DecodedLen(len(b.Data)))

	// base64 decode the encrypted repo hash
	n, err := base64.StdEncoding.Decode(dst, b.Data)
	if err != nil {
		return err
	}

	dst = dst[:n]

	// decrypt the base64 decoded executable data
	decrypted, err := util.Decrypt(key, dst)
	if err != nil {
		return err
	}

	// set the decrypted executable
	b.Data = decrypted

	return nil
}

// Encrypt will manipulate the existing build executable by
// creating a AES-256 cipher block from the encryption
// key in order to encrypt the build executable. Then, the
// build executable is base64 encoded for transport across
// network boundaries.
func (b *BuildExecutable) Encrypt(key string) error {
	// encrypt the executable data
	encrypted, err := util.Encrypt(key, b.Data)
	if err != nil {
		return err
	}

	// base64 encode the encrypted executable to make it network safe
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(encrypted)))
	base64.StdEncoding.Encode(dst, encrypted)

	b.Data = dst

	return nil
}

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the BuildExecutable type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (b *BuildExecutable) Nullify() *BuildExecutable {
	if b == nil {
		return nil
	}

	// check if the ID field should be false
	if b.ID.Int64 == 0 {
		b.ID.Valid = false
	}

	// check if the BuildID field should be false
	if b.BuildID.Int64 == 0 {
		b.BuildID.Valid = false
	}

	return b
}

// ToAPI converts the BuildExecutable type
// to a API BuildExecutable type.
func (b *BuildExecutable) ToAPI() *api.BuildExecutable {
	buildExecutable := new(api.BuildExecutable)

	buildExecutable.SetID(b.ID.Int64)
	buildExecutable.SetBuildID(b.BuildID.Int64)
	buildExecutable.SetData(b.Data)

	return buildExecutable
}

// Validate verifies the necessary fields for
// the BuildExecutable type are populated correctly.
func (b *BuildExecutable) Validate() error {
	// verify the BuildID field is populated
	if b.BuildID.Int64 <= 0 {
		return ErrEmptyBuildExecutableBuildID
	}

	return nil
}

// BuildExecutableFromAPI converts the API BuildExecutable type
// to a database BuildExecutable type.
func BuildExecutableFromAPI(c *api.BuildExecutable) *BuildExecutable {
	buildExecutable := &BuildExecutable{
		ID:      sql.NullInt64{Int64: c.GetID(), Valid: true},
		BuildID: sql.NullInt64{Int64: c.GetBuildID(), Valid: true},
		Data:    c.GetData(),
	}

	return buildExecutable.Nullify()
}
