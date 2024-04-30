// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
)

func TestJWK_Engine_ListJWKs(t *testing.T) {
	// setup types
	_jwkOne := testutils.APIJWK()
	_jwkOne.Kid = "c8da1302-07d6-11ea-882f-4893bca275b8"
	_jwkOne.Algorithm = "RS256"
	_jwkOne.Kty = "rsa"
	_jwkOne.Use = "sig"
	_jwkOne.N = "123456"
	_jwkOne.E = "123"

	_jwkTwo := testutils.APIJWK()
	_jwkTwo.Kid = "c8da1302-07d6-11ea-882f-4893bca275b9"
	_jwkTwo.Algorithm = "RS256"
	_jwkTwo.Kty = "rsa"
	_jwkTwo.Use = "sig"
	_jwkTwo.N = "123789"
	_jwkTwo.E = "456"

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "active", "key"}).
		AddRow("c8da1302-07d6-11ea-882f-4893bca275b8", true, []byte(`{"alg":"RS256","use":"sig","x5t":"","kid":"c8da1302-07d6-11ea-882f-4893bca275b8","kty":"rsa","x5c":null,"n":"123456","e":"123"}`)).
		AddRow("c8da1302-07d6-11ea-882f-4893bca275b8", true, []byte(`{"alg":"RS256","use":"sig","x5t":"","kid":"c8da1302-07d6-11ea-882f-4893bca275b9","kty":"rsa","x5c":null,"n":"123789","e":"456"}`))

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "jwks"`).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateJWK(context.TODO(), _jwkOne)
	if err != nil {
		t.Errorf("unable to create test jwk for sqlite: %v", err)
	}

	err = _sqlite.CreateJWK(context.TODO(), _jwkTwo)
	if err != nil {
		t.Errorf("unable to create test jwk for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     []api.JWK
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     []api.JWK{_jwkOne, _jwkTwo},
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     []api.JWK{_jwkOne, _jwkTwo},
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
