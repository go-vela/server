// SPDX-License-Identifier: Apache-2.0

package constants

// Repo visibility types.
const (
	// VisibilityPublic defines the visibility type for allowing any
	// users in Vela to access their repo regardless of the access
	// defined in the source control system.
	VisibilityPublic = "public"

	// VisibilityPrivate defines the visibility type for only allowing
	// users in Vela with pre-defined access in the source control
	// system to access their repo.
	VisibilityPrivate = "private"
)
