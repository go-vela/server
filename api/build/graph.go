// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package build

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/org"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
)

// Graph contains nodes, and relationships between nodes, or edges.
//
//	a node is a pipeline stage and its relevant steps.
//	an edge is a relationship between nodes, defined by the 'needs' tag.
//
// swagger:model Graph
type Graph struct {
	BuildID int64         `json:"build_id"`
	Nodes   map[int]*node `json:"nodes"`
	Edges   []*edge       `json:"edges"`
}

// node represents a pipeline stage and its relevant steps.
type node struct {
	ID      int    `json:"id"`
	Cluster int    `json:"cluster"`
	Name    string `json:"name"`

	// vela metadata
	Status     string          `json:"status"`
	StartedAt  int             `json:"started_at"`
	FinishedAt int             `json:"finished_at"`
	Steps      []*library.Step `json:"steps"`

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
	canceled   int
	skipped    int
	errored    int
	startedAt  int
	finishedAt int
}

const (
	// clusters determine graph orientation
	BuiltInCluster       = 2
	PipelineCluster      = 1
	ServiceCluster       = 0
	GraphComplexityLimit = 1000 // arbitrary value to limit render complexity
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
//       type: json
//       items:
//         "$ref": "#/definitions/Graph"
//   '401':
//     description: Unable to retrieve graph for the build — unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Unable to retrieve graph for the build — not found
//     schema:
//       "$ref": "#/definitions/Error"
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

	// update engine logger with API metadata
	//
	// https://pkg.go.dev/github.com/sirupsen/logrus?tab=doc#Entry.WithFields
	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())
	logger := logrus.WithFields(logrus.Fields{
		"build": b.GetNumber(),
		"org":   o,
		"repo":  r.GetName(),
		"user":  u.GetName(),
	})

	baseErr := "unable to retrieve graph"

	logger.Infof("constructing graph for build %s", entry)

	// retrieve the steps for the build from the step table
	steps := []*library.Step{}
	page := 1
	perPage := 100
	for page > 0 {
		// retrieve build steps (per page) from the database
		stepsPart, stepsCount, err := database.FromContext(c).ListStepsForBuild(b, nil, page, perPage)
		if err != nil {
			retErr := fmt.Errorf("unable to retrieve steps for build %s: %w", entry, err)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		// add page of steps to list steps
		steps = append(steps, stepsPart...)

		// assume no more pages exist if under 100 results are returned
		if stepsCount < 100 {
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
		servicesPart, servicesCount, err := database.FromContext(c).ListServicesForBuild(ctx, b, nil, page, perPage)
		if err != nil {
			retErr := fmt.Errorf("unable to retrieve services for build %s: %w", entry, err)

			util.HandleError(c, http.StatusNotFound, retErr)

			return
		}

		// add page of services to list services
		services = append(services, servicesPart...)

		// assume no more pages exist if under 100 results are returned
		if servicesCount < 100 {
			page = 0
		} else {
			page++
		}
	}

	logger.Info("retrieving pipeline configuration")

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

	logger.Info("compiling pipeline configuration")

	// parse and compile the pipeline configuration file
	p, _, err := compiler.FromContext(c).
		Duplicate().
		WithBuild(b).
		WithFiles(files).
		WithCommit(b.GetCommit()).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u).
		Compile(config)
	if err != nil {
		// format the error message with extra information
		err = fmt.Errorf("unable to compile pipeline configuration for %s: %v", r.GetFullName(), err)

		logger.Error(err.Error())

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

	// this is a simple check
	// but it will save on processing a massive build that should not be rendered
	complexity := len(steps) + len(p.Stages) + len(services)
	if complexity > GraphComplexityLimit {
		c.JSON(http.StatusInternalServerError, "build is too complex, too many resources")
		return
	}

	logger.Info("generating build graph")

	// create nodes from pipeline stages
	nodes := make(map[int]*node)

	// create edges from nodes
	//   an edge is as a relationship between two nodes
	//   that is defined by the 'needs' tag
	edges := []*edge{}

	// initialize a map for grouping steps by stage name
	//   and tracking stage information
	stageMap := map[string]*stg{}
	for _, step := range steps {
		name := step.GetStage()
		if len(name) == 0 {
			name = step.GetName()
		}

		// initialize a stage tracker
		if _, ok := stageMap[name]; !ok {
			stageMap[name] = &stg{
				steps:      []*library.Step{},
				success:    0,
				running:    0,
				failure:    0,
				killed:     0,
				canceled:   0,
				skipped:    0,
				errored:    0,
				startedAt:  0,
				finishedAt: 0,
			}
		}

		// retrieve the stage to update
		s := stageMap[name]

		// count each step status
		switch step.GetStatus() {
		case constants.StatusRunning:
			s.running++
		case constants.StatusSuccess:
			s.success++
		case constants.StatusFailure:
			s.failure++
		case constants.StatusKilled:
			s.killed++
		case constants.StatusCanceled:
			s.canceled++
		case constants.StatusSkipped:
			s.skipped++
		case constants.StatusError:
			s.errored++
		default:
		}

		// keep track of the widest time window possible
		if s.finishedAt == 0 || s.finishedAt < int(step.GetFinished()) {
			s.finishedAt = int(step.GetFinished())
		}
		if s.startedAt == 0 || s.startedAt > int(step.GetStarted()) {
			s.startedAt = int(step.GetStarted())
		}

		s.steps = append(s.steps, step)
	}

	// construct services nodes separately
	// services are grouped via cluster and manually-constructed edges
	for _, service := range services {
		// create the node
		nodeID := len(nodes)
		node := nodeFromService(nodeID, service)
		nodes[nodeID] = node

		// create a sequential edge for nodes after the first
		if nodeID > 0 {
			edge := &edge{
				Cluster:     ServiceCluster,
				Source:      nodeID - 1,
				Destination: nodeID,
				Status:      service.GetStatus(),
			}
			edges = append(edges, edge)
		}
	}

	// construct pipeline stages nodes when stages exist
	for _, stage := range p.Stages {
		// scrub the environment
		for _, step := range stage.Steps {
			step.Environment = nil
		}

		// create the node
		nodeID := len(nodes)

		cluster := PipelineCluster

		// override the cluster id for built-in nodes
		// this allows for better layout control because all stages need 'clone'
		if stage.Name == "init" {
			cluster = BuiltInCluster
		}
		if stage.Name == "clone" {
			cluster = BuiltInCluster
		}

		node := nodeFromStage(nodeID, cluster, stage, stageMap[stage.Name])
		nodes[nodeID] = node
	}

	// create single-step stages when no stages exist
	if len(p.Stages) == 0 {
		for _, step := range p.Steps {
			// scrub the environment
			step.Environment = nil

			// mock stage for edge creation
			stage := &pipeline.Stage{
				Name:  step.Name,
				Needs: []string{},
			}

			// create the node
			nodeID := len(nodes)

			// no built-in step separation for graphs without stages
			cluster := PipelineCluster

			node := nodeFromStage(nodeID, cluster, stage, stageMap[stage.Name])
			nodes[nodeID] = node
		}
	}

	// loop over all nodes and create edges based on 'needs'
	for _, destinationNode := range nodes {
		// if theres no stage, skip because the edge is already created
		if destinationNode.Stage == nil {
			continue
		}

		// compare all nodes against all nodes
		for _, sourceNode := range nodes {
			if sourceNode.Stage == nil {
				continue
			}

			// create a sequential edge for built-in nodes
			if sourceNode.Cluster == BuiltInCluster &&
				destinationNode.Cluster == BuiltInCluster &&
				sourceNode.ID < destinationNode.ID &&
				destinationNode.ID-sourceNode.ID == 1 {
				edge := &edge{
					Cluster:     sourceNode.Cluster,
					Source:      sourceNode.ID,
					Destination: destinationNode.ID,
					Status:      sourceNode.Status,
				}
				edges = append(edges, edge)
			}

			// skip normal processing for built-in nodes
			if destinationNode.Cluster == BuiltInCluster ||
				sourceNode.Cluster == BuiltInCluster {
				continue
			}

			// dont compare the same node
			if destinationNode.ID == sourceNode.ID {
				continue
			}

			// use needs to create an edge
			if len((*destinationNode.Stage).Needs) > 0 {
				// check destination node needs
				for _, need := range (*destinationNode.Stage).Needs {
					// all stages need "clone"
					if need == "clone" {
						continue
					}

					// check destination needs source stage
					if sourceNode.Stage.Name != need {
						continue
					}

					// create an edge for source->destination
					edge := &edge{
						Cluster:     sourceNode.Cluster,
						Source:      sourceNode.ID,
						Destination: destinationNode.ID,
						Status:      sourceNode.Status,
					}
					edges = append(edges, edge)
				}
			} else {
				// create an edge for source->next
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

	// validate the generated graph's complexity is beneath the limit
	if len(nodes)+len(edges) > GraphComplexityLimit {
		c.JSON(http.StatusInternalServerError, "graph is too complex, too many nodes and edges")
		return
	}

	// construct the response
	graph := Graph{
		BuildID: b.GetID(),
		Nodes:   nodes,
		Edges:   edges,
	}

	c.JSON(http.StatusOK, graph)
}

// nodeFromStage returns a new node from a stage
func nodeFromStage(nodeID, cluster int, stage *pipeline.Stage, s *stg) *node {
	return &node{
		ID:         nodeID,
		Cluster:    cluster,
		Name:       stage.Name,
		Stage:      stage,
		Steps:      s.steps,
		Status:     s.GetOverallStatus(),
		StartedAt:  int(s.startedAt),
		FinishedAt: int(s.finishedAt),
	}
}

// nodeFromService returns a new node from a service
func nodeFromService(nodeID int, service *library.Service) *node {
	return &node{
		ID:         nodeID,
		Cluster:    ServiceCluster,
		Name:       service.GetName(),
		Stage:      nil,
		Steps:      []*library.Step{},
		Status:     service.GetStatus(),
		StartedAt:  int(service.GetStarted()),
		FinishedAt: int(service.GetFinished()),
	}
}

// GetOverallStatus determines the "status" for a stage based on the steps within it.
// this could potentially get complicated with ruleset logic (continue/detach)
func (s *stg) GetOverallStatus() string {
	if s.running > 0 {
		return constants.StatusRunning
	}

	if s.failure > 0 {
		return constants.StatusFailure
	}

	if s.success > 0 {
		return constants.StatusSuccess
	}

	if s.killed > 0 {
		return constants.StatusKilled
	}

	if s.skipped > 0 {
		return constants.StatusSkipped
	}

	if s.canceled > 0 {
		return constants.StatusCanceled
	}

	if s.errored > 0 {
		return constants.StatusError
	}

	return constants.StatusPending
}
