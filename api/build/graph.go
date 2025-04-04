// SPDX-License-Identifier: Apache-2.0

package build

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/internal"
	"github.com/go-vela/server/router/middleware/build"
	"github.com/go-vela/server/router/middleware/repo"
	"github.com/go-vela/server/router/middleware/user"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
)

// Graph contains nodes, and relationships between nodes, or edges.
//
//	a node is a pipeline stage and its relevant steps.
//	an edge is a relationship between nodes, defined by the 'needs' tag.
//
// swagger:model Graph
type Graph struct {
	BuildID     int64         `json:"build_id"`
	BuildNumber int64         `json:"build_number"`
	Org         string        `json:"org"`
	Repo        string        `json:"repo"`
	Nodes       map[int]*node `json:"nodes"`
	Edges       []*edge       `json:"edges"`
}

// node represents a pipeline stage and its relevant steps.
type node struct {
	ID      int    `json:"id"`
	Cluster int    `json:"cluster"`
	Name    string `json:"name"`

	// vela metadata
	Status     string        `json:"status"`
	StartedAt  int           `json:"started_at"`
	FinishedAt int           `json:"finished_at"`
	Steps      []*types.Step `json:"steps"`

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

// stg represents a stage's steps and some metadata for producing node/edge information.
type stg struct {
	steps []*types.Step
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
	// clusters determine graph orientation.
	BuiltInCluster       = 2
	PipelineCluster      = 1
	ServiceCluster       = 0
	GraphComplexityLimit = 1000 // arbitrary value to limit render complexity.
)

// swagger:operation GET /api/v1/repos/{org}/{repo}/builds/{build}/graph builds GetBuildGraph
//
// Get directed a-cyclical graph for a build
//
// ---
// produces:
// - application/json
// parameters:
// - in: path
//   name: org
//   description: Name of the organization
//   required: true
//   type: string
// - in: path
//   name: repo
//   description: Name of the repository
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
//     type: json
//     schema:
//       "$ref": "#/definitions/Graph"
//   '400':
//     description: Invalid request payload or path
//     schema:
//       "$ref": "#/definitions/Error"
//   '401':
//     description: Unauthorized
//     schema:
//       "$ref": "#/definitions/Error"
//   '404':
//     description: Not found
//     schema:
//       "$ref": "#/definitions/Error"
//   '500':
//     description: Unexpected server error
//     schema:
//       "$ref": "#/definitions/Error"

// GetBuildGraph represents the API handler to get a
// directed a-cyclical graph for a build.
//
//nolint:funlen,goconst,gocyclo // ignore function length and constants
func GetBuildGraph(c *gin.Context) {
	// capture middleware values
	m := c.MustGet("metadata").(*internal.Metadata)
	l := c.MustGet("logger").(*logrus.Entry)
	b := build.Retrieve(c)
	r := repo.Retrieve(c)
	u := user.Retrieve(c)
	ctx := c.Request.Context()

	entry := fmt.Sprintf("%s/%d", r.GetFullName(), b.GetNumber())

	baseErr := "unable to retrieve graph"

	l.Debugf("constructing graph for build %s and retrieving pipeline configuration", entry)

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
			files, err = scm.FromContext(c).Changeset(ctx, r, b.GetCommit())
			if err != nil {
				retErr := fmt.Errorf("%s: failed to get changeset for %s: %w", baseErr, r.GetFullName(), err)

				util.HandleError(c, http.StatusInternalServerError, retErr)

				return
			}
		}
	}

	l.Debug("compiling pipeline configuration")

	// parse and compile the pipeline configuration file
	p, _, err := compiler.FromContext(c).
		Duplicate().
		WithBuild(b).
		WithFiles(files).
		WithCommit(b.GetCommit()).
		WithMetadata(m).
		WithRepo(r).
		WithUser(u).
		Compile(ctx, config)
	if err != nil {
		// format the error message with extra information
		err = fmt.Errorf("unable to compile pipeline configuration for %s: %w", r.GetFullName(), err)

		l.Error(err.Error())

		retErr := fmt.Errorf("%s: %w", baseErr, err)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	if p == nil {
		retErr := fmt.Errorf("unable to compile pipeline configuration for %s: pipeline is nil", r.GetFullName())

		l.Error(retErr)

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// skip the build if only the init or clone steps are found
	skip := SkipEmptyBuild(p)
	if skip != "" {
		c.JSON(http.StatusOK, skip)
		return
	}

	// retrieve the steps for the build from the step table
	steps := []*types.Step{}
	page := 1
	perPage := 100

	// only fetch steps when necessary
	if len(p.Stages) > 0 || len(p.Steps) > 0 {
		for page > 0 {
			// retrieve build steps (per page) from the database
			stepsPart, err := database.FromContext(c).ListStepsForBuild(ctx, b, nil, page, perPage)
			if err != nil {
				retErr := fmt.Errorf("unable to retrieve steps for build %s: %w", entry, err)

				util.HandleError(c, http.StatusNotFound, retErr)

				return
			}

			// add page of steps to list steps
			steps = append(steps, stepsPart...)

			// assume no more pages exist if under 100 results are returned
			if len(stepsPart) < perPage {
				page = 0
			} else {
				page++
			}
		}
	}

	if len(steps) == 0 {
		retErr := fmt.Errorf("no steps found for build %s", entry)

		util.HandleError(c, http.StatusNotFound, retErr)

		return
	}

	// retrieve the services for the build from the service table
	services := []*types.Service{}
	page = 1
	perPage = 100

	// only fetch services when necessary
	if len(p.Services) > 0 {
		for page > 0 {
			// retrieve build services (per page) from the database
			servicesPart, err := database.FromContext(c).ListServicesForBuild(ctx, b, nil, page, perPage)
			if err != nil {
				retErr := fmt.Errorf("unable to retrieve services for build %s: %w", entry, err)

				util.HandleError(c, http.StatusNotFound, retErr)

				return
			}

			// add page of services to list services
			services = append(services, servicesPart...)

			// assume no more pages exist if under 100 results are returned
			if len(servicesPart) < perPage {
				page = 0
			} else {
				page++
			}
		}
	}

	// this is a simple check
	// but it will save on processing a massive build that should not be rendered
	complexity := len(steps) + len(p.Stages) + len(services)
	if complexity > GraphComplexityLimit {
		retErr := fmt.Errorf("build is too complex, too many resources")

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	l.Debug("generating build graph")

	// create nodes from pipeline stages
	nodes := make(map[int]*node)

	// create edges from nodes
	//   an edge is as a relationship between two nodes
	//   that is defined by the 'needs' tag
	edges := []*edge{}

	// initialize a map for grouping steps by stage name
	//   and tracking stage information
	stageMap := map[string]*stg{}

	// build a map for step_id to pipeline step
	stepMap := map[int32]*pipeline.Container{}

	for _, pStep := range p.Steps {
		stepMap[pStep.Number] = pStep
	}

	for _, pStage := range p.Stages {
		for _, pStep := range pStage.Steps {
			stepMap[pStep.Number] = pStep
		}
	}

	for _, step := range steps {
		if step == nil {
			continue
		}

		name := step.GetStage()
		if len(name) == 0 {
			name = step.GetName()
		}

		// initialize a stage tracker
		if _, ok := stageMap[name]; !ok {
			stageMap[name] = &stg{
				steps: []*types.Step{},
			}
		}

		// retrieve the stage to update
		s, ok := stageMap[name]
		if !ok {
			continue
		}

		// count each step status
		switch step.GetStatus() {
		case constants.StatusRunning:
			s.running++
		case constants.StatusSuccess:
			s.success++
		case constants.StatusFailure:
			stp, ok := stepMap[step.GetNumber()]
			if ok && stp.Ruleset.Continue {
				s.success++
			} else {
				s.failure++
			}
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
		if service == nil {
			continue
		}

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
		if stage == nil {
			continue
		}

		// skip steps/stages that were not present in the build
		// this fixes the scenario where mutable templates are updated
		s, ok := stageMap[stage.Name]
		if !ok {
			continue
		}

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

		node := nodeFromStage(nodeID, cluster, stage, s)
		nodes[nodeID] = node
	}

	// create single-step stages when no stages exist
	if len(p.Stages) == 0 {
		// sort steps by number
		sort.Slice(steps, func(i, j int) bool {
			return steps[i].GetNumber() < steps[j].GetNumber()
		})

		for _, step := range steps {
			// mock stage for edge creation
			stage := &pipeline.Stage{
				Name:  step.GetName(),
				Needs: []string{},
			}

			s, ok := stageMap[stage.Name]
			if !ok {
				continue
			}

			// create the node
			nodeID := len(nodes)

			// no built-in step separation for graphs without stages
			cluster := PipelineCluster

			node := nodeFromStage(nodeID, cluster, stage, s)
			nodes[nodeID] = node
		}
	}

	// loop over all nodes and create edges based on 'needs'
	for _, destinationNode := range nodes {
		if destinationNode == nil {
			continue
		}

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

			if destinationNode.Stage == nil {
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
		retErr := fmt.Errorf("build is too complex, too many resources")

		util.HandleError(c, http.StatusInternalServerError, retErr)

		return
	}

	// construct the response
	graph := Graph{
		BuildID:     b.GetID(),
		BuildNumber: b.GetNumber(),
		Org:         r.GetOrg(),
		Repo:        r.GetName(),
		Nodes:       nodes,
		Edges:       edges,
	}

	c.JSON(http.StatusOK, graph)
}

// nodeFromStage returns a new node from a stage.
func nodeFromStage(nodeID, cluster int, stage *pipeline.Stage, s *stg) *node {
	return &node{
		ID:         nodeID,
		Cluster:    cluster,
		Name:       stage.Name,
		Stage:      stage,
		Steps:      s.steps,
		Status:     s.GetOverallStatus(),
		StartedAt:  s.startedAt,
		FinishedAt: s.finishedAt,
	}
}

// nodeFromService returns a new node from a service.
func nodeFromService(nodeID int, service *types.Service) *node {
	return &node{
		ID:         nodeID,
		Cluster:    ServiceCluster,
		Name:       service.GetName(),
		Stage:      nil,
		Steps:      []*types.Step{},
		Status:     service.GetStatus(),
		StartedAt:  int(service.GetStarted()),
		FinishedAt: int(service.GetFinished()),
	}
}

// GetOverallStatus determines the "status" for a stage based on the steps within it.
// this could potentially get complicated with ruleset logic (continue/detach).
func (s *stg) GetOverallStatus() string {
	if s.running > 0 {
		return constants.StatusRunning
	}

	if s.failure > 0 {
		return constants.StatusFailure
	}

	if s.errored > 0 {
		return constants.StatusError
	}

	if s.killed >= len(s.steps) {
		return constants.StatusKilled
	}

	if s.skipped > 0 {
		return constants.StatusSkipped
	}

	if s.success > 0 {
		return constants.StatusSuccess
	}

	if s.canceled > 0 {
		return constants.StatusCanceled
	}

	return constants.StatusPending
}
