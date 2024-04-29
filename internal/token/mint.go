// SPDX-License-Identifier: Apache-2.0

package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
)

// Claims struct is an extension of the JWT standard claims. It
// includes information about the user.
type Claims struct {
	BuildID     int64  `json:"build_id,omitempty"`
	BuildNumber int    `json:"build_number,omitempty"`
	IsActive    bool   `json:"is_active,omitempty"`
	IsAdmin     bool   `json:"is_admin,omitempty"`
	Repo        string `json:"repo,omitempty"`
	TokenType   string `json:"token_type,omitempty"`
	jwt.RegisteredClaims
}

// MintTokenOpts is a type to inform the token minter how to construct
// the token.
type MintTokenOpts struct {
	BuildID       int64
	BuildNumber   int
	Hostname      string
	Repo          string
	TokenDuration time.Duration
	TokenType     string
	User          *api.User
	Audience      []string
	Commit        string
}

// MintToken mints a Vela JWT Token given a set of options.
func (tm *Manager) MintToken(mto *MintTokenOpts) (string, error) {
	// initialize claims struct
	var claims = new(Claims)

	// apply claims based on token type
	switch mto.TokenType {
	case constants.UserAccessTokenType, constants.UserRefreshTokenType:
		if mto.User == nil {
			return "", fmt.Errorf("no user provided for user access token")
		}

		claims.IsActive = mto.User.GetActive()
		claims.IsAdmin = mto.User.GetAdmin()
		claims.Subject = mto.User.GetName()

	case constants.WorkerBuildTokenType:
		if mto.BuildID == 0 {
			return "", errors.New("missing build id for build token")
		}

		if len(mto.Repo) == 0 {
			return "", errors.New("missing repo for build token")
		}

		if len(mto.Hostname) == 0 {
			return "", errors.New("missing host name for build token")
		}

		claims.BuildID = mto.BuildID
		claims.Repo = mto.Repo
		claims.Subject = mto.Hostname

	case constants.WorkerAuthTokenType, constants.WorkerRegisterTokenType:
		if len(mto.Hostname) == 0 {
			return "", fmt.Errorf("missing host name for %s token", mto.TokenType)
		}

		claims.Subject = mto.Hostname

	case constants.IDRequestTokenType:
		if len(mto.Repo) == 0 {
			return "", errors.New("missing repo for ID request token")
		}

		if len(mto.Commit) == 0 {
			return "", errors.New("missing commit for ID request token")
		}

		if mto.BuildID == 0 {
			return "", errors.New("missing build id for ID request token")
		}

		claims.Repo = mto.Repo
		claims.Subject = fmt.Sprintf("%s/%s", mto.Repo, mto.Commit)
		claims.BuildID = mto.BuildID
		claims.BuildNumber = mto.BuildNumber

	default:
		return "", errors.New("invalid token type")
	}

	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(mto.TokenDuration))
	claims.TokenType = mto.TokenType

	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//sign token with configured private signing key
	token, err := tk.SignedString([]byte(tm.PrivateKeyHMAC))
	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}

	return token, nil
}

func (tm *Manager) MintIDToken(mto *MintTokenOpts, db database.Interface) (string, error) {
	// initialize claims struct
	var claims = new(Claims)

	if len(mto.Repo) == 0 {
		return "", errors.New("missing repo for ID token")
	}

	if len(mto.Commit) == 0 {
		return "", errors.New("missing commit for ID token")
	}

	if mto.BuildNumber == 0 {
		return "", errors.New("missing build id for ID token")
	}

	claims.BuildNumber = mto.BuildNumber
	claims.Repo = mto.Repo
	claims.Subject = fmt.Sprintf("%s/%s", mto.Repo, mto.Commit)
	claims.Audience = mto.Audience

	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(mto.TokenDuration))
	claims.TokenType = mto.TokenType
	claims.Issuer = tm.Issuer

	tk := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	_, err := db.GetActiveKeySet(context.TODO(), tm.RSAKeySet.KID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("unable to get active key set: %w", err)
		}

		err = tm.GenerateRSA(db)
		if err != nil {
			return "", fmt.Errorf("unable to generate RSA key set: %w", err)
		}
	}

	tk.Header["kid"] = tm.RSAKeySet.KID

	//sign token with configured private signing key
	token, err := tk.SignedString(tm.RSAKeySet.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}

	return token, nil
}
