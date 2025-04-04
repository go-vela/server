// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"

	api "github.com/go-vela/server/api/types"
)

// ServiceInterface represents the Vela interface for service
// functions with the supported Database backends.
//
//nolint:revive // ignore name stutter
type ServiceInterface interface {
	// Service Data Definition Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_definition_language

	// CreateServiceTable defines a function that creates the services table.
	CreateServiceTable(context.Context, string) error

	// Service Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CleanServices defines a function that sets running or pending services to error status before a given created time.
	CleanServices(context.Context, string, int64) (int64, error)
	// CountServices defines a function that gets the count of all services.
	CountServices(context.Context) (int64, error)
	// CountServicesForBuild defines a function that gets the count of services by build ID.
	CountServicesForBuild(context.Context, *api.Build, map[string]interface{}) (int64, error)
	// CreateService defines a function that creates a new service.
	CreateService(context.Context, *api.Service) (*api.Service, error)
	// DeleteService defines a function that deletes an existing service.
	DeleteService(context.Context, *api.Service) error
	// GetService defines a function that gets a service by ID.
	GetService(context.Context, int64) (*api.Service, error)
	// GetServiceForBuild defines a function that gets a service by number and build ID.
	GetServiceForBuild(context.Context, *api.Build, int32) (*api.Service, error)
	// ListServices defines a function that gets a list of all services.
	ListServices(context.Context) ([]*api.Service, error)
	// ListServicesForBuild defines a function that gets a list of services by build ID.
	ListServicesForBuild(context.Context, *api.Build, map[string]interface{}, int, int) ([]*api.Service, error)
	// ListServiceImageCount defines a function that gets a list of all service images and the count of their occurrence.
	ListServiceImageCount(context.Context) (map[string]float64, error)
	// ListServiceStatusCount defines a function that gets a list of all service statuses and the count of their occurrence.
	ListServiceStatusCount(context.Context) (map[string]float64, error)
	// UpdateService defines a function that updates an existing service.
	UpdateService(context.Context, *api.Service) (*api.Service, error)
}
