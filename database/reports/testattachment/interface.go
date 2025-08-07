// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// TestAttachmentInterface represents the Vela interface for testattachments
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type TestAttachmentInterface interface {
	// TestAttachment Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateTestAttachmentIndexes defines a function that creates the indexes for the testattachments table.
	CreateTestAttachmentIndexes(context.Context) error
	// CreateTestAttachmentTable defines a function that creates the testattachments table.
	CreateTestAttachmentTable(context.Context, string) error

	// TestAttachment Management Functions

	// CountTestAttachments returns the count of all test attachments.
	CountTestAttachments(context.Context) (int64, error)

	// CreateTestAttachment creates a new test attachment.
	CreateTestAttachment(context.Context, *api.TestAttachment) (*api.TestAttachment, error)

	// DeleteTestAttachment removes a test attachment by ID.
	DeleteTestAttachment(context.Context, *api.TestAttachment) error

	// GetTestAttachment returns a test attachment by ID.
	GetTestAttachment(context.Context, int64) (*api.TestAttachment, error)

	// GetTestAttachmentForBuild defines a function that gets a test report by number and build ID.
	GetTestAttachmentForBuild(context.Context, *api.Build) (*api.TestAttachment, error)

	// ListTestAttachments returns a list of all test attachments.
	ListTestAttachments(context.Context) ([]*api.TestAttachment, error)

	// UpdateTestAttachment updates a test attachment by ID.
	UpdateTestAttachment(context.Context, *api.TestAttachment) (*api.TestAttachment, error)
}
