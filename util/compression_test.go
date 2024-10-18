// SPDX-License-Identifier: Apache-2.0

package util

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/constants"
)

func TestDatabase_compress(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		failure bool
		level   int
		data    []byte
		want    []byte
	}{
		{
			name:    "compression level -1",
			failure: false,
			level:   constants.CompressionNegOne,
			data:    []byte("foo"),
			want:    []byte{120, 156, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 0",
			failure: false,
			level:   constants.CompressionZero,
			data:    []byte("foo"),
			want:    []byte{120, 1, 0, 3, 0, 252, 255, 102, 111, 111, 1, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 1",
			failure: false,
			level:   constants.CompressionOne,
			data:    []byte("foo"),
			want:    []byte{120, 1, 0, 3, 0, 252, 255, 102, 111, 111, 1, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 2",
			failure: false,
			level:   constants.CompressionTwo,
			data:    []byte("foo"),
			want:    []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 3",
			failure: false,
			level:   constants.CompressionThree,
			data:    []byte("foo"),
			want:    []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 4",
			failure: false,
			level:   constants.CompressionFour,
			data:    []byte("foo"),
			want:    []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 5",
			failure: false,
			level:   constants.CompressionFive,
			data:    []byte("foo"),
			want:    []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 6",
			failure: false,
			level:   constants.CompressionSix,
			data:    []byte("foo"),
			want:    []byte{120, 156, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 7",
			failure: false,
			level:   constants.CompressionSeven,
			data:    []byte("foo"),
			want:    []byte{120, 218, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 8",
			failure: false,
			level:   constants.CompressionEight,
			data:    []byte("foo"),
			want:    []byte{120, 218, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
		{
			name:    "compression level 9",
			failure: false,
			level:   constants.CompressionNine,
			data:    []byte("foo"),
			want:    []byte{120, 218, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Compress(test.level, test.data)

			if test.failure {
				if err == nil {
					t.Errorf("compress for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("compress for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("compress for %s is %v, want %v", test.name, string(got), string(test.want))
			}
		})
	}
}

func TestDatabase_decompress(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		failure bool
		data    []byte
		want    []byte
	}{
		{
			name:    "compression level -1",
			failure: false,
			data:    []byte{120, 156, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 0",
			failure: false,
			data:    []byte{120, 1, 0, 3, 0, 252, 255, 102, 111, 111, 1, 0, 0, 255, 255, 2, 130, 1, 69},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 1",
			failure: false,
			data:    []byte{120, 1, 0, 3, 0, 252, 255, 102, 111, 111, 1, 0, 0, 255, 255, 2, 130, 1, 69},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 2",
			failure: false,
			data:    []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 3",
			failure: false,
			data:    []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 4",
			failure: false,
			data:    []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 5",
			failure: false,
			data:    []byte{120, 94, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 6",
			failure: false,
			data:    []byte{120, 156, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 7",
			failure: false,
			data:    []byte{120, 218, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 8",
			failure: false,
			data:    []byte{120, 218, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
			want:    []byte("foo"),
		},
		{
			name:    "compression level 9",
			failure: false,
			data:    []byte{120, 218, 74, 203, 207, 7, 4, 0, 0, 255, 255, 2, 130, 1, 69},
			want:    []byte("foo"),
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Decompress(test.data)

			if test.failure {
				if err == nil {
					t.Errorf("decompress for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("decompressm for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("decompress for %s is %v, want %v", test.name, string(got), string(test.want))
			}
		})
	}
}
