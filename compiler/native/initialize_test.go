// SPDX-License-Identifier: Apache-2.0

package native

import (
	"flag"
	"reflect"
	"testing"

	"github.com/urfave/cli/v2"

	"github.com/go-vela/server/compiler/types/yaml/yaml"
)

func TestNative_InitStage(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)

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

	want := &yaml.Build{
		Version: "v1",
		Stages: yaml.StageSlice{
			&yaml.Stage{
				Name: "init",
				Steps: yaml.StepSlice{
					&yaml.Step{
						Image: "#init",
						Name:  "init",
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
	}

	// run test
	compiler, err := FromCLIContext(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got, err := compiler.InitStage(p)
	if err != nil {
		t.Errorf("InitStage returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("InitStage is %v, want %v", got, want)
	}
}

func TestNative_InitStep(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	set.String("clone-image", defaultCloneImage, "doc")
	c := cli.NewContext(nil, set, nil)

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

	want := &yaml.Build{
		Version: "v1",
		Steps: yaml.StepSlice{
			&yaml.Step{
				Image: "#init",
				Name:  "init",
				Pull:  "not_present",
			},
			&yaml.Step{
				Image: "alpine",
				Name:  str,
				Pull:  "not_present",
			},
		},
	}
	// run test
	compiler, err := FromCLIContext(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got, err := compiler.InitStep(p)
	if err != nil {
		t.Errorf("InitStep returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("InitStep is %v, want %v", got, want)
	}
}
