// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"testing"

	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/storage"
)

func TestSetupStorage(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantNil bool
		wantErr bool
	}{
		{
			name:    "storage disabled",
			args:    []string{"test"},
			wantNil: true,
			wantErr: false,
		},
		{
			name: "storage enabled with static credentials",
			args: []string{
				"test",
				"--storage.enable=true",
				"--storage.driver=minio",
				"--storage.addr=http://localhost:9000",
				"--storage.access.key=access",
				"--storage.secret.key=secret",
				"--storage.bucket.name=bucket",
			},
			wantNil: false,
			wantErr: false,
		},
		{
			name: "storage enabled with IAM credentials",
			args: []string{
				"test",
				"--storage.enable=true",
				"--storage.driver=minio",
				"--storage.addr=http://localhost:9000",
				"--storage.bucket.name=bucket",
				"--storage.use.iam=true",
			},
			wantNil: false,
			wantErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var (
				got    storage.Storage
				gotErr error
			)

			cmd := &cli.Command{
				Name:  "test",
				Flags: storage.Flags,
				Action: func(ctx context.Context, c *cli.Command) error {
					got, gotErr = setupStorage(ctx, c)
					return nil
				},
			}

			_ = cmd.Run(context.Background(), tc.args)

			if tc.wantErr {
				if gotErr == nil {
					t.Error("expected error, got nil")
				}

				return
			}

			if gotErr != nil {
				t.Errorf("unexpected error: %v", gotErr)
			}

			if tc.wantNil && got != nil {
				t.Errorf("expected nil storage, got %v", got)
			}

			if !tc.wantNil && got == nil {
				t.Error("expected storage client, got nil")
			}
		})
	}
}
