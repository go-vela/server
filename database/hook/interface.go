// SPDX-License-Identifier: Apache-2.0

package hook

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// HookInterface represents the Vela interface for hook
// functions with the supported Database backends.
//

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
	CountHooksForRepo(context.Context, *api.Repo) (int64, error)
	// CreateHook defines a function that creates a new hook.
	CreateHook(context.Context, *api.Hook) (*api.Hook, error)
	// DeleteHook defines a function that deletes an existing hook.
	DeleteHook(context.Context, *api.Hook) error
	// GetHook defines a function that gets a hook by ID.
	GetHook(context.Context, int64) (*api.Hook, error)
	// GetHookByWebhookID defines a function that gets any hook with a matching webhook_id.
	GetHookByWebhookID(context.Context, int64) (*api.Hook, error)
	// GetHookForRepo defines a function that gets a hook by repo ID and number.
	GetHookForRepo(context.Context, *api.Repo, int64) (*api.Hook, error)
	// ListHooks defines a function that gets a list of all hooks.
	ListHooks(context.Context) ([]*api.Hook, error)
	// ListHooksForRepo defines a function that gets a list of hooks by repo ID.
	ListHooksForRepo(context.Context, *api.Repo, int, int) ([]*api.Hook, error)
	// UpdateHook defines a function that updates an existing hook.
	UpdateHook(context.Context, *api.Hook) (*api.Hook, error)
}
