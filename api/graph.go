// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"fmt"
	"net/http"
	"sort"
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
//   a graphs is a collection of subgraphs and edges.
//   a subgraph is a cluster of nodes.
//   a node is a pipeline step.
//   an edge is a connection between two nodes on the graph.
//   a connection between two nodes is defined by the 'needs' tag.
type graph struct {
	Subgraphs  map[int]*subgraph `json:"subgraphs"`
	StageNodes []*stagenode      `json:"stage_nodes"`
	StageEdges []*edge           `json:"stage_edges"`
}

// subgraph represents is a pipeline stage and its relevant steps.
type subgraph struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	StepNodes []*stepnode     `json:"step_nodes"`
	StepEdges []*edge         `json:"step_edges"`
	Stage     *pipeline.Stage `json:"stage,omitempty"`
}

// stagenode represents is a pipeline stage and its relevant steps.
type stagenode struct {
	ID    int             `json:"id"`
	Name  string          `json:"name"`
	Stage *pipeline.Stage `json:"stage,omitempty"`
	Steps []*library.Step `json:"steps,omitempty"`
}

// stepnode represents is a pipeline step and its relevant info.
type stepnode struct {
	ID   int           `json:"id"`
	Name string        `json:"name"`
	Step *library.Step `json:"step,omitempty"`
}

// an edge points between two stagenodes.
type edge struct {
	SourceID   int    `json:"source_id"`
	SourceName string `json:"source_name"`

	DestinationID   int    `json:"destination_id"`
	DestinationName string `json:"destination_name"`
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
		stepsPart, err := database.FromContext(c).GetBuildStepList(b, page, perPage)
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
	p, err := compiler.FromContext(c).
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
	skip := skipEmptyBuild(p)
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

	// create subgraphs from pipeline stages
	subgraphs := make(map[int]*subgraph)
	stageNodes := make([]*stagenode, 0)
	for _, stage := range p.Stages {
		steps := stages[stage.Name]
		if len(steps) == 0 {
			// somehow we have a stage with no steps
			break
		}

		// sort by step number
		sort.Slice(steps, func(i int, j int) bool {
			return steps[i].GetNumber() < steps[j].GetNumber()
		})

		for _, step := range stage.Steps {
			// scrub the container environment
			step.Environment = nil
		}

		stepNodes := make([]*stepnode, 0)
		for _, step := range steps {
			// build a stepnode
			stepNode := stepnode{
				ID:   int(step.GetID()),
				Name: step.GetName(),
				Step: step,
			}
			stepNodes = append(stepNodes, &stepNode)
		}

		stepPairs := getSequencePairs(steps, [][]*library.Step{})

		stepEdges := make([]*edge, 0)
		for _, stepPair := range stepPairs {
			if len(stepPair) == 2 {
				// build a stepedge
				sourceID, sourceName := int(stepPair[0].GetID()), stepPair[0].GetName()
				destinationID, destinationName := int(stepPair[1].GetID()), stepPair[1].GetName()

				stepEdge := edge{
					SourceID:        sourceID,
					SourceName:      sourceName,
					DestinationID:   destinationID,
					DestinationName: destinationName,
				}
				stepEdges = append(stepEdges, &stepEdge)
			}
		}

		// subgraph ID is the front step ID
		subgraphID := int(steps[0].GetID())
		subgraph := subgraph{
			ID:        subgraphID,
			Name:      stage.Name,
			StepNodes: stepNodes,
			StepEdges: stepEdges,
			Stage:     stage,
		}
		subgraphs[subgraphID] = &subgraph

		stageNode := stagenode{
			ID:    subgraphID,
			Name:  stage.Name,
			Steps: steps,
		}
		stageNodes = append(stageNodes, &stageNode)
	}

	// create edges from nodes
	//   an edge is as a relationship between two nodes
	//   that is defined by the 'needs' tag
	stageEdges := []*edge{}
	// loop over nodes
	for _, destinationSubgraph := range subgraphs {
		// compare all nodes against all nodes
		for _, sourceSubgraph := range subgraphs {
			// dont compare the same node
			if destinationSubgraph.ID != sourceSubgraph.ID {
				// check destination node needs
				for _, need := range (*destinationSubgraph.Stage).Needs {
					// check if destination needs source stage
					if sourceSubgraph.Stage.Name == need {
						stageEdge := edge{
							SourceID:        sourceSubgraph.ID,
							SourceName:      sourceSubgraph.Name,
							DestinationID:   destinationSubgraph.ID,
							DestinationName: destinationSubgraph.Name,
						}
						stageEdges = append(stageEdges, &stageEdge)
					}
				}
			}
		}
	}

	// construct the response
	graph := graph{
		Subgraphs:  subgraphs,
		StageNodes: stageNodes,
		StageEdges: stageEdges,
	}

	c.JSON(http.StatusOK, graph)
}

func getSequencePairs(slice []*library.Step, pairs [][]*library.Step) [][]*library.Step {
	if len(slice) > 1 {
		pair := []*library.Step{slice[0], slice[1]}
		return getSequencePairs(slice[1:], append(pairs, pair))
	}
	return pairs
}
