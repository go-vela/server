// SPDX-License-Identifier: Apache-2.0

package native

import (
	"fmt"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/compiler/types/yaml"
)

// TransformStages converts a yaml configuration with stages into an executable pipeline.
func (c *Client) TransformStages(r *pipeline.RuleData, p *yaml.Build) (*pipeline.Build, error) {
	// create new executable pipeline
	pipeline := &pipeline.Build{
		Version:    p.Version,
		Metadata:   *p.Metadata.ToPipeline(),
		Stages:     *p.Stages.ToPipeline(),
		Secrets:    *p.Secrets.ToPipeline(),
		Services:   *p.Services.ToPipeline(),
		Worker:     *p.Worker.ToPipeline(),
		Deployment: *p.Deployment.ToPipeline(),
	}

	if c.netrc != nil {
		pipeline.Token = *c.netrc
		pipeline.TokenExp = c.netrcExp
	}

	build, err := pipeline.Purge(r)
	if err != nil {
		return nil, fmt.Errorf("unable to purge pipeline: %w", err)
	}

	return build, nil
}

// TransformSteps converts a yaml configuration with steps into an executable pipeline.
func (c *Client) TransformSteps(r *pipeline.RuleData, p *yaml.Build) (*pipeline.Build, error) {
	// create new executable pipeline
	pipeline := &pipeline.Build{
		Version:    p.Version,
		Metadata:   *p.Metadata.ToPipeline(),
		Deployment: *p.Deployment.ToPipeline(),
		Steps:      *p.Steps.ToPipeline(),
		Secrets:    *p.Secrets.ToPipeline(),
		Services:   *p.Services.ToPipeline(),
		Worker:     *p.Worker.ToPipeline(),
	}

	if c.netrc != nil {
		pipeline.Token = *c.netrc
		pipeline.TokenExp = c.netrcExp
	}

	build, err := pipeline.Purge(r)
	if err != nil {
		return nil, fmt.Errorf("unable to purge pipeline: %w", err)
	}

	return build, nil
}
