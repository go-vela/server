// SPDX-License-Identifier: Apache-2.0

package scm

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
				"scm.client":               "myClientID",
				"scm.secret":               "myClientSecret",
				"scm.app.id":               "42",
				"scm.app.private-key.path": "/secrets/key.pem",
				"scm.app.webhook-secret":   "myWebhookSecret",
			},
			wantErr: false,
		},
		{
			name: "custom addr",
			flags: map[string]string{
				"scm.addr":                 "https://git.company.com",
				"scm.client":               "myClientID",
				"scm.secret":               "myClientSecret",
				"scm.app.id":               "42",
				"scm.app.private-key.path": "/secrets/key.pem",
				"scm.app.webhook-secret":   "myWebhookSecret",
			},
			wantErr: false,
		},
		{
			name: "empty scm client",
			flags: map[string]string{
				"scm.secret":               "myClientSecret",
				"scm.app.id":               "42",
				"scm.app.private-key.path": "/secrets/key.pem",
				"scm.app.webhook-secret":   "myWebhookSecret",
			},
			wantErr: true,
		},
		{
			name: "empty scm secret",
			flags: map[string]string{
				"scm.client":               "myClientID",
				"scm.app.id":               "42",
				"scm.app.private-key.path": "/secrets/key.pem",
				"scm.app.webhook-secret":   "myWebhookSecret",
			},
			wantErr: true,
		},
		{
			name: "invalid addr",
			flags: map[string]string{
				"scm.client":               "myClientID",
				"scm.secret":               "myClientSecret",
				"scm.app.id":               "42",
				"scm.app.private-key.path": "/secrets/key.pem",
				"scm.app.webhook-secret":   "myWebhookSecret",
				"scm.addr":                 "github",
			},
			wantErr: true,
		},
		{
			name: "invalid addr - trailing slash",
			flags: map[string]string{
				"scm.client":               "myClientID",
				"scm.secret":               "myClientSecret",
				"scm.app.id":               "42",
				"scm.app.private-key.path": "/secrets/key.pem",
				"scm.app.webhook-secret":   "myWebhookSecret",
				"scm.addr":                 "https://github.com/",
			},
			wantErr: true,
		},
		{
			name: "app id with no private key",
			flags: map[string]string{
				"scm.client":             "myClientID",
				"scm.secret":             "myClientSecret",
				"scm.app.id":             "42",
				"scm.app.webhook-secret": "myWebhookSecret",
			},
			wantErr: true,
		},
		{
			name: "app id with a private key and a private key path",
			flags: map[string]string{
				"scm.client":               "myClientID",
				"scm.secret":               "myClientSecret",
				"scm.app.id":               "42",
				"scm.app.webhook-secret":   "myWebhookSecret",
				"scm.app.private-key":      "base64-encoded-key",
				"scm.app.private-key.path": "/secrets/key.pem",
			},
			wantErr: true,
		},
		{
			name: "empty webhook app secret",
			flags: map[string]string{
				"scm.client":               "myClientID",
				"scm.secret":               "myClientSecret",
				"scm.app.id":               "42",
				"scm.app.private-key.path": "/secrets/key.pem",
			},
			wantErr: true,
		},
		{
			name: "empty webhook app secret but no app",
			flags: map[string]string{
				"scm.client": "myClientID",
				"scm.secret": "myClientSecret",
			},
			wantErr: false,
		},
		{
			name: "empty webhook app secret but disabled webhook validation",
			flags: map[string]string{
				"scm.client":                      "myClientID",
				"scm.secret":                      "myClientSecret",
				"scm.app.id":                      "42",
				"scm.app.private-key.path":        "/secrets/key.pem",
				"vela-disable-webhook-validation": "true",
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
