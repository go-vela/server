// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

// +build !cgo

package database

import (
	// Load in the Gorm Postgres dialect for
	// integrating with a Postgres instance.
	//
	// https://gorm.io/docs/dialects.html
	_ "github.com/jinzhu/gorm/dialects/postgres"
)
