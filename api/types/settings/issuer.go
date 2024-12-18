// SPDX-License-Identifier: Apache-2.0

package settings

import "fmt"

// OIDCIssuerSlice is an alias for a slice of OIDCIssuer types.
type OIDCIssuerSlice []*OIDCIssuer

// OIDCIssuer is the API representation of the OIDCIssuer setting.
type OIDCIssuer struct {
	Issuer      *string `json:"issuer,omitempty"         yaml:"issuer,omitempty"`
	UsernameMap *string `json:"username_map,omitempty"      yaml:"username_map,omitempty"`
	Redirect    *string `json:"redirect,omitempty" yaml:"redirect,omitempty"`
}

// GetIssuer returns the Issuer field.
//
// When the provided OIDCIssuer type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (oi *OIDCIssuer) GetIssuer() string {
	// return zero value if OIDCIssuer type or Issuer field is nil
	if oi == nil || oi.Issuer == nil {
		return ""
	}

	return *oi.Issuer
}

// GetUsernameMap returns the UsernameMap field.
//
// When the provided OIDCIssuer type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (oi *OIDCIssuer) GetUsernameMap() string {
	// return zero value if OIDCIssuer type or UsernameMap field is nil
	if oi == nil || oi.UsernameMap == nil {
		return ""
	}

	return *oi.UsernameMap
}

// GetRedirect returns the Redirect field.
//
// When the provided OIDCIssuer type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (oi *OIDCIssuer) GetRedirect() string {
	// return zero value if OIDCIssuer type or Redirect field is nil
	if oi == nil || oi.Redirect == nil {
		return ""
	}

	return *oi.Redirect
}

// SetIssuer sets the Issuer field.
//
// When the provided OIDCIssuer type is nil, it
// will set nothing and immediately return.
func (oi *OIDCIssuer) SetIssuer(v string) {
	// return if OIDCIssuer type is nil
	if oi == nil {
		return
	}

	oi.Issuer = &v
}

// SetUsernameMap sets the UsernameMap field.
//
// When the provided OIDCIssuer type is nil, it
// will set nothing and immediately return.
func (oi *OIDCIssuer) SetUsernameMap(v string) {
	// return if OIDCIssuer type is nil
	if oi == nil {
		return
	}

	oi.UsernameMap = &v
}

// SetRedirect sets the Redirect field.
//
// When the provided OIDCIssuer type is nil, it
// will set nothing and immediately return.
func (oi *OIDCIssuer) SetRedirect(v string) {
	// return if OIDCIssuer type is nil
	if oi == nil {
		return
	}

	oi.Redirect = &v
}

// String implements the Stringer interface for the OIDCIssuer type.
func (oi *OIDCIssuer) String() string {
	return fmt.Sprintf(`{
  Issuer: %s,
  UsernameMap: %s,
  Redirect: %s,
}`,
		oi.GetIssuer(),
		oi.GetUsernameMap(),
		oi.GetRedirect(),
	)
}

// OIDCIssuerMockEmpty returns an empty OIDCIssuer type.
func OIDCIssuerMockEmpty() OIDCIssuer {
	oi := OIDCIssuer{}
	oi.SetIssuer("")
	oi.SetUsernameMap("")
	oi.SetRedirect("")

	return oi
}
