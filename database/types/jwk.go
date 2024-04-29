// SPDX-License-Identifier: Apache-2.0

package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	api "github.com/go-vela/server/api/types"
)

var (
	// ErrEmptyKID defines the error type when a
	// JWK type has an empty ID field provided.
	ErrEmptyKID = errors.New("empty key identifier provided")
)

type (
	// JWK is the database representation of a jwk.
	JWK struct {
		ID     uuid.UUID    `gorm:"type:uuid"`
		Active sql.NullBool `sql:"active"`
		Key    KeyJSON      `sql:"key"`
	}

	KeyJSON api.JWK
)

// Value - Implementation of valuer for database/sql for KeyJSON.
func (k KeyJSON) Value() (driver.Value, error) {
	valueString, err := json.Marshal(k)
	return string(valueString), err
}

// Scan - Implement the database/sql scanner interface for KeyJSON.
func (k *KeyJSON) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &k)
	case string:
		return json.Unmarshal([]byte(v), &k)
	default:
		return fmt.Errorf("wrong type for key: %T", v)
	}
}

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
func (j *JWK) ToAPI() api.JWK {
	return api.JWK(j.Key)
}

// Validate verifies the necessary fields for
// the JWK type are populated correctly.
func (j *JWK) Validate() error {
	// verify the Name field is populated
	if len(j.ID.String()) == 0 {
		return ErrEmptyKID
	}

	return nil
}

// JWKFromAPI converts the API JWK type
// to a database JWK type.
func JWKFromAPI(j api.JWK) *JWK {
	var (
		id  uuid.UUID
		err error
	)

	id, err = uuid.Parse(j.Kid)
	if err != nil {
		return nil
	}

	dashboard := &JWK{
		ID:  id,
		Key: KeyJSON(j),
	}

	return dashboard.Nullify()
}
