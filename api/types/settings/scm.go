// SPDX-License-Identifier: Apache-2.0

package settings

import "fmt"

type SCM struct {
	RepoRoleMap map[string]string `json:"repo_role_map,omitempty" yaml:"repo_role_map,omitempty"`
	OrgRoleMap  map[string]string `json:"org_role_map,omitempty"  yaml:"org_role_map,omitempty"`
	TeamRoleMap map[string]string `json:"team_role_map,omitempty" yaml:"team_role_map,omitempty"`
}

// GetRepoRoleMap returns the RepoRoleMap field.
//
// When the provided SCM type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *SCM) GetRepoRoleMap() map[string]string {
	// return zero value if SCM type or RepoRoleMap field is nil
	if s == nil || s.RepoRoleMap == nil {
		return map[string]string{}
	}

	return s.RepoRoleMap
}

// GetOrgRoleMap returns the OrgRoleMap field.
//
// When the provided SCM type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *SCM) GetOrgRoleMap() map[string]string {
	// return zero value if SCM type or OrgRoleMap field is nil
	if s == nil || s.OrgRoleMap == nil {
		return map[string]string{}
	}

	return s.OrgRoleMap
}

// GetTeamRoleMap returns the TeamRoleMap field.
//
// When the provided SCM type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (s *SCM) GetTeamRoleMap() map[string]string {
	// return zero value if SCM type or TeamRoleMap field is nil
	if s == nil || s.TeamRoleMap == nil {
		return map[string]string{}
	}

	return s.TeamRoleMap
}

// SetRepoRoleMap sets the RepoRoleMap field.
//
// When the provided SCM type is nil, it
// will set nothing and immediately return.
func (s *SCM) SetRepoRoleMap(v map[string]string) {
	// return if SCM type is nil
	if s == nil {
		return
	}

	s.RepoRoleMap = v
}

// SetOrgRoleMap sets the OrgRoleMap field.
//
// When the provided SCM type is nil, it
// will set nothing and immediately return.
func (s *SCM) SetOrgRoleMap(v map[string]string) {
	// return if SCM type is nil
	if s == nil {
		return
	}

	s.OrgRoleMap = v
}

// SetTeamRoleMap sets the TeamRoleMap field.
//
// When the provided SCM type is nil, it
// will set nothing and immediately return.
func (s *SCM) SetTeamRoleMap(v map[string]string) {
	// return if SCM type is nil
	if s == nil {
		return
	}

	s.TeamRoleMap = v
}

// String implements the Stringer interface for the SCM type.
func (s *SCM) String() string {
	return fmt.Sprintf(`{
  RepoRoleMap: %v,
  OrgRoleMap: %v,
  TeamRoleMap: %v,
}`,
		s.GetRepoRoleMap(),
		s.GetOrgRoleMap(),
		s.GetTeamRoleMap(),
	)
}

// SCMMockEmpty returns an empty SCM type.
func SCMMockEmpty() SCM {
	s := SCM{}
	s.SetRepoRoleMap(map[string]string{})
	s.SetOrgRoleMap(map[string]string{})
	s.SetTeamRoleMap(map[string]string{})

	return s
}
