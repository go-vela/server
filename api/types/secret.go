// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"strings"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
)

// Secret is the API representation of a secret.
//
// swagger:model Secret
type Secret struct {
	ID                *int64    `json:"id,omitempty"`
	Org               *string   `json:"org,omitempty"`
	OrgSCMID          *int64    `json:"org_scm_id,omitempty"`
	Repo              *string   `json:"repo,omitempty"`
	RepoSCMID         *int64    `json:"repo_scm_id,omitempty"`
	Team              *string   `json:"team,omitempty"`
	TeamSCMID         *int64    `json:"team_scm_id,omitempty"`
	Name              *string   `json:"name,omitempty"`
	Value             *string   `json:"value,omitempty"`
	Type              *string   `json:"type,omitempty"`
	Images            *[]string `json:"images,omitempty"`
	AllowEvents       *Events   `json:"allow_events,omitempty"       yaml:"allow_events"`
	AllowCommand      *bool     `json:"allow_command,omitempty"`
	AllowSubstitution *bool     `json:"allow_substitution,omitempty"`
	CreatedAt         *int64    `json:"created_at,omitempty"`
	CreatedBy         *string   `json:"created_by,omitempty"`
	UpdatedAt         *int64    `json:"updated_at,omitempty"`
	UpdatedBy         *string   `json:"updated_by,omitempty"`
}

// UnmarshalYAML implements the Unmarshaler interface for the Secret type.
// This allows custom fields in the Secret type to be read from a YAML file, like AllowEvents.
func (s *Secret) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// create an alias to perform a normal unmarshal and avoid an infinite loop
	type jsonSecret Secret

	tmp := &jsonSecret{}

	err := unmarshal(tmp)
	if err != nil {
		return err
	}

	// overwrite existing secret
	*s = Secret(*tmp)

	return nil
}

// Sanitize creates a duplicate of the Secret without the value.
func (s *Secret) Sanitize() *Secret {
	// create a variable since constants can not be addressable
	//
	// https://golang.org/ref/spec#Address_operators
	value := constants.SecretMask

	return &Secret{
		ID:                s.ID,
		Org:               s.Org,
		Repo:              s.Repo,
		Team:              s.Team,
		Name:              s.Name,
		Value:             &value,
		Type:              s.Type,
		Images:            s.Images,
		AllowEvents:       s.AllowEvents,
		AllowCommand:      s.AllowCommand,
		AllowSubstitution: s.AllowSubstitution,
		CreatedAt:         s.CreatedAt,
		CreatedBy:         s.CreatedBy,
		UpdatedAt:         s.UpdatedAt,
		UpdatedBy:         s.UpdatedBy,
	}
}

// Match returns true when the provided container matches
// the conditions to inject a secret into a pipeline container
// resource.
func (s *Secret) Match(from *pipeline.Container) bool {
	eACL, iACL := false, false
	images, commands := s.GetImages(), s.GetAllowCommand()

	// check if commands are utilized when not allowed
	if !commands && len(from.Commands) > 0 {
		return false
	}

	// check if a custom entrypoint is utilized when not allowed
	if !commands && len(from.Commands) == 0 && len(from.Entrypoint) > 0 {
		return false
	}

	eACL = s.GetAllowEvents().Allowed(
		from.Environment["VELA_BUILD_EVENT"],
		from.Environment["VELA_BUILD_EVENT_ACTION"],
	)

	// check images whitelist
	for _, i := range images {
		if strings.HasPrefix(from.Image, i) && (len(i) != 0) {
			iACL = true
			break
		}
	}

	// inject secrets into environment
	switch {
	case eACL && (len(images) == 0):
		return true
	case eACL && iACL:
		return true
	}

	// return false if not match is found
	return false
}

// GetID returns the ID field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetID() int64 {
	// return zero value if Secret type or ID field is nil
	if s == nil || s.ID == nil {
		return 0
	}

	return *s.ID
}

// GetOrg returns the Org field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetOrg() string {
	// return zero value if Secret type or Org field is nil
	if s == nil || s.Org == nil {
		return ""
	}

	return *s.Org
}

// GetOrgSCMID returns the OrgSCMID field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetOrgSCMID() int64 {
	// return zero value if Secret type or OrgSCMID field is nil
	if s == nil || s.OrgSCMID == nil {
		return 0
	}

	return *s.OrgSCMID
}

// GetRepo returns the Repo field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetRepo() string {
	// return zero value if Secret type or Repo field is nil
	if s == nil || s.Repo == nil {
		return ""
	}

	return *s.Repo
}

// GetRepoSCMID returns the RepoSCMID field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetRepoSCMID() int64 {
	// return zero value if Secret type or RepoSCMID field is nil
	if s == nil || s.RepoSCMID == nil {
		return 0
	}

	return *s.RepoSCMID
}

// GetTeam returns the Team field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetTeam() string {
	// return zero value if Secret type or Team field is nil
	if s == nil || s.Team == nil {
		return ""
	}

	return *s.Team
}

// GetTeamSCMID returns the TeamSCMID field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetTeamSCMID() int64 {
	// return zero value if Secret type or TeamSCMID field is nil
	if s == nil || s.TeamSCMID == nil {
		return 0
	}

	return *s.TeamSCMID
}

// GetName returns the Name field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetName() string {
	// return zero value if Secret type or Name field is nil
	if s == nil || s.Name == nil {
		return ""
	}

	return *s.Name
}

// GetValue returns the Value field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetValue() string {
	// return zero value if Secret type or Value field is nil
	if s == nil || s.Value == nil {
		return ""
	}

	return *s.Value
}

// GetType returns the Type field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetType() string {
	// return zero value if Secret type or Type field is nil
	if s == nil || s.Type == nil {
		return ""
	}

	return *s.Type
}

// GetImages returns the Images field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetImages() []string {
	// return zero value if Secret type or Images field is nil
	if s == nil || s.Images == nil {
		return []string{}
	}

	return *s.Images
}

// GetAllowEvents returns the AllowEvents field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetAllowEvents() *Events {
	// return zero value if Secret type or AllowEvents field is nil
	if s == nil || s.AllowEvents == nil {
		return new(Events)
	}

	return s.AllowEvents
}

// GetAllowCommand returns the AllowCommand field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetAllowCommand() bool {
	// return zero value if Secret type or Images field is nil
	if s == nil || s.AllowCommand == nil {
		return false
	}

	return *s.AllowCommand
}

// GetAllowSubstitution returns the AllowSubstitution field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetAllowSubstitution() bool {
	// return zero value if Secret type or AllowSubstitution field is nil
	if s == nil || s.AllowSubstitution == nil {
		return false
	}

	return *s.AllowSubstitution
}

// GetCreatedAt returns the CreatedAt field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetCreatedAt() int64 {
	// return zero value if Secret type or CreatedAt field is nil
	if s == nil || s.CreatedAt == nil {
		return 0
	}

	return *s.CreatedAt
}

// GetCreatedBy returns the CreatedBy field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetCreatedBy() string {
	// return zero value if Secret type or CreatedBy field is nil
	if s == nil || s.CreatedBy == nil {
		return ""
	}

	return *s.CreatedBy
}

// GetUpdatedAt returns the UpdatedAt field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetUpdatedAt() int64 {
	// return zero value if Secret type or UpdatedAt field is nil
	if s == nil || s.UpdatedAt == nil {
		return 0
	}

	return *s.UpdatedAt
}

// GetUpdatedBy returns the UpdatedBy field.
//
// When the provided Secret type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *Secret) GetUpdatedBy() string {
	// return zero value if Secret type or UpdatedBy field is nil
	if s == nil || s.UpdatedBy == nil {
		return ""
	}

	return *s.UpdatedBy
}

// SetID sets the ID field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetID(v int64) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.ID = &v
}

// SetOrg sets the Org field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetOrg(v string) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.Org = &v
}

// SetOrgSCMID sets the OrgSCMID field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetOrgSCMID(v int64) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.OrgSCMID = &v
}

// SetRepo sets the Repo field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetRepo(v string) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.Repo = &v
}

// SetRepoSCMID sets the RepoSCMID field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetRepoSCMID(v int64) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.RepoSCMID = &v
}

// SetTeam sets the Team field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetTeam(v string) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.Team = &v
}

// SetTeamSCMID sets the TeamSCMID field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetTeamSCMID(v int64) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.TeamSCMID = &v
}

// SetName sets the Name field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetName(v string) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.Name = &v
}

// SetValue sets the Value field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetValue(v string) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.Value = &v
}

// SetType sets the Type field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetType(v string) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.Type = &v
}

// SetImages sets the Images field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetImages(v []string) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.Images = &v
}

// SetAllowEvents sets the AllowEvents field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetAllowEvents(v *Events) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.AllowEvents = v
}

// SetAllowCommand sets the AllowCommand field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetAllowCommand(v bool) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.AllowCommand = &v
}

// SetAllowSubstitution sets the AllowSubstitution field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetAllowSubstitution(v bool) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.AllowSubstitution = &v
}

// SetCreatedAt sets the CreatedAt field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetCreatedAt(v int64) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.CreatedAt = &v
}

// SetCreatedBy sets the CreatedBy field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetCreatedBy(v string) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.CreatedBy = &v
}

// SetUpdatedAt sets the UpdatedAt field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetUpdatedAt(v int64) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.UpdatedAt = &v
}

// SetUpdatedBy sets the UpdatedBy field.
//
// When the provided Secret type is nil, it
// will set nothing and immediately return.
func (s *Secret) SetUpdatedBy(v string) {
	// return if Secret type is nil
	if s == nil {
		return
	}

	s.UpdatedBy = &v
}

// String implements the Stringer interface for the Secret type.
func (s *Secret) String() string {
	return fmt.Sprintf(`{
	AllowCommand: %t,
	AllowEvents: %s,
	AllowSubstitution: %t,
	ID: %d,
	Images: %s,
	Name: %s,
	Org: %s,
	Repo: %s,
	Team: %s,
	Type: %s,
	Value: %s,
	CreatedAt: %d,
	CreatedBy: %s,
	UpdatedAt: %d,
	UpdatedBy: %s,
}`,
		s.GetAllowCommand(),
		s.GetAllowEvents().List(),
		s.GetAllowSubstitution(),
		s.GetID(),
		s.GetImages(),
		s.GetName(),
		s.GetOrg(),
		s.GetRepo(),
		s.GetTeam(),
		s.GetType(),
		s.GetValue(),
		s.GetCreatedAt(),
		s.GetCreatedBy(),
		s.GetUpdatedAt(),
		s.GetUpdatedBy(),
	)
}
