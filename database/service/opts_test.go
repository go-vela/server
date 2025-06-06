// SPDX-License-Identifier: Apache-2.0

package service

import (
	"context"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func TestService_EngineOpt_WithClient(t *testing.T) {
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

func TestService_EngineOpt_WithLogger(t *testing.T) {
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

func TestService_EngineOpt_WithSkipCreation(t *testing.T) {
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

func TestService_EngineOpt_WithContext(t *testing.T) {
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
