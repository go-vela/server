// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/testutils"
)

func TestLog_Engine_CleanLogs(t *testing.T) {
	// setup types
	cutoffTime := time.Now().Add(-48 * time.Hour).Unix()

	// setup the test database client
	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	// mock PostgreSQL advisory lock acquisition and release
	_mock.ExpectQuery("SELECT pg_try_advisory_lock($1)").
		WithArgs(int64(123456789)).
		WillReturnRows(sqlmock.NewRows([]string{"pg_try_advisory_lock"}).AddRow(true))

	// setup mock expectations for batch deletion using single DELETE queries
	// first batch: delete 1000 records
	_mock.ExpectExec("DELETE FROM logs WHERE id IN (SELECT id FROM logs WHERE created_at < $1 ORDER BY created_at ASC LIMIT $2)").
		WithArgs(testutils.AnyArgument{}, testutils.AnyArgument{}).
		WillReturnResult(sqlmock.NewResult(1, 1000))

	// second batch: delete 500 records (less than batch size, indicating we're done)
	_mock.ExpectExec("DELETE FROM logs WHERE id IN (SELECT id FROM logs WHERE created_at < $1 ORDER BY created_at ASC LIMIT $2)").
		WithArgs(testutils.AnyArgument{}, testutils.AnyArgument{}).
		WillReturnResult(sqlmock.NewResult(1, 500))

	// mock VACUUM operation
	_mock.ExpectExec(`VACUUM ANALYZE logs`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// mock PostgreSQL advisory lock release
	_mock.ExpectQuery("SELECT pg_advisory_unlock($1)").
		WithArgs(int64(123456789)).
		WillReturnRows(sqlmock.NewRows([]string{"pg_advisory_unlock"}).AddRow(true))

	_sqlite := testSqlite(t)

	defer func() { _sql, _ := _sqlite.client.DB(); _sql.Close() }()

	// create test logs in sqlite for integration testing
	for i := 1; i <= 10; i++ {
		_log := testutils.APILog()
		_log.SetID(int64(i))
		_log.SetBuildID(1)
		_log.SetRepoID(1)
		_log.SetStepID(int64(i))
		_log.SetData([]byte("test log data"))
		_log.SetCreatedAt(cutoffTime - int64(i*3600)) // logs older than cutoff

		err := _sqlite.CreateLog(context.TODO(), _log)
		if err != nil {
			t.Errorf("unable to create test log for sqlite: %v", err)
		}
	}

	// setup tests
	tests := []struct {
		failure     bool
		name        string
		database    *Engine
		before      int64
		batchSize   int
		withVacuum  bool
		driver      string
		wantDeleted int64
	}{
		{
			failure:     false,
			name:        "postgres with vacuum",
			database:    _postgres,
			before:      cutoffTime,
			batchSize:   1000,
			withVacuum:  true,
			driver:      constants.DriverPostgres,
			wantDeleted: 1500,
		},
		{
			failure:     false,
			name:        "sqlite without vacuum",
			database:    _sqlite,
			before:      cutoffTime,
			batchSize:   5,
			withVacuum:  false,
			driver:      constants.DriverSqlite,
			wantDeleted: 10,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.database.CleanLogs(context.TODO(), test.before, test.batchSize, test.withVacuum, test.driver)

			if test.failure {
				if err == nil {
					t.Errorf("CleanLogs for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("CleanLogs for %s returned err: %v", test.name, err)
			}

			if result == nil {
				t.Errorf("CleanLogs for %s returned nil result", test.name)
				return
			}

			if !reflect.DeepEqual(result.DeletedCount, test.wantDeleted) {
				t.Errorf("CleanLogs for %s returned %d deleted, want %d", test.name, result.DeletedCount, test.wantDeleted)
			}

			// For traditional mode (non-partitioned), affected partitions should be empty
			if len(result.AffectedPartitions) != 0 {
				t.Errorf("CleanLogs for %s returned %d affected partitions, want 0 for traditional mode", test.name, len(result.AffectedPartitions))
			}
		})
	}
}

func TestLog_Engine_CleanLogs_EmptyDatabase(t *testing.T) {
	// setup the test database client
	_postgres, _mock := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	cutoffTime := time.Now().Add(-48 * time.Hour).Unix()

	// mock PostgreSQL advisory lock acquisition
	_mock.ExpectQuery("SELECT pg_try_advisory_lock($1)").
		WithArgs(int64(123456789)).
		WillReturnRows(sqlmock.NewRows([]string{"pg_try_advisory_lock"}).AddRow(true))

	// mock the single DELETE query that returns 0 rows (empty database)
	_mock.ExpectExec("DELETE FROM logs WHERE id IN (SELECT id FROM logs WHERE created_at < $1 ORDER BY created_at ASC LIMIT $2)").
		WithArgs(testutils.AnyArgument{}, testutils.AnyArgument{}).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// mock PostgreSQL advisory lock release
	_mock.ExpectQuery("SELECT pg_advisory_unlock($1)").
		WithArgs(int64(123456789)).
		WillReturnRows(sqlmock.NewRows([]string{"pg_advisory_unlock"}).AddRow(true))

	deleted, err := _postgres.CleanLogs(context.TODO(), cutoffTime, 1000, false, constants.DriverPostgres)
	if err != nil {
		t.Errorf("CleanLogs should not have returned err: %v", err)
	}

	if deleted == nil {
		t.Errorf("CleanLogs returned nil result")
		return
	}

	if deleted.DeletedCount != 0 {
		t.Errorf("CleanLogs returned %d deleted, want 0", deleted.DeletedCount)
	}

	if len(deleted.AffectedPartitions) != 0 {
		t.Errorf("CleanLogs returned %d affected partitions, want 0", len(deleted.AffectedPartitions))
	}
}

func TestLog_Engine_CleanLogs_ContextCancellation(t *testing.T) {
	// setup the test database client
	_postgres, _ := testPostgres(t)

	defer func() { _sql, _ := _postgres.client.DB(); _sql.Close() }()

	cutoffTime := time.Now().Add(-48 * time.Hour).Unix()

	// create a canceled context
	ctx, cancel := context.WithCancel(t.Context())
	cancel() // immediately cancel

	// the function should return early due to context cancellation
	result, err := _postgres.CleanLogs(ctx, cutoffTime, 1000, false, constants.DriverPostgres)
	if err == nil {
		t.Errorf("CleanLogs should have returned context cancellation error")
	}

	if result != nil && result.DeletedCount != 0 {
		t.Errorf("CleanLogs returned %d deleted, want 0", result.DeletedCount)
	}
}
