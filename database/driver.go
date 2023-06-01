// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package database

// Driver outputs the configured database driver.
func (e *engine) Driver() string {
	return e.Config.Driver
}
