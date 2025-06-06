// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"reflect"
	"testing"

	"github.com/urfave/cli/v3"

	"github.com/go-vela/server/compiler/types/yaml/yaml"
)

const defaultCloneImage = "target/vela-git-slim:latest"

func TestNative_CloneStage(t *testing.T) {
	// setup types
	cmd := new(cli.Command)
	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "clone-image",
			Value: defaultCloneImage,
		},
	}

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Stages: yaml.StageSlice{
			&yaml.Stage{
				Name: str,
				Steps: yaml.StepSlice{
					&yaml.Step{
						Image: "alpine",
						Name:  str,
						Pull:  "not_present",
					},
				},
			},
		},
	}

	// setup tests
	tests := []struct {
		failure  bool
		local    bool
		pipeline *yaml.Build
		want     *yaml.Build
	}{
		{
			failure:  false,
			local:    false,
			pipeline: p,
			want: &yaml.Build{
				Version: "v1",
				Stages: yaml.StageSlice{
					&yaml.Stage{
						Name: "clone",
						Steps: yaml.StepSlice{
							&yaml.Step{
								Image: defaultCloneImage,
								Name:  "clone",
								Pull:  "not_present",
							},
						},
					},
					&yaml.Stage{
						Name: str,
						Steps: yaml.StepSlice{
							&yaml.Step{
								Image: "alpine",
								Name:  str,
								Pull:  "not_present",
							},
						},
					},
				},
			},
		},
		{
			failure:  false,
			local:    true,
			pipeline: p,
			want:     p,
		},
	}

	// run tests
	for _, test := range tests {
		compiler, err := FromCLICommand(context.Background(), cmd)
		if err != nil {
			t.Errorf("unable to create new compiler: %v", err)
		}

		// set the local field for the test
		compiler.WithLocal(test.local)

		got, err := compiler.CloneStage(test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("CloneStage should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CloneStage returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("CloneStage is %v, want %v", got, test.want)
		}
	}
}

func TestNative_CloneStep(t *testing.T) {
	// setup types
	cmd := new(cli.Command)
	cmd.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "clone-image",
			Value: defaultCloneImage,
		},
	}

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Steps: yaml.StepSlice{
			&yaml.Step{
				Image: "alpine",
				Name:  str,
				Pull:  "not_present",
			},
		},
	}

	// setup tests
	tests := []struct {
		failure  bool
		local    bool
		pipeline *yaml.Build
		want     *yaml.Build
	}{
		{
			failure:  false,
			local:    false,
			pipeline: p,
			want: &yaml.Build{
				Version: "v1",
				Steps: yaml.StepSlice{
					&yaml.Step{
						Image: defaultCloneImage,
						Name:  "clone",
						Pull:  "not_present",
					},
					&yaml.Step{
						Image: "alpine",
						Name:  str,
						Pull:  "not_present",
					},
				},
			},
		},
		{
			failure:  false,
			local:    true,
			pipeline: p,
			want:     p,
		},
	}

	// run tests
	for _, test := range tests {
		compiler, err := FromCLICommand(context.Background(), cmd)
		if err != nil {
			t.Errorf("Unable to create new compiler: %v", err)
		}

		// set the local field for the test
		compiler.WithLocal(test.local)

		got, err := compiler.CloneStep(test.pipeline)

		if test.failure {
			if err == nil {
				t.Errorf("CloneStep should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("CloneStep returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("CloneStep is %v, want %v", got, test.want)
		}
	}
}
