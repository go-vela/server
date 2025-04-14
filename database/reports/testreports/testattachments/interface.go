// SPDX-License-Identifier: Apache-2.0

package testattachments

import (
	"context"
)

// TestAttachmentsInterface represents the Vela interface for testattachments
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type TestAttachmentsInterface interface {
	// TestAttachments Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateTestAttachmentsIndexes defines a function that creates the indexes for the testattachments table.
	CreateTestAttachmentsIndexes(context.Context) error
	// CreateTestAttachmentsTable defines a function that creates the testattachments table.
	CreateTestAttachmentsTable(context.Context, string) error
}
