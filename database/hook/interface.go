// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package hook

import (
	"context"

	"github.com/go-vela/types/library"
)

// HookInterface represents the Vela interface for hook
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type HookInterface interface {
	// Hook Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateHookIndexes defines a function that creates the indexes for the hooks table.
	CreateHookIndexes(context.Context) error
	// CreateHookTable defines a function that creates the hooks table.
	CreateHookTable(context.Context, string) error

	// Hook Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CountHooks defines a function that gets the count of all hooks.
	CountHooks(context.Context) (int64, error)
	// CountHooksForRepo defines a function that gets the count of hooks by repo ID.
	CountHooksForRepo(context.Context, *library.Repo) (int64, error)
	// CreateHook defines a function that creates a new hook.
	CreateHook(context.Context, *library.Hook) (*library.Hook, error)
	// DeleteHook defines a function that deletes an existing hook.
	DeleteHook(context.Context, *library.Hook) error
	// GetHook defines a function that gets a hook by ID.
	GetHook(context.Context, int64) (*library.Hook, error)
	// GetHookByWebhook defines a function that gets any hook with a matching webhook_id.
	GetHookByWebhook(context.Context, int64) (*library.Hook, error)
	// GetHookForRepo defines a function that gets a hook by repo ID and number.
	GetHookForRepo(context.Context, *library.Repo, int) (*library.Hook, error)
	// LastHookForRepo defines a function that gets the last hook by repo ID.
	LastHookForRepo(context.Context, *library.Repo) (*library.Hook, error)
	// ListHooks defines a function that gets a list of all hooks.
	ListHooks(context.Context) ([]*library.Hook, error)
	// ListHooksForRepo defines a function that gets a list of hooks by repo ID.
	ListHooksForRepo(context.Context, *library.Repo, int, int) ([]*library.Hook, int64, error)
	// UpdateHook defines a function that updates an existing hook.
	UpdateHook(context.Context, *library.Hook) (*library.Hook, error)
}
