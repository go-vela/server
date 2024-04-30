// SPDX-License-Identifier: Apache-2.0

package token

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
)

// GenerateRSA creates an RSA key pair and sets it in the token manager and saves the JWK in the database.
func (tm *Manager) GenerateRSA(db database.Interface) error {
	// generate key pair
	privateRSAKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// assign KID to key pair
	kid, err := uuid.NewV7()
	if err != nil {
		return err
	}

	// abstract the JWK from the public key information
	key := api.JWK{
		Algorithm: jwt.SigningMethodRS256.Name,
		Kid:       kid.String(),
		Use:       "sig",
		Kty:       "RSA",
		N:         base64.RawURLEncoding.EncodeToString(privateRSAKey.PublicKey.N.Bytes()),
		E:         base64.RawURLEncoding.EncodeToString([]byte(strconv.Itoa(privateRSAKey.PublicKey.E))),
	}

	// create the JWK in the database
	err = db.CreateJWK(context.TODO(), key)
	if err != nil {
		return err
	}

	// create the RSA key set for token manager
	keySet := RSAKeySet{
		PrivateKey: privateRSAKey,
		KID:        kid.String(),
	}

	tm.RSAKeySet = keySet

	return nil
}
