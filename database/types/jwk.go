// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwk"
)

var (
	// ErrInvalidKID defines the error type when a
	// JWK type has an invalid ID field provided.
	ErrInvalidKID = errors.New("invalid key identifier provided")
)

type (
	// JWK is the database representation of a jwk.
	JWK struct {
		ID     uuid.UUID    `gorm:"type:uuid"`
		Active sql.NullBool `sql:"active"`
		Key    []byte       `sql:"key"`
	}
)

// Nullify ensures the valid flag for
// the sql.Null types are properly set.
//
// When a field within the JWK type is the zero
// value for the field, the valid flag is set to
// false causing it to be NULL in the database.
func (j *JWK) Nullify() *JWK {
	if j == nil {
		return nil
	}

	return j
}

// ToAPI converts the JWK type
// to an API JWK type.
func (j *JWK) ToAPI() jwk.RSAPublicKey {
	parsedKey, _ := jwk.ParseKey(j.Key)

	switch jwk := parsedKey.(type) {
	case jwk.RSAPublicKey:
		return jwk
	default:
		return nil
	}
}

// JWKFromAPI converts the API JWK type
// to a database JWK type.
func JWKFromAPI(j jwk.RSAPublicKey) *JWK {
	var (
		id  uuid.UUID
		err error
	)

	id, err = uuid.Parse(j.KeyID())
	if err != nil {
		return nil
	}

	bytesKey, _ := json.Marshal(j)

	jwk := &JWK{
		ID:  id,
		Key: bytesKey,
	}

	return jwk.Nullify()
}
