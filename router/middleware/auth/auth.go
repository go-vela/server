// SPDX-License-Identifier: Apache-2.0

package auth

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5/request"

	"github.com/go-vela/server/constants"
)

// RetrieveAccessToken gets the passed in access token from the header in the request.
func RetrieveAccessToken(r *http.Request) (accessToken string, err error) {
	return request.AuthorizationHeaderExtractor.ExtractToken(r)
}

// RetrieveTokenHeader gets the passed in install token from the 'Token' header in the request.
//
// this is only used for builds that have app installation tokens that are used for status updates.
// it is not required unless the repository has installed the Vela app.
func RetrieveTokenHeader(r *http.Request) string {
	tkn, ok := r.Header["Token"]
	if !ok || len(tkn) == 0 {
		return ""
	}

	return tkn[0]
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
