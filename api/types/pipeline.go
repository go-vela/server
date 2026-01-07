// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
)

// Pipeline is the API representation of a Pipeline.
//
// swagger:model Pipeline
type Pipeline struct {
	ID              *int64    `json:"id,omitempty"`
	Repo            *Repo     `json:"repo,omitempty"`
	Commit          *string   `json:"commit,omitempty"`
	Flavor          *string   `json:"flavor,omitempty"`
	Platform        *string   `json:"platform,omitempty"`
	Ref             *string   `json:"ref,omitempty"`
	Type            *string   `json:"type,omitempty"`
	Version         *string   `json:"version,omitempty"`
	ExternalSecrets *bool     `json:"external_secrets,omitempty"`
	InternalSecrets *bool     `json:"internal_secrets,omitempty"`
	Services        *bool     `json:"services,omitempty"`
	Stages          *bool     `json:"stages,omitempty"`
	Steps           *bool     `json:"steps,omitempty"`
	Templates       *bool     `json:"templates,omitempty"`
	TestReport      *bool     `json:"test_report,omitempty"`
	Warnings        *[]string `json:"warnings,omitempty"`
	// swagger:strfmt base64
	Data *[]byte `json:"data,omitempty"`
}

// GetID returns the ID field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetID() int64 {
	// return zero value if Pipeline type or ID field is nil
	if p == nil || p.ID == nil {
		return 0
	}

	return *p.ID
}

// GetRepo returns the Repo field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetRepo() *Repo {
	// return zero value if Pipeline type or Repo field is nil
	if p == nil || p.Repo == nil {
		return new(Repo)
	}

	return p.Repo
}

// GetCommit returns the Commit field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetCommit() string {
	// return zero value if Pipeline type or Commit field is nil
	if p == nil || p.Commit == nil {
		return ""
	}

	return *p.Commit
}

// GetFlavor returns the Flavor field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetFlavor() string {
	// return zero value if Pipeline type or Flavor field is nil
	if p == nil || p.Flavor == nil {
		return ""
	}

	return *p.Flavor
}

// GetPlatform returns the Platform field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetPlatform() string {
	// return zero value if Pipeline type or Platform field is nil
	if p == nil || p.Platform == nil {
		return ""
	}

	return *p.Platform
}

// GetRef returns the Ref field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetRef() string {
	// return zero value if Pipeline type or Ref field is nil
	if p == nil || p.Ref == nil {
		return ""
	}

	return *p.Ref
}

// GetType returns the Type field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetType() string {
	// return zero value if Pipeline type or Type field is nil
	if p == nil || p.Type == nil {
		return ""
	}

	return *p.Type
}

// GetVersion returns the Version field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetVersion() string {
	// return zero value if Pipeline type or Version field is nil
	if p == nil || p.Version == nil {
		return ""
	}

	return *p.Version
}

// GetExternalSecrets returns the ExternalSecrets field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetExternalSecrets() bool {
	// return zero value if Pipeline type or ExternalSecrets field is nil
	if p == nil || p.ExternalSecrets == nil {
		return false
	}

	return *p.ExternalSecrets
}

// GetInternalSecrets returns the InternalSecrets field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetInternalSecrets() bool {
	// return zero value if Pipeline type or InternalSecrets field is nil
	if p == nil || p.InternalSecrets == nil {
		return false
	}

	return *p.InternalSecrets
}

// GetServices returns the Services field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetServices() bool {
	// return zero value if Pipeline type or Services field is nil
	if p == nil || p.Services == nil {
		return false
	}

	return *p.Services
}

// GetStages returns the Stages field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetStages() bool {
	// return zero value if Pipeline type or Stages field is nil
	if p == nil || p.Stages == nil {
		return false
	}

	return *p.Stages
}

// GetSteps returns the Steps field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetSteps() bool {
	// return zero value if Pipeline type or Steps field is nil
	if p == nil || p.Steps == nil {
		return false
	}

	return *p.Steps
}

// GetTemplates returns the Templates field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetTemplates() bool {
	// return zero value if Pipeline type or Templates field is nil
	if p == nil || p.Templates == nil {
		return false
	}

	return *p.Templates
}

// GetWarnings returns the Warnings field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetWarnings() []string {
	// return zero value if Pipeline type or Warnings field is nil
	if p == nil || p.Warnings == nil {
		return []string{}
	}

	return *p.Warnings
}

// GetData returns the Data field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetData() []byte {
	// return zero value if Pipeline type or Data field is nil
	if p == nil || p.Data == nil {
		return []byte{}
	}

	return *p.Data
}

// GetTestReport returns the TestReport results field.
//
// When the provided Pipeline type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (p *Pipeline) GetTestReport() bool {
	// return zero value if Pipeline type or Artifacts field is nil
	if p == nil || p.TestReport == nil {
		return false
	}

	return *p.TestReport
}

// SetID sets the ID field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetID(v int64) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.ID = &v
}

// SetRepo sets the Repo field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetRepo(v *Repo) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.Repo = v
}

// SetCommit sets the Commit field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetCommit(v string) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.Commit = &v
}

// SetFlavor sets the Flavor field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetFlavor(v string) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.Flavor = &v
}

// SetPlatform sets the Platform field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetPlatform(v string) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.Platform = &v
}

// SetRef sets the Ref field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetRef(v string) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.Ref = &v
}

// SetType sets the Type field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetType(v string) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.Type = &v
}

// SetVersion sets the Version field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetVersion(v string) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.Version = &v
}

// SetExternalSecrets sets the ExternalSecrets field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetExternalSecrets(v bool) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.ExternalSecrets = &v
}

// SetInternalSecrets sets the InternalSecrets field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetInternalSecrets(v bool) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.InternalSecrets = &v
}

// SetServices sets the Services field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetServices(v bool) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.Services = &v
}

// SetStages sets the Stages field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetStages(v bool) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.Stages = &v
}

// SetSteps sets the Steps field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetSteps(v bool) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.Steps = &v
}

// SetTemplates sets the Templates field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetTemplates(v bool) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.Templates = &v
}

// SetTestReport sets the TestReport field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetTestReport(v bool) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.TestReport = &v
}

// SetWarnings sets the Warnings field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetWarnings(v []string) {
	// return if Pipeline type is nil
	if p == nil {
		return
	}

	p.Warnings = &v
}

// SetData sets the Data field.
//
// When the provided Pipeline type is nil, it
// will set nothing and immediately return.
func (p *Pipeline) SetData(v []byte) {
	// return if Log type is nil
	if p == nil {
		return
	}

	p.Data = &v
}

// String implements the Stringer interface for the Pipeline type.
func (p *Pipeline) String() string {
	return fmt.Sprintf(`{
  Commit: %s,
  Data: %s,
  Flavor: %s,
  ID: %d,
  Platform: %s,
  Ref: %s,
  Repo: %v,
  ExternalSecrets: %t,
  InternalSecrets: %t,
  Services: %t,
  Stages: %t,
  Steps: %t,
  Templates: %t,
  Artifacts: %t,
  Type: %s,
  Version: %s,
  Warnings: %v,
}`,
		p.GetCommit(),
		p.GetData(),
		p.GetFlavor(),
		p.GetID(),
		p.GetPlatform(),
		p.GetRef(),
		p.GetRepo(),
		p.GetExternalSecrets(),
		p.GetInternalSecrets(),
		p.GetServices(),
		p.GetStages(),
		p.GetSteps(),
		p.GetTemplates(),
		p.GetTestReport(),
		p.GetType(),
		p.GetVersion(),
		p.GetWarnings(),
	)
}
