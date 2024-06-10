// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
)

// PlanServices is a helper function to plan all services
// in the build for execution. This creates the services
// for the build.
func PlanServices(ctx context.Context, database database.Interface, p *pipeline.Build, b *types.Build) ([]*library.Service, error) {
	// variable to store planned services
	services := []*library.Service{}

	// iterate through all pipeline services
	for _, service := range p.Services {
		// create the service object
		s := new(library.Service)
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
			"repo":       b.GetRepo().GetFullName(),
		}).Info("service created")

		// populate environment variables from service library
		//
		// https://pkg.go.dev/github.com/go-vela/types/library#Service.Environment
		err = service.MergeEnv(s.Environment())
		if err != nil {
			return services, err
		}

		// create the log object
		l := new(library.Log)
		l.SetServiceID(s.GetID())
		l.SetBuildID(b.GetID())
		l.SetRepoID(b.GetRepo().GetID())
		l.SetData([]byte{})

		// send API call to create the service logs
		err = database.CreateLog(ctx, l)
		if err != nil {
			return services, fmt.Errorf("unable to create service logs for service %s: %w", s.GetName(), err)
		}

		logrus.WithFields(logrus.Fields{
			"service":    s.GetName(),
			"service_id": s.GetID(),
			"log_id":     l.GetID(),
			"repo":       b.GetRepo().GetFullName(),
		}).Info("log for service created")
	}

	return services, nil
}
