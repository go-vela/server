// SPDX-License-Identifier: Apache-2.0

package build

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-vela/server/compiler"
	"github.com/go-vela/server/database"
	"github.com/go-vela/server/queue"
	"github.com/go-vela/server/scm"
	"github.com/go-vela/server/util"
	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/sirupsen/logrus"
)

// CompileAndPublishConfig is a struct that contains information for the CompileAndPublish function.
type CompileAndPublishConfig struct {
	Build    *library.Build
	Repo     *library.Repo
	Metadata *types.Metadata
	BaseErr  string
	Source   string
	Comment  string
	Retries  int
}

// CompileAndPublish is a helper function to generate the queue items for a build. It takes a form
// as well as the database, scm, compiler, and queue services as arguments. It is used in webhook handling,
// schedule processing, and API build creation.
//
//nolint:funlen,gocyclo // ignore function length due to comments, error handling, and general complexity of function
func CompileAndPublish(
	c context.Context,
	cfg CompileAndPublishConfig,
	database database.Interface,
	scm scm.Service,
	compiler compiler.Engine,
	queue queue.Service,
) (bool, *pipeline.Build, *types.Item, error) {
	logrus.Debugf("generating queue items for build %s/%d", cfg.Repo.GetFullName(), cfg.Build.GetNumber())

	// assign variables from form for readibility
	r := cfg.Repo
	b := cfg.Build
	baseErr := cfg.BaseErr

	// send API call to capture repo owner
	logrus.Debugf("capturing owner of repository %s", cfg.Repo.GetFullName())

	u, err := database.GetUser(c, r.GetUserID())
	if err != nil {
		retErr := fmt.Errorf("%s: failed to get owner for %s: %w", baseErr, r.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return false, nil, nil, retErr
	}

	// confirm current repo owner has at least write access to repo (needed for status update later)
	_, err = scm.RepoAccess(c, u.GetName(), u.GetToken(), r.GetOrg(), r.GetName())
	if err != nil {
		retErr := fmt.Errorf("unable to publish build to queue: repository owner %s no longer has write access to repository %s", u.GetName(), r.GetFullName())
		util.HandleError(c, http.StatusUnauthorized, retErr)

		return false, nil, nil, retErr
	}

	// get pull request number from build if event is pull_request or issue_comment
	var prNum int
	if strings.EqualFold(b.GetEvent(), constants.EventPull) || strings.EqualFold(b.GetEvent(), constants.EventComment) {
		prNum, err = getPRNumberFromBuild(b)
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get pull request number for %s: %w", baseErr, r.GetFullName(), err)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return false, nil, nil, retErr
		}
	}

	// if the event is issue_comment and the issue is a pull request,
	// call SCM for more data not provided in webhook payload
	if strings.EqualFold(cfg.Source, "webhook") && strings.EqualFold(b.GetEvent(), constants.EventComment) {
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get pull request number for %s: %w", baseErr, r.GetFullName(), err)
			util.HandleError(c, http.StatusBadRequest, retErr)

			return false, nil, nil, retErr
		}

		commit, branch, baseref, headref, err := scm.GetPullRequest(c, u, r, prNum)
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get pull request info for %s: %w", baseErr, r.GetFullName(), err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			return false, nil, nil, retErr
		}

		b.SetCommit(commit)
		b.SetBranch(strings.ReplaceAll(branch, "refs/heads/", ""))
		b.SetBaseRef(baseref)
		b.SetHeadRef(headref)
	}

	// if the source is from a schedule, fetch the commit sha from schedule branch (same as build branch at this moment)
	if strings.EqualFold(cfg.Source, "schedule") {
		// send API call to capture the commit sha for the branch
		_, commit, err := scm.GetBranch(c, u, r, b.GetBranch())
		if err != nil {
			retErr := fmt.Errorf("failed to get commit for repo %s on %s branch: %w", r.GetFullName(), r.GetBranch(), err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			return false, nil, nil, retErr
		}

		b.SetCommit(commit)
	}

	// create SQL filters for querying pending and running builds for repo
	filters := map[string]interface{}{
		"status": []string{constants.StatusPending, constants.StatusRunning},
	}

	// send API call to capture the number of pending or running builds for the repo
	builds, err := database.CountBuildsForRepo(c, r, filters)
	if err != nil {
		retErr := fmt.Errorf("%s: unable to get count of builds for repo %s", baseErr, r.GetFullName())
		util.HandleError(c, http.StatusBadRequest, retErr)

		return false, nil, nil, retErr
	}

	logrus.Debugf("currently %d builds running on repo %s", builds, r.GetFullName())

	// check if the number of pending and running builds exceeds the limit for the repo
	if builds >= r.GetBuildLimit() {
		retErr := fmt.Errorf("%s: repo %s has exceeded the concurrent build limit of %d", baseErr, r.GetFullName(), r.GetBuildLimit())
		util.HandleError(c, http.StatusBadRequest, retErr)

		return false, nil, nil, retErr
	}

	// update fields in build object
	// this is necessary in case source is restart and the build is prepopulated with these values
	b.SetID(0)
	b.SetCreated(time.Now().UTC().Unix())
	b.SetEnqueued(0)
	b.SetStarted(0)
	b.SetFinished(0)
	b.SetStatus(constants.StatusPending)
	b.SetError("")
	b.SetHost("")
	b.SetRuntime("")
	b.SetDistribution("")

	// variable to store changeset files
	var files []string

	// check if the build event is not issue_comment or pull_request
	if !strings.EqualFold(b.GetEvent(), constants.EventComment) &&
		!strings.EqualFold(b.GetEvent(), constants.EventPull) &&
		!strings.EqualFold(b.GetEvent(), constants.EventDelete) {
		// send API call to capture list of files changed for the commit
		files, err = scm.Changeset(c, u, r, b.GetCommit())
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get changeset for %s: %w", baseErr, r.GetFullName(), err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			return false, nil, nil, retErr
		}
	}

	// check if the build event is a pull_request
	if strings.EqualFold(b.GetEvent(), constants.EventPull) && prNum > 0 {
		// send API call to capture list of files changed for the pull request
		files, err = scm.ChangesetPR(c, u, r, prNum)
		if err != nil {
			retErr := fmt.Errorf("%s: failed to get changeset for %s: %w", baseErr, r.GetFullName(), err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			return false, nil, nil, retErr
		}
	}

	var (
		// variable to store the raw pipeline configuration
		pipelineFile []byte
		// variable to store executable pipeline
		p *pipeline.Build
		// variable to store pipeline configuration
		pipeline *library.Pipeline
		// variable to store the pipeline type for the repository
		pipelineType = r.GetPipelineType()
		// variable to store updated repository record
		repo *library.Repo
	)

	// implement a loop to process asynchronous operations with a retry limit
	//
	// Some operations taken during the webhook workflow can lead to race conditions
	// failing to successfully process the request. This logic ensures we attempt our
	// best efforts to handle these cases gracefully.
	for i := 0; i < cfg.Retries; i++ {
		logrus.Debugf("compilation loop - attempt %d", i+1)
		// check if we're on the first iteration of the loop
		if i > 0 {
			// incrementally sleep in between retries
			time.Sleep(time.Duration(i) * time.Second)
		}

		// send database call to attempt to capture the pipeline if we already processed it before
		pipeline, err = database.GetPipelineForRepo(c, b.GetCommit(), r)
		if err != nil { // assume the pipeline doesn't exist in the database yet
			// send API call to capture the pipeline configuration file
			pipelineFile, err = scm.ConfigBackoff(c, u, r, b.GetCommit())
			if err != nil {
				retErr := fmt.Errorf("%s: unable to get pipeline configuration for %s: %w", baseErr, r.GetFullName(), err)

				util.HandleError(c, http.StatusNotFound, retErr)

				return false, nil, nil, retErr
			}
		} else {
			pipelineFile = pipeline.GetData()
		}

		// send API call to capture repo for the counter (grabbing repo again to ensure counter is correct)
		repo, err = database.GetRepoForOrg(c, r.GetOrg(), r.GetName())
		if err != nil {
			retErr := fmt.Errorf("%s: unable to get repo %s: %w", baseErr, r.GetFullName(), err)

			// check if the retry limit has been exceeded
			if i < cfg.Retries-1 {
				logrus.WithError(retErr).Warningf("retrying #%d", i+1)

				// continue to the next iteration of the loop
				continue
			}

			util.HandleError(c, http.StatusBadRequest, retErr)

			return false, nil, nil, retErr
		}

		// update DB record of repo (repo) with any changes captured from webhook payload (r)
		repo.SetTopics(r.GetTopics())
		repo.SetBranch(r.GetBranch())

		// update the build numbers based off repo counter
		inc := repo.GetCounter() + 1
		repo.SetCounter(inc)
		b.SetNumber(inc)

		// populate the build link if a web address is provided
		if len(cfg.Metadata.Vela.WebAddress) > 0 {
			b.SetLink(
				fmt.Sprintf("%s/%s/%d", cfg.Metadata.Vela.WebAddress, repo.GetFullName(), b.GetNumber()),
			)
		}

		// ensure we use the expected pipeline type when compiling
		//
		// The pipeline type for a repo can change at any time which can break compiling
		// existing pipelines in the system for that repo. To account for this, we update
		// the repo pipeline type to match what was defined for the existing pipeline
		// before compiling. After we're done compiling, we reset the pipeline type.
		if len(pipeline.GetType()) > 0 {
			repo.SetPipelineType(pipeline.GetType())
		}

		var compiled *library.Pipeline
		// parse and compile the pipeline configuration file
		p, compiled, err = compiler.
			Duplicate().
			WithBuild(b).
			WithComment(cfg.Comment).
			WithCommit(b.GetCommit()).
			WithFiles(files).
			WithMetadata(cfg.Metadata).
			WithRepo(repo).
			WithUser(u).
			Compile(pipelineFile)
		if err != nil {
			// format the error message with extra information
			err = fmt.Errorf("unable to compile pipeline configuration for %s: %w", repo.GetFullName(), err)

			// log the error for traceability
			logrus.Error(err.Error())

			retErr := fmt.Errorf("%s: %w", baseErr, err)
			util.HandleError(c, http.StatusInternalServerError, retErr)

			return false, nil, nil, retErr
		}

		// reset the pipeline type for the repo
		//
		// The pipeline type for a repo can change at any time which can break compiling
		// existing pipelines in the system for that repo. To account for this, we update
		// the repo pipeline type to match what was defined for the existing pipeline
		// before compiling. After we're done compiling, we reset the pipeline type.
		repo.SetPipelineType(pipelineType)

		// skip the build if pipeline compiled to only the init and clone steps
		skip := SkipEmptyBuild(p)
		if skip != "" {
			// set build to successful status
			b.SetStatus(constants.StatusSkipped)

			// send API call to set the status on the commit
			err = scm.Status(c, u, b, repo.GetOrg(), repo.GetName())
			if err != nil {
				logrus.Errorf("unable to set commit status for %s/%d: %v", repo.GetFullName(), b.GetNumber(), err)
			}

			return false,
				nil,
				&types.Item{
					Build: b,
				},
				nil
		}

		// check if the pipeline did not already exist in the database
		if pipeline == nil {
			pipeline = compiled
			pipeline.SetRepoID(repo.GetID())
			pipeline.SetCommit(b.GetCommit())
			pipeline.SetRef(b.GetRef())

			// send API call to create the pipeline
			pipeline, err = database.CreatePipeline(c, pipeline)
			if err != nil {
				retErr := fmt.Errorf("%s: failed to create pipeline for %s: %w", baseErr, repo.GetFullName(), err)

				// check if the retry limit has been exceeded
				if i < cfg.Retries-1 {
					logrus.WithError(retErr).Warningf("retrying #%d", i+1)

					// continue to the next iteration of the loop
					continue
				}

				util.HandleError(c, http.StatusBadRequest, retErr)

				return false, nil, nil, retErr
			}
		}

		b.SetPipelineID(pipeline.GetID())

		// create the objects from the pipeline in the database
		// TODO:
		// - if a build gets created and something else fails midway,
		//   the next loop will attempt to create the same build,
		//   using the same Number and thus create a constraint
		//   conflict; consider deleting the partially created
		//   build object in the database
		err = PlanBuild(c, database, scm, p, b, repo)
		if err != nil {
			retErr := fmt.Errorf("%s: %w", baseErr, err)

			// check if the retry limit has been exceeded
			if i < cfg.Retries-1 {
				logrus.WithError(retErr).Warningf("retrying #%d", i+1)

				// reset fields set by cleanBuild for retry
				b.SetError("")
				b.SetStatus(constants.StatusPending)
				b.SetFinished(0)

				// continue to the next iteration of the loop
				continue
			}

			util.HandleError(c, http.StatusInternalServerError, retErr)

			return false, nil, nil, retErr
		}

		// break the loop because everything was successful
		break
	} // end of retry loop

	// send API call to update repo for ensuring counter is incremented
	repo, err = database.UpdateRepo(c, repo)
	if err != nil {
		retErr := fmt.Errorf("%s: failed to update repo %s: %w", baseErr, repo.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return false, nil, nil, retErr
	}

	// return error if pipeline didn't get populated
	if p == nil {
		retErr := fmt.Errorf("%s: failed to set pipeline for %s: %w", baseErr, repo.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return false, nil, nil, retErr
	}

	// return error if build didn't get populated
	if b == nil {
		retErr := fmt.Errorf("%s: failed to set build for %s: %w", baseErr, repo.GetFullName(), err)
		util.HandleError(c, http.StatusBadRequest, retErr)

		return false, nil, nil, retErr
	}

	// send API call to capture the triggered build
	b, err = database.GetBuildForRepo(c, repo, b.GetNumber())
	if err != nil {
		retErr := fmt.Errorf("%s: failed to get new build %s/%d: %w", baseErr, repo.GetFullName(), b.GetNumber(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return false, nil, nil, retErr
	}

	// determine queue route
	route, err := queue.Route(&p.Worker)
	if err != nil {
		retErr := fmt.Errorf("unable to set route for build %d for %s: %w", b.GetNumber(), r.GetFullName(), err)

		// error out the build
		CleanBuild(c, database, b, nil, nil, retErr)

		return false, nil, nil, retErr
	}

	// temporarily set host to the route before it gets picked up by a worker
	b.SetHost(route)

	// publish the pipeline.Build to the build_executables table to be requested by a worker
	err = PublishBuildExecutable(c, database, p, b)
	if err != nil {
		retErr := fmt.Errorf("unable to publish build executable for %s/%d: %w", repo.GetFullName(), b.GetNumber(), err)
		util.HandleError(c, http.StatusInternalServerError, retErr)

		return false, nil, nil, retErr
	}

	return true, p, types.ToItem(b, repo, u), nil
}

// getPRNumberFromBuild is a helper function to
// extract the pull request number from a Build.
func getPRNumberFromBuild(b *library.Build) (int, error) {
	// parse out pull request number from base ref
	//
	// pattern: refs/pull/1/head
	var parts []string
	if strings.HasPrefix(b.GetRef(), "refs/pull/") {
		parts = strings.Split(b.GetRef(), "/")
	}

	// just being safe to avoid out of range index errors
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid ref: %s", b.GetRef())
	}

	// return the results of converting number to string
	return strconv.Atoi(parts[2])
}
