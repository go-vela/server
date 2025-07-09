// SPDX-License-Identifier: Apache-2.0

package testattachment

import (
	"context"

	"github.com/sirupsen/logrus"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
)

// GetTestAttachmentForBuild gets a test attachment by number and build ID from the database.
func (e *Engine) GetTestAttachmentForBuild(ctx context.Context, b *api.Build) (*api.TestAttachment, error) {
	e.logger.WithFields(logrus.Fields{
		"build": b.GetNumber(),
	}).Tracef("getting testattachment")

	// variable to store query results
	tr := new(types.TestAttachment)

	// send query to the database and store result in variable
	err := e.client.
		WithContext(ctx).
		Table(constants.TableTestAttachment).
		Where("build_id = ?", b.GetID()).
		Take(tr).
		Error
	if err != nil {
		return nil, err
	}

	return tr.ToAPI(), nil
}
