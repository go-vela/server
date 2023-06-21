// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

import (
	"fmt"
	"time"
)

// Ping sends a "ping" request with backoff to the database.
func (e *engine) Ping() error {
	e.logger.Tracef("sending ping request to the %s database", e.Driver())

	// create a loop to attempt ping requests 5 times
	for i := 0; i < 5; i++ {
		// capture database/sql database from gorm.io/gorm database
		_sql, err := e.client.DB()
		if err != nil {
			return err
		}

		// send ping request to database
		err = _sql.Ping()
		if err != nil {
			// create the duration of time to sleep for before attempting to retry
			duration := time.Duration(i+1) * time.Second

			e.logger.Warnf("unable to ping %s database - retrying in %v", e.Driver(), duration)

			// sleep for loop iteration in seconds
			time.Sleep(duration)

			continue
		}

		return nil
	}

	return fmt.Errorf("unable to successfully ping %s database", e.Driver())
}
