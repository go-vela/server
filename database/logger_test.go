// SPDX-License-Identifier: Apache-2.0

package database

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sirupsen/logrus"
)

func TestNewGormLogger(t *testing.T) {
	logger := logrus.NewEntry(logrus.New())

	type args struct {
		logger        *logrus.Entry
		slowThreshold time.Duration
		skipNotFound  bool
		showSQL       bool
	}

	tests := []struct {
		name string
		args args
		want *GormLogger
	}{
		{
			name: "logger set",
			args: args{
				logger:        logger,
				slowThreshold: time.Second,
				skipNotFound:  false,
				showSQL:       true,
			},
			want: &GormLogger{
				slowThreshold:         time.Second,
				skipErrRecordNotFound: false,
				showSQL:               true,
				entry:                 logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if diff := cmp.Diff(NewGormLogger(tt.args.logger, tt.args.slowThreshold, tt.args.skipNotFound, tt.args.showSQL), tt.want, cmpopts.EquateComparable(GormLogger{})); diff != "" {
				t.Errorf("NewGormLogger() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
