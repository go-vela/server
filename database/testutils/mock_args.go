// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"database/sql/driver"
	"reflect"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/lestrrat-go/jwx/jwk"
)

// This will be used with the github.com/DATA-DOG/go-sqlmock library to compare values
// that are otherwise not easily compared. These typically would be values generated
// before adding or updating them in the database.
//
// https://github.com/DATA-DOG/go-sqlmock#matching-arguments-like-timetime
type AnyArgument struct{}

// Match satisfies sqlmock.Argument interface.
func (a AnyArgument) Match(_ driver.Value) bool {
	return true
}

// NowTimestamp is used to test whether timestamps get updated correctly to the current time with lenience.
type NowTimestamp struct{}

// Match satisfies sqlmock.Argument interface.
func (t NowTimestamp) Match(v driver.Value) bool {
	ts, ok := v.(int64)
	if !ok {
		return false
	}

	now := time.Now().Unix()

	return now-ts < 10
}

var JwkKeyOpts = cmp.Options{
	cmp.FilterValues(func(x, y interface{}) bool {
		_, xOk := x.(jwk.RSAPublicKey)
		_, yOk := y.(jwk.RSAPublicKey)
		return xOk && yOk
	}, cmp.Comparer(func(x, y interface{}) bool {
		xJWK := x.(jwk.RSAPublicKey)
		yJWK := y.(jwk.RSAPublicKey)

		var rawXKey, rawYKey interface{}

		if err := xJWK.Raw(&rawXKey); err != nil {
			return false
		}

		if err := yJWK.Raw(&rawYKey); err != nil {
			return false
		}

		return reflect.DeepEqual(rawXKey, rawYKey) && xJWK.KeyID() == yJWK.KeyID()
	})),
}
