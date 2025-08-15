// SPDX-License-Identifier: Apache-2.0

package queue

import (
	"context"
	"testing"

	"github.com/urfave/cli/v3"
)

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
			case *cli.StringSliceFlag:
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
				"queue.driver":     "redis",
				"queue.addr":       "redis://redis.example.com",
				"queue.routes":     "vela,worker",
				"queue.public-key": "CuS+EQAzofbk3tVFS3bt5f2tIb4YiJJC4nVMFQYQElg=",
			},
			wantErr: false,
		},
		{
			name: "empty driver",
			flags: map[string]string{
				"queue.addr":       "redis://redis.example.com",
				"queue.routes":     "vela,worker",
				"queue.public-key": "CuS+EQAzofbk3tVFS3bt5f2tIb4YiJJC4nVMFQYQElg=",
			},
			wantErr: true,
		},
		{
			name: "empty addr - validated in opts",
			flags: map[string]string{
				"queue.driver":     "redis",
				"queue.routes":     "vela,worker",
				"queue.public-key": "CuS+EQAzofbk3tVFS3bt5f2tIb4YiJJC4nVMFQYQElg=",
			},
		},
		{
			name: "invalid addr",
			flags: map[string]string{
				"queue.driver":     "redis",
				"queue.addr":       "redis.example.com",
				"queue.routes":     "vela,worker",
				"queue.public-key": "CuS+EQAzofbk3tVFS3bt5f2tIb4YiJJC4nVMFQYQElg=",
			},
			wantErr: true,
		},
		{
			name: "invalid addr - trailing slash",
			flags: map[string]string{
				"queue.driver":     "redis",
				"queue.addr":       "redis://redis.example.com/",
				"queue.routes":     "vela,worker",
				"queue.public-key": "CuS+EQAzofbk3tVFS3bt5f2tIb4YiJJC4nVMFQYQElg=",
			},
			wantErr: true,
		},
		{
			name: "empty public key - validated in opts",
			flags: map[string]string{
				"queue.driver": "redis",
				"queue.addr":   "redis://redis.example.com",
				"queue.routes": "vela,worker",
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
