// SPDX-License-Identifier: Apache-2.0

package log

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/constants"
)

// CleanupResult represents the result of a log cleanup operation.
type CleanupResult struct {
	DeletedCount       int64
	AffectedPartitions []string
}

// CleanLogs deletes log records created before the specified timestamp in batches.
// This function processes deletions in configurable batches with sleep intervals
// to reduce database load and prevent timeouts on large datasets.
// It uses a database-based lock to prevent multiple concurrent cleanup operations.
// If partitioned mode is enabled and the database is PostgreSQL, it will use
// partition-aware cleanup for better performance.
func (e *Engine) CleanLogs(ctx context.Context, before int64, batchSize int, withVacuum bool, driver string) (*CleanupResult, error) {
	logrus.Tracef("cleaning logs created before %d in batches of %d", before, batchSize)

	// try to acquire a distributed lock to prevent concurrent cleanup operations
	acquired, err := e.acquireCleanupLock(ctx, driver)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire cleanup lock: %w", err)
	}

	if !acquired {
		logrus.Debug("cleanup lock is already held by another operation")
		return nil, fmt.Errorf("another cleanup operation is already in progress")
	}

	logrus.Debug("successfully acquired cleanup lock")

	defer func() {
		if acquired {
			if releaseErr := e.releaseCleanupLock(ctx, driver); releaseErr != nil {
				logrus.Warnf("failed to release cleanup lock: %v", releaseErr)
			} else {
				logrus.Debug("successfully released cleanup lock")
			}
		}
	}()

	if e.isPartitionedModeEnabled() && driver == constants.DriverPostgres {
		logrus.Debug("attempting partition-aware log cleanup")

		result, err := e.cleanLogsPartitioned(ctx, before, batchSize, withVacuum)
		if err == nil {
			logrus.Infof("partition-aware cleanup completed successfully, deleted %d logs", result.DeletedCount)
			return result, nil
		}
		// fall back to traditional cleanup
		logrus.Warnf("partition-aware cleanup failed (%v), falling back to traditional cleanup", err)
	}

	return e.cleanLogs(ctx, before, batchSize, withVacuum, driver)
}

// cleanLogs implements the original non-partitioned cleanup logic.
func (e *Engine) cleanLogs(ctx context.Context, before int64, batchSize int, withVacuum bool, driver string) (*CleanupResult, error) {
	var totalDeleted int64

	// process deletions in batches
	for {
		// check for context cancellation
		select {
		case <-ctx.Done():
			return &CleanupResult{
				DeletedCount:       totalDeleted,
				AffectedPartitions: []string{}, // Traditional mode doesn't use partitions
			}, ctx.Err()
		default:
		}

		// delete a batch of logs using a single query
		// this combines the SELECT and DELETE into one operation for better performance
		// both PostgreSQL and SQLite support this DELETE with subquery pattern
		deleteResult := e.client.
			WithContext(ctx).
			Exec("DELETE FROM logs WHERE id IN (SELECT id FROM logs WHERE created_at < ? ORDER BY created_at ASC LIMIT ?)", before, batchSize)

		if deleteResult.Error != nil {
			return &CleanupResult{
				DeletedCount:       totalDeleted,
				AffectedPartitions: []string{}, // Traditional mode doesn't use partitions
			}, fmt.Errorf("failed to delete batch of logs: %w", deleteResult.Error)
		}

		batchDeleted := deleteResult.RowsAffected

		// if no records were deleted, we're done
		if batchDeleted == 0 {
			break
		}

		totalDeleted += batchDeleted

		logrus.Debugf("deleted batch of %d logs (total: %d)", batchDeleted, totalDeleted)

		// if we deleted fewer records than the batch size, we're done
		if batchDeleted < int64(batchSize) {
			break
		}

		// sleep between batches to reduce database load
		time.Sleep(100 * time.Millisecond)
	}

	logrus.Infof("cleaned %d logs created before %d", totalDeleted, before)

	// optionally run VACUUM to reclaim space
	if withVacuum && totalDeleted > 0 {
		err := e.vacuumLogs(ctx, driver)
		if err != nil {
			// don't fail the entire operation if vacuum fails
			logrus.Warnf("failed to vacuum logs table after cleanup: %v", err)
		} else {
			logrus.Info("successfully vacuumed logs table after cleanup")
		}
	}

	return &CleanupResult{
		DeletedCount:       totalDeleted,
		AffectedPartitions: []string{}, // no partitions in traditional mode
	}, nil
}

// vacuumLogs runs VACUUM on the logs table to reclaim space.
func (e *Engine) vacuumLogs(ctx context.Context, driver string) error {
	switch driver {
	case constants.DriverPostgres:
		// PostgreSQL VACUUM ANALYZE
		return e.client.
			WithContext(ctx).
			Exec("VACUUM ANALYZE logs").Error
	case constants.DriverSqlite:
		fallthrough
	default:
		// SQLite VACUUM (note: VACUUM cannot be run inside a transaction in SQLite)
		return e.client.
			WithContext(ctx).
			Exec("VACUUM").Error
	}
}

// acquireCleanupLock attempts to acquire a distributed advisory lock for cleanup operations.
// returns true if the lock was acquired, false if another operation is already running.
func (e *Engine) acquireCleanupLock(ctx context.Context, driver string) (bool, error) {
	// use database-specific advisory locking mechanisms
	// Lock ID: 123456789 (arbitrary number for log cleanup operations)
	lockID := int64(123456789)

	switch driver {
	case constants.DriverPostgres:
		// PostgreSQL advisory lock - non-blocking attempt
		var acquired bool
		err := e.client.WithContext(ctx).
			Raw("SELECT pg_try_advisory_lock(?)", lockID).
			Scan(&acquired).Error

		if err != nil {
			return false, fmt.Errorf("failed to acquire PostgreSQL advisory lock: %w", err)
		}

		return acquired, nil

	case constants.DriverSqlite:
		fallthrough
	default:
		// SQLite doesn't have advisory locks, so we use a BEGIN IMMEDIATE transaction
		// to get an exclusive lock on the database for coordination
		tx := e.client.WithContext(ctx).Begin()
		if tx.Error != nil {
			return false, fmt.Errorf("failed to begin SQLite transaction: %w", tx.Error)
		}

		// try to acquire an exclusive lock by attempting to write to a coordination table
		// if this fails, another cleanup is likely in progress
		result := tx.Exec("CREATE TABLE IF NOT EXISTS log_cleanup_coordination (lock_holder TEXT PRIMARY KEY)")
		if result.Error != nil {
			tx.Rollback()
			return false, fmt.Errorf("failed to create coordination table: %w", result.Error)
		}

		result = tx.Exec("INSERT OR FAIL INTO log_cleanup_coordination (lock_holder) VALUES ('active')")
		if result.Error != nil {
			tx.Rollback()
			// if insert fails, another cleanup is in progress
			return false, nil
		}

		// keep the transaction open to hold the lock
		tx.Commit()

		return true, nil
	}
}

// releaseCleanupLock releases the distributed advisory lock.
func (e *Engine) releaseCleanupLock(ctx context.Context, driver string) error {
	lockID := int64(123456789)

	switch driver {
	case constants.DriverPostgres:
		// PostgreSQL advisory lock release
		var released bool
		err := e.client.WithContext(ctx).
			Raw("SELECT pg_advisory_unlock(?)", lockID).
			Scan(&released).Error

		if err != nil {
			return fmt.Errorf("failed to release PostgreSQL advisory lock: %w", err)
		}

		if !released {
			// this is a warning condition, not an error - log it but don't fail
			logrus.Debugf("PostgreSQL advisory lock was not held during release (this can happen if cleanup was interrupted)")
			return nil
		}

		return nil

	case constants.DriverSqlite:
		fallthrough
	default:
		// for SQLite, clean up the coordination table
		result := e.client.WithContext(ctx).Exec("DELETE FROM log_cleanup_coordination WHERE lock_holder = 'active'")
		if result.Error != nil {
			return fmt.Errorf("failed to release SQLite coordination lock: %w", result.Error)
		}

		return nil
	}
}

// isPartitionedModeEnabled checks if partition-aware cleanup is enabled.
func (e *Engine) isPartitionedModeEnabled() bool {
	return e.config.LogPartitioned
}

// cleanLogsPartitioned implements partition-aware cleanup for PostgreSQL partitioned tables.
func (e *Engine) cleanLogsPartitioned(ctx context.Context, before int64, batchSize int, withVacuum bool) (*CleanupResult, error) {
	// discover partitions matching the configured pattern
	partitions, err := e.discoverLogPartitions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to discover log partitions: %w", err)
	}

	if len(partitions) == 0 {
		return nil, fmt.Errorf("no log partitions found matching pattern %s", e.config.LogPartitionPattern)
	}

	logrus.Debugf("found %d log partitions for cleanup", len(partitions))

	var totalDeleted int64

	affectedPartitions := []string{}

	for _, partition := range partitions {
		// check if partition has records to delete
		count, err := e.countLogsInPartition(ctx, partition, before)
		if err != nil {
			logrus.Warnf("failed to count logs in partition %s: %v", partition, err)
			continue
		}

		if count == 0 {
			logrus.Debugf("partition %s has no logs to delete", partition)
			continue
		}

		logrus.Debugf("partition %s has %d logs to delete", partition, count)

		// delete from this partition in batches
		deleted, err := e.cleanLogsFromPartition(ctx, partition, before, batchSize)
		if err != nil {
			return &CleanupResult{
				DeletedCount:       totalDeleted,
				AffectedPartitions: affectedPartitions,
			}, fmt.Errorf("failed to clean partition %s: %w", partition, err)
		}

		totalDeleted += deleted

		if deleted > 0 {
			affectedPartitions = append(affectedPartitions, partition)
			logrus.Debugf("deleted %d logs from partition %s", deleted, partition)
		}
	}

	logrus.Infof("partition-aware cleanup deleted %d logs from %d partitions", totalDeleted, len(affectedPartitions))

	// run VACUUM on affected partitions if requested
	if withVacuum && len(affectedPartitions) > 0 {
		for _, partition := range affectedPartitions {
			if err := e.vacuumPartition(ctx, partition); err != nil {
				logrus.Warnf("failed to vacuum partition %s: %v", partition, err)
			} else {
				logrus.Debugf("successfully vacuumed partition %s", partition)
			}
		}
	}

	return &CleanupResult{
		DeletedCount:       totalDeleted,
		AffectedPartitions: affectedPartitions,
	}, nil
}

// discoverLogPartitions queries the PostgreSQL system catalogs to find log table partitions.
func (e *Engine) discoverLogPartitions(ctx context.Context) ([]string, error) {
	var partitions []string

	query := `
		SELECT schemaname, tablename 
		FROM pg_tables 
		WHERE tablename LIKE $1 
		  AND schemaname = $2
		  AND tablename != 'logs'
		ORDER BY tablename`

	rows, err := e.client.WithContext(ctx).Raw(query, e.config.LogPartitionPattern, e.config.LogPartitionSchema).Rows()
	if err != nil {
		return nil, fmt.Errorf("failed to query log partitions: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var schema, table string
		if err := rows.Scan(&schema, &table); err != nil {
			return nil, fmt.Errorf("failed to scan partition row: %w", err)
		}

		// use schema.table format for qualified table names
		partitionName := fmt.Sprintf("%s.%s", schema, table)
		partitions = append(partitions, partitionName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating partition rows: %w", err)
	}

	return partitions, nil
}

// countLogsInPartition counts logs in a specific partition that would be deleted.
func (e *Engine) countLogsInPartition(ctx context.Context, partition string, before int64) (int64, error) {
	var count int64

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE created_at < ?", partition)

	err := e.client.WithContext(ctx).Raw(query, before).Scan(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count logs in partition %s: %w", partition, err)
	}

	return count, nil
}

// cleanLogsFromPartition deletes logs from a specific partition in batches.
func (e *Engine) cleanLogsFromPartition(ctx context.Context, partition string, before int64, batchSize int) (int64, error) {
	var totalDeleted int64

	// process deletions in batches for this partition
	for {
		// check for context cancellation
		select {
		case <-ctx.Done():
			return totalDeleted, ctx.Err()
		default:
		}

		// delete a batch from this specific partition
		query := fmt.Sprintf("DELETE FROM %s WHERE id IN (SELECT id FROM %s WHERE created_at < ? ORDER BY created_at ASC LIMIT ?)", partition, partition)
		result := e.client.WithContext(ctx).Exec(query, before, batchSize)

		if result.Error != nil {
			return totalDeleted, fmt.Errorf("failed to delete batch from partition %s: %w", partition, result.Error)
		}

		batchDeleted := result.RowsAffected

		// if no records were deleted, we're done with this partition
		if batchDeleted == 0 {
			break
		}

		totalDeleted += batchDeleted

		// if we deleted fewer records than the batch size, we're done with this partition
		if batchDeleted < int64(batchSize) {
			break
		}

		logrus.Debugf("deleted batch of %d logs (total: %d) from partition %s", batchDeleted, totalDeleted, partition)

		// sleep between batches to reduce database load
		time.Sleep(100 * time.Millisecond)
	}

	return totalDeleted, nil
}

// vacuumPartition runs VACUUM on a specific partition.
func (e *Engine) vacuumPartition(ctx context.Context, partition string) error {
	query := fmt.Sprintf("VACUUM ANALYZE %s", partition)
	return e.client.WithContext(ctx).Exec(query).Error
}
