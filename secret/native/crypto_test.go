// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"testing"
)

func TestNative_encrypt(t *testing.T) {
	// setup types
	tests := []struct {
		name       string
		data       []byte
		passphrase string
		wantErr    bool
	}{
		{
			name:       "success when encrypting a value with a passphrase",
			data:       []byte("hello, world"),
			passphrase: "C639A572E14D5075C526FDDD43E4ECF6",
			wantErr:    false,
		},
		{
			name:       "success when encrypting a value with a special characters",
			data:       []byte("!@#$%^&*()"),
			passphrase: "C639A572E14D5075C526FDDD43E4ECF6",
			wantErr:    false,
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := encrypt(test.data, test.passphrase)
			if (err != nil) != test.wantErr {
				t.Errorf("encrypt() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if len(got) == 0 {
				t.Errorf("encrypt() = %d, want %d", len(got), 0)
			}
		})
	}
}

func TestNative_decrypt(t *testing.T) {
	// setup types
	tests := []struct {
		name       string
		data       []byte
		passphrase string
		want       string
		wantErr    bool
	}{
		{
			name:       "success when decrypting a value with a passphrase",
			data:       []byte("hello, world"),
			passphrase: "C639A572E14D5075C526FDDD43E4ECF6",
			wantErr:    false,
		},
		{
			name:       "success when decrypting a value with a special characters",
			data:       []byte("!@#$%^&*()"),
			passphrase: "C639A572E14D5075C526FDDD43E4ECF6",
			wantErr:    false,
		},
	}

	// run test
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			val, _ := encrypt(test.data, test.passphrase)

			got, err := decrypt([]byte(val), test.passphrase)
			if (err != nil) != test.wantErr {
				t.Errorf("decrypt() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if len(got) == 0 {
				t.Errorf("decrypt() = %d, want %d", len(got), 0)
			}
		})
	}
}
