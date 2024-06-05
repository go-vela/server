// SPDX-License-Identifier: Apache-2.0

package token

import (
	"context"
	"crypto/rand"
	"crypto/rsa"

	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/v2/jwk"

	"github.com/go-vela/server/database"
)

// GenerateRSA creates an RSA key pair and sets it in the token manager and saves the JWK in the database.
func (tm *Manager) GenerateRSA(ctx context.Context, db database.Interface) error {
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

	j := jwk.NewRSAPublicKey()

	err = j.FromRaw(&privateRSAKey.PublicKey)
	if err != nil {
		return err
	}

	err = j.Set(jwk.KeyIDKey, kid.String())
	if err != nil {
		return err
	}

	// create the JWK in the database
	err = db.CreateJWK(context.TODO(), j)
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
