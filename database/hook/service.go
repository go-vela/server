// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

import (
	"github.com/go-vela/types/library"
)

// HookService represents the Vela interface for hook
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type HookService interface {
	// Hook Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateHookIndexes defines a function that creates the indexes for the hooks table.
	CreateHookIndexes() error
	// CreateHookTable defines a function that creates the hooks table.
	CreateHookTable(string) error

	// Hook Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountHooks defines a function that gets the count of all hooks.
	CountHooks() (int64, error)
	// CountHooksForRepo defines a function that gets the count of hooks by repo ID.
	CountHooksForRepo(*library.Repo) (int64, error)
	// CreateHook defines a function that creates a new hook.
	CreateHook(*library.Hook) error
	// DeleteHook defines a function that deletes an existing hook.
	DeleteHook(*library.Hook) error
	// GetHook defines a function that gets a hook by ID.
	GetHook(int64) (*library.Hook, error)
	// GetHookForRepo defines a function that gets a hook by repo ID and number.
	GetHookForRepo(*library.Repo, int) (*library.Hook, error)
	// LatestHookForRepo defines a function that gets the latest hook by repo ID.
	LatestHookForRepo(*library.Repo, int) (*library.Hook, error)
	// ListHooks defines a function that gets a list of all hooks.
	ListHooks() ([]*library.Hook, error)
	// ListHooksForRepo defines a function that gets a list of hooks by repo ID.
	ListHooksForRepo(*library.Repo, int, int, int) (*library.Hook, error)
	// UpdateHook defines a function that updates an existing hook.
	UpdateHook(*library.Hook) error
}
