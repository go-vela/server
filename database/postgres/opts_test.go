// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package postgres

import (
	"reflect"
	"testing"
	"time"
)

func TestPostgres_ClientOpt_WithAddress(t *testing.T) {
	// setup types
	c := new(client)
	c.config = new(config)

	// setup tests
	tests := []struct {
		failure bool
		address string
		want    string
	}{
		{
			failure: false,
			address: "postgres://foo:bar@localhost:5432/vela",
			want:    "postgres://foo:bar@localhost:5432/vela",
		},
		{
			failure: true,
			address: "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		err := WithAddress(test.address)(c)

		if test.failure {
			if err == nil {
				t.Errorf("WithAddress should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithAddress returned err: %v", err)
		}

		if !reflect.DeepEqual(c.config.Address, test.want) {
			t.Errorf("WithAddress is %v, want %v", c.config.Address, test.want)
		}
	}
}

func TestPostgres_ClientOpt_WithCompressionLevel(t *testing.T) {
	// setup types
	c := new(client)
	c.config = new(config)

	// setup tests
	tests := []struct {
		level int
		want  int
	}{
		{
			level: 3,
			want:  3,
		},
		{
			level: 0,
			want:  0,
		},
	}

	// run tests
	for _, test := range tests {
		err := WithCompressionLevel(test.level)(c)

		if err != nil {
			t.Errorf("WithCompressionLevel returned err: %v", err)
		}

		if !reflect.DeepEqual(c.config.CompressionLevel, test.want) {
			t.Errorf("WithCompressionLevel is %v, want %v", c.config.CompressionLevel, test.want)
		}
	}
}

func TestPostgres_ClientOpt_WithConnectionLife(t *testing.T) {
	// setup types
	c := new(client)
	c.config = new(config)

	// setup tests
	tests := []struct {
		duration time.Duration
		want     time.Duration
	}{
		{
			duration: 10 * time.Second,
			want:     10 * time.Second,
		},
		{
			duration: 0,
			want:     0,
		},
	}

	// run tests
	for _, test := range tests {
		err := WithConnectionLife(test.duration)(c)

		if err != nil {
			t.Errorf("WithConnectionLife returned err: %v", err)
		}

		if !reflect.DeepEqual(c.config.ConnectionLife, test.want) {
			t.Errorf("WithConnectionLife is %v, want %v", c.config.ConnectionLife, test.want)
		}
	}
}

func TestPostgres_ClientOpt_WithConnectionIdle(t *testing.T) {
	// setup types
	c := new(client)
	c.config = new(config)

	// setup tests
	tests := []struct {
		idle int
		want int
	}{
		{
			idle: 5,
			want: 5,
		},
		{
			idle: 0,
			want: 0,
		},
	}

	// run tests
	for _, test := range tests {
		err := WithConnectionIdle(test.idle)(c)

		if err != nil {
			t.Errorf("WithConnectionIdle returned err: %v", err)
		}

		if !reflect.DeepEqual(c.config.ConnectionIdle, test.want) {
			t.Errorf("WithConnectionIdle is %v, want %v", c.config.ConnectionIdle, test.want)
		}
	}
}

func TestPostgres_ClientOpt_WithConnectionOpen(t *testing.T) {
	// setup types
	c := new(client)
	c.config = new(config)

	// setup tests
	tests := []struct {
		open int
		want int
	}{
		{
			open: 10,
			want: 10,
		},
		{
			open: 0,
			want: 0,
		},
	}

	// run tests
	for _, test := range tests {
		err := WithConnectionOpen(test.open)(c)

		if err != nil {
			t.Errorf("WithConnectionOpen returned err: %v", err)
		}

		if !reflect.DeepEqual(c.config.ConnectionOpen, test.want) {
			t.Errorf("WithConnectionOpen is %v, want %v", c.config.ConnectionOpen, test.want)
		}
	}
}

func TestPostgres_ClientOpt_WithEncryptionKey(t *testing.T) {
	// setup types
	c := new(client)
	c.config = new(config)

	// setup tests
	tests := []struct {
		failure bool
		key     string
		want    string
	}{
		{
			failure: false,
			key:     "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
			want:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
		},
		{
			failure: true,
			key:     "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		err := WithEncryptionKey(test.key)(c)

		if test.failure {
			if err == nil {
				t.Errorf("WithEncryptionKey should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("WithEncryptionKey returned err: %v", err)
		}

		if !reflect.DeepEqual(c.config.EncryptionKey, test.want) {
			t.Errorf("WithEncryptionKey is %v, want %v", c.config.EncryptionKey, test.want)
		}
	}
}

func TestPostgres_ClientOpt_WithSkipCreation(t *testing.T) {
	// setup types
	c := new(client)
	c.config = new(config)

	// setup tests
	tests := []struct {
		skipCreation bool
		want         bool
	}{
		{
			skipCreation: true,
			want:         true,
		},
		{
			skipCreation: false,
			want:         false,
		},
	}

	// run tests
	for _, test := range tests {
		err := WithSkipCreation(test.skipCreation)(c)

		if err != nil {
			t.Errorf("WithSkipCreation returned err: %v", err)
		}

		if !reflect.DeepEqual(c.config.SkipCreation, test.want) {
			t.Errorf("WithSkipCreation is %v, want %v", c.config.SkipCreation, test.want)
		}
	}
}
