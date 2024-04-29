// SPDX-License-Identifier: Apache-2.0

package token

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"strconv"

	api "github.com/go-vela/server/api/types"
	"github.com/go-vela/server/database"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// MintToken mints a Vela JWT Token given a set of options.
func (tm *Manager) GenerateRSA(db database.Interface) error {
	privateRSAKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	kid, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	key := api.JWK{
		Algorithm: jwt.SigningMethodRS256.Name,
		Kid:       kid.String(),
		N:         base64.RawURLEncoding.EncodeToString(privateRSAKey.PublicKey.N.Bytes()),
		E:         base64.RawURLEncoding.EncodeToString([]byte(strconv.Itoa(privateRSAKey.PublicKey.E))),
	}

	err = db.CreateKeySet(context.TODO(), key)
	if err != nil {
		return err
	}

	keySet := RSAKeySet{
		PrivateKey: privateRSAKey,
		KID:        kid.String(),
	}

	tm.RSAKeySet = keySet

	return nil
}
