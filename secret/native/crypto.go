// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

// encrypt represents a function that is used for encrypting the value
// of a secret before it gets written into the database. The Go standard
// library package is used for encryption with AES ciphers:
//
// https://golang.org/pkg/crypto/aes/
//
// This function is based on methods that exist in the Vault project:
// https://github.com/hashicorp/vault/blob/v1.6.0/helper/dhutil/dhutil.go#L99
func encrypt(data []byte, key string) (string, error) {
	// within the validate process we enforce a 64 bit key which
	// ensures all secrets are encrypted with AES-256:
	// https://en.wikipedia.org/wiki/Advanced_Encryption_Standard
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// nonce is an arbitrary number used to to ensure that
	// old communications cannot be reused in replay attacks.
	// https://en.wikipedia.org/wiki/Cryptographic_nonce
	nonce := make([]byte, gcm.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	// encrypt the data with the randomly generated nonce
	encData := gcm.Seal(nonce, nonce, data, nil)

	// encode the encrypt data to make it network safe
	sEnc := base64.StdEncoding.EncodeToString(encData)

	return sEnc, nil
}

// decrypt represents a function that is used for decrypting the value
// of a secret when a user is retrieving it from the database. The Go standard
// library package is used for encryption with AES ciphers:
//
// https://golang.org/pkg/crypto/aes/
//
// This function is based on methods that exist in the Vault project:
// https://github.com/hashicorp/vault/blob/v1.6.0/helper/dhutil/dhutil.go#L132
func decrypt(data []byte, key string) (string, error) {
	// decode the encrypted data
	dDec, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return "", err
	}

	// within the validate process we enforce a 64 bit key which
	// ensures all secrets are decrypted with AES-256:
	// https://en.wikipedia.org/wiki/Advanced_Encryption_Standard
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// nonce is an arbitrary number used to to ensure that
	// old communications cannot be reused in replay attacks.
	// https://en.wikipedia.org/wiki/Cryptographic_nonce
	nonceSize := gcm.NonceSize()

	nonce, ciphertext := dDec[:nonceSize], dDec[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
