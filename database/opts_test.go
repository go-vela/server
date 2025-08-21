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

func TestDatabase_EngineOpt_WithLogLevel(t *testing.T) {
	e := &engine{config: new(config)}

	tests := []struct {
		failure  bool
		name     string
		logLevel string
		want     string
	}{
		{
			failure:  false,
			name:     "log level set to debug",
			logLevel: "debug",
			want:     "debug",
		},
		{
			failure:  false,
			name:     "log level set to info",
			logLevel: "info",
			want:     "info",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithLogLevel(test.logLevel)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithLogLevel for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("WithLogLevel returned err: %v", err)
			}

			if !reflect.DeepEqual(e.config.LogLevel, test.want) {
				t.Errorf("WithLogLevel is %v, want %v", e.config.SkipCreation, test.want)
			}
		})
	}
}

func TestDatabase_EngineOpt_WithLogSkipNotFound(t *testing.T) {
	e := &engine{config: new(config)}

	tests := []struct {
		failure bool
		name    string
		skip    bool
		want    bool
	}{
		{
			failure: false,
			name:    "log skip not found set to true",
			skip:    true,
			want:    true,
		},
		{
			failure: false,
			name:    "log skip not found set to false",
			skip:    false,
			want:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithLogSkipNotFound(test.skip)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithLogSkipNotFound for %s should have returned err", test.name)
				}

				if err != nil {
					t.Errorf("WithLogSkipNotFound for %s returned err: %v", test.name, err)
				}

				if !reflect.DeepEqual(e.config.LogSkipNotFound, test.want) {
					t.Errorf("WithLogSkipNotFound for %s is %v, want %v", test.name, e.config.LogSkipNotFound, test.want)
				}
			}
		})
	}
}

func TestDatabase_EngineOpt_WithLogSlowThreshold(t *testing.T) {
	e := &engine{config: new(config)}

	tests := []struct {
		failure   bool
		name      string
		threshold time.Duration
		want      time.Duration
	}{
		{
			failure:   false,
			name:      "log slow threshold set to 1ms",
			threshold: 1 * time.Millisecond,
			want:      1 * time.Millisecond,
		},
		{
			failure:   false,
			name:      "log slow threshold set to 2ms",
			threshold: 2 * time.Millisecond,
			want:      2 * time.Millisecond,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithLogSlowThreshold(test.threshold)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithLogSlowThreshold for %s should have returned err", test.name)
				}

				if err != nil {
					t.Errorf("WithLogSlowThreshold for %s returned err: %v", test.name, err)
				}

				if !reflect.DeepEqual(e.config.LogSlowThreshold, test.want) {
					t.Errorf("WithLogSlowThreshold for %s is %v, want %v", test.name, e.config.LogSlowThreshold, test.want)
				}
			}
		})
	}
}

func TestDatabase_EngineOpt_WithLogShowSQL(t *testing.T) {
	e := &engine{config: new(config)}

	tests := []struct {
		failure bool
		name    string
		show    bool
		want    bool
	}{
		{
			failure: false,
			name:    "log show SQL set to true",
			show:    true,
			want:    true,
		},
		{
			failure: false,
			name:    "log show SQL set to false",
			show:    false,
			want:    false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithLogShowSQL(test.show)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithLogShowSQL for %s should have returned err", test.name)
				}

				if err != nil {
					t.Errorf("WithLogShowSQL for %s returned err: %v", test.name, err)
				}

				if !reflect.DeepEqual(e.config.LogShowSQL, test.want) {
					t.Errorf("WithLogShowSQL for %s is %v, want %v", test.name, e.config.LogShowSQL, test.want)
				}
			}
		})
	}
}

func TestDatabase_EngineOpt_WithLogPartitioned(t *testing.T) {
	// setup types
	e := &engine{config: new(config)}

	// setup tests
	tests := []struct {
		failure     bool
		name        string
		partitioned bool
		want        bool
	}{
		{
			failure:     false,
			name:        "log partitioned set to true",
			partitioned: true,
			want:        true,
		},
		{
			failure:     false,
			name:        "log partitioned set to false",
			partitioned: false,
			want:        false,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithLogPartitioned(test.partitioned)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithLogPartitioned for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("WithLogPartitioned for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(e.config.LogPartitioned, test.want) {
				t.Errorf("WithLogPartitioned for %s is %v, want %v", test.name, e.config.LogPartitioned, test.want)
			}
		})
	}
}

func TestDatabase_EngineOpt_WithLogPartitionPattern(t *testing.T) {
	// setup types
	e := &engine{config: new(config)}

	// setup tests
	tests := []struct {
		failure bool
		name    string
		pattern string
		want    string
	}{
		{
			failure: false,
			name:    "pattern set to logs_%",
			pattern: "logs_%",
			want:    "logs_%",
		},
		{
			failure: false,
			name:    "pattern set to custom_logs_%",
			pattern: "custom_logs_%",
			want:    "custom_logs_%",
		},
		{
			failure: false,
			name:    "pattern set to empty string",
			pattern: "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithLogPartitionPattern(test.pattern)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithLogPartitionPattern for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("WithLogPartitionPattern for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(e.config.LogPartitionPattern, test.want) {
				t.Errorf("WithLogPartitionPattern for %s is %v, want %v", test.name, e.config.LogPartitionPattern, test.want)
			}
		})
	}
}

func TestDatabase_EngineOpt_WithLogPartitionSchema(t *testing.T) {
	// setup types
	e := &engine{config: new(config)}

	// setup tests
	tests := []struct {
		failure bool
		name    string
		schema  string
		want    string
	}{
		{
			failure: false,
			name:    "schema set to public",
			schema:  "public",
			want:    "public",
		},
		{
			failure: false,
			name:    "schema set to custom_schema",
			schema:  "custom_schema",
			want:    "custom_schema",
		},
		{
			failure: false,
			name:    "schema set to empty string",
			schema:  "",
			want:    "",
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithLogPartitionSchema(test.schema)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithLogPartitionSchema for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("WithLogPartitionSchema for %s returned err: %v", test.name, err)
			}

			if !reflect.DeepEqual(e.config.LogPartitionSchema, test.want) {
				t.Errorf("WithLogPartitionSchema for %s is %v, want %v", test.name, e.config.LogPartitionSchema, test.want)
			}
		})
	}
}
