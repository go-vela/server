// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package token

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

// ParseToken scans the signed JWT token as a string and extracts
// the user login from the claims to be looked up in the database.
// This function will return an error for a few different reasons:
//
// * the token signature doesn't match what is expected
// * the token signing method doesn't match what is expected
// * the token is invalid (potentially expired or improper).
func (tm *Manager) ParseToken(token string) (*Claims, error) {
	var claims = new(Claims)

	// create a new JWT parser
	p := &jwt.Parser{
		// explicitly only allow these signing methods
		ValidMethods: []string{jwt.SigningMethodHS256.Name},
	}

	// parse and validate given token
	tkn, err := p.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		var err error

		// extract the claims from the token
		claims = t.Claims.(*Claims)
		name := claims.Subject

		// according to JWT RFC, the iat field is optional for security purposes and is purely informational.
		// setting it to nil avoids any worries of race conditions.
		claims.IssuedAt = nil

		// check if subject has a value in claims;
		// we can save a db lookup attempt
		if len(name) == 0 {
			return nil, errors.New("no subject defined")
		}

		// ParseWithClaims will skip expiration check
		// if expiration has default value;
		// forcing a check and exiting if not set
		if claims.ExpiresAt == nil {
			return nil, errors.New("token has no expiration")
		}

		return []byte(tm.PrivateKey), err
	})

	if err != nil {
		return nil, errors.New("failed parsing: " + err.Error())
	}

	if !tkn.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
