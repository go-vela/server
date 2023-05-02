package postgres

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

// Transaction defines a function that executes a given
// database transaction with the given details.
func (c *client) Transaction(ctx context.Context, fc func(tx *gorm.DB) error, opts ...*sql.TxOptions) error {
	return c.Postgres.WithContext(ctx).Transaction(fc, opts...)
}

// TryTransactionLock defines a function that executes a postgres query
// to attempt to lock the given transaction ID.
func (c *client) TryTransactionLock(txID int, tx *gorm.DB) error {
	// use transaction db if provided
	//  allows for nested transactions
	db := c.Postgres
	if tx != nil {
		db = tx
	}

	err := db.Exec("SELECT PG_TRY_ADVISORY_XACT_LOCK(?);", txID).Error
	if err != nil {
		return err
	}

	return nil
}
