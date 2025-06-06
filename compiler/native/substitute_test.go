// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/server/compiler/types/yaml/yaml"
)

func Test_client_SubstituteStages(t *testing.T) {
	type args struct {
		stages yaml.StageSlice
	}

	tests := []struct {
		name    string
		args    args
		want    yaml.StageSlice
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				stages: yaml.StageSlice{
					{
						Name: "simple",
						Steps: yaml.StepSlice{
							{
								Commands:    []string{"echo ${FOO}", "echo $${BAR}"},
								Environment: map[string]string{"FOO": "baz", "BAR": "baz"},
								Image:       "alpine:latest",
								Name:        "simple",
								Pull:        "always",
							},
						},
					},
					{
						Name: "advanced",
						Steps: yaml.StepSlice{
							{
								Commands:    []string{"echo ${COMPLEX}"},
								Environment: map[string]string{"COMPLEX": "{\"hello\":\n  \"world\"}"},
								Image:       "alpine:latest",
								Name:        "advanced",
								Pull:        "always",
							},
						},
					},
					{
						Name: "not_found",
						Steps: yaml.StepSlice{
							{
								Commands:    []string{"echo $NOT_FOUND", "echo ${NOT_FOUND}", "echo $${NOT_FOUND}"},
								Environment: map[string]string{"FOO": "baz", "BAR": "baz"},
								Image:       "alpine:latest",
								Name:        "not_found",
								Pull:        "always",
							},
						},
					},
				},
			},
			want: yaml.StageSlice{
				{
					Name: "simple",
					Steps: yaml.StepSlice{
						{
							Commands:    []string{"echo baz", "echo ${BAR}"},
							Environment: map[string]string{"FOO": "baz", "BAR": "baz"},
							Image:       "alpine:latest",
							Name:        "simple",
							Pull:        "always",
						},
					},
				},
				{
					Name: "advanced",
					Steps: yaml.StepSlice{
						{
							Commands:    []string{"echo \"{\\\"hello\\\":\\n  \\\"world\\\"}\""},
							Environment: map[string]string{"COMPLEX": "{\"hello\":\n  \"world\"}"},
							Image:       "alpine:latest",
							Name:        "advanced",
							Pull:        "always",
						},
					},
				},
				{
					Name: "not_found",
					Steps: yaml.StepSlice{
						{
							Commands:    []string{"echo $NOT_FOUND", "echo ${NOT_FOUND}", "echo ${NOT_FOUND}"},
							Environment: map[string]string{"FOO": "baz", "BAR": "baz"},
							Image:       "alpine:latest",
							Name:        "not_found",
							Pull:        "always",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
			if err != nil {
				t.Errorf("Creating compiler returned err: %v", err)
			}

			got, err := compiler.SubstituteStages(tt.args.stages)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubstituteStages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("SubstituteStages() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func Test_client_SubstituteSteps(t *testing.T) {
	type args struct {
		steps yaml.StepSlice
	}

	tests := []struct {
		name    string
		args    args
		want    yaml.StepSlice
		wantErr bool
	}{
		{
			name: "steps",
			args: args{
				steps: yaml.StepSlice{
					{
						Commands:    []string{"echo ${FOO}", "echo $${BAR}"},
						Environment: map[string]string{"FOO": "baz", "BAR": "baz"},
						Image:       "alpine:latest",
						Name:        "simple",
						Pull:        "always",
					},
					{
						Commands:    []string{"echo ${COMPLEX}"},
						Environment: map[string]string{"COMPLEX": "{\"hello\":\n  \"world\"}"},
						Image:       "alpine:latest",
						Name:        "advanced",
						Pull:        "always",
					},
					{
						Commands:    []string{"echo $NOT_FOUND", "echo ${NOT_FOUND}", "echo $${NOT_FOUND}"},
						Environment: map[string]string{"FOO": "baz", "BAR": "baz"},
						Image:       "alpine:latest",
						Name:        "not_found",
						Pull:        "always",
					},
				},
			},
			want: yaml.StepSlice{
				{
					Commands:    []string{"echo baz", "echo ${BAR}"},
					Environment: map[string]string{"FOO": "baz", "BAR": "baz"},
					Image:       "alpine:latest",
					Name:        "simple",
					Pull:        "always",
				},
				{
					Commands:    []string{"echo \"{\\\"hello\\\":\\n  \\\"world\\\"}\""},
					Environment: map[string]string{"COMPLEX": "{\"hello\":\n  \"world\"}"},
					Image:       "alpine:latest",
					Name:        "advanced",
					Pull:        "always",
				},
				{
					Commands:    []string{"echo $NOT_FOUND", "echo ${NOT_FOUND}", "echo ${NOT_FOUND}"},
					Environment: map[string]string{"FOO": "baz", "BAR": "baz"},
					Image:       "alpine:latest",
					Name:        "not_found",
					Pull:        "always",
				},
			},
			wantErr: false,
		},
		{
			name: "template",
			args: args{
				steps: yaml.StepSlice{
					{
						Name: "sample",
						Template: yaml.StepTemplate{
							Name: "go",
							Variables: map[string]interface{}{
								"build_author": "${BUILD_AUTHOR}",
								"unknown":      "${DEPLOYMENT_PARAMETER_API_IMAGE}",
							},
						},
						Environment: map[string]string{
							"BUILD_AUTHOR": "testauthor",
						},
					},
				},
			},
			want: yaml.StepSlice{
				{
					Name: "sample",
					Template: yaml.StepTemplate{
						Name: "go",
						Variables: map[string]interface{}{
							"build_author": "testauthor",
							"unknown":      "${DEPLOYMENT_PARAMETER_API_IMAGE}",
						},
					},
					Environment: map[string]string{
						"BUILD_AUTHOR": "testauthor",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
			if err != nil {
				t.Errorf("Creating compiler returned err: %v", err)
			}

			got, err := compiler.SubstituteSteps(tt.args.steps)
			if (err != nil) != tt.wantErr {
				t.Errorf("SubstituteSteps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("SubstituteSteps() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
