// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"flag"
	"testing"

	"github.com/go-vela/types/raw"
	"github.com/go-vela/types/yaml"

	"github.com/urfave/cli/v2"
)

func TestNative_Validate_NoVersion(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	p := &yaml.Build{}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_NoStagesOrSteps(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	p := &yaml.Build{
		Version: "v1",
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_StagesAndSteps(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Stages: yaml.StageSlice{
			&yaml.Stage{
				Name: str,
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands: raw.StringSlice{"echo hello"},
						Image:    "alpine",
						Name:     str,
						Pull:     "always",
					},
				},
			},
		},
		Steps: yaml.StepSlice{
			&yaml.Step{
				Commands: raw.StringSlice{"echo hello"},
				Image:    "alpine",
				Name:     str,
				Pull:     "always",
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Services(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Services: yaml.ServiceSlice{
			&yaml.Service{
				Image: "postgres",
				Name:  str,
				Ports: raw.StringSlice{"8080:8080"},
			},
		},
		Steps: yaml.StepSlice{
			&yaml.Step{
				Commands: raw.StringSlice{"echo hello"},
				Image:    "alpine",
				Name:     str,
				Pull:     "always",
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestNative_Validate_Services_NoName(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Services: yaml.ServiceSlice{
			&yaml.Service{
				Image: "postgres",
				Name:  "",
				Ports: raw.StringSlice{"8080:8080"},
			},
		},
		Steps: yaml.StepSlice{
			&yaml.Step{
				Commands: raw.StringSlice{"echo hello"},
				Image:    "alpine",
				Name:     str,
				Pull:     "always",
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Services_NoImage(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Services: yaml.ServiceSlice{
			&yaml.Service{
				Image: "",
				Name:  str,
				Ports: raw.StringSlice{"8080:8080"},
			},
		},
		Steps: yaml.StepSlice{
			&yaml.Step{
				Commands: raw.StringSlice{"echo hello"},
				Image:    "alpine",
				Name:     str,
				Pull:     "always",
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Stages(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Stages: yaml.StageSlice{
			&yaml.Stage{
				Name: str,
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands: raw.StringSlice{"echo hello"},
						Image:    "alpine",
						Name:     str,
						Pull:     "always",
					},
				},
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestNative_Validate_Stages_NoName(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Stages: yaml.StageSlice{
			&yaml.Stage{
				Name: "",
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands: raw.StringSlice{"echo hello"},
						Name:     str,
						Pull:     "always",
					},
				},
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Stages_NoStepName(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Stages: yaml.StageSlice{
			&yaml.Stage{
				Name: str,
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands: raw.StringSlice{"echo hello"},
						Name:     "",
						Pull:     "always",
					},
				},
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Stages_NoImage(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Stages: yaml.StageSlice{
			&yaml.Stage{
				Name: str,
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands: raw.StringSlice{"echo hello"},
						Name:     str,
						Pull:     "always",
					},
				},
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Stages_NoCommands(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
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
						Pull:  "always",
					},
				},
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Stages_NeedsSelfReference(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Stages: yaml.StageSlice{
			&yaml.Stage{
				Name:  str,
				Needs: raw.StringSlice{str},
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands: raw.StringSlice{"echo hello"},
						Image:    "alpine",
						Name:     str,
						Pull:     "always",
					},
				},
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Steps(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Steps: yaml.StepSlice{
			&yaml.Step{
				Commands: raw.StringSlice{"echo hello"},
				Image:    "alpine",
				Name:     str,
				Pull:     "always",
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestNative_Validate_Steps_NoName(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	p := &yaml.Build{
		Version: "v1",
		Steps: yaml.StepSlice{
			&yaml.Step{
				Commands: raw.StringSlice{"echo hello"},
				Name:     "",
				Pull:     "always",
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Steps_NoImage(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Steps: yaml.StepSlice{
			&yaml.Step{
				Commands: raw.StringSlice{"echo hello"},
				Name:     str,
				Pull:     "always",
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Steps_NoCommands(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Steps: yaml.StepSlice{
			&yaml.Step{
				Image: "alpine",
				Name:  str,
				Pull:  "always",
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.Validate(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
