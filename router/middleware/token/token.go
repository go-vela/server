// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package token

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-vela/server/database"
	"github.com/go-vela/types/library"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

// Compose generates the token unique to the provided user.
// It uses a secret hash, which is unique for every user, to
// sign the token which guarantees the signature is unique
// per token. The signed token is returned as a string.
func Compose(u *library.User) (string, error) {
	// generate a token using dgrijalva/jwt-go
	t := jwt.New(jwt.SigningMethodHS256)

	// extract the claims from the token
	claims := t.Claims.(jwt.MapClaims)

	// append extra metadata to token claims
	claims["active"] = u.Active
	claims["admin"] = u.Admin
	claims["name"] = u.Name

	// sign the token using our secret key
	token, err := t.SignedString([]byte(u.GetHash()))
	if err != nil {
		return "", err
	}

	return token, nil
}

// Parse scans the signed JWT token as a string and extracts
// the user login from the claims to be looked up in the database.
// This function will return an error for a few different reasons:
//
// * the token signature doesn't match what is expected
// * the token signing method doesn't match what is expected
// * the token is invalid (potentially expired or improper)
func Parse(t string, db database.Service) (*library.User, error) {
	u := new(library.User)

	// parse the signed JWT token string
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		var err error

		// extract the claims from the token
		claims := token.Claims.(jwt.MapClaims)
		name := claims["name"].(string)

		// lookup the user in the database
		logrus.Debugf("Reading user %s", name)
		u, err = db.GetUserName(name)
		return []byte(u.GetHash()), err
	})

	if err != nil {
		return nil, fmt.Errorf("unable to parse JWT token: %v", err)
	}

	// validate the correct signing method is being used
	if token.Method != jwt.SigningMethodHS256 {
		return nil, jwt.ErrSignatureInvalid
	}

	// ensure the token is valid
	if !token.Valid {
		return nil, fmt.Errorf("invalid JWT token provided for %s", u.GetName())
	}

	return u, nil
}

// Retrieve gets the token from the provided request http.Request
// to be parsed and validated. This is called on every request
// to enable capturing the user making the request and validating
// they have the proper access. The following methods of providing
// authentication to Vela are supported:
//
// * Bearer token in `Authorization` header
func Retrieve(r *http.Request) (string, error) {
	// get the token from the `Authorization` header
	token := r.Header.Get("Authorization")
	if len(token) > 0 {
		if strings.Contains(token, "Bearer") {
			return strings.Split(token, "Bearer ")[1], nil
		}
	}

	// This code is commented out to reduce the amount of methods
	// Vela supports for authentication. If these other methods of
	// providing authentication are found valuable, we're willing
	// to enable those use cases.
	//
	// // get the token from the `access_token` query parameter
	// token = r.FormValue("access_token")
	// if len(token) > 0 {
	// 	return token, nil
	// }
	//
	// // get the token from the `vela_session` browser cookie
	// cookie, err := r.Cookie("vela_session")
	// if err == nil {
	// 	return cookie.Value, nil
	// }

	return "", fmt.Errorf("No token provided in Authorization header")
}
