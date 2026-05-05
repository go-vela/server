// SPDX-License-Identifier: Apache-2.0

package native

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-vela/server/api/types/settings"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/raw"
	"github.com/go-vela/server/compiler/types/yaml"
)

func TestNative_ValidateYAML_NoVersion(t *testing.T) {
	// setup types
	p := &yaml.Build{}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidateYAML(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_ValidateYAML_NoStagesOrSteps(t *testing.T) {
	// setup types
	p := &yaml.Build{
		Version: "v1",
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidateYAML(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_ValidateYAML_StagesAndSteps(t *testing.T) {
	// setup types
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
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidateYAML(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_ValidateYAML_RenderInLineStepTemplate(t *testing.T) {
	// setup types
	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Metadata: yaml.Metadata{
			RenderInline: true,
		},
		Steps: yaml.StepSlice{
			&yaml.Step{
				Commands: raw.StringSlice{"echo hello"},
				Image:    "alpine",
				Name:     str,
				Pull:     "always",
			},
			&yaml.Step{
				Template: yaml.StepTemplate{
					Name: "foo",
				},
				Name: str,
			},
		},
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidateYAML(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_ValidatePipeline_Services(t *testing.T) {
	// setup types
	str := "foo"
	p := &pipeline.Build{
		Version: "v1",
		Services: pipeline.ContainerSlice{
			&pipeline.Container{
				Image: "postgres",
				Name:  str,
				Ports: raw.StringSlice{"8080:8080"},
			},
		},
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				Commands: raw.StringSlice{"echo hello"},
				Image:    "alpine",
				Name:     str,
				Pull:     "always",
			},
		},
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidatePipeline(p)
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestNative_ValidateYAML_Services_NoName(t *testing.T) {
	// setup types
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
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidateYAML(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Services_NoImage(t *testing.T) {
	// setup types
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
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidateYAML(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Stages(t *testing.T) {
	// setup types
	str := "foo"
	p := &pipeline.Build{
		Version: "v1",
		Stages: pipeline.StageSlice{
			&pipeline.Stage{
				Name: str,
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
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
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidatePipeline(p)
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestNative_Validate_StagesSameName(t *testing.T) {
	// setup types
	strFoo := "foo"
	strBar := "bar"

	p := &pipeline.Build{
		Version: "v1",
		Stages: pipeline.StageSlice{
			&pipeline.Stage{
				Name: strFoo,
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						Commands: raw.StringSlice{"echo hello"},
						Image:    "alpine",
						Name:     strFoo,
						Pull:     "always",
					},
				},
			},
			&pipeline.Stage{
				Name: strFoo,
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						Commands: raw.StringSlice{"echo hello"},
						Image:    "alpine",
						Name:     strBar,
						Pull:     "always",
					},
				},
			},
		},
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidatePipeline(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Stages_NoName(t *testing.T) {
	// setup types
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
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidateYAML(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Stages_NoStepName(t *testing.T) {
	// setup types
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
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidateYAML(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Stages_NoImage(t *testing.T) {
	// setup types
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
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidateYAML(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Stages_NoCommands(t *testing.T) {
	// setup types
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
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidateYAML(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Stages_StepNameConflict(t *testing.T) {
	// setup types
	str := "foo"
	p := &pipeline.Build{
		Version: "v1",
		Stages: pipeline.StageSlice{
			&pipeline.Stage{
				Name: str,
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						Commands: raw.StringSlice{"echo hello"},
						Image:    "alpine",
						Name:     str,
						Pull:     "always",
					},
				},
			},
			&pipeline.Stage{
				Name: str,
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
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
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidatePipeline(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Stages_NeedsSelfReference(t *testing.T) {
	// setup types
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
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidateYAML(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Steps(t *testing.T) {
	// setup types
	str := "foo"
	p := &pipeline.Build{
		Version: "v1",
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				Commands: raw.StringSlice{"echo hello"},
				Image:    "alpine",
				Name:     str,
				Pull:     "always",
			},
		},
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidatePipeline(p)
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestNative_Validate_Steps_NoName(t *testing.T) {
	// setup types
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
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidateYAML(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Services_NameCollision(t *testing.T) {
	// setup types
	str := "foo"
	p := &pipeline.Build{
		Version: "v1",
		Services: pipeline.ContainerSlice{
			&pipeline.Container{
				Environment: raw.StringSliceMap{
					"FOO": "bar",
				},
				Image: "postgres",
				Name:  str,
				Pull:  "always",
			},
			&pipeline.Container{
				Environment: raw.StringSliceMap{
					"FOO": "bar",
				},
				Image: "kafka",
				Name:  str,
				Pull:  "always",
			},
		},
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidatePipeline(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Steps_NoImage(t *testing.T) {
	// setup types
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
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidateYAML(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Steps_NoCommands(t *testing.T) {
	// setup types
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
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidateYAML(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Steps_ExceedReportAs(t *testing.T) {
	// setup types
	str := "foo"

	reportSteps := pipeline.ContainerSlice{}

	for i := range 12 {
		reportStep := &pipeline.Container{
			Commands: raw.StringSlice{"echo hello"},
			Image:    "alpine",
			Name:     fmt.Sprintf("%s-%d", str, i),
			Pull:     "always",
			ReportAs: fmt.Sprintf("step-%d", i),
		}
		reportSteps = append(reportSteps, reportStep)
	}

	p := &pipeline.Build{
		Version: "v1",
		Steps:   reportSteps,
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidatePipeline(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_MultiReportAs(t *testing.T) {
	// setup types
	str := "foo"
	p := &pipeline.Build{
		Version: "v1",
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				Commands: raw.StringSlice{"echo hello"},
				Image:    "alpine",
				Name:     str,
				Pull:     "always",
				ReportAs: "bar",
			},
			&pipeline.Container{
				Commands: raw.StringSlice{"echo hello"},
				Image:    "alpine",
				Name:     str + "-2",
				Pull:     "always",
				ReportAs: "bar",
			},
		},
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidatePipeline(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_Validate_Steps_StepNameConflict(t *testing.T) {
	// setup types
	str := "foo"
	p := &pipeline.Build{
		Version: "v1",
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				Commands: raw.StringSlice{"echo hello"},
				Image:    "alpine",
				Name:     str,
				Pull:     "always",
			},
			&pipeline.Container{
				Commands: raw.StringSlice{"echo goodbye"},
				Image:    "alpine",
				Name:     str,
				Pull:     "always",
			},
		},
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidatePipeline(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
func TestNative_Validate_Artifact(t *testing.T) {
	// setup types
	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Steps: yaml.StepSlice{
			&yaml.Step{
				Commands: raw.StringSlice{"echo hello"},
				Image:    "alpine",
				Name:     str,
				Pull:     "always",
				Artifacts: yaml.Artifacts{
					Paths: []string{"results.xml", "artifacts.png"},
				},
			},
		},
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidateYAML(p)
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}
func TestNative_Validate_Secrets_SecretOriginNameConflict(t *testing.T) {
	// setup types
	str := "foo"
	p := &pipeline.Build{
		Version: "v1",
		Secrets: pipeline.SecretSlice{
			&pipeline.Secret{
				Origin: &pipeline.Container{
					Name:  "secrets",
					Image: "vault",
				},
			},
			&pipeline.Secret{
				Origin: &pipeline.Container{
					Name:  "secrets",
					Image: "vault",
				},
			},
		},
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				Commands: raw.StringSlice{"echo hello"},
				Image:    "alpine",
				Name:     str,
				Pull:     "always",
			},
		},
	}

	// run test
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	err = compiler.ValidatePipeline(p)
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestNative_CheckImageRestrictions_BlockedImage(t *testing.T) {
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	compiler.Compiler = settings.Compiler{
		BlockedImages: &[]settings.ImageRestriction{
			{Image: new("docker.io/blocked/image:latest"), Reason: new("this image is not allowed")},
		},
	}

	p := &pipeline.Build{
		Steps: pipeline.ContainerSlice{
			{Name: "blocked-step", Image: "blocked/image:latest"},
			{Name: "allowed-step", Image: "alpine:latest"},
		},
	}

	warnings, err := compiler.checkImageRestrictions(p)
	if err == nil {
		t.Errorf("checkImageRestrictions should have returned err for blocked image")
	}

	if len(warnings) != 0 {
		t.Errorf("checkImageRestrictions should not have returned warnings, got: %v", warnings)
	}
}

func TestNative_CheckImageRestrictions_WarnImage(t *testing.T) {
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	compiler.Compiler = settings.Compiler{
		WarnImages: &[]settings.ImageRestriction{
			{Image: new("docker.io/deprecated/image:latest"), Reason: new("this image is deprecated")},
		},
	}

	p := &pipeline.Build{
		Steps: pipeline.ContainerSlice{
			{Name: "warn-step", Image: "deprecated/image:latest"},
			{Name: "fine-step", Image: "alpine:latest"},
		},
	}

	warnings, err := compiler.checkImageRestrictions(p)
	if err != nil {
		t.Errorf("checkImageRestrictions returned unexpected err: %v", err)
	}

	if len(warnings) != 1 {
		t.Errorf("checkImageRestrictions should have returned 1 warning, got: %d", len(warnings))
	}
}

func TestNative_CheckImageRestrictions_WildcardPattern(t *testing.T) {
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	compiler.Compiler = settings.Compiler{
		BlockedImages: &[]settings.ImageRestriction{
			{Image: new("docker.io/blocked/*"), Reason: new("entire namespace is blocked")},
		},
	}

	p := &pipeline.Build{
		Steps: pipeline.ContainerSlice{
			{Name: "step-a", Image: "blocked/image-one:latest"},
			{Name: "step-b", Image: "blocked/image-two:v2"},
			{Name: "step-c", Image: "allowed/image:latest"},
		},
	}

	_, err = compiler.checkImageRestrictions(p)
	if err == nil {
		t.Errorf("checkImageRestrictions should have returned err for wildcard blocked images")
	}
}

func TestNative_CheckImageRestrictions_NoMatch(t *testing.T) {
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	compiler.Compiler = settings.Compiler{
		BlockedImages: &[]settings.ImageRestriction{
			{Image: new("docker.io/blocked/image:latest"), Reason: new("this image is not allowed")},
		},
		WarnImages: &[]settings.ImageRestriction{
			{Image: new("docker.io/deprecated/image:latest"), Reason: new("this image is deprecated")},
		},
	}

	p := &pipeline.Build{
		Steps: pipeline.ContainerSlice{
			{Name: "fine-step", Image: "alpine:latest"},
		},
	}

	warnings, err := compiler.checkImageRestrictions(p)
	if err != nil {
		t.Errorf("checkImageRestrictions returned unexpected err: %v", err)
	}

	if len(warnings) != 0 {
		t.Errorf("checkImageRestrictions should not have returned warnings, got: %v", warnings)
	}
}

func TestNative_CheckImageRestrictions_EmptyLists(t *testing.T) {
	compiler, err := FromCLICommand(context.Background(), testCommand(t, "http://foo.example.com"))
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	p := &pipeline.Build{
		Steps: pipeline.ContainerSlice{
			{Name: "step", Image: "alpine:latest"},
		},
	}

	warnings, err := compiler.checkImageRestrictions(p)
	if err != nil {
		t.Errorf("checkImageRestrictions returned unexpected err: %v", err)
	}

	if len(warnings) != 0 {
		t.Errorf("checkImageRestrictions should not have returned warnings, got: %v", warnings)
	}
}

func TestNative_MatchesImagePattern(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		image   string
		want    bool
	}{
		{
			name:    "exact match normalized",
			pattern: "docker.io/library/alpine:latest",
			image:   "alpine:latest",
			want:    true,
		},
		{
			name:    "wildcard tag",
			pattern: "docker.io/org/image:*",
			image:   "org/image:v1.2.3",
			want:    true,
		},
		{
			name:    "wildcard org",
			pattern: "docker.io/blocked/*",
			image:   "blocked/tool:latest",
			want:    true,
		},
		{
			name:    "no match",
			pattern: "docker.io/blocked/image:latest",
			image:   "allowed/image:latest",
			want:    false,
		},
		{
			name:    "empty pattern",
			pattern: "",
			image:   "alpine:latest",
			want:    false,
		},
		{
			name:    "empty image",
			pattern: "alpine:latest",
			image:   "",
			want:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := matchesImagePattern(test.pattern, test.image)

			if got != test.want {
				t.Errorf("matchesImagePattern(%q, %q) = %v, want %v", test.pattern, test.image, got, test.want)
			}
		})
	}
}
