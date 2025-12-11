// SPDX-License-Identifier: Apache-2.0

package types

import "fmt"

// Token is the API representation of a token response from server.
//
// swagger:model Token
type Token struct {
	Token      *string `json:"token,omitempty"`
	Expiration *int64  `json:"expiration,omitempty"`
}

// GetToken returns the Token field.
//
// When the provided Token type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (l *Token) GetToken() string {
	// return zero value if Token type or Token field is nil
	if l == nil || l.Token == nil {
		return ""
	}

	return *l.Token
}

// GetExpiration returns the Expiration field.
//
// When the provided Expiration type is nil, or the field within
// the type is nil, it returns the zero value for the field.
func (l *Token) GetExpiration() int64 {
	// return zero value if Expiration type or Expiration field is nil
	if l == nil || l.Expiration == nil {
		return 0
	}

	return *l.Expiration
}

// SetToken sets the Token field.
//
// When the provided Token type is nil, it
// will set nothing and immediately return.
func (l *Token) SetToken(v string) {
	// return if Token type is nil
	if l == nil {
		return
	}

	l.Token = &v
}

// SetExpiration sets the Expiration field.
//
// When the provided Expiration type is nil, it
// will set nothing and immediately return.
func (l *Token) SetExpiration(v int64) {
	// return if Expiration type is nil
	if l == nil {
		return
	}

	l.Expiration = &v
}

// String implements the Stringer interface for the Token type.
func (l *Token) String() string {
	return fmt.Sprintf(`{
  Token: %s,
  Expiration: %d
}`,
		l.GetToken(),
		l.GetExpiration(),
	)
}
