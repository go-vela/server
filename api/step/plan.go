// SPDX-License-Identifier: Apache-2.0

package step

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/scm"
)

// PlanSteps is a helper function to plan all steps
// in the build for execution. This creates the steps
// for the build.
func PlanSteps(ctx context.Context, database database.Interface, scm scm.Service, p *pipeline.Build, b *types.Build) ([]*types.Step, error) {
	// variable to store planned steps
	steps := []*types.Step{}

	// iterate through all pipeline stages
	for _, stage := range p.Stages {
		// iterate through all steps for each pipeline stage
		for _, step := range stage.Steps {
			// create the step object
			s, err := planStep(ctx, database, scm, b, step, stage.Name)
			if err != nil {
				return steps, err
			}

			steps = append(steps, s)
		}
	}

	// iterate through all pipeline steps
	for _, step := range p.Steps {
		s, err := planStep(ctx, database, scm, b, step, "")
		if err != nil {
			return steps, err
		}

		steps = append(steps, s)
	}

	return steps, nil
}

func planStep(ctx context.Context, database database.Interface, scm scm.Service, b *types.Build, c *pipeline.Container, stage string) (*types.Step, error) {
	// create the step object
	s := new(types.Step)
	s.SetBuildID(b.GetID())
	s.SetRepoID(b.GetRepo().GetID())
	s.SetNumber(c.Number)
	s.SetName(c.Name)
	s.SetImage(c.Image)
	s.SetStage(stage)
	s.SetStatus(constants.StatusPending)
	s.SetReportAs(c.ReportAs)
	s.SetCreated(time.Now().UTC().Unix())

	// send API call to create the step
	s, err := database.CreateStep(ctx, s)
	if err != nil {
		return nil, fmt.Errorf("unable to create step %s: %w", s.GetName(), err)
	}

	logrus.WithFields(logrus.Fields{
		"step":    s.GetName(),
		"step_id": s.GetID(),
		"org":     b.GetRepo().GetOrg(),
		"repo":    b.GetRepo().GetName(),
		"repo_id": b.GetRepo().GetID(),
	}).Info("step created")

	// populate environment variables from step
	err = c.MergeEnv(s.Environment())
	if err != nil {
		return nil, err
	}

	// create the log object
	l := new(types.Log)
	l.SetStepID(s.GetID())
	l.SetBuildID(b.GetID())
	l.SetRepoID(b.GetRepo().GetID())
	l.SetData([]byte{})

	// send API call to create the step logs
	err = database.CreateLog(ctx, l)
	if err != nil {
		return nil, fmt.Errorf("unable to create logs for step %s: %w", s.GetName(), err)
	}

	logrus.WithFields(logrus.Fields{
		"step":    s.GetName(),
		"step_id": s.GetID(),
		"log_id":  l.GetID(), // it won't have an ID here
		"org":     b.GetRepo().GetOrg(),
		"repo":    b.GetRepo().GetName(),
		"repo_id": b.GetRepo().GetID(),
	}).Info("log for step created")

	if len(s.GetReportAs()) > 0 {
		// send API call to set the status on the commit
		err = scm.StepStatus(ctx, b.GetRepo().GetOwner(), b, s, b.GetRepo().GetOrg(), b.GetRepo().GetName())
		if err != nil {
			logrus.Errorf("unable to set commit status for build: %v", err)
		}
	}

	return s, nil
}
