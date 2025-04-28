package testattachments

import (
	"context"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database/types"
	"github.com/sirupsen/logrus"
)

// CreateTestReportAttachment creates a new test attachment in the database.
func (e *Engine) CreateTestReportAttachment(ctx context.Context, r *api.TestReportAttachments) (*api.TestReportAttachments, error) {
	e.logger.WithFields(logrus.Fields{
		"test_report": r.GetID(),
	}).Tracef("creating test report attachment %d", r.GetID())

	attachment := types.TestReportAttachmentFromAPI(r)

	err := attachment.Validate()
	if err != nil {
		return nil, err
	}

	// send query to the database
	err = e.client.
		WithContext(ctx).
		Table(constants.TableAttachments).
		Create(attachment).Error
	if err != nil {
		return nil, err
	}

	result := attachment.ToAPI()
	result.SetTestReportID(r.GetTestReportID())

	return result, nil
}
