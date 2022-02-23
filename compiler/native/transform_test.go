// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"flag"
	"reflect"
	"testing"

	"github.com/go-vela/types"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/types/yaml"

	"github.com/urfave/cli/v2"
)

func TestNative_TransformStages(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	m := &types.Metadata{
		Database: &types.Database{
			Driver: "foo",
			Host:   "foo",
		},
		Queue: &types.Queue{
			Channel: "foo",
			Driver:  "foo",
			Host:    "foo",
		},
		Source: &types.Source{
			Driver: "foo",
			Host:   "foo",
		},
		Vela: &types.Vela{
			Address:    "foo",
			WebAddress: "foo",
		},
	}

	p := &yaml.Build{
		Version: "v1",
		Services: yaml.ServiceSlice{
			&yaml.Service{
				Ports: []string{"5432:5432"},
				Name:  "postgres backend",
				Image: "postgres:latest",
			},
		},
		Worker: yaml.Worker{
			Flavor:   "16cpu8gb",
			Platform: "gcp",
		},
		Stages: yaml.StageSlice{
			&yaml.Stage{
				Name: "install deps",
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands:    []string{"./gradlew downloadDependencies"},
						Environment: environment(nil, nil, nil, nil),
						Image:       "openjdk:latest",
						Name:        "install",
						Pull:        "always",
					},
				},
			},
			&yaml.Stage{
				Name:  "test",
				Needs: []string{"install"},
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands:    []string{"./gradlew check"},
						Environment: environment(nil, nil, nil, nil),
						Image:       "openjdk:latest",
						Name:        "test",
						Pull:        "always",
						Ruleset: yaml.Ruleset{
							If: yaml.Rules{
								Event: []string{"push"},
							},
							Operator: "and",
						},
					},
				},
			},
		},
		Secrets: yaml.SecretSlice{
			&yaml.Secret{
				Name: "foobar",
				Origin: yaml.Origin{
					Image: "vault:latest",
					Name:  "vault",
					Pull:  "always",
				},
			},
		},
	}

	// setup tests
	tests := []struct {
		failure  bool
		local    bool
		pipeline *yaml.Build
		want     *pipeline.Build
	}{
		{
			failure:  false,
			local:    false,
			pipeline: p,
			want: &pipeline.Build{
				ID:      "__0",
				Version: "v1",
				Metadata: pipeline.Metadata{
					Clone: true,
				},
				Services: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:     "service___0_postgres backend",
						Ports:  []string{"5432:5432"},
						Name:   "postgres backend",
						Image:  "postgres:latest",
						Number: 1,
						Detach: true,
					},
				},
				Worker: pipeline.Worker{
					Flavor:   "16cpu8gb",
					Platform: "gcp",
				},
				Stages: pipeline.StageSlice{
					&pipeline.Stage{
						Name: "install deps",
						Steps: pipeline.ContainerSlice{
							&pipeline.Container{
								ID:          "__0_install deps_install",
								Commands:    []string{"./gradlew downloadDependencies"},
								Directory:   "/vela/src",
								Environment: environment(nil, nil, nil, nil),
								Image:       "openjdk:latest",
								Name:        "install",
								Number:      1,
								Pull:        "always",
							},
						},
					},
				},
				Secrets: pipeline.SecretSlice{
					&pipeline.Secret{
						Name: "foobar",
						Origin: &pipeline.Container{
							ID:     "secret___0_vault",
							Name:   "vault",
							Image:  "vault:latest",
							Pull:   "always",
							Number: 1,
						},
					},
				},
			},
		},
		{
			failure:  false,
			local:    true,
			pipeline: p,
			want: &pipeline.Build{
				ID:      "localOrg_localRepo_1",
				Version: "v1",
				Metadata: pipeline.Metadata{
					Clone: true,
				},
				Services: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:     "service_localOrg_localRepo_1_postgres backend",
						Ports:  []string{"5432:5432"},
						Name:   "postgres backend",
						Image:  "postgres:latest",
						Number: 1,
						Detach: true,
					},
				},
				Worker: pipeline.Worker{
					Flavor:   "16cpu8gb",
					Platform: "gcp",
				},
				Stages: pipeline.StageSlice{
					&pipeline.Stage{
						Name: "install deps",
						Steps: pipeline.ContainerSlice{
							&pipeline.Container{
								ID:          "localOrg_localRepo_1_install deps_install",
								Commands:    []string{"./gradlew downloadDependencies"},
								Directory:   "/vela/src",
								Environment: environment(nil, nil, nil, nil),
								Image:       "openjdk:latest",
								Name:        "install",
								Number:      1,
								Pull:        "always",
							},
						},
					},
				},
				Secrets: pipeline.SecretSlice{
					&pipeline.Secret{
						Name: "foobar",
						Origin: &pipeline.Container{
							ID:     "secret_localOrg_localRepo_1_vault",
							Name:   "vault",
							Image:  "vault:latest",
							Pull:   "always",
							Number: 1,
						},
					},
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		compiler, err := New(c)
		if err != nil {
			t.Errorf("unable to create new compiler: %v", err)
		}

		// set the metadata field for the test
		compiler.WithMetadata(m)

		// set the local field for the test
		compiler.WithLocal(test.local)

		got, err := compiler.TransformStages(new(pipeline.RuleData), test.pipeline)
		if err != nil {
			t.Errorf("TransformStages returned err: %v", err)
		}

		// WARNING: hack to compare stages
		//
		// Channel values can only be compared for equality.
		// Two channel values are considered equal if they
		// originated from the same make call meaning they
		// refer to the same channel value in memory.
		for i, stage := range got.Stages {
			tmp := test.want.Stages

			tmp[i].Done = stage.Done
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("TransformStages is %v, want %v", got, test.want)
		}
	}
}

func TestNative_TransformSteps(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	m := &types.Metadata{
		Database: &types.Database{
			Driver: "foo",
			Host:   "foo",
		},
		Queue: &types.Queue{
			Channel: "foo",
			Driver:  "foo",
			Host:    "foo",
		},
		Source: &types.Source{
			Driver: "foo",
			Host:   "foo",
		},
		Vela: &types.Vela{
			Address:    "foo",
			WebAddress: "foo",
		},
	}

	p := &yaml.Build{
		Version: "v1",
		Services: yaml.ServiceSlice{
			&yaml.Service{
				Ports: []string{"5432:5432"},
				Name:  "postgres backend",
				Image: "postgres:latest",
			},
		},
		Worker: yaml.Worker{
			Flavor:   "16cpu8gb",
			Platform: "gcp",
		},
		Steps: yaml.StepSlice{
			&yaml.Step{
				Commands:    []string{"./gradlew downloadDependencies"},
				Environment: environment(nil, nil, nil, nil),
				Image:       "openjdk:latest",
				Name:        "install deps",
				Pull:        "always",
			},
			&yaml.Step{
				Commands:    []string{"./gradlew check"},
				Environment: environment(nil, nil, nil, nil),
				Image:       "openjdk:latest",
				Name:        "test",
				Pull:        "always",
				Ruleset: yaml.Ruleset{
					If: yaml.Rules{
						Event: []string{"push"},
					},
					Operator: "and",
				},
			},
		},
		Secrets: yaml.SecretSlice{
			&yaml.Secret{
				Name: "foobar",
				Origin: yaml.Origin{
					Image: "vault:latest",
					Name:  "vault",
					Pull:  "always",
				},
			},
		},
	}

	// setup tests
	tests := []struct {
		failure  bool
		local    bool
		pipeline *yaml.Build
		want     *pipeline.Build
	}{
		{
			failure:  false,
			local:    false,
			pipeline: p,
			want: &pipeline.Build{
				ID:      "__0",
				Version: "v1",
				Metadata: pipeline.Metadata{
					Clone: true,
				},
				Services: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:     "service___0_postgres backend",
						Ports:  []string{"5432:5432"},
						Name:   "postgres backend",
						Image:  "postgres:latest",
						Number: 1,
						Detach: true,
					},
				},
				Worker: pipeline.Worker{
					Flavor:   "16cpu8gb",
					Platform: "gcp",
				},
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "step___0_install deps",
						Commands:    []string{"./gradlew downloadDependencies"},
						Directory:   "/vela/src",
						Environment: environment(nil, nil, nil, nil),
						Image:       "openjdk:latest",
						Name:        "install deps",
						Number:      1,
						Pull:        "always",
					},
				},
				Secrets: pipeline.SecretSlice{
					&pipeline.Secret{
						Name: "foobar",
						Origin: &pipeline.Container{
							ID:     "secret___0_vault",
							Name:   "vault",
							Image:  "vault:latest",
							Pull:   "always",
							Number: 1,
						},
					},
				},
			},
		},
		{
			failure:  false,
			local:    true,
			pipeline: p,
			want: &pipeline.Build{
				ID:      "localOrg_localRepo_1",
				Version: "v1",
				Metadata: pipeline.Metadata{
					Clone: true,
				},
				Services: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:     "service_localOrg_localRepo_1_postgres backend",
						Ports:  []string{"5432:5432"},
						Name:   "postgres backend",
						Image:  "postgres:latest",
						Number: 1,
						Detach: true,
					},
				},
				Worker: pipeline.Worker{
					Flavor:   "16cpu8gb",
					Platform: "gcp",
				},
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:          "step_localOrg_localRepo_1_install deps",
						Commands:    []string{"./gradlew downloadDependencies"},
						Directory:   "/vela/src",
						Environment: environment(nil, nil, nil, nil),
						Image:       "openjdk:latest",
						Name:        "install deps",
						Number:      1,
						Pull:        "always",
					},
				},
				Secrets: pipeline.SecretSlice{
					&pipeline.Secret{
						Name: "foobar",
						Origin: &pipeline.Container{
							ID:     "secret_localOrg_localRepo_1_vault",
							Name:   "vault",
							Image:  "vault:latest",
							Pull:   "always",
							Number: 1,
						},
					},
				},
			},
		},
	}

	// run tests
	for _, test := range tests {
		compiler, err := New(c)
		if err != nil {
			t.Errorf("unable to create new compiler: %v", err)
		}

		// set the metadata field for the test
		compiler.WithMetadata(m)

		// set the local field for the test
		compiler.WithLocal(test.local)

		got, err := compiler.TransformSteps(new(pipeline.RuleData), test.pipeline)
		if err != nil {
			t.Errorf("TransformSteps returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("TransformSteps is %v, want %v", got, test.want)
		}
	}
}
