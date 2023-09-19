// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-vela/server/database"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

// PlanServices is a helper function to plan all services
// in the build for execution. This creates the services
// for the build in the configured backend.
func PlanServices(ctx context.Context, database database.Interface, p *pipeline.Build, b *library.Build) ([]*library.Service, error) {
	// variable to store planned services
	services := []*library.Service{}

	// iterate through all pipeline services
	for _, service := range p.Services {
		// create the service object
		s := new(library.Service)
		s.SetBuildID(b.GetID())
		s.SetRepoID(b.GetRepoID())
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
		l.SetRepoID(b.GetRepoID())
		l.SetData([]byte{})

		// send API call to create the service logs
		err = database.CreateLog(ctx, l)
		if err != nil {
			return services, fmt.Errorf("unable to create service logs for service %s: %w", s.GetName(), err)
		}
	}

	return services, nil
}
