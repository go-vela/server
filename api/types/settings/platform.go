// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"fmt"

	"github.com/urfave/cli/v3"
)

// Platform is the API representation of platform settingps.
//
// swagger:model Platform
type Platform struct {
	ID                *int32 `json:"id"`
	*Compiler         `json:"compiler,omitempty"            yaml:"compiler,omitempty"`
	*Queue            `json:"queue,omitempty"               yaml:"queue,omitempty"`
	*SCM              `json:"scm,omitempty"                 yaml:"scm,omitempty"`
	RepoAllowlist     *[]string `json:"repo_allowlist,omitempty"      yaml:"repo_allowlist,omitempty"`
	ScheduleAllowlist *[]string `json:"schedule_allowlist,omitempty"  yaml:"schedule_allowlist,omitempty"`
	MaxDashboardRepos *int32    `json:"max_dashboard_repos,omitempty" yaml:"max_dashboard_repos,omitempty"`
	QueueRestartLimit *int32    `json:"queue_restart_limit,omitempty" yaml:"queue_restart_limit,omitempty"`
	CreatedAt         *int64    `json:"created_at,omitempty"          yaml:"created_at,omitempty"`
	UpdatedAt         *int64    `json:"updated_at,omitempty"          yaml:"updated_at,omitempty"`
	UpdatedBy         *string   `json:"updated_by,omitempty"          yaml:"updated_by,omitempty"`
}

// FromCLICommand returns a new Platform record from a cli command.
func FromCLICommand(c *cli.Command) *Platform {
	ps := new(Platform)

	// set repos permitted to be added
	ps.SetRepoAllowlist(c.StringSlice("vela-repo-allowlist"))

	// set repos permitted to use schedules
	ps.SetScheduleAllowlist(c.StringSlice("vela-schedule-allowlist"))

	// set max repos per dashboard
	ps.SetMaxDashboardRepos(c.Int32("max-dashboard-repos"))

	// set queue restart limit
	ps.SetQueueRestartLimit(c.Int32("queue-restart-limit"))

	return ps
}

// GetID returns the ID field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetID() int32 {
	// return zero value if Platform type or ID field is nil
	if ps == nil || ps.ID == nil {
		return 0
	}

	return *ps.ID
}

// GetCompiler returns the Compiler field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetCompiler() Compiler {
	// return zero value if Platform type or Compiler field is nil
	if ps == nil || ps.Compiler == nil {
		return Compiler{}
	}

	return *ps.Compiler
}

// GetQueue returns the Queue field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetQueue() Queue {
	// return zero value if Platform type or Queue field is nil
	if ps == nil || ps.Queue == nil {
		return Queue{}
	}

	return *ps.Queue
}

// GetSCM returns the SCM field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetSCM() SCM {
	// return zero value if Platform type or SCM field is nil
	if ps == nil || ps.SCM == nil {
		return SCM{}
	}

	return *ps.SCM
}

// GetRepoAllowlist returns the RepoAllowlist field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetRepoAllowlist() []string {
	// return zero value if Platform type or RepoAllowlist field is nil
	if ps == nil || ps.RepoAllowlist == nil {
		return []string{}
	}

	return *ps.RepoAllowlist
}

// GetScheduleAllowlist returns the ScheduleAllowlist field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetScheduleAllowlist() []string {
	// return zero value if Platform type or ScheduleAllowlist field is nil
	if ps == nil || ps.ScheduleAllowlist == nil {
		return []string{}
	}

	return *ps.ScheduleAllowlist
}

// GetMaxDashboardRepos returns the MaxDashboardRepos field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetMaxDashboardRepos() int32 {
	// return zero value if Platform type or MaxDashboardRepos field is nil
	if ps == nil || ps.MaxDashboardRepos == nil {
		return 0
	}

	return *ps.MaxDashboardRepos
}

// GetQueueRestartLimit returns the QueueRestartLimit field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetQueueRestartLimit() int32 {
	// return zero value if Platform type or QueueRestartLimit field is nil
	if ps == nil || ps.QueueRestartLimit == nil {
		return 0
	}

	return *ps.QueueRestartLimit
}

// GetCreatedAt returns the CreatedAt field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetCreatedAt() int64 {
	// return zero value if Platform type or CreatedAt field is nil
	if ps == nil || ps.CreatedAt == nil {
		return 0
	}

	return *ps.CreatedAt
}

// GetUpdatedAt returns the UpdatedAt field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetUpdatedAt() int64 {
	// return zero value if Platform type or UpdatedAt field is nil
	if ps == nil || ps.UpdatedAt == nil {
		return 0
	}

	return *ps.UpdatedAt
}

// GetUpdatedBy returns the UpdatedBy field.
//
// When the provided Platform type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (ps *Platform) GetUpdatedBy() string {
	// return zero value if Platform type or UpdatedBy field is nil
	if ps == nil || ps.UpdatedBy == nil {
		return ""
	}

	return *ps.UpdatedBy
}

// SetID sets the ID field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetID(v int32) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.ID = &v
}

// SetCompiler sets the Compiler field.
//
// When the provided Compiler type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetCompiler(cs Compiler) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.Compiler = &cs
}

// SetQueue sets the Queue field.
//
// When the provided Queue type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetQueue(qs Queue) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.Queue = &qs
}

// SetSCM sets the SCM field.
//
// When the provided SCM type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetSCM(scm SCM) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.SCM = &scm
}

// SetRepoAllowlist sets the RepoAllowlist field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetRepoAllowlist(v []string) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.RepoAllowlist = &v
}

// SetScheduleAllowlist sets the RepoAllowlist field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetScheduleAllowlist(v []string) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.ScheduleAllowlist = &v
}

// SetMaxDashboardRepos sets the MaxDashboardRepos field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetMaxDashboardRepos(v int32) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.MaxDashboardRepos = &v
}

// SetQueueRestartLimit sets the QueueRestartLimit field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetQueueRestartLimit(v int32) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.QueueRestartLimit = &v
}

// SetCreatedAt sets the CreatedAt field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetCreatedAt(v int64) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.CreatedAt = &v
}

// SetUpdatedAt sets the UpdatedAt field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetUpdatedAt(v int64) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.UpdatedAt = &v
}

// SetUpdatedBy sets the UpdatedBy field.
//
// When the provided Platform type is nil, it
// will set nothing and immediately return.
func (ps *Platform) SetUpdatedBy(v string) {
	// return if Platform type is nil
	if ps == nil {
		return
	}

	ps.UpdatedBy = &v
}

// FromSettings takes another settings record and updates the internal fields,
// used when the updating settings and refreshing the record shared across the server.
func (ps *Platform) FromSettings(_ps *Platform) {
	if ps == nil {
		return
	}

	if _ps == nil {
		return
	}

	ps.SetCompiler(_ps.GetCompiler())
	ps.SetQueue(_ps.GetQueue())
	ps.SetSCM(_ps.GetSCM())
	ps.SetRepoAllowlist(_ps.GetRepoAllowlist())
	ps.SetScheduleAllowlist(_ps.GetScheduleAllowlist())
	ps.SetMaxDashboardRepos(_ps.GetMaxDashboardRepos())
	ps.SetQueueRestartLimit(_ps.GetQueueRestartLimit())

	ps.SetCreatedAt(_ps.GetCreatedAt())
	ps.SetUpdatedAt(_ps.GetUpdatedAt())
	ps.SetUpdatedBy(_ps.GetUpdatedBy())
}

// String implements the Stringer interface for the Platform type.
func (ps *Platform) String() string {
	cs := ps.GetCompiler()
	qs := ps.GetQueue()
	scms := ps.GetSCM()

	return fmt.Sprintf(`{
  ID: %d,
  Compiler: %v,
  Queue: %v,
  SCM: %v,
  RepoAllowlist: %v,
  ScheduleAllowlist: %v,
  MaxDashboardRepos: %d,
  QueueRestartLimit: %d,
  CreatedAt: %d,
  UpdatedAt: %d,
  UpdatedBy: %s,
}`,
		ps.GetID(),
		cs.String(),
		qs.String(),
		scms.String(),
		ps.GetRepoAllowlist(),
		ps.GetScheduleAllowlist(),
		ps.GetMaxDashboardRepos(),
		ps.GetQueueRestartLimit(),
		ps.GetCreatedAt(),
		ps.GetUpdatedAt(),
		ps.GetUpdatedBy(),
	)
}

// PlatformMockEmpty returns an empty Platform type.
func PlatformMockEmpty() Platform {
	ps := Platform{}

	ps.SetCompiler(CompilerMockEmpty())
	ps.SetQueue(QueueMockEmpty())
	ps.SetSCM(SCMMockEmpty())

	ps.SetRepoAllowlist([]string{})
	ps.SetScheduleAllowlist([]string{})
	ps.SetMaxDashboardRepos(0)
	ps.SetQueueRestartLimit(0)

	return ps
}
