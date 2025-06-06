// SPDX-License-Identifier: Apache-2.0

package executable

import (
	"context"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func TestExecutable_EngineOpt_WithClient(t *testing.T) {
	// setup types
	e := &Engine{client: new(gorm.DB)}

	// setup tests
	tests := []struct {
		failure bool
		name    string
		client  *gorm.DB
		want    *gorm.DB
	}{
		{
			failure: false,
			name:    "client set to new database",
			client:  new(gorm.DB),
			want:    new(gorm.DB),
		},
		{
			failure: false,
			name:    "client set to nil",
			client:  nil,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithClient(test.client)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithClient for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("WithClient returned err: %v", err)
			}

			if !reflect.DeepEqual(e.client, test.want) {
				t.Errorf("WithClient is %v, want %v", e.client, test.want)
			}
		})
	}
}

func TestExecutable_EngineOpt_WithCompressionLevel(t *testing.T) {
	// setup types
	e := &Engine{config: new(config)}

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
				t.Errorf("WithCompressionLevel returned err: %v", err)
			}

			if !reflect.DeepEqual(e.config.CompressionLevel, test.want) {
				t.Errorf("WithCompressionLevel is %v, want %v", e.config.CompressionLevel, test.want)
			}
		})
	}
}

func TestExecutable_EngineOpt_WithEncryptionKey(t *testing.T) {
	// setup types
	e := &Engine{config: new(config)}

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
				t.Errorf("WithEncryptionKey returned err: %v", err)
			}

			if !reflect.DeepEqual(e.config.EncryptionKey, test.want) {
				t.Errorf("WithEncryptionKey is %v, want %v", e.config.EncryptionKey, test.want)
			}
		})
	}
}

func TestExecutable_EngineOpt_WithLogger(t *testing.T) {
	// setup types
	e := &Engine{logger: new(logrus.Entry)}

	// setup tests
	tests := []struct {
		failure bool
		name    string
		logger  *logrus.Entry
		want    *logrus.Entry
	}{
		{
			failure: false,
			name:    "logger set to new entry",
			logger:  new(logrus.Entry),
			want:    new(logrus.Entry),
		},
		{
			failure: false,
			name:    "logger set to nil",
			logger:  nil,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithLogger(test.logger)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithLogger for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("WithLogger returned err: %v", err)
			}

			if !reflect.DeepEqual(e.logger, test.want) {
				t.Errorf("WithLogger is %v, want %v", e.logger, test.want)
			}
		})
	}
}

func TestExecutable_EngineOpt_WithSkipCreation(t *testing.T) {
	// setup types
	e := &Engine{config: new(config)}

	// setup tests
	tests := []struct {
		failure      bool
		name         string
		skipCreation bool
		want         bool
	}{
		{
			failure:      false,
			name:         "skip creation set to true",
			skipCreation: true,
			want:         true,
		},
		{
			failure:      false,
			name:         "skip creation set to false",
			skipCreation: false,
			want:         false,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithSkipCreation(test.skipCreation)(e)

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

func TestExecutable_EngineOpt_WithContext(t *testing.T) {
	// setup types
	e := &Engine{config: new(config)}

	// setup tests
	tests := []struct {
		failure bool
		name    string
		ctx     context.Context
		want    context.Context
	}{
		{
			failure: false,
			name:    "context set to TODO",
			ctx:     context.TODO(),
			want:    context.TODO(),
		},
		{
			failure: false,
			name:    "context set to nil",
			ctx:     nil,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := WithContext(test.ctx)(e)

			if test.failure {
				if err == nil {
					t.Errorf("WithContext for %s should have returned err", test.name)
				}

				return
			}

			if err != nil {
				t.Errorf("WithContext returned err: %v", err)
			}

			if !reflect.DeepEqual(e.ctx, test.want) {
				t.Errorf("WithContext is %v, want %v", e.ctx, test.want)
			}
		})
	}
}
