// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/compiler/types/yaml"
	"github.com/go-vela/server/constants"
)

func TestNative_ScriptStages(t *testing.T) {
	// setup types
	baseEnv := environment(nil, nil, nil, nil, nil)

	s := yaml.StageSlice{
		&yaml.Stage{
			Name: "install",
			Steps: yaml.StepSlice{
				&yaml.Step{
					Commands:    []string{"./gradlew downloadDependencies"},
					Environment: baseEnv,
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
					Environment: baseEnv,
					Image:       "openjdk:latest",
					Name:        "test",
					Pull:        "always",
				},
			},
		},
	}

	baseEnv["HOME"] = constants.DefaultHomeDir
	baseEnv["SHELL"] = constants.DefaultShell

	installEnv := baseEnv
	installEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew downloadDependencies"})
	testEnv := baseEnv
	testEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew check"})

	want := yaml.StageSlice{
		&yaml.Stage{
			Name: "install",
			Steps: yaml.StepSlice{
				&yaml.Step{
					Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
					Entrypoint:  []string{constants.DefaultShell, "-c"},
					Environment: installEnv,
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
					Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
					Entrypoint:  []string{constants.DefaultShell, "-c"},
					Environment: testEnv,
					Image:       "openjdk:latest",
					Name:        "test",
					Pull:        "always",
				},
			},
		},
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	got, err := compiler.ScriptStages(s)
	if err != nil {
		t.Errorf("ScriptStages returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ScriptStages is %v, want %v", got, want)
	}
}

func TestNative_ScriptSteps(t *testing.T) {
	// setup types
	emptyEnv := environment(nil, nil, nil, nil, nil)

	baseEnv := emptyEnv
	baseEnv["HOME"] = constants.DefaultHomeDir
	baseEnv["SHELL"] = constants.DefaultShell

	installEnv := baseEnv
	installEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew downloadDependencies"})

	testEnv := baseEnv
	testEnv["VELA_BUILD_SCRIPT"] = generateScriptPosix([]string{"./gradlew check"})

	newHomeEnv := baseEnv
	newHomeEnv["HOME"] = "/usr/share/test"

	type args struct {
		s yaml.StepSlice
	}

	tests := []struct {
		name    string
		args    args
		want    yaml.StepSlice
		wantErr bool
	}{
		{"no user defined", args{s: yaml.StepSlice{
			&yaml.Step{
				Commands:    []string{"./gradlew downloadDependencies"},
				Environment: baseEnv,
				Image:       "openjdk:latest",
				Name:        "install",
				Pull:        "always",
			},
			&yaml.Step{
				Commands:    []string{"./gradlew check"},
				Environment: baseEnv,
				Image:       "openjdk:latest",
				Name:        "test",
				Pull:        "always",
			},
		}}, yaml.StepSlice{
			&yaml.Step{
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{constants.DefaultShell, "-c"},
				Environment: installEnv,
				Image:       "openjdk:latest",
				Name:        "install",
				Pull:        "always",
			},
			&yaml.Step{
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{constants.DefaultShell, "-c"},
				Environment: testEnv,
				Image:       "openjdk:latest",
				Name:        "test",
				Pull:        "always",
			},
		}, false},
		{"root user defined", args{s: yaml.StepSlice{
			&yaml.Step{
				Commands:    []string{"./gradlew downloadDependencies"},
				Environment: baseEnv,
				Image:       "openjdk:latest",
				Name:        "install",
				User:        "root",
				Pull:        "always",
			},
			&yaml.Step{
				Commands:    []string{"./gradlew check"},
				Environment: baseEnv,
				Image:       "openjdk:latest",
				Name:        "test",
				User:        "root",
				Pull:        "always",
			},
		}}, yaml.StepSlice{
			&yaml.Step{
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{constants.DefaultShell, "-c"},
				Environment: installEnv,
				Image:       "openjdk:latest",
				Name:        "install",
				User:        "root",
				Pull:        "always",
			},
			&yaml.Step{
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{constants.DefaultShell, "-c"},
				Environment: testEnv,
				Image:       "openjdk:latest",
				Name:        "test",
				User:        "root",
				Pull:        "always",
			},
		}, false},
		{"foo user defined", args{s: yaml.StepSlice{
			&yaml.Step{
				Commands:    []string{"./gradlew downloadDependencies"},
				Environment: baseEnv,
				Image:       "openjdk:latest",
				Name:        "install",
				User:        "foo",
				Pull:        "always",
			},
			&yaml.Step{
				Commands:    []string{"./gradlew check"},
				Environment: baseEnv,
				Image:       "openjdk:latest",
				Name:        "test",
				User:        "foo",
				Pull:        "always",
			},
		}}, yaml.StepSlice{
			&yaml.Step{
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{constants.DefaultShell, "-c"},
				Environment: installEnv,
				Image:       "openjdk:latest",
				Name:        "install",
				User:        "foo",
				Pull:        "always",
			},
			&yaml.Step{
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{constants.DefaultShell, "-c"},
				Environment: testEnv,
				Image:       "openjdk:latest",
				Name:        "test",
				User:        "foo",
				Pull:        "always",
			},
		}, false},
		{"user with home dir override", args{s: yaml.StepSlice{
			&yaml.Step{
				Commands:    []string{"./gradlew downloadDependencies"},
				Environment: newHomeEnv,
				Image:       "openjdk:latest",
				Name:        "install",
				User:        "test",
				Pull:        "always",
			},
			&yaml.Step{
				Commands:    []string{"./gradlew check"},
				Environment: newHomeEnv,
				Image:       "openjdk:latest",
				Name:        "test",
				User:        "test",
				Pull:        "always",
			},
		}}, yaml.StepSlice{
			&yaml.Step{
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{constants.DefaultShell, "-c"},
				Environment: newHomeEnv,
				Image:       "openjdk:latest",
				Name:        "install",
				User:        "test",
				Pull:        "always",
			},
			&yaml.Step{
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{constants.DefaultShell, "-c"},
				Environment: newHomeEnv,
				Image:       "openjdk:latest",
				Name:        "test",
				User:        "test",
				Pull:        "always",
			},
		}, false},
		{"empty env - no user", args{s: yaml.StepSlice{
			&yaml.Step{
				Commands:    []string{"./gradlew downloadDependencies"},
				Environment: emptyEnv,
				Image:       "openjdk:latest",
				Name:        "install",
				Pull:        "always",
			},
			&yaml.Step{
				Commands:    []string{"./gradlew check"},
				Environment: emptyEnv,
				Image:       "openjdk:latest",
				Name:        "test",
				Pull:        "always",
			},
		}}, yaml.StepSlice{
			&yaml.Step{
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{constants.DefaultShell, "-c"},
				Environment: installEnv,
				Image:       "openjdk:latest",
				Name:        "install",
				Pull:        "always",
			},
			&yaml.Step{
				Commands:    []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"},
				Entrypoint:  []string{constants.DefaultShell, "-c"},
				Environment: testEnv,
				Image:       "openjdk:latest",
				Name:        "test",
				Pull:        "always",
			},
		}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
			if err != nil {
				t.Errorf("Creating compiler returned err: %v", err)
			}

			got, err := compiler.ScriptSteps(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ScriptSteps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ScriptSteps() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
