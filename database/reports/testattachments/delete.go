package testattachments

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
	"github.com/sirupsen/logrus"
)

// DeleteByID deletes an existing test report attachment from the database.
func (e *Engine) DeleteByID(ctx context.Context, r *api.TestReportAttachments) error {
	e.logger.WithFields(logrus.Fields{
		"test_report_attachment": r.GetID(),
	}).Tracef("deleting test report attachment %d", r.GetID())

	// cast the API type to database type
	attachment := types.TestReportAttachmentFromAPI(r)

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableAttachments).
		Delete(attachment).
		Error
}
