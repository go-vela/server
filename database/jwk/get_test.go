// SPDX-License-Identifier: Apache-2.0

package jwk

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"
	"github.com/lestrrat-go/jwx/v3/jwk"

	"github.com/go-vela/server/database/testutils"
)

func TestJWK_Engine_GetJWK(t *testing.T) {
	// setup types
	_jwk := testutils.JWK()
	_jwkBytes, err := json.Marshal(_jwk)
	if err != nil {
		t.Errorf("unable to marshal JWK: %v", err)
	}

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "active", "key"},
	).AddRow(_jwk.KeyID(), true, _jwkBytes)

	// ensure the mock expects the query
	_mock.ExpectQuery(`SELECT * FROM "jwks" WHERE id = $1 AND active = $2 LIMIT $3`).WithArgs(_jwk.KeyID(), true, 1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err = _sqlite.CreateJWK(context.TODO(), _jwk)
	if err != nil {
		t.Errorf("unable to create test repo for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     jwk.RSAPublicKey
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
			got, err := test.database.GetActiveJWK(context.TODO(), _jwk.KeyID())

			if test.failure {
				if err == nil {
					t.Errorf("GetActiveJWK for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("GetActiveJWK for %s returned err: %v", test.name, err)
			}

			if diff := cmp.Diff(test.want, got, testutils.JwkKeyOpts); diff != "" {
				t.Errorf("GetActiveJWK mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
