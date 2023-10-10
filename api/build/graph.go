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
	Cluster int    `json:"cluster"`
	ID      int    `json:"id"`
	Name    string `json:"name"`

	// vela metadata
	Status   string          `json:"status"`
	Duration int             `json:"duration"`
	Steps    []*library.Step `json:"steps"`

	// unexported data used for building edges
	Stage *pipeline.Stage `json:"-"`
}

// edge represents a relationship between nodes, primarily defined by a stage 'needs' tag.
type edge struct {
	Cluster     int `json:"cluster"`
	Source      int `json:"source"`
	Destination int `json:"destination"`

	// vela metadata
	Status string `json:"status"`
}

// stg represents a stage's steps and some metadata for producing node/edge information
type stg struct {
	steps []*library.Step
	// used for tracking stage status
	success    int
	running    int
	failure    int
	killed     int
	startedAt  int
	finishedAt int
}

const (
	// clusters determine graph orientation
	BuiltInCluster  = 2
	PipelineCluster = 1
	ServiceCluster  = 0
)

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

	baseErr := "unable to generate build graph"

	logrus.Info("retrieving pipeline configuration")
	var config []byte

	lp, err := database.FromContext(c).GetPipelineForRepo(ctx, b.GetCommit(), r)
	if err != nil { // assume the pipeline doesn't exist in the database yet (before pipeline support was added)
		// send API call to capture the pipeline configuration file
		config, err = scm.FromContext(c).ConfigBackoff(ctx, u, r, b.GetCommit())
		if err != nil {
			retErr := fmt.Errorf("unable to get pipeline configuration for %s: %w", r.GetFullName(), err)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}
	} else {
		config = lp.GetData()
	}

	if config == nil {
		retErr := fmt.Errorf("unable to get pipeline configuration for %s: config is nil", r.GetFullName())

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

	logrus.Info("compiling pipeline configuration")

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

	stages := p.Stages

	// create nodes from pipeline stages
	nodes := make(map[int]*node)

	// create edges from nodes
	//   an edge is as a relationship between two nodes
	//   that is defined by the 'needs' tag
	edges := []*edge{}

	// initialize a map for grouping steps by stage name
	//   and tracking stage information
	stageMap := map[string]*stg{}
	for _, _step := range steps {
		name := _step.GetStage()
		if len(name) == 0 {
			name = _step.GetName()
		}
		if _, ok := stageMap[name]; !ok {
			stageMap[name] = &stg{
				steps:      []*library.Step{},
				success:    0,
				running:    0,
				failure:    0,
				killed:     0,
				startedAt:  0,
				finishedAt: 0,
			}
		}
		stageMap[name].updateStgTracker(_step)
		stageMap[name].steps = append(stageMap[name].steps, _step)
	}

	for _, service := range services {
		nodeID := len(nodes)
		node := node{
			// set service cluster
			Cluster: ServiceCluster,
			ID:      nodeID,

			Name:   service.GetName(),
			Status: service.GetStatus(),
			// todo: runtime
			Duration: int(service.GetFinished() - service.GetStarted()),
			// todo: is this right?
			// Steps: []string{},
		}
		nodes[nodeID] = &node

		// group services using invisible edges
		if nodeID > 0 {
			edge := &edge{
				// set service cluster
				Cluster: ServiceCluster,
				// link them together for visual effect
				Source:      nodeID - 1,
				Destination: nodeID,

				Status: service.GetStatus(),
			}
			edges = append(edges, edge)
		}
	}

	for _, stage := range stages {
		for _, step := range stage.Steps {
			// scrub the environment
			step.Environment = nil
		}

		nodeID := len(nodes)

		// determine the "status" for a stage based on the steps within it.
		// this could potentially get complicated with ruleset logic (continue/detach)
		status := stageStatus(stageMap[stage.Name])

		node := node{
			Cluster: PipelineCluster,
			ID:      nodeID,

			Name:     stage.Name,
			Status:   status,
			Duration: int(stageMap[stage.Name].finishedAt - stageMap[stage.Name].startedAt),
			Steps:    stageMap[stage.Name].steps,

			// used for edge creation
			Stage: stage,
		}

		// override the id for built-in nodes
		// this allows for better layout control
		if stage.Name == "init" {
			node.Cluster = BuiltInCluster
		}
		if stage.Name == "clone" {
			node.Cluster = BuiltInCluster
		}

		nodes[nodeID] = &node
	}

	// no stages so use steps
	if len(p.Stages) == 0 {
		for _, step := range p.Steps {
			// scrub the environment
			step.Environment = nil

			// mock stage for edge creation
			stage := &pipeline.Stage{
				Name: step.Name,
			}

			// determine the "status" for a stage based on the steps within it.
			// this could potentially get complicated with ruleset logic (continue/detach)
			status := stageStatus(stageMap[stage.Name])

			nodeID := len(nodes)

			node := node{
				Cluster: PipelineCluster,
				ID:      nodeID,

				Name:     stage.Name,
				Status:   status,
				Duration: int(stageMap[stage.Name].finishedAt - stageMap[stage.Name].startedAt),
				Steps:    stageMap[stage.Name].steps,

				// used for edge creation
				Stage: stage,
			}
			nodes[nodeID] = &node
		}
	}
	// done building nodes

	// loop over nodes and create edges based on 'needs'
	for _, destinationNode := range nodes {
		// if theres no stage, skip because the edge is already created?
		if destinationNode.Stage == nil {
			continue
		}

		// compare all nodes against all nodes
		for _, sourceNode := range nodes {
			if sourceNode.Stage == nil {
				continue
			}

			if sourceNode.Cluster == BuiltInCluster && destinationNode.Cluster == BuiltInCluster && sourceNode.ID < destinationNode.ID && destinationNode.ID-sourceNode.ID == 1 {
				edge := &edge{
					Cluster:     sourceNode.Cluster,
					Source:      sourceNode.ID,
					Destination: destinationNode.ID,
					Status:      sourceNode.Status,
				}
				edges = append(edges, edge)
			}

			// skip normal processing for built-in nodes
			if destinationNode.Cluster == BuiltInCluster || sourceNode.Cluster == BuiltInCluster {
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
								Cluster:     sourceNode.Cluster,
								Source:      sourceNode.ID,
								Destination: destinationNode.ID,
								Status:      sourceNode.Status,
							}
							edges = append(edges, edge)
						}
					}
				} else {
					edge := &edge{
						Cluster:     sourceNode.Cluster,
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

func (s *stg) updateStgTracker(step *library.Step) {
	switch step.GetStatus() {
	case constants.StatusRunning:
		s.running++
	case constants.StatusSuccess:
		s.success++
	case constants.StatusFailure:
		// check if ruleset has 'continue' ?
		s.failure++
	case constants.StatusKilled:
		s.killed++
	default:
	}

	if s.finishedAt == 0 || s.finishedAt < int(step.GetFinished()) {
		s.finishedAt = int(step.GetFinished())
	}
	if s.startedAt == 0 || s.startedAt > int(step.GetStarted()) {
		s.startedAt = int(step.GetStarted())
	}
}

func stageStatus(s *stg) string {
	status := constants.StatusPending
	if s.running > 0 {
		status = constants.StatusRunning
	} else if s.failure > 0 {
		status = constants.StatusFailure
	} else if s.success > 0 {
		status = constants.StatusSuccess
	} else if s.killed > 0 {
		status = constants.StatusKilled
	}
	return status
}
