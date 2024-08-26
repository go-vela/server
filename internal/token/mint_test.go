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
			got, got1, err := imageParse(tt.args.image)
			if (err != nil) != tt.wantErr {
				t.Errorf("imageParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantName {
				t.Errorf("imageParse() got = %v, wantName %v", got, tt.wantName)
			}
			if got1 != tt.wantTag {
				t.Errorf("imageParse() got1 = %v, wantName %v", got1, tt.wantTag)
			}
		})
	}
}
