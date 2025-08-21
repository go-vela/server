// SPDX-License-Identifier: Apache-2.0

package storage

import (
	"context"
	"maps"
	"testing"

	"github.com/urfave/cli/v3"
)

func TestStorage_Flags(t *testing.T) {
	// deep copy flags since they are global variables and will hold onto modifications during testing
	deepCopyFlags := func(flags []cli.Flag) []cli.Flag {
		copiedFlags := make([]cli.Flag, len(flags))
		for i, flag := range flags {
			switch f := flag.(type) {
			case *cli.StringFlag:
				copyFlag := *f
				copiedFlags[i] = &copyFlag
			case *cli.BoolFlag:
				copyFlag := *f
				copiedFlags[i] = &copyFlag
			default:
				t.Fatalf("unsupported flag type: %T", f)
			}
		}
		return copiedFlags
	}

	validFlags := map[string]string{
		"storage.enable":      "true",
		"storage.driver":      "s3",
		"storage.addr":        "https://s3.amazonaws.com",
		"storage.access.key":  "test-access-key",
		"storage.secret.key":  "test-secret-key",
		"storage.bucket.name": "test-bucket",
		"storage.use.ssl":     "true",
	}

	// Define test cases
	tests := []struct {
		name     string
		override map[string]string
		wantErr  bool
	}{
		{
			name:    "happy path",
			wantErr: false,
		},
		{
			name: "invalid storage addr - no scheme",
			override: map[string]string{
				"storage.addr": "s3.amazonaws.com",
			},
			wantErr: true,
		},
		{
			name: "invalid storage addr - trailing slash",
			override: map[string]string{
				"storage.addr": "https://s3.amazonaws.com/",
			},
			wantErr: true,
		},
		{
			name: "valid storage addr with port",
			override: map[string]string{
				"storage.addr": "https://localhost:9000",
			},
			wantErr: false,
		},
		{
			name: "valid storage addr with http scheme",
			override: map[string]string{
				"storage.addr": "http://minio.local:9000",
			},
			wantErr: false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create a new command with a deep copy of the Flags slice
			cmd := cli.Command{
				Name: "test",
				Action: func(_ context.Context, _ *cli.Command) error {
					return nil
				},
				Flags: deepCopyFlags(Flags),
			}

			copyMap := maps.Clone(validFlags)

			maps.Copy(copyMap, test.override)

			args := []string{"test"}
			// Set command line arguments
			for key, value := range copyMap {
				if len(value) == 0 {
					continue
				}
				args = append(args, `--`+key+"="+value)
			}

			// Run command
			err := cmd.Run(context.Background(), args)

			// Check the result
			if (err != nil) != test.wantErr {
				t.Errorf("error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}
