// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package auth

import (
	"fmt"
	"net/http"

	"github.com/go-vela/types/constants"

	"github.com/golang-jwt/jwt/v5/request"
)

// RetrieveAccessToken gets the passed in access token from the header in the request.
func RetrieveAccessToken(r *http.Request) (accessToken string, err error) {
	return request.AuthorizationHeaderExtractor.ExtractToken(r)
}

// RetrieveRefreshToken gets the refresh token sent along with the request as a cookie.
func RetrieveRefreshToken(r *http.Request) (string, error) {
	refreshToken, err := r.Cookie(constants.RefreshTokenName)

	if refreshToken == nil || len(refreshToken.Value) == 0 {
		// cookie will not be sent if it has expired
		return "", fmt.Errorf("refresh token expired or not provided")
	}

	return refreshToken.Value, err
}
