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
	BuildID int64         `json:"build_id"`
	Nodes   map[int]*node `json:"nodes"`
	Edges   []*edge       `json:"edges"`
}

// node represents a pipeline stage and its relevant steps.
type node struct {
	Name    string           `json:"name"`
	Stage   *pipeline.Stage  `json:"stage"`
	Steps   []*library.Step  `json:"steps"`
	ID      int              `json:"id"`
	Status  string           `json:"status"`
	Service *library.Service `json:"service"`
}

// an edge is a relationship between nodes, defined by the 'needs' tag.
type edge struct {
	Source      int    `json:"source"`
	Destination int    `json:"destination"`
	Status      string `json:"status"`
	Service     bool   `json:"service"`
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
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  u.GetName(),
	}).Infof("getting constructing graph for build %s", entry)

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

	// retrieve the services for the build from the service table
	services := []*library.Service{}
	page = 1
	perPage = 100
	for page > 0 {
		// retrieve build services (per page) from the database
		servicesPart, _, err := database.FromContext(c).ListServicesForBuild(ctx, b, nil, page, perPage)
		if err != nil {
			retErr := fmt.Errorf("unable to retrieve services for build %s: %w", entry, err)
			util.HandleError(c, http.StatusNotFound, retErr)
			return
		}

		// add page of services to list services
		services = append(services, servicesPart...)

		// assume no more pages exist if under 100 results are returned
		//
		// nolint: gomnd // ignore magic number
		if len(servicesPart) < 100 {
			page = 0
		} else {
			page++
		}
	}

	logrus.Info("retrieving pipeline configuration file")

	baseErr := "unable to generate build graph"

	// send API call to capture the pipeline configuration file
	// todo: pipeline table
	config, err := scm.FromContext(c).ConfigBackoff(ctx, u, r, b.GetCommit())
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
			files, err = scm.FromContext(c).Changeset(ctx, u, r, b.GetCommit())
			if err != nil {
				retErr := fmt.Errorf("%s: failed to get changeset for %s: %v", baseErr, r.GetFullName(), err)
				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}
		}
	}

	// todo: get pipeline from db instead?

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

	logrus.Info("generating build graph")

	type stg struct {
		steps []*library.Step
		// used for tracking stage status
		success int
		running int
		failure int
	}
	// group library steps by stage name
	stages := map[string]*stg{}
	for _, _step := range steps {
		name := _step.GetStage()
		if len(name) == 0 {
			name = _step.GetName()
		}
		if _, ok := stages[name]; !ok {
			stages[name] = &stg{
				steps:   []*library.Step{},
				success: 0,
				running: 0,
				failure: 0,
			}
		}
		switch _step.GetStatus() {
		case constants.StatusRunning:
			stages[name].running++
		case constants.StatusSuccess:
			stages[name].success++
		case constants.StatusFailure:
			// check if ruleset has 'continue' ?
			stages[name].failure++
		default:
		}
		stages[name].steps = append(stages[name].steps, _step)
	}

	// create nodes from pipeline stages
	nodes := make(map[int]*node)

	// create edges from nodes
	//   an edge is as a relationship between two nodes
	//   that is defined by the 'needs' tag
	edges := []*edge{}

	for _, service := range services {
		// scrub the environment
		nodeID := len(nodes)
		node := node{
			Name:    service.GetName(),
			ID:      nodeID,
			Status:  service.GetStatus(),
			Service: service,
		}
		nodes[nodeID] = &node

		// group services using invisible edges
		if nodeID > 0 {
			edge := &edge{
				Source:      nodeID - 1,
				Destination: nodeID,
				Status:      service.GetStatus(),
				Service:     true,
			}
			edges = append(edges, edge)
		}
	}

	for _, stage := range p.Stages {
		for _, step := range stage.Steps {
			// scrub the environment
			step.Environment = nil
		}

		// determine the "status" for a stage based on the steps within it.
		// this could potentially get complicated with ruleset logic (continue/detach)
		status := constants.StatusPending
		if stages[stage.Name].running > 0 {
			status = constants.StatusRunning
		} else if stages[stage.Name].failure > 0 {
			status = constants.StatusFailure
		} else if stages[stage.Name].success > 0 {
			status = constants.StatusSuccess
		}

		nodeID := len(nodes)

		// override the id for built-in nodes
		// this allows for better layout control
		if stage.Name == "init" {
			nodeID = -3
		}
		if stage.Name == "clone" {
			nodeID = -2
		}

		node := node{
			Name:   stage.Name,
			Stage:  stage,
			Steps:  stages[stage.Name].steps,
			ID:     nodeID,
			Status: status,
		}
		nodes[nodeID] = &node
	}

	// no stages so use steps
	if len(p.Stages) == 0 {
		for _, step := range p.Steps {
			// scrub the environment
			step.Environment = nil

			stage := &pipeline.Stage{
				Name: step.Name,
			}

			// determine the "status" for a stage based on the steps within it.
			// this could potentially get complicated with ruleset logic (continue/detach)
			status := constants.StatusPending
			if stages[stage.Name].running > 0 {
				status = constants.StatusRunning
			} else if stages[stage.Name].failure > 0 {
				status = constants.StatusFailure
			} else if stages[stage.Name].success > 0 {
				status = constants.StatusSuccess
			}

			nodeID := len(nodes)

			// override the id for built-in nodes
			// this allows for better layout control
			// if stage.Name == "init" {
			// 	nodeID = -3
			// }
			// if stage.Name == "clone" {
			// 	nodeID = -2
			// }

			node := node{
				Name:   stage.Name,
				Stage:  stage,
				Steps:  stages[stage.Name].steps,
				ID:     nodeID,
				Status: status,
			}
			nodes[nodeID] = &node
		}
	}

	// loop over nodes
	for _, destinationNode := range nodes {
		if destinationNode.Stage == nil {
			continue
		}
		// compare all nodes against all nodes
		for _, sourceNode := range nodes {
			if sourceNode.Stage == nil {
				continue
			}
			if sourceNode.ID < 0 && destinationNode.ID < 0 && sourceNode.ID < destinationNode.ID && destinationNode.ID-sourceNode.ID == 1 {
				edge := &edge{
					Source:      sourceNode.ID,
					Destination: destinationNode.ID,
					Status:      sourceNode.Status,
				}
				edges = append(edges, edge)
			}

			// skip normal processing for built-in nodes
			if destinationNode.ID < 0 || sourceNode.ID < 0 {
				continue
			}

			// dont compare the same node
			if destinationNode.ID != sourceNode.ID {
				if len((*destinationNode.Stage).Needs) > 0 {
					// check destination node needs
					for _, need := range (*destinationNode.Stage).Needs {
						// check if destination needs source stage
						if sourceNode.Stage.Name == need && need != "clone" {
							edge := &edge{
								Source:      sourceNode.ID,
								Destination: destinationNode.ID,
								Status:      sourceNode.Status,
							}
							edges = append(edges, edge)
						}
					}
				} else {
					edge := &edge{
						Source:      sourceNode.ID,
						Destination: sourceNode.ID + 1,
						Status:      sourceNode.Status,
					}
					edges = append(edges, edge)
				}
			}
		}
	}

	// for loop over edges, and collapse same parent edge
	// todo: move this check above the processing ?
	if len(nodes) > 5000 {
		c.JSON(http.StatusInternalServerError, "too many nodes on this graph")
	}
	if len(edges) > 5000 {
		c.JSON(http.StatusInternalServerError, "too many edges on this graph")
	}

	// construct the response
	graph := graph{
		BuildID: b.GetID(),
		Nodes:   nodes,
		Edges:   edges,
	}

	// todo: cli to generate digraph? output in format that can be used in other visualizers?
	c.JSON(http.StatusOK, graph)
}
