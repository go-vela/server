// SPDX-License-Identifier: Apache-2.0

package token

import "testing"

func Test_imageParse(t *testing.T) {
	type args struct {
		image string
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantTag  string
		wantErr  bool
	}{
		{
			name: "image with tag",
			args: args{
				image: "alpine:1.20",
			},
			wantName: "alpine",
			wantTag:  "1.20",
			wantErr:  false,
		},
		{
			name: "image with tag and sha",
			args: args{
				image: "alpine:1.20@sha:fc0d4410fd2343cf6f7a75d5819001a34ca3b549fbab0c231b7aab49b57e9e43",
			},
			wantName: "alpine",
			wantTag:  "1.20",
			wantErr:  false,
		},
		{
			name: "image without latest tag",
			args: args{
				image: "alpine:latest",
			},
			wantName: "alpine",
			wantTag:  "latest",
			wantErr:  false,
		},
		{
			name: "image without tag",
			args: args{
				image: "alpine",
			},
			wantName: "alpine",
			wantTag:  "latest",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotTag, err := imageParse(tt.args.image)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("imageParse() gotName = %v, wantName %v", gotName, tt.wantName)
			}
			if gotTag != tt.wantTag {
				t.Errorf("imageParse() gotTag = %v, wantTag %v", gotTag, tt.wantTag)
			}
		})
	}
}
