// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
)

// Service is the API representation of a service in a build.
//
// swagger:model Service
type Service struct {
	ID           *int64  `json:"id,omitempty"`
	BuildID      *int64  `json:"build_id,omitempty"`
	RepoID       *int64  `json:"repo_id,omitempty"`
	Number       *int32  `json:"number,omitempty"`
	Name         *string `json:"name,omitempty"`
	Image        *string `json:"image,omitempty"`
	Status       *string `json:"status,omitempty"`
	Error        *string `json:"error,omitempty"`
	ExitCode     *int32  `json:"exit_code,omitempty"`
	Created      *int64  `json:"created,omitempty"`
	Started      *int64  `json:"started,omitempty"`
	Finished     *int64  `json:"finished,omitempty"`
	Host         *string `json:"host,omitempty"`
	Runtime      *string `json:"runtime,omitempty"`
	Distribution *string `json:"distribution,omitempty"`
}

// Duration calculates and returns the total amount of
// time the service ran for in a human-readable format.
func (s *Service) Duration() string {
	// check if the service doesn't have a started timestamp
	if s.GetStarted() == 0 {
		return constants.ErrorEmptyDuration
	}

	// capture started unix timestamp from the service
	started := time.Unix(s.GetStarted(), 0)

	// check if the service doesn't have a finished timestamp
	if s.GetFinished() == 0 {
		// return the duration in a human-readable form by
		// subtracting the service started time from the
		// current time rounded to the nearest second
		return time.Since(started).Round(time.Second).String()
	}

	// capture finished unix timestamp from the service
	finished := time.Unix(s.GetFinished(), 0)

	// calculate the duration by subtracting the service
	// started time from the service finished time
	duration := finished.Sub(started)

	// return the duration in a human-readable form
	return duration.String()
}

// Environment returns a list of environment variables
// provided from the fields of the Service type.
func (s *Service) Environment() map[string]string {
	return map[string]string{
		"VELA_SERVICE_CREATED":      ToString(s.GetCreated()),
		"VELA_SERVICE_DISTRIBUTION": ToString(s.GetDistribution()),
		"VELA_SERVICE_EXIT_CODE":    ToString(s.GetExitCode()),
		"VELA_SERVICE_HOST":         ToString(s.GetHost()),
		"VELA_SERVICE_IMAGE":        ToString(s.GetImage()),
		"VELA_SERVICE_NAME":         ToString(s.GetName()),
		"VELA_SERVICE_NUMBER":       ToString(s.GetNumber()),
		"VELA_SERVICE_RUNTIME":      ToString(s.GetRuntime()),
		"VELA_SERVICE_STARTED":      ToString(s.GetStarted()),
		"VELA_SERVICE_STATUS":       ToString(s.GetStatus()),
	}
}

// GetID returns the ID field.
//
// When the provided Service type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Service) GetID() int64 {
	// return zero value if Service type or ID field is nil
	if s == nil || s.ID == nil {
		return 0
	}

	return *s.ID
}

// GetBuildID returns the BuildID field.
//
// When the provided Service type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Service) GetBuildID() int64 {
	// return zero value if Service type or BuildID field is nil
	if s == nil || s.BuildID == nil {
		return 0
	}

	return *s.BuildID
}

// GetRepoID returns the RepoID field.
//
// When the provided Service type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Service) GetRepoID() int64 {
	// return zero value if Service type or RepoID field is nil
	if s == nil || s.RepoID == nil {
		return 0
	}

	return *s.RepoID
}

// GetNumber returns the Number field.
//
// When the provided Service type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Service) GetNumber() int32 {
	// return zero value if Service type or Number field is nil
	if s == nil || s.Number == nil {
		return 0
	}

	return *s.Number
}

// GetName returns the Name field.
//
// When the provided Service type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Service) GetName() string {
	// return zero value if Service type or Name field is nil
	if s == nil || s.Name == nil {
		return ""
	}

	return *s.Name
}

// GetImage returns the Image field.
//
// When the provided Service type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Service) GetImage() string {
	// return zero value if Service type or Image field is nil
	if s == nil || s.Image == nil {
		return ""
	}

	return *s.Image
}

// GetStatus returns the Status field.
//
// When the provided Service type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Service) GetStatus() string {
	// return zero value if Service type or Status field is nil
	if s == nil || s.Status == nil {
		return ""
	}

	return *s.Status
}

// GetError returns the Error field.
//
// When the provided Service type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Service) GetError() string {
	// return zero value if Service type or Error field is nil
	if s == nil || s.Error == nil {
		return ""
	}

	return *s.Error
}

// GetExitCode returns the ExitCode field.
//
// When the provided Service type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Service) GetExitCode() int32 {
	// return zero value if Service type or ExitCode field is nil
	if s == nil || s.ExitCode == nil {
		return 0
	}

	return *s.ExitCode
}

// GetCreated returns the Created field.
//
// When the provided Service type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Service) GetCreated() int64 {
	// return zero value if Service type or Created field is nil
	if s == nil || s.Created == nil {
		return 0
	}

	return *s.Created
}

// GetStarted returns the Started field.
//
// When the provided Service type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Service) GetStarted() int64 {
	// return zero value if Service type or Started field is nil
	if s == nil || s.Started == nil {
		return 0
	}

	return *s.Started
}

// GetFinished returns the Finished field.
//
// When the provided Service type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Service) GetFinished() int64 {
	// return zero value if Service type or Finished field is nil
	if s == nil || s.Finished == nil {
		return 0
	}

	return *s.Finished
}

// GetHost returns the Host field.
//
// When the provided Service type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Service) GetHost() string {
	// return zero value if Service type or Host field is nil
	if s == nil || s.Host == nil {
		return ""
	}

	return *s.Host
}

// GetRuntime returns the Runtime field.
//
// When the provided Service type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Service) GetRuntime() string {
	// return zero value if Service type or Runtime field is nil
	if s == nil || s.Runtime == nil {
		return ""
	}

	return *s.Runtime
}

// GetDistribution returns the Runtime field.
//
// When the provided Service type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Service) GetDistribution() string {
	// return zero value if Service type or Distribution field is nil
	if s == nil || s.Distribution == nil {
		return ""
	}

	return *s.Distribution
}

// SetID sets the ID field.
//
// When the provided Service type is nil, it
// will set nothing and immediately return.
func (s *Service) SetID(v int64) {
	// return if Service type is nil
	if s == nil {
		return
	}

	s.ID = &v
}

// SetBuildID sets the BuildID field.
//
// When the provided Service type is nil, it
// will set nothing and immediately return.
func (s *Service) SetBuildID(v int64) {
	// return if Service type is nil
	if s == nil {
		return
	}

	s.BuildID = &v
}

// SetRepoID sets the RepoID field.
//
// When the provided Service type is nil, it
// will set nothing and immediately return.
func (s *Service) SetRepoID(v int64) {
	// return if Service type is nil
	if s == nil {
		return
	}

	s.RepoID = &v
}

// SetNumber sets the Number field.
//
// When the provided Service type is nil, it
// will set nothing and immediately return.
func (s *Service) SetNumber(v int32) {
	// return if Service type is nil
	if s == nil {
		return
	}

	s.Number = &v
}

// SetName sets the Name field.
//
// When the provided Service type is nil, it
// will set nothing and immediately return.
func (s *Service) SetName(v string) {
	// return if Service type is nil
	if s == nil {
		return
	}

	s.Name = &v
}

// SetImage sets the Image field.
//
// When the provided Service type is nil, it
// will set nothing and immediately return.
func (s *Service) SetImage(v string) {
	// return if Service type is nil
	if s == nil {
		return
	}

	s.Image = &v
}

// SetStatus sets the Status field.
//
// When the provided Service type is nil, it
// will set nothing and immediately return.
func (s *Service) SetStatus(v string) {
	// return if Service type is nil
	if s == nil {
		return
	}

	s.Status = &v
}

// SetError sets the Error field.
//
// When the provided Service type is nil, it
// will set nothing and immediately return.
func (s *Service) SetError(v string) {
	// return if Service type is nil
	if s == nil {
		return
	}

	s.Error = &v
}

// SetExitCode sets the ExitCode field.
//
// When the provided Service type is nil, it
// will set nothing and immediately return.
func (s *Service) SetExitCode(v int32) {
	// return if Service type is nil
	if s == nil {
		return
	}

	s.ExitCode = &v
}

// SetCreated sets the Created field.
//
// When the provided Service type is nil, it
// will set nothing and immediately return.
func (s *Service) SetCreated(v int64) {
	// return if Service type is nil
	if s == nil {
		return
	}

	s.Created = &v
}

// SetStarted sets the Started field.
//
// When the provided Service type is nil, it
// will set nothing and immediately return.
func (s *Service) SetStarted(v int64) {
	// return if Service type is nil
	if s == nil {
		return
	}

	s.Started = &v
}

// SetFinished sets the Finished field.
//
// When the provided Service type is nil, it
// will set nothing and immediately return.
func (s *Service) SetFinished(v int64) {
	// return if Service type is nil
	if s == nil {
		return
	}

	s.Finished = &v
}

// SetHost sets the Host field.
//
// When the provided Service type is nil, it
// will set nothing and immediately return.
func (s *Service) SetHost(v string) {
	// return if Service type is nil
	if s == nil {
		return
	}

	s.Host = &v
}

// SetRuntime sets the Runtime field.
//
// When the provided Service type is nil, it
// will set nothing and immediately return.
func (s *Service) SetRuntime(v string) {
	// return if Service type is nil
	if s == nil {
		return
	}

	s.Runtime = &v
}

// SetDistribution sets the Runtime field.
//
// When the provided Service type is nil, it
// will set nothing and immediately return.
func (s *Service) SetDistribution(v string) {
	// return if Service type is nil
	if s == nil {
		return
	}

	s.Distribution = &v
}

// String implements the Stringer interface for the Service type.
func (s *Service) String() string {
	return fmt.Sprintf(`{
  BuildID: %d,
  Created: %d,
  Distribution: %s,
  Error: %s,
  ExitCode: %d,
  Finished: %d,
  Host: %s,
  ID: %d,
  Image: %s,
  Name: %s,
  Number: %d,
  RepoID: %d,
  Runtime: %s,
  Started: %d,
  Status: %s,
}`,
		s.GetBuildID(),
		s.GetCreated(),
		s.GetDistribution(),
		s.GetError(),
		s.GetExitCode(),
		s.GetFinished(),
		s.GetHost(),
		s.GetID(),
		s.GetImage(),
		s.GetName(),
		s.GetNumber(),
		s.GetRepoID(),
		s.GetRuntime(),
		s.GetStarted(),
		s.GetStatus(),
	)
}

// ServiceFromBuildContainer creates a new Service based on a Build and pipeline Container.
func ServiceFromBuildContainer(build *Build, ctn *pipeline.Container) *Service {
	// create new service type we want to return
	s := new(Service)

	// default status to Pending
	s.SetStatus(constants.StatusPending)

	// copy fields from build
	if build != nil {
		// set values from the build
		s.SetHost(build.GetHost())
		s.SetRuntime(build.GetRuntime())
		s.SetDistribution(build.GetDistribution())
	}

	// copy fields from container
	if ctn != nil && ctn.Name != "" {
		// set values from the container
		s.SetName(ctn.Name)
		s.SetNumber(ctn.Number)
		s.SetImage(ctn.Image)
	}

	return s
}

// ServiceFromContainerEnvironment converts the pipeline Container
// to an API Service using the container's Environment.
func ServiceFromContainerEnvironment(ctn *pipeline.Container) *Service {
	// check if container or container environment are nil
	if ctn == nil || ctn.Environment == nil {
		return nil
	}

	// create new service type we want to return
	s := new(Service)

	// check if the VELA_SERVICE_DISTRIBUTION environment variable exists
	value, ok := ctn.Environment["VELA_SERVICE_DISTRIBUTION"]
	if ok {
		// set the Distribution field to the value from environment variable
		s.SetDistribution(value)
	}

	// check if the VELA_SERVICE_HOST environment variable exists
	value, ok = ctn.Environment["VELA_SERVICE_HOST"]
	if ok {
		// set the Host field to the value from environment variable
		s.SetHost(value)
	}

	// check if the VELA_SERVICE_IMAGE environment variable exists
	value, ok = ctn.Environment["VELA_SERVICE_IMAGE"]
	if ok {
		// set the Image field to the value from environment variable
		s.SetImage(value)
	}

	// check if the VELA_SERVICE_NAME environment variable exists
	value, ok = ctn.Environment["VELA_SERVICE_NAME"]
	if ok {
		// set the Name field to the value from environment variable
		s.SetName(value)
	}

	// check if the VELA_SERVICE_RUNTIME environment variable exists
	value, ok = ctn.Environment["VELA_SERVICE_RUNTIME"]
	if ok {
		// set the Runtime field to the value from environment variable
		s.SetRuntime(value)
	}

	// check if the VELA_SERVICE_STATUS environment variable exists
	value, ok = ctn.Environment["VELA_SERVICE_STATUS"]
	if ok {
		// set the Status field to the value from environment variable
		s.SetStatus(value)
	}

	// check if the VELA_SERVICE_CREATED environment variable exists
	value, ok = ctn.Environment["VELA_SERVICE_CREATED"]
	if ok {
		// parse the environment variable value into an int64
		i, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			// set the Created field to the parsed int64
			s.SetCreated(i)
		}
	}

	// check if the VELA_SERVICE_EXIT_CODE environment variable exists
	value, ok = ctn.Environment["VELA_SERVICE_EXIT_CODE"]
	if ok {
		// parse the environment variable value into an int
		i, err := strconv.ParseInt(value, 10, 32)
		if err == nil {
			// set the ExitCode field to the parsed int
			s.SetExitCode(int32(i))
		}
	}

	// check if the VELA_SERVICE_FINISHED environment variable exists
	value, ok = ctn.Environment["VELA_SERVICE_FINISHED"]
	if ok {
		// parse the environment variable value into an int64
		i, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			// set the Finished field to the parsed int64
			s.SetFinished(i)
		}
	}

	// check if the VELA_SERVICE_NUMBER environment variable exists
	value, ok = ctn.Environment["VELA_SERVICE_NUMBER"]
	if ok {
		// parse the environment variable value into an int
		i, err := strconv.ParseInt(value, 10, 16)
		if err == nil {
			// set the Number field to the parsed int
			s.SetNumber(int32(i))
		}
	}

	// check if the VELA_SERVICE_STARTED environment variable exists
	value, ok = ctn.Environment["VELA_SERVICE_STARTED"]
	if ok {
		// parse the environment variable value into an int64
		i, err := strconv.ParseInt(value, 10, 64)
		if err == nil {
			// set the Started field to the parsed int64
			s.SetStarted(i)
		}
	}

	return s
}
