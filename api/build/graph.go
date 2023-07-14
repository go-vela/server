// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"

	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// graph contains nodes, and relationships between nodes, or edges.
//
//	a node is a pipeline stage and its relevant steps.
//	an edge is a relationship between nodes, defined by the 'needs' tag.
type graph struct {
	Nodes map[int]*node `json:"nodes"`
	Edges []*edge       `json:"edges"`
}

// node represents is a pipeline stage and its relevant steps.
type node struct {
	Name  string          `json:"name"`
	Stage *pipeline.Stage `json:"stage"`
	Steps []*library.Step `json:"steps"`
	ID    int             `json:"id"`
}

type edge struct {
	Source      int `json:"source"`
	Destination int `json:"destination"`
}

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/graph builds GetBuildGraph
//
// Get directed a-cyclical graph for a build in the configured backend
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the org
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repo
//   required: true
//   type: string
// - in: path
//   name: build
//   description: Build number
//   required: true
//   type: integer
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: Successfully retrieved graph for the build
//     schema:
//       type: array
//       items:
//         "$ref": "#/definitions/Graph"
//   '500':
//     description: Unable to retrieve graph for the build
//     schema:
//       "$ref": "#/definitions/Error"

// GetBuildGraph represents the API handler to capture a
// directed a-cyclical graph for a build from the configured backend.
func GetBuildGraph(c *gin.Context) {
	// capture middleware values
	b := build.Retrieve(c)
	o := org.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	m := c.MustGet("metadata").(*types.Metadata)

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  u.GetName(),
	}).Infof("getting all steps for build %s", entry)

	// retrieve the steps for the build from the step table
	steps := []*library.Step{}
	page := 1
	perPage := 100
	for page > 0 {
		// retrieve build steps (per page) from the database
		stepsPart, _, err := database.FromContext(c).ListStepsForBuild(b, nil, page, perPage)
		if err != nil {
			retErr := fmt.Errorf("unable to retrieve steps for build %s: %w", entry, err)
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}

		// add page of steps to list steps
		steps = append(steps, stepsPart...)

		// assume no more pages exist if under 100 results are returned
		//
		// nolint: gomnd // ignore magic number
		if len(stepsPart) < 100 {
			page = 0
		} else {
			page++
		}
	}

	if len(steps) == 0 {
		retErr := fmt.Errorf("no steps found for build %s", entry)
		util.HandleError(c, http.StatusNotFound, retErr)
		return
	}

	logrus.Info("retrieving pipeline configuration file")

	baseErr := "unable to generate build graph"

	// send API call to capture the pipeline configuration file
	config, err := scm.FromContext(c).ConfigBackoff(u, r, b.GetCommit())
	if err != nil {
		// nolint: lll // ignore long line length due to error message
		retErr := fmt.Errorf("%s: failed to get pipeline configuration for %s: %v", baseErr, r.GetFullName(), err)
		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// variable to store changeset files
	var files []string
	// check if the build event is not issue_comment
	if !strings.EqualFold(b.GetEvent(), constants.EventComment) {
		// check if the build event is not pull_request
		if !strings.EqualFold(b.GetEvent(), constants.EventPull) {
			// send API call to capture list of files changed for the commit
			files, err = scm.FromContext(c).Changeset(u, r, b.GetCommit())
			if err != nil {
				retErr := fmt.Errorf("%s: failed to get changeset for %s: %v", baseErr, r.GetFullName(), err)
				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}
		}
	}

	logrus.Info("compiling pipeline")
	// parse and compile the pipeline configuration file
	p, _, err := compiler.FromContext(c).
		Duplicate().
		WithBuild(b).
		WithFiles(files).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u).
		Compile(config)
	if err != nil {
		// format the error message with extra information
		err = fmt.Errorf("unable to compile pipeline configuration for %s: %v", r.GetFullName(), err)

		// log the error for traceability
		logrus.Error(err.Error())

		retErr := fmt.Errorf("%s: %v", baseErr, err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// skip the build if only the init or clone steps are found
	skip := SkipEmptyBuild(p)
	if skip != "" {
		c.JSON(http.StatusOK, skip)
		return
	}

	logrus.Info("creating dag using 'needs'")

	// group library steps by stage name
	stages := map[string][]*library.Step{}
	for _, _step := range steps {
		if _, ok := stages[_step.GetStage()]; !ok {
			stages[_step.GetStage()] = []*library.Step{}
		}
		stages[_step.GetStage()] = append(stages[_step.GetStage()], _step)
	}

	// create nodes from pipeline stages
	nodes := make(map[int]*node)
	for _, stage := range p.Stages {
		// scrub the environment
		for _, step := range stage.Steps {
			step.Environment = nil
		}
		nodeID := len(nodes)
		node := node{
			Name:  stage.Name,
			Stage: stage,
			Steps: stages[stage.Name],
			ID:    nodeID,
		}
		nodes[nodeID] = &node
	}

	// create edges from nodes
	//   an edge is as a relationship between two nodes
	//   that is defined by the 'needs' tag
	edges := []*edge{}
	// loop over nodes
	for _, destinationNode := range nodes {
		// compare all nodes against all nodes
		for _, sourceNode := range nodes {
			// dont compare the same node
			if destinationNode.ID != sourceNode.ID {
				// check destination node needs
				for _, need := range (*destinationNode.Stage).Needs {
					// check if destination needs source stage
					if sourceNode.Stage.Name == need && need != "clone" {
						edge := &edge{
							Source:      sourceNode.ID,
							Destination: destinationNode.ID,
						}
						edges = append(edges, edge)
					}
				}
			}
		}
	}

	// for loop over edges, and collapse same parent edge

	if len(nodes) > 5000 {
		c.JSON(http.StatusInternalServerError, "too many nodes on this graph")
	}
	if len(edges) > 5000 {
		c.JSON(http.StatusInternalServerError, "too many edges on this graph")
	}

	// construct the response
	graph := graph{
		Nodes: nodes,
		Edges: edges,
	}

	c.JSON(http.StatusOK, graph)
}
