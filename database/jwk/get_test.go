// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database/testutils"
)

func TestJWK_Engine_GetJWK(t *testing.T) {
	// setup types
	_jwk := testutils.APIJWK()
	_jwk.Kid = "c8da1302-07d6-11ea-882f-4893bca275b8"
	_jwk.Algorithm = "RS256"
	_jwk.Kty = "rsa"
	_jwk.Use = "sig"
	_jwk.N = "123456"
	_jwk.E = "123"

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "active", "key"},
	).AddRow("c8da1302-07d6-11ea-882f-4893bca275b8", true, []byte(`{"alg":"RS256","use":"sig","x5t":"","kid":"c8da1302-07d6-11ea-882f-4893bca275b8","kty":"rsa","x5c":null,"n":"123456","e":"123"}`))

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "jwks" WHERE id = $1 AND active = $2 LIMIT $3`).WithArgs("c8da1302-07d6-11ea-882f-4893bca275b8", true, 1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateJWK(context.TODO(), _jwk)
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     api.JWK
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _jwk,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _jwk,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.GetActiveJWK(context.TODO(), "c8da1302-07d6-11ea-882f-4893bca275b8")

			if test.failure {
				if err == nil {
					t.Errorf("GetActiveJWK for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetActiveJWK for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(got, test.want); diff != "" {
				t.Errorf("GetActiveJWK mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
