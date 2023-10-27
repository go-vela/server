// SPDX-License-Identifier: Apache-2.0

package build

import (
	"github.com/go-vela/types/pipeline"
)

// SkipEmptyBuild checks if the build should be skipped due to it
// not containing any steps besides init or clone.
//
//nolint:goconst // ignore init and clone constants
func SkipEmptyBuild(p *pipeline.Build) string {
	if len(p.Stages) == 1 {
		if p.Stages[0].Name == "init" {
			return "skipping build since only init stage found — it is likely no rulesets matched for the webhook payload"
		}
	}

	if len(p.Stages) == 2 {
		if p.Stages[0].Name == "init" && p.Stages[1].Name == "clone" {
			return "skipping build since only init and clone stages found — it is likely no rulesets matched for the webhook payload"
		}
	}

	if len(p.Steps) == 1 {
		if p.Steps[0].Name == "init" {
			return "skipping build since only init step found — it is likely no rulesets matched for the webhook payload"
		}
	}

	if len(p.Steps) == 2 {
		if p.Steps[0].Name == "init" && p.Steps[1].Name == "clone" {
			return "skipping build since only init and clone steps found — it is likely no rulesets matched for the webhook payload"
		}
	}

	return ""
}
