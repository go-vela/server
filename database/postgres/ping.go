// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"fmt"
	"time"
)

// Ping sends a "ping" request with backoff to the database.
func (c *client) Ping() error {
	c.Logger.Trace("sending ping requests to the postgres database")

	// create a loop to attempt ping requests 5 times
	for i := 0; i < 5; i++ {
		// capture database/sql database from gorm database
		//
		// https://pkg.go.dev/gorm.io/gorm#DB.DB
		_sql, err := c.Postgres.DB()
		if err != nil {
			return err
		}

		// send ping request to database
		//
		// https://pkg.go.dev/database/sql#DB.Ping
		err = _sql.Ping()
		if err != nil {
			c.Logger.Debugf("unable to ping database - retrying in %v", time.Duration(i)*time.Second)

			// sleep for loop iteration in seconds
			time.Sleep(time.Duration(i) * time.Second)

			// continue to next iteration of the loop
			continue
		}

		// able to ping database so return with no error
		return nil
	}

	return fmt.Errorf("unable to successfully ping database")
}
