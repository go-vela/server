// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package schedule

import (
	"testing"
	"time"
)

func Test_validateEntry(t *testing.T) {
	type args struct {
		minimum time.Duration
		entry   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "exceeds minimum frequency",
			args: args{
				minimum: 30 * time.Minute,
				entry:   "* * * * *",
			},
			wantErr: true,
		},
		{
			name: "exceeds minimum frequency with tag",
			args: args{
				minimum: 30 * time.Minute,
				entry:   "@15minutes",
			},
			wantErr: true,
		},
		{
			name: "exceeds minimum frequency with scalene entry pattern",
			args: args{
				minimum: 30 * time.Minute,
				entry:   "1,2,45 * * * *",
			},
			wantErr: true,
		},
		{
			name: "meets minimum frequency",
			args: args{
				minimum: 30 * time.Second,
				entry:   "* * * * *",
			},
			wantErr: false,
		},
		{
			name: "meets minimum frequency with tag",
			args: args{
				minimum: 30 * time.Second,
				entry:   "@hourly",
			},
			wantErr: false,
		},
		{
			name: "meets minimum frequency with comma entry pattern",
			args: args{
				minimum: 15 * time.Minute,
				entry:   "0,15,30,45 * * * *",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateEntry(tt.args.minimum, tt.args.entry); (err != nil) != tt.wantErr {
				t.Errorf("validateEntry() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
