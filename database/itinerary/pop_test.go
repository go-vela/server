// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package itinerary

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-vela/types/library"
)

func TestItinerary_Engine_PopBuildItinerary(t *testing.T) {
	// setup types
	_bItinerary := testBuildItinerary()
	_bItinerary.SetID(1)
	_bItinerary.SetBuildID(1)
	_bItinerary.SetData([]byte("foo"))

	_postgres, _mock := testPostgres(t)
	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// create expected result in mock
	_rows := sqlmock.NewRows(
		[]string{"id", "build_id", "data"}).
		AddRow(1, 1, []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69})

	// ensure the mock expects the query
	_mock.ExpectQuery(`DELETE FROM "build_itineraries" WHERE build_id = $1 RETURNING *`).WithArgs(1).WillReturnRows(_rows)

	_sqlite := testSqlite(t)
	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	err := _sqlite.CreateBuildItinerary(_bItinerary)
	if err != nil {
		t.Errorf("unable to create test build itinerary for sqlite: %v", err)
	}

	// setup tests
	tests := []struct {
		failure  bool
		name     string
		database *engine
		want     *library.BuildItinerary
	}{
		{
			failure:  false,
			name:     "postgres",
			database: _postgres,
			want:     _bItinerary,
		},
		{
			failure:  false,
			name:     "sqlite3",
			database: _sqlite,
			want:     _bItinerary,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := test.database.PopBuildItinerary(1)

			if test.failure {
				if err == nil {
					t.Errorf("PopBuildItinerary for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("PopBuildItinerary for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("PopBuildItinerary for %s is %v, want %v", test.name, got, test.want)
			}
		})
	}
}
