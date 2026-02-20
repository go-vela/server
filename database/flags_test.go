// SPDX-License-Identifier: Apache-2.0

package database

import (
	"context"
	"testing"

	"github.com/urfave/cli/v3"
)

//nolint:gosec // ignore hardcoded credentials in tests
func TestDatabase_Flags(t *testing.T) {
	// deep copy flags since they are global variables and will hold onto modifications during testing
	deepCopyFlags := func(flags []cli.Flag) []cli.Flag {
		copiedFlags := make([]cli.Flag, len(flags))
		for i, flag := range flags {
			switch f := flag.(type) {
			case *cli.StringFlag:
				copyFlag := *f
				copiedFlags[i] = &copyFlag
			case *cli.Int64Flag:
				copyFlag := *f
				copiedFlags[i] = &copyFlag
			case *cli.Int32Flag:
				copyFlag := *f
				copiedFlags[i] = &copyFlag
			case *cli.IntFlag:
				copyFlag := *f
				copiedFlags[i] = &copyFlag
			case *cli.DurationFlag:
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

	// Define test cases
	tests := []struct {
		name    string
		flags   map[string]string
		wantErr bool
	}{
		{
			name: "happy path",
			flags: map[string]string{
				"database.driver":            "postgres",
				"database.addr":              "postgres://vela:pass@postgres:5432/vela?sslmode=disable",
				"database.encryption.key":    "C639A572E14D5075C526FDDD43E4ECF6",
				"database.compression.level": "3",
			},
			wantErr: false,
		},
		{
			name: "empty driver",
			flags: map[string]string{
				"database.driver": "",
				"database.addr":   "postgres://vela:pass@postgres:5432/vela?sslmode=disable",
			},
			wantErr: true,
		},
		{
			name: "empty addr",
			flags: map[string]string{
				"database.driver": "postgres",
				"database.addr":   "",
			},
			wantErr: true,
		},
		{
			name: "invalid addr",
			flags: map[string]string{
				"database.driver": "postgres",
				"database.addr":   "invalid-url/",
			},
			wantErr: true,
		},
		{
			name: "empty key",
			flags: map[string]string{
				"database.driver":         "postgres",
				"database.addr":           "postgres://vela:pass@postgres:5432/vela?sslmode=disable",
				"database.encryption.key": "",
			},
			wantErr: true,
		},
		{
			name: "invalid compression level",
			flags: map[string]string{
				"database.driver":            "postgres",
				"database.addr":              "postgres://vela:pass@postgres:5432/vela?sslmode=disable",
				"database.encryption.key":    "C639A572E14D5075C526FDDD43E4ECF6",
				"database.compression.level": "10",
			},
			wantErr: true,
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

			args := []string{"test"}
			// Set environment variables
			for key, value := range test.flags {
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
