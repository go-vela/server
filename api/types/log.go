// SPDX-License-Identifier: Apache-2.0

package types

import (
	"bytes"
	"fmt"

	"github.com/go-vela/server/constants"
)

// Log is the API representation of a log for a step in a build.
//
// swagger:model Log
type Log struct {
	ID        *int64 `json:"id,omitempty"`
	BuildID   *int64 `json:"build_id,omitempty"`
	RepoID    *int64 `json:"repo_id,omitempty"`
	ServiceID *int64 `json:"service_id,omitempty"`
	StepID    *int64 `json:"step_id,omitempty"`
	// swagger:strfmt base64
	Data      *[]byte `json:"data,omitempty"`
	CreatedAt *int64  `json:"created_at,omitempty"`
}

// AppendData adds the provided data to the end of
// the Data field for the Log type. If the Data
// field is empty, then the function overwrites
// the entire Data field.
func (l *Log) AppendData(data []byte) {
	// check if Data field is empty
	if len(l.GetData()) == 0 {
		// overwrite the Data field
		l.SetData(data)

		return
	}

	// add the data to the Data field
	l.SetData(append(l.GetData(), data...))
}

// MaskData reads through the log data and masks
// all values provided in the string slice. If the
// log is empty, we do nothing.
func (l *Log) MaskData(secrets []string) {
	data := l.GetData()

	// early exit on empty log or secret list
	if len(data) == 0 || len(secrets) == 0 {
		return
	}

	// byte replace data with masked logs
	for _, secret := range secrets {
		data = bytes.ReplaceAll(data, []byte(secret), []byte(constants.SecretLogMask))
	}

	// update data field to masked logs
	l.SetData(data)
}

// GetID returns the ID field.
//
// When the provided Log type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (l *Log) GetID() int64 {
	// return zero value if Log type or ID field is nil
	if l == nil || l.ID == nil {
		return 0
	}

	return *l.ID
}

// GetBuildID returns the BuildID field.
//
// When the provided Log type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (l *Log) GetBuildID() int64 {
	// return zero value if Log type or BuildID field is nil
	if l == nil || l.BuildID == nil {
		return 0
	}

	return *l.BuildID
}

// GetRepoID returns the RepoID field.
//
// When the provided Log type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (l *Log) GetRepoID() int64 {
	// return zero value if Log type or RepoID field is nil
	if l == nil || l.RepoID == nil {
		return 0
	}

	return *l.RepoID
}

// GetServiceID returns the ServiceID field.
//
// When the provided Log type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (l *Log) GetServiceID() int64 {
	// return zero value if Log type or ServiceID field is nil
	if l == nil || l.ServiceID == nil {
		return 0
	}

	return *l.ServiceID
}

// GetStepID returns the StepID field.
//
// When the provided Log type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (l *Log) GetStepID() int64 {
	// return zero value if Log type or StepID field is nil
	if l == nil || l.StepID == nil {
		return 0
	}

	return *l.StepID
}

// GetData returns the Data field.
//
// When the provided Log type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (l *Log) GetData() []byte {
	// return zero value if Log type or Data field is nil
	if l == nil || l.Data == nil {
		return []byte{}
	}

	return *l.Data
}

// GetCreatedAt returns the CreatedAt field.
//
// When the provided log type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (l *Log) GetCreatedAt() int64 {
	// return zero value if log type or CreatedAt field is nil
	if l == nil || l.CreatedAt == nil {
		return 0
	}

	return *l.CreatedAt
}

// SetID sets the ID field.
//
// When the provided Log type is nil, it
// will set nothing and immediately return.
func (l *Log) SetID(v int64) {
	// return if Log type is nil
	if l == nil {
		return
	}

	l.ID = &v
}

// SetBuildID sets the BuildID field.
//
// When the provided Log type is nil, it
// will set nothing and immediately return.
func (l *Log) SetBuildID(v int64) {
	// return if Log type is nil
	if l == nil {
		return
	}

	l.BuildID = &v
}

// SetRepoID sets the RepoID field.
//
// When the provided Log type is nil, it
// will set nothing and immediately return.
func (l *Log) SetRepoID(v int64) {
	// return if Log type is nil
	if l == nil {
		return
	}

	l.RepoID = &v
}

// SetServiceID sets the ServiceID field.
//
// When the provided Log type is nil, it
// will set nothing and immediately return.
func (l *Log) SetServiceID(v int64) {
	// return if Log type is nil
	if l == nil {
		return
	}

	l.ServiceID = &v
}

// SetStepID sets the StepID field.
//
// When the provided Log type is nil, it
// will set nothing and immediately return.
func (l *Log) SetStepID(v int64) {
	// return if Log type is nil
	if l == nil {
		return
	}

	l.StepID = &v
}

// SetData sets the Data field.
//
// When the provided Log type is nil, it
// will set nothing and immediately return.
func (l *Log) SetData(v []byte) {
	// return if Log type is nil
	if l == nil {
		return
	}

	l.Data = &v
}

// SetCreatedAt sets the CreatedAt field.
//
// When the provided log type is nil, it
// will set nothing and immediately return.
func (l *Log) SetCreatedAt(v int64) {
	// return if log type is nil
	if l == nil {
		return
	}

	l.CreatedAt = &v
}

// String implements the Stringer interface for the Log type.
func (l *Log) String() string {
	return fmt.Sprintf(`{
  BuildID: %d,
  Data: %s,
  ID: %d,
  RepoID: %d,
  ServiceID: %d,
  StepID: %d,
  CreatedAt: %d,
}`,
		l.GetBuildID(),
		l.GetData(),
		l.GetID(),
		l.GetRepoID(),
		l.GetServiceID(),
		l.GetStepID(),
		l.GetCreatedAt(),
	)
}

// LogCleanupResponse represents the response body for log cleanup operations.
//
// swagger:model LogCleanupResponse
type LogCleanupResponse struct {
	DeletedCount       int64    `json:"deleted_count"`
	BatchesProcessed   int64    `json:"batches_processed"`
	DurationSeconds    float64  `json:"duration_seconds"`
	VacuumPerformed    bool     `json:"vacuum_performed"`
	PartitionedMode    bool     `json:"partitioned_mode"`
	AffectedPartitions []string `json:"affected_partitions,omitempty"`
	Message            string   `json:"message"`
}
