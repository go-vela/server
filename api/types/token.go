// SPDX-License-Identifier: Apache-2.0

package types

import "fmt"

// Token is the API representation of a token response from server.
//
// swagger:model Token
type Token struct {
	Token *string `json:"token,omitempty"`
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

// String implements the Stringer interface for the Token type.
func (l *Token) String() string {
	return fmt.Sprintf(`{
  Token: %s,
}`,
		l.GetToken(),
	)
}
