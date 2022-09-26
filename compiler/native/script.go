// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-vela/types/yaml"
)

// ScriptStages injects the script for each step in every stage in a yaml configuration.
func (c *client) ScriptStages(s yaml.StageSlice) (yaml.StageSlice, error) {
	// iterate through all stages
	for _, stage := range s {
		// inject the scripts into the steps for the stage
		steps, err := c.ScriptSteps(stage.Steps)
		if err != nil {
			return nil, err
		}

		stage.Steps = steps
	}

	return s, nil
}

// ScriptSteps injects the script for each step in a yaml configuration.
func (c *client) ScriptSteps(s yaml.StepSlice) (yaml.StepSlice, error) {
	// iterate through all steps
	for _, step := range s {
		// skip if no commands block for the step
		if len(step.Commands) == 0 {
			continue
		}

		// set the default home
		//nolint:goconst // ignore making this a constant for now
		home := "/root"
		// override the home value if user is defined
		// TODO:
		// - add ability to override user home directory
		if step.User != "" {
			home = fmt.Sprintf("/home/%s", step.User)
		}

		// generate script from commands
		script := generateScriptPosix(step.Commands)

		// set the entrypoint for the step
		step.Entrypoint = []string{"/bin/sh", "-c"}

		// set the commands for the step
		step.Commands = []string{"echo $VELA_BUILD_SCRIPT | base64 -d | /bin/sh -e"}

		// set the environment variables for the step
		step.Environment["VELA_BUILD_SCRIPT"] = script
		step.Environment["HOME"] = home
		//nolint:goconst // ignore making this a constant for now
		step.Environment["SHELL"] = "/bin/sh"
	}

	return s, nil
}

// generateScriptPosix is a helper function that generates a build script
// for a linux container using the given commands.
func generateScriptPosix(commands []string) string {
	var buf bytes.Buffer

	// iterate through each command provided
	for _, command := range commands {
		// safely escape entire command
		escaped := fmt.Sprintf("%q", command)

		// safely escape trace character
		escaped = strings.Replace(escaped, "$", `\$`, -1)

		// write escaped lines to buffer
		buf.WriteString(fmt.Sprintf(
			traceScript,
			escaped,
			command,
		))
	}

	// create build script with netrc and buffer information
	script := fmt.Sprintf(
		setupScript,
		buf.String(),
	)

	return base64.StdEncoding.EncodeToString([]byte(script))
}

// setupScript is a helper script this is added to the build to ensure
// a minimum set of environment variables are set correctly.
const setupScript = `
cat <<EOF > $HOME/.netrc
machine $VELA_NETRC_MACHINE
login $VELA_NETRC_USERNAME
password $VELA_NETRC_PASSWORD
EOF
chmod 0600 $HOME/.netrc
unset VELA_NETRC_MACHINE
unset VELA_NETRC_USERNAME
unset VELA_NETRC_PASSWORD
unset VELA_BUILD_SCRIPT
%s
`

// traceScript is a helper script that is added to the build script
// to trace a command.
const traceScript = `
echo $ %s
%s
`
