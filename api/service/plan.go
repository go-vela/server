// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
)

// PlanServices is a helper function to plan all services
// in the build for execution. This creates the services
// for the build.
func PlanServices(ctx context.Context, database database.Interface, p *pipeline.Build, b *types.Build) ([]*types.Service, error) {
	// variable to store planned services
	services := []*types.Service{}

	// iterate through all pipeline services
	for _, service := range p.Services {
		// create the service object
		s := new(types.Service)
		s.SetBuildID(b.GetID())
		s.SetRepoID(b.GetRepo().GetID())
		s.SetName(service.Name)
		s.SetImage(service.Image)
		s.SetNumber(service.Number)
		s.SetStatus(constants.StatusPending)
		s.SetCreated(time.Now().UTC().Unix())

		// send API call to create the service
		s, err := database.CreateService(ctx, s)
		if err != nil {
			return services, fmt.Errorf("unable to create service %s: %w", s.GetName(), err)
		}

		logrus.WithFields(logrus.Fields{
			"service":    s.GetName(),
			"service_id": s.GetID(),
			"org":        b.GetRepo().GetOrg(),
			"repo":       b.GetRepo().GetName(),
			"repo_id":    b.GetRepo().GetID(),
		}).Info("service created")

		err = service.MergeEnv(s.Environment())
		if err != nil {
			return services, err
		}

		// create the log object
		l := new(types.Log)
		l.SetServiceID(s.GetID())
		l.SetBuildID(b.GetID())
		l.SetRepoID(b.GetRepo().GetID())
		l.SetData([]byte{})
		l.SetCreatedAt(time.Now().UTC().Unix())

		// send API call to create the service logs
		err = database.CreateLog(ctx, l)
		if err != nil {
			return services, fmt.Errorf("unable to create service logs for service %s: %w", s.GetName(), err)
		}

		logrus.WithFields(logrus.Fields{
			"service":    s.GetName(),
			"service_id": s.GetID(),
			"log_id":     l.GetID(), // it won't have an ID here, because CreateLog doesn't return the created log
			"org":        b.GetRepo().GetOrg(),
			"repo":       b.GetRepo().GetName(),
			"repo_id":    b.GetRepo().GetID(),
		}).Info("log for service created")
	}

	return services, nil
}
