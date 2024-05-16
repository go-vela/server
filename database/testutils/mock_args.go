// SPDX-License-Identifier: Apache-2.0

package testutils

import (
	"database/sql/driver"
	"time"
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
