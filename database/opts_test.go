// SPDX-License-Identifier: Apache-2.0

package database

import (
	"reflect"
	"testing"
	"time"
)

func TestDatabase_EngineOpt_WithAddress(t *testing.T) {
	// setup types
	e := &engine{config: new(config)}

	// setup tests
	tests := []struct {
		failure bool
		name    string
		address string
		want    string
	}{
		{
			failure: false,
			name:    "address set",
			address: "file::memory:?cache=shared",
			want:    "file::memory:?cache=shared",
		},
		{
			failure: false,
			name:    "address not set",
			address: "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithAddress(test.address)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithAddress for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("WithAddress for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(e.config.Address, test.want) {
				t.Errorf("WithAddress for %s is %v, want %v", test.name, e.config.Address, test.want)
			}
		})
	}
}

func TestDatabase_EngineOpt_WithCompressionLevel(t *testing.T) {
	// setup types
	e := &engine{config: new(config)}

	// setup tests
	tests := []struct {
		failure bool
		name    string
		level   int
		want    int
	}{
		{
			failure: false,
			name:    "compression level set to -1",
			level:   -1,
			want:    -1,
		},
		{
			failure: false,
			name:    "compression level set to 0",
			level:   0,
			want:    0,
		},
		{
			failure: false,
			name:    "compression level set to 1",
			level:   1,
			want:    1,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithCompressionLevel(test.level)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithCompressionLevel for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("WithCompressionLevel for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(e.config.CompressionLevel, test.want) {
				t.Errorf("WithCompressionLevel for %s is %v, want %v", test.name, e.config.CompressionLevel, test.want)
			}
		})
	}
}

func TestDatabase_EngineOpt_WithConnectionLife(t *testing.T) {
	// setup types
	e := &engine{config: new(config)}

	// setup tests
	tests := []struct {
		failure bool
		name    string
		life    time.Duration
		want    time.Duration
	}{
		{
			failure: false,
			name:    "life of connections set",
			life:    30 * time.Minute,
			want:    30 * time.Minute,
		},
		{
			failure: false,
			name:    "life of connections not set",
			life:    0,
			want:    0,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithConnectionLife(test.life)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithConnectionLife for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("WithConnectionLife for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(e.config.ConnectionLife, test.want) {
				t.Errorf("WithConnectionLife for %s is %v, want %v", test.name, e.config.ConnectionLife, test.want)
			}
		})
	}
}

func TestDatabase_EngineOpt_WithConnectionIdle(t *testing.T) {
	// setup types
	e := &engine{config: new(config)}

	// setup tests
	tests := []struct {
		failure bool
		name    string
		open    int
		want    int
	}{
		{
			failure: false,
			name:    "idle connections set",
			open:    2,
			want:    2,
		},
		{
			failure: false,
			name:    "idle connections not set",
			open:    0,
			want:    0,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithConnectionIdle(test.open)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithConnectionIdle for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("WithConnectionIdle for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(e.config.ConnectionIdle, test.want) {
				t.Errorf("WithConnectionIdle for %s is %v, want %v", test.name, e.config.ConnectionIdle, test.want)
			}
		})
	}
}

func TestDatabase_EngineOpt_WithConnectionOpen(t *testing.T) {
	// setup types
	e := &engine{config: new(config)}

	// setup tests
	tests := []struct {
		failure bool
		name    string
		open    int
		want    int
	}{
		{
			failure: false,
			name:    "open connections set",
			open:    2,
			want:    2,
		},
		{
			failure: false,
			name:    "open connections not set",
			open:    0,
			want:    0,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithConnectionOpen(test.open)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithConnectionOpen for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("WithConnectionOpen for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(e.config.ConnectionOpen, test.want) {
				t.Errorf("WithConnectionOpen for %s is %v, want %v", test.name, e.config.ConnectionOpen, test.want)
			}
		})
	}
}

func TestDatabase_EngineOpt_WithDriver(t *testing.T) {
	// setup types
	e := &engine{config: new(config)}

	// setup tests
	tests := []struct {
		failure bool
		name    string
		driver  string
		want    string
	}{
		{
			failure: false,
			name:    "driver set",
			driver:  "sqlite3",
			want:    "sqlite3",
		},
		{
			failure: false,
			name:    "driver not set",
			driver:  "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithDriver(test.driver)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithDriver for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("WithDriver for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(e.config.Driver, test.want) {
				t.Errorf("WithDriver for %s is %v, want %v", test.name, e.config.Driver, test.want)
			}
		})
	}
}

func TestDatabase_EngineOpt_WithEncryptionKey(t *testing.T) {
	// setup types
	e := &engine{config: new(config)}

	// setup tests
	tests := []struct {
		failure bool
		name    string
		key     string
		want    string
	}{
		{
			failure: false,
			name:    "encryption key set",
			key:     "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
			want:    "A1B2C3D4E5G6H7I8J9K0LMNOPQRSTUVW",
		},
		{
			failure: false,
			name:    "encryption key not set",
			key:     "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithEncryptionKey(test.key)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithEncryptionKey for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("WithEncryptionKey for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(e.config.EncryptionKey, test.want) {
				t.Errorf("WithEncryptionKey for %s is %v, want %v", test.name, e.config.EncryptionKey, test.want)
			}
		})
	}
}

func TestDatabase_EngineOpt_WithSkipCreation(t *testing.T) {
	// setup types
	e := &engine{config: new(config)}

	// setup tests
	tests := []struct {
		failure bool
		name    string
		skip    bool
		want    bool
	}{
		{
			failure: false,
			name:    "skip creation set to true",
			skip:    true,
			want:    true,
		},
		{
			failure: false,
			name:    "skip creation set to false",
			skip:    false,
			want:    false,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithSkipCreation(test.skip)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithSkipCreation for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("WithSkipCreation returned err: %v", err)
			}

			if !reflect.DeepEqual(e.config.SkipCreation, test.want) {
				t.Errorf("WithSkipCreation is %v, want %v", e.config.SkipCreation, test.want)
			}
		})
	}
}
