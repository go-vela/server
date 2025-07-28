// SPDX-License-Identifier: Apache-2.0

package database

// IsLogPartitioned returns whether log partitioning is enabled in the database engine.
func (e *engine) IsLogPartitioned() bool {
	return e.config.LogPartitioned
}
