// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"maps"
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

	validFlags := map[string]string{
		"log-level":                   "debug",
		"server-addr":                 "http://localhost:8080",
		"webui-addr":                  "http://localhost:8888",
		"vela-server-private-key":     "123abc",
		"clone-image":                 "target/vela-git-slim:latest",
		"default-build-limit":         "10",
		"max-build-limit":             "100",
		"default-repo-events":         "push,pull_request,tag,deployment,comment",
		"default-repo-approve-build":  "fork-always",
		"user-refresh-token-duration": "24h",
		"build-token-buffer-duration": "1h",
		"oidc-issuer":                 "http://localhost:8080",
		"github-driver":               "true",
		"github-url":                  "https://github.com",
		"github-token":                "123abc",
		"max-template-depth":          "5",
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
			name: "invalid server addr",
			override: map[string]string{
				"server-addr": "vela.example.com",
			},
			wantErr: true,
		},
		{
			name: "invalid server addr - trailing slash",
			override: map[string]string{
				"server-addr": "http://vela.example.com/",
			},
			wantErr: true,
		},
		{
			name: "no error on missing webui addr",
			override: map[string]string{
				"webui-addr": "",
			},
			wantErr: false,
		},
		{
			name: "invalid webui addr",
			override: map[string]string{
				"webui-addr": "vela.example.com",
			},
			wantErr: true,
		},
		{
			name: "invalid webui addr - trailing slash",
			override: map[string]string{
				"webui-addr": "http://vela.example.com/",
			},
			wantErr: true,
		},
		{
			name: "invalid clone image",
			override: map[string]string{
				"clone-image": "not-an-image:{",
			},
			wantErr: true,
		},
		{
			name: "0 build limit",
			override: map[string]string{
				"default-build-limit": "0",
			},
			wantErr: true,
		},
		{
			name: "0 max build limit",
			override: map[string]string{
				"max-build-limit": "0",
			},
			wantErr: true,
		},
		{
			name: "max build limit less than default",
			override: map[string]string{
				"max-build-limit": "2",
			},
			wantErr: true,
		},
		{
			name: "bad event type for default events",
			override: map[string]string{
				"default-repo-events": "not_an_event",
			},
			wantErr: true,
		},
		{
			name: "bad policy type for default approve build",
			override: map[string]string{
				"default-repo-approve-build": "not_a_policy",
			},
			wantErr: true,
		},
		{
			name: "refresh less than access duration",
			override: map[string]string{
				"user-refresh-token-duration": "10s",
			},
			wantErr: true,
		},
		{
			name: "zero build token buffer",
			override: map[string]string{
				"build-token-buffer-duration": "0m",
			},
			wantErr: true,
		},
		{
			name: "invalid url for oidc issuer",
			override: map[string]string{
				"oidc-issuer": "http://exa mple.com",
			},
			wantErr: true,
		},
		{
			name: "no github url for github driver",
			override: map[string]string{
				"github-url": "",
			},
			wantErr: true,
		},
		{
			name: "no github token for github driver",
			override: map[string]string{
				"github-token": "",
			},
			wantErr: true,
		},
		{
			name: "0 template depth",
			override: map[string]string{
				"max-template-depth": "0",
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

			copyMap := maps.Clone(validFlags)

			maps.Copy(copyMap, test.override)

			args := []string{"test"}
			// Set environment variables
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
