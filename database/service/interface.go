// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package service

import (
	"github.com/go-vela/types/library"
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
	CreateServiceTable(string) error

	// Service Data Manipulation Language Functions
	//
	// https://en.wikipedia.org/wiki/Data_manipulation_language

	// CleanServices defines a function that sets running or pending services to error status before a given created time.
	CleanServices(string, int64) (int64, error)
	// CountServices defines a function that gets the count of all services.
	CountServices() (int64, error)
	// CountServicesForBuild defines a function that gets the count of services by build ID.
	CountServicesForBuild(*library.Build, map[string]interface{}) (int64, error)
	// CreateService defines a function that creates a new service.
	CreateService(*library.Service) (*library.Service, error)
	// DeleteService defines a function that deletes an existing service.
	DeleteService(*library.Service) error
	// GetService defines a function that gets a service by ID.
	GetService(int64) (*library.Service, error)
	// GetServiceForBuild defines a function that gets a service by number and build ID.
	GetServiceForBuild(*library.Build, int) (*library.Service, error)
	// ListServices defines a function that gets a list of all services.
	ListServices() ([]*library.Service, error)
	// ListServicesForBuild defines a function that gets a list of services by build ID.
	ListServicesForBuild(*library.Build, map[string]interface{}, int, int) ([]*library.Service, int64, error)
	// ListServiceImageCount defines a function that gets a list of all service images and the count of their occurrence.
	ListServiceImageCount() (map[string]float64, error)
	// ListServiceStatusCount defines a function that gets a list of all service statuses and the count of their occurrence.
	ListServiceStatusCount() (map[string]float64, error)
	// UpdateService defines a function that updates an existing service.
	UpdateService(*library.Service) (*library.Service, error)
}
