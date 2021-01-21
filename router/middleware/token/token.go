// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package token

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/database"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/sirupsen/logrus"
)

type Claims struct {
	IsAdmin  bool `json:"is_admin"`
	IsActive bool `json:"is_active"`
	jwt.StandardClaims
}

// Compose generates an refresh and access token pair unique
// to the provided user and sets a secure cookie.
// It uses a secret hash, which is unique for every user.
// The hash signs the token to guarantee the signature is unique
// per token. The refresh token is returned to store with the user
// in the database.
// nolint:lll // reference links cause long lines
func Compose(c *gin.Context, u *library.User) (string, string, error) {
	// grab the metadata from the context to pull in provided
	// cookie duration information
	m := c.MustGet("metadata").(*types.Metadata)

	// create a refresh with the provided duration
	refreshToken, refreshExpiry, err := CreateRefreshToken(u, m.Vela.RefreshTokenDuration)
	if err != nil {
		return "", "", err
	}

	// create an access token with the provided duration
	accessToken, err := CreateAccessToken(u, m.Vela.AccessTokenDuration)
	if err != nil {
		return "", "", err
	}

	// parse the address for the backend server
	// so we can set it for the cookie domain
	addr, err := url.Parse(m.Vela.Address)
	if err != nil {
		return "", "", err
	}

	// set the SameSite value for the cookie
	// https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html#samesite-attribute
	// We set to Lax because we will have links from source provider web UI.
	// Setting this to Strict would force a login when navigating via source provider web UI links.
	c.SetSameSite(http.SameSiteLaxMode)
	// set the cookie with the refresh token as a HttpOnly, Secure cookie
	// https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html#httponly-attribute
	// https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html#secure-attribute
	c.SetCookie(constants.RefreshTokenName, refreshToken, refreshExpiry, "/", addr.Hostname(), c.Value("securecookie").(bool), true)

	// return the refresh and access tokens
	return refreshToken, accessToken, nil
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

	// create a new JWT parser
	p := &jwt.Parser{
		// explicitly only allow these signing methods
		ValidMethods: []string{jwt.SigningMethodHS256.Name},
	}

	// parse the signed JWT token string
	// parse also validates the claims and token by default.
	_, err := p.ParseWithClaims(t, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		var err error

		// extract the claims from the token
		claims := token.Claims.(*Claims)
		name := claims.Subject

		// lookup the user in the database
		logrus.Debugf("Reading user %s", name)
		u, err = db.GetUserName(name)
		return []byte(u.GetHash()), err
	})

	// there will be an error if we're not able to parse
	// the token, eg. due to expiration, invalid signature, etc
	if err != nil {
		return nil, fmt.Errorf("invalid token provided for %s: %w", u.GetName(), err)
	}

	return u, nil
}

// RetrieveAccessToken gets the passed in access token from the header in the request
func RetrieveAccessToken(r *http.Request) (accessToken string, err error) {
	accessToken, err = request.AuthorizationHeaderExtractor.ExtractToken(r)

	return
}

// RetrieveRefreshToken gets the refresh token sent along with the request as a cookie
func RetrieveRefreshToken(r *http.Request) (string, error) {
	refreshToken, err := r.Cookie(constants.RefreshTokenName)

	if refreshToken == nil || len(refreshToken.Value) == 0 {
		// cookie will not be sent if it has expired
		return "", fmt.Errorf("refresh token expired or not provided")
	}

	return refreshToken.Value, err
}

// CreateAccessToken creates a new access token for the given user and duration
func CreateAccessToken(u *library.User, d time.Duration) (string, error) {
	now := time.Now()
	exp := now.Add(d)

	claims := &Claims{
		IsActive: u.GetActive(),
		IsAdmin:  u.GetAdmin(),
		StandardClaims: jwt.StandardClaims{
			Subject:   u.GetName(),
			IssuedAt:  now.Unix(),
			ExpiresAt: exp.Unix(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := t.SignedString([]byte(u.GetHash()))
	if err != nil {
		return "", err
	}

	return token, nil
}

// CreateCreateRefreshToken creates a new refresh token for the given user and duration
func CreateRefreshToken(u *library.User, d time.Duration) (string, int, error) {
	exp := time.Now().Add(d)

	claims := jwt.StandardClaims{}
	claims.Subject = u.GetName()
	claims.ExpiresAt = exp.Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	refreshToken, err := t.SignedString([]byte(u.GetHash()))
	if err != nil {
		return "", 0, err
	}

	return refreshToken, int(d.Seconds()), nil
}

// Refresh returns a new access token, if the provided refreshToken is valid
func Refresh(c *gin.Context, refreshToken string) (string, error) {
	// get the metadata
	m := c.MustGet("metadata").(*types.Metadata)
	// get a reference to the database
	db := database.FromContext(c)

	// check to see if a user exists with that refresh token
	// we are comparing with db to allow for leverage in
	// invalidating a refresh token in the DB
	u, err := db.GetUserRefreshToken(refreshToken)
	if err != nil {
		return "", fmt.Errorf("refresh token not valid - please log in")
	}

	// parse (which also validates) the token
	_, err = Parse(refreshToken, db)
	if err != nil {
		return "", err
	}

	// create a new access token
	at, err := CreateAccessToken(u, m.Vela.AccessTokenDuration)
	if err != nil {
		return "", err
	}

	return at, nil
}
