// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lestrrat-go/jwx/v3/jwk"

	"github.com/go-vela/server/database/testutils"
)

func TestJWK_Engine_ListJWKs(t *testing.T) {
	// setup types
	_jwkOne := testutils.JWK()
	_jwkOneBytes, err := json.Marshal(_jwkOne)
	if err != nil {
		t.Errorf("unable to marshal JWK: %v", err)
	}

	_jwkTwo := testutils.JWK()
	_jwkTwoBytes, err := json.Marshal(_jwkTwo)
	if err != nil {
		t.Errorf("unable to marshal JWK: %v", err)
	}

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	kidOne, ok := _jwkOne.KeyID()
	if !ok {
		t.Errorf("unable to get key ID for jwk")
	}

	kidTwo, ok := _jwkTwo.KeyID()
	if !ok {
		t.Errorf("unable to get key ID for jwk")
	}

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "active", "key"}).
		AddRow(kidOne, true, _jwkOneBytes).
		AddRow(kidTwo, true, _jwkTwoBytes)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "jwks"`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err = _sqlite.CreateJWK(context.TODO(), _jwkOne)
	if err != nil {
		t.Errorf("unable to create test jwk for sqlite: %v", err)
	}

	err = _sqlite.CreateJWK(context.TODO(), _jwkTwo)
	if err != nil {
		t.Errorf("unable to create test jwk for sqlite: %v", err)
	}

	wantSet := jwk.NewSet()

	err = wantSet.AddKey(_jwkOne)
	if err != nil {
		t.Errorf("unable to add jwk to set: %v", err)
	}

	err = wantSet.AddKey(_jwkTwo)
	if err != nil {
		t.Errorf("unable to add jwk to set: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     jwk.Set
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     wantSet,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     wantSet,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.ListJWKs(context.TODO())

			if test.failure {
				if err == nil {
					t.Errorf("ListJWKs for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("ListJWKs for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("ListJWKs for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
