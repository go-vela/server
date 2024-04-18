// SPDX-License-Identifier: Apache-2.0

package models

import (
	api "github.com/go-vela/server/api/types"
)

// ItemVersion allows the worker to detect items that were queued before an Vela server
// upgrade or downgrade, so it can handle such stale data gracefully.
// For example, the worker could fail a stale build or ask the server to recompile it.
// This is not a public API and is unrelated to the version key in pipeline yaml.
const ItemVersion uint64 = 3

// Item is the queue representation of an item to publish to the queue.
type Item struct {
	Build *api.Build `json:"build"`
	// The 0-value is the implicit ItemVersion for queued Items that pre-date adding the field.
	ItemVersion uint64 `json:"item_version"`
}

// ToItem creates a queue item from a build, repo and user.
func ToItem(b *api.Build) *Item {
	return &Item{
		Build:       b,
		ItemVersion: ItemVersion,
	}
}
