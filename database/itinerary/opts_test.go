// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package itinerary

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/constants"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

func TestItinerary_EngineOpt_WithClient(t *testing.T) {
	// setup types
	e := &engine{client: new(gorm.DB)}

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

func TestItinerary_EngineOpt_WithCompressionLevel(t *testing.T) {
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
				t.Errorf("WithCompressionLevel returned err: %v", err)
			}

			if !reflect.DeepEqual(e.config.CompressionLevel, test.want) {
				t.Errorf("WithCompressionLevel is %v, want %v", e.config.CompressionLevel, test.want)
			}
		})
	}
}

func TestItinerary_EngineOpt_WithLogger(t *testing.T) {
	// setup types
	e := &engine{logger: new(logrus.Entry)}

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

func TestItinerary_EngineOpt_WithSkipCreation(t *testing.T) {
	// setup types
	e := &engine{config: new(config)}

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

func TestItinerary_EngineOpt_WithDriver(t *testing.T) {
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
			name:    "postgres",
			driver:  constants.DriverPostgres,
			want:    constants.DriverPostgres,
		},
		{
			failure: false,
			name:    "sqlite",
			driver:  constants.DriverSqlite,
			want:    constants.DriverSqlite,
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
				t.Errorf("WithDriver returned err: %v", err)
			}

			if !reflect.DeepEqual(e.config.CompressionLevel, test.want) {
				t.Errorf("WithDriver is %v, want %v", e.config.CompressionLevel, test.want)
			}
		})
	}
}
