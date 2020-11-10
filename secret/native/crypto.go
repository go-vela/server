// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// encrypt represents a function that is used for encrypting the value
// of a secret before it gets written into the database. The Go standard
// library package is used for encryption with AES ciphers:
//
// https://golang.org/pkg/crypto/aes/
func encrypt(data []byte, passphrase string) (string, error) {
	// create new md5 type hasher
	hasher := md5.New()

	_, err := hasher.Write([]byte(passphrase))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(hex.EncodeToString(hasher.Sum(nil))))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())

	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	return string(gcm.Seal(nonce, nonce, data, nil)), nil
}

// decrypt represents a function that is used for decrypting the value
// of a secret when a user is retrieving it from the database. The Go standard
// library package is used for encryption with AES ciphers:
//
// https://golang.org/pkg/crypto/aes/
func decrypt(data []byte, passphrase string) (string, error) {
	// create new md5 type hasher
	hasher := md5.New()

	_, err := hasher.Write([]byte(passphrase))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(hex.EncodeToString(hasher.Sum(nil))))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
