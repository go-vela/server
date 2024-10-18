// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"

	"github.com/go-vela/server/api/types/actions"
	"github.com/go-vela/server/constants"
)

// Events is the API representation of the various events that generate a
// webhook from the SCM.
type Events struct {
	Push        *actions.Push     `json:"push"`
	PullRequest *actions.Pull     `json:"pull_request"`
	Deployment  *actions.Deploy   `json:"deployment"`
	Comment     *actions.Comment  `json:"comment"`
	Schedule    *actions.Schedule `json:"schedule"`
}

// NewEventsFromMask is an instatiation function for the Events type that
// takes in an event mask integer value and populates the nested Events struct.
func NewEventsFromMask(mask int64) *Events {
	pushActions := new(actions.Push).FromMask(mask)
	pullActions := new(actions.Pull).FromMask(mask)
	deployActions := new(actions.Deploy).FromMask(mask)
	commentActions := new(actions.Comment).FromMask(mask)
	scheduleActions := new(actions.Schedule).FromMask(mask)

	e := new(Events)

	e.SetPush(pushActions)
	e.SetPullRequest(pullActions)
	e.SetDeployment(deployActions)
	e.SetComment(commentActions)
	e.SetSchedule(scheduleActions)

	return e
}

// NewEventsFromSlice is an instantiation function for the Events type that
// takes in a slice of event strings and populates the nested Events struct.
func NewEventsFromSlice(events []string) (*Events, error) {
	mask := int64(0)

	// iterate through all events provided
	for _, event := range events {
		switch event {
		// push actions
		case constants.EventPush, constants.EventPush + ":branch":
			mask = mask | constants.AllowPushBranch
		case constants.EventTag, constants.EventPush + ":" + constants.EventTag:
			mask = mask | constants.AllowPushTag
		case constants.EventDelete + ":" + constants.ActionBranch:
			mask = mask | constants.AllowPushDeleteBranch
		case constants.EventDelete + ":" + constants.ActionTag:
			mask = mask | constants.AllowPushDeleteTag
		case constants.EventDelete:
			mask = mask | constants.AllowPushDeleteBranch | constants.AllowPushDeleteTag

		// pull_request actions
		case constants.EventPull, constants.EventPullAlternate:
			mask = mask | constants.AllowPullOpen | constants.AllowPullSync | constants.AllowPullReopen
		case constants.EventPull + ":" + constants.ActionOpened:
			mask = mask | constants.AllowPullOpen
		case constants.EventPull + ":" + constants.ActionEdited:
			mask = mask | constants.AllowPullEdit
		case constants.EventPull + ":" + constants.ActionSynchronize:
			mask = mask | constants.AllowPullSync
		case constants.EventPull + ":" + constants.ActionReopened:
			mask = mask | constants.AllowPullReopen
		case constants.EventPull + ":" + constants.ActionLabeled:
			mask = mask | constants.AllowPullLabel
		case constants.EventPull + ":" + constants.ActionUnlabeled:
			mask = mask | constants.AllowPullUnlabel

		// deployment actions
		case constants.EventDeploy, constants.EventDeployAlternate, constants.EventDeploy + ":" + constants.ActionCreated:
			mask = mask | constants.AllowDeployCreate

		// comment actions
		case constants.EventComment:
			mask = mask | constants.AllowCommentCreate | constants.AllowCommentEdit
		case constants.EventComment + ":" + constants.ActionCreated:
			mask = mask | constants.AllowCommentCreate
		case constants.EventComment + ":" + constants.ActionEdited:
			mask = mask | constants.AllowCommentEdit

		// schedule actions
		case constants.EventSchedule, constants.EventSchedule + ":" + constants.ActionRun:
			mask = mask | constants.AllowSchedule

		default:
			return nil, fmt.Errorf("invalid event provided: %s", event)
		}
	}

	return NewEventsFromMask(mask), nil
}

// Allowed determines whether or not an event + action is allowed based on whether
// its event:action is set to true in the Events struct.
func (e *Events) Allowed(event, action string) bool {
	allowed := false

	// if there is an action, create `event:action` comparator string
	if len(action) > 0 {
		event = event + ":" + action
	}

	switch event {
	case constants.EventPush:
		allowed = e.GetPush().GetBranch()
	case constants.EventPull + ":" + constants.ActionOpened:
		allowed = e.GetPullRequest().GetOpened()
	case constants.EventPull + ":" + constants.ActionSynchronize:
		allowed = e.GetPullRequest().GetSynchronize()
	case constants.EventPull + ":" + constants.ActionEdited:
		allowed = e.GetPullRequest().GetEdited()
	case constants.EventPull + ":" + constants.ActionReopened:
		allowed = e.GetPullRequest().GetReopened()
	case constants.EventPull + ":" + constants.ActionLabeled:
		allowed = e.GetPullRequest().GetLabeled()
	case constants.EventPull + ":" + constants.ActionUnlabeled:
		allowed = e.GetPullRequest().GetUnlabeled()
	case constants.EventTag:
		allowed = e.GetPush().GetTag()
	case constants.EventComment + ":" + constants.ActionCreated:
		allowed = e.GetComment().GetCreated()
	case constants.EventComment + ":" + constants.ActionEdited:
		allowed = e.GetComment().GetEdited()
	case constants.EventDeploy + ":" + constants.ActionCreated:
		allowed = e.GetDeployment().GetCreated()
	case constants.EventSchedule:
		allowed = e.GetSchedule().GetRun()
	case constants.EventDelete + ":" + constants.ActionBranch:
		allowed = e.GetPush().GetDeleteBranch()
	case constants.EventDelete + ":" + constants.ActionTag:
		allowed = e.GetPush().GetDeleteTag()
	}

	return allowed
}

// List is an Events method that generates a comma-separated list of event:action
// combinations that are allowed for the repo.
func (e *Events) List() []string {
	eventSlice := []string{}

	if e.GetPush().GetBranch() {
		eventSlice = append(eventSlice, constants.EventPush)
	}

	if e.GetPullRequest().GetOpened() {
		eventSlice = append(eventSlice, constants.EventPull+":"+constants.ActionOpened)
	}

	if e.GetPullRequest().GetSynchronize() {
		eventSlice = append(eventSlice, constants.EventPull+":"+constants.ActionSynchronize)
	}

	if e.GetPullRequest().GetEdited() {
		eventSlice = append(eventSlice, constants.EventPull+":"+constants.ActionEdited)
	}

	if e.GetPullRequest().GetReopened() {
		eventSlice = append(eventSlice, constants.EventPull+":"+constants.ActionReopened)
	}

	if e.GetPullRequest().GetLabeled() {
		eventSlice = append(eventSlice, constants.EventPull+":"+constants.ActionLabeled)
	}

	if e.GetPullRequest().GetUnlabeled() {
		eventSlice = append(eventSlice, constants.EventPull+":"+constants.ActionUnlabeled)
	}

	if e.GetPush().GetTag() {
		eventSlice = append(eventSlice, constants.EventTag)
	}

	if e.GetDeployment().GetCreated() {
		eventSlice = append(eventSlice, constants.EventDeploy)
	}

	if e.GetComment().GetCreated() {
		eventSlice = append(eventSlice, constants.EventComment+":"+constants.ActionCreated)
	}

	if e.GetComment().GetEdited() {
		eventSlice = append(eventSlice, constants.EventComment+":"+constants.ActionEdited)
	}

	if e.GetSchedule().GetRun() {
		eventSlice = append(eventSlice, constants.EventSchedule)
	}

	if e.GetPush().GetDeleteBranch() {
		eventSlice = append(eventSlice, constants.EventDelete+":"+constants.ActionBranch)
	}

	if e.GetPush().GetDeleteTag() {
		eventSlice = append(eventSlice, constants.EventDelete+":"+constants.ActionTag)
	}

	return eventSlice
}

// ToDatabase is an Events method that converts a nested Events struct into an integer event mask.
func (e *Events) ToDatabase() int64 {
	return 0 |
		e.GetPush().ToMask() |
		e.GetPullRequest().ToMask() |
		e.GetComment().ToMask() |
		e.GetDeployment().ToMask() |
		e.GetSchedule().ToMask()
}

// GetPush returns the Push field from the provided Events. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (e *Events) GetPush() *actions.Push {
	// return zero value if Events type or Push field is nil
	if e == nil || e.Push == nil {
		return new(actions.Push)
	}

	return e.Push
}

// GetPullRequest returns the PullRequest field from the provided Events. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (e *Events) GetPullRequest() *actions.Pull {
	// return zero value if Events type or PullRequest field is nil
	if e == nil || e.PullRequest == nil {
		return new(actions.Pull)
	}

	return e.PullRequest
}

// GetDeployment returns the Deployment field from the provided Events. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (e *Events) GetDeployment() *actions.Deploy {
	// return zero value if Events type or Deployment field is nil
	if e == nil || e.Deployment == nil {
		return new(actions.Deploy)
	}

	return e.Deployment
}

// GetComment returns the Comment field from the provided Events. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (e *Events) GetComment() *actions.Comment {
	// return zero value if Events type or Comment field is nil
	if e == nil || e.Comment == nil {
		return new(actions.Comment)
	}

	return e.Comment
}

// GetSchedule returns the Schedule field from the provided Events. If the object is nil,
// or the field within the object is nil, it returns the zero value instead.
func (e *Events) GetSchedule() *actions.Schedule {
	// return zero value if Events type or Schedule field is nil
	if e == nil || e.Schedule == nil {
		return new(actions.Schedule)
	}

	return e.Schedule
}

// SetPush sets the Events Push field.
//
// When the provided Events type is nil, it
// will set nothing and immediately return.
func (e *Events) SetPush(v *actions.Push) {
	// return if Events type is nil
	if e == nil {
		return
	}

	e.Push = v
}

// SetPullRequest sets the Events PullRequest field.
//
// When the provided Events type is nil, it
// will set nothing and immediately return.
func (e *Events) SetPullRequest(v *actions.Pull) {
	// return if Events type is nil
	if e == nil {
		return
	}

	e.PullRequest = v
}

// SetDeployment sets the Events Deployment field.
//
// When the provided Events type is nil, it
// will set nothing and immediately return.
func (e *Events) SetDeployment(v *actions.Deploy) {
	// return if Events type is nil
	if e == nil {
		return
	}

	e.Deployment = v
}

// SetComment sets the Events Comment field.
//
// When the provided Events type is nil, it
// will set nothing and immediately return.
func (e *Events) SetComment(v *actions.Comment) {
	// return if Events type is nil
	if e == nil {
		return
	}

	e.Comment = v
}

// SetSchedule sets the Events Schedule field.
//
// When the provided Events type is nil, it
// will set nothing and immediately return.
func (e *Events) SetSchedule(v *actions.Schedule) {
	// return if Events type is nil
	if e == nil {
		return
	}

	e.Schedule = v
}
