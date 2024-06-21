// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"strings"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
)

var (
	skipDirectiveMsg = "skip ci directive found in commit title/message"
)

// PullRequest defines the data pulled from PRs while
// processing a webhook.
type PullRequest struct {
	Comment    string
	Number     int
	IsFromFork bool
	Labels     []string
}

// Webhook defines a struct that is used to return
// the required data when processing webhook event
// a for a source provider event.
type Webhook struct {
	Hook        *library.Hook
	Repo        *api.Repo
	Build       *api.Build
	PullRequest PullRequest
	Deployment  *library.Deployment
}

// ShouldSkip uses the build information
// associated with the given hook to determine
// whether the hook should be skipped.
func (w *Webhook) ShouldSkip() (bool, string) {
	// push or tag event
	if strings.EqualFold(constants.EventPush, w.Build.GetEvent()) || strings.EqualFold(constants.EventTag, w.Build.GetEvent()) {
		// check for skip ci directive in message or title
		if hasSkipDirective(w.Build.GetMessage()) ||
			hasSkipDirective(w.Build.GetTitle()) {
			return true, skipDirectiveMsg
		}
	}

	return false, ""
}

// hasSkipDirective is a small helper function
// to check a string for a number of patterns
// that signal to vela that the hook should
// be skipped from processing.
func hasSkipDirective(s string) bool {
	sl := strings.ToLower(s)

	switch {
	case strings.Contains(sl, "[skip ci]"),
		strings.Contains(sl, "[ci skip]"),
		strings.Contains(sl, "[skip vela]"),
		strings.Contains(sl, "[vela skip]"),
		strings.Contains(sl, "***no_ci***"):
		return true
	default:
		return false
	}
}
