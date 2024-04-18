// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
)

// User is the API representation of a user.
//
// swagger:model User
type User struct {
	ID           *int64    `json:"id,omitempty"`
	Name         *string   `json:"name,omitempty"`
	RefreshToken *string   `json:"-"`
	Token        *string   `json:"-"`
	Favorites    *[]string `json:"favorites,omitempty"`
	Active       *bool     `json:"active,omitempty"`
	Admin        *bool     `json:"admin,omitempty"`
	Dashboards   *[]string `json:"dashboards,omitempty"`
}

// CropPreferences removes personal info from a user.
func (u *User) CropPreferences() *User {
	return &User{
		ID:     u.ID,
		Name:   u.Name,
		Token:  u.Token,
		Active: u.Active,
	}
}

// Environment returns a list of environment variables
// provided from the fields of the User type.
func (u *User) Environment() map[string]string {
	return map[string]string{
		"VELA_USER_ACTIVE":    ToString(u.GetActive()),
		"VELA_USER_ADMIN":     ToString(u.GetAdmin()),
		"VELA_USER_FAVORITES": ToString(u.GetFavorites()),
		"VELA_USER_NAME":      ToString(u.GetName()),
	}
}

// GetID returns the ID field.
//
// When the provided User type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (u *User) GetID() int64 {
	// return zero value if User type or ID field is nil
	if u == nil || u.ID == nil {
		return 0
	}

	return *u.ID
}

// GetName returns the Name field.
//
// When the provided User type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (u *User) GetName() string {
	// return zero value if User type or Name field is nil
	if u == nil || u.Name == nil {
		return ""
	}

	return *u.Name
}

// GetRefreshToken returns the RefreshToken field.
//
// When the provided User type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (u *User) GetRefreshToken() string {
	// return zero value if User type or RefreshToken field is nil
	if u == nil || u.RefreshToken == nil {
		return ""
	}

	return *u.RefreshToken
}

// GetToken returns the Token field.
//
// When the provided User type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (u *User) GetToken() string {
	// return zero value if User type or Token field is nil
	if u == nil || u.Token == nil {
		return ""
	}

	return *u.Token
}

// GetActive returns the Active field.
//
// When the provided User type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (u *User) GetActive() bool {
	// return zero value if User type or Active field is nil
	if u == nil || u.Active == nil {
		return false
	}

	return *u.Active
}

// GetAdmin returns the Admin field.
//
// When the provided User type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (u *User) GetAdmin() bool {
	// return zero value if User type or Admin field is nil
	if u == nil || u.Admin == nil {
		return false
	}

	return *u.Admin
}

// GetFavorites returns the Favorites field.
//
// When the provided User type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (u *User) GetFavorites() []string {
	// return zero value if User type or Favorites field is nil
	if u == nil || u.Favorites == nil {
		return []string{}
	}

	return *u.Favorites
}

// GetDashboards returns the Dashboards field.
//
// When the provided User type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (u *User) GetDashboards() []string {
	// return zero value if User type or Favorites field is nil
	if u == nil || u.Dashboards == nil {
		return []string{}
	}

	return *u.Dashboards
}

// SetID sets the ID field.
//
// When the provided User type is nil, it
// will set nothing and immediately return.
func (u *User) SetID(v int64) {
	// return if User type is nil
	if u == nil {
		return
	}

	u.ID = &v
}

// SetName sets the Name field.
//
// When the provided User type is nil, it
// will set nothing and immediately return.
func (u *User) SetName(v string) {
	// return if User type is nil
	if u == nil {
		return
	}

	u.Name = &v
}

// SetRefreshToken sets the RefreshToken field.
//
// When the provided User type is nil, it
// will set nothing and immediately return.
func (u *User) SetRefreshToken(v string) {
	// return if User type is nil
	if u == nil {
		return
	}

	u.RefreshToken = &v
}

// SetToken sets the Token field.
//
// When the provided User type is nil, it
// will set nothing and immediately return.
func (u *User) SetToken(v string) {
	// return if User type is nil
	if u == nil {
		return
	}

	u.Token = &v
}

// SetActive sets the Active field.
//
// When the provided User type is nil, it
// will set nothing and immediately return.
func (u *User) SetActive(v bool) {
	// return if User type is nil
	if u == nil {
		return
	}

	u.Active = &v
}

// SetAdmin sets the Admin field.
//
// When the provided User type is nil, it
// will set nothing and immediately return.
func (u *User) SetAdmin(v bool) {
	// return if User type is nil
	if u == nil {
		return
	}

	u.Admin = &v
}

// SetFavorites sets the Favorites field.
//
// When the provided User type is nil, it
// will set nothing and immediately return.
func (u *User) SetFavorites(v []string) {
	// return if User type is nil
	if u == nil {
		return
	}

	u.Favorites = &v
}

// SetDashboard sets the Dashboard field.
//
// When the provided User type is nil, it
// will set nothing and immediately return.
func (u *User) SetDashboards(v []string) {
	// return if User type is nil
	if u == nil {
		return
	}

	u.Dashboards = &v
}

// String implements the Stringer interface for the User type.
func (u *User) String() string {
	return fmt.Sprintf(`{
  Active: %t,
  Admin: %t,
  Dashboards: %s,
  Favorites: %s,
  ID: %d,
  Name: %s,
  Token: %s,
}`,
		u.GetActive(),
		u.GetAdmin(),
		u.GetDashboards(),
		u.GetFavorites(),
		u.GetID(),
		u.GetName(),
		u.GetToken(),
	)
}
