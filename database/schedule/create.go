// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

//nolint:dupl // ignore similar code with update.go
package schedule

import (
	"context"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/database"
	"github.com/go-vela/types/library"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

// CreateSchedule creates a new schedule in the database.
func (e *engine) CreateSchedule(ctx context.Context, s *library.Schedule) error {
	e.logger.WithFields(logrus.Fields{
		"schedule": s.GetName(),
		"span_id":  trace.SpanFromContext(ctx).SpanContext().SpanID(),
		"trace_id": trace.SpanFromContext(ctx).SpanContext().TraceID(),
	}).Tracef("creating schedule %s in the database", s.GetName())

	// cast the library type to database type
	schedule := database.ScheduleFromLibrary(s)

	// validate the necessary fields are populated
	err := schedule.Validate()
	if err != nil {
		return err
	}

	// send query to the database
	return e.client.
		WithContext(ctx).
		Table(constants.TableSchedule).
		Create(schedule).
		Error
}
