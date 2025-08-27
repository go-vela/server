// SPDX-License-Identifier: Apache-2.0

package build

import (
	"testing"

	"github.com/go-vela/server/api/types"
	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
)

func Test_isCancelable(t *testing.T) {
	// setup types
	pushEvent := constants.EventPush
	pullEvent := constants.EventPull
	mergeEvent := constants.EventMergeGroup
	tagEvent := constants.EventTag

	branchDev := "dev"
	branchPatch := "patch-1"
	branchMergeGroupA := "merge/123abc"
	branchMergeGroupB := "merge/456def"

	actionOpened := constants.ActionOpened
	actionSync := constants.ActionSynchronize
	actionEdited := constants.ActionEdited

	tests := []struct {
		name    string
		target  *types.Build
		current *types.Build
		want    bool
	}{
		{
			name: "Wrong Event",
			target: &types.Build{
				Event:  &tagEvent,
				Branch: &branchDev,
			},
			current: &types.Build{
				Event:  &pushEvent,
				Branch: &branchDev,
			},
			want: false,
		},
		{
			name: "Cancelable Push",
			target: &types.Build{
				Event:  &pushEvent,
				Branch: &branchDev,
			},
			current: &types.Build{
				Event:  &pushEvent,
				Branch: &branchDev,
			},
			want: true,
		},
		{
			name: "Push Branch Mismatch",
			target: &types.Build{
				Event:  &pushEvent,
				Branch: &branchDev,
			},
			current: &types.Build{
				Event:  &pushEvent,
				Branch: &branchPatch,
			},
			want: false,
		},
		{
			name: "Event Mismatch",
			target: &types.Build{
				Event:  &pushEvent,
				Branch: &branchDev,
			},
			current: &types.Build{
				Event:   &pullEvent,
				Branch:  &branchDev,
				HeadRef: &branchPatch,
			},
			want: false,
		},
		{
			name: "Cancelable Pull",
			target: &types.Build{
				Event:       &pullEvent,
				Branch:      &branchDev,
				HeadRef:     &branchPatch,
				EventAction: &actionOpened,
			},
			current: &types.Build{
				Event:       &pullEvent,
				Branch:      &branchDev,
				HeadRef:     &branchPatch,
				EventAction: &actionSync,
			},
			want: true,
		},
		{
			name: "Pull Head Ref Mismatch",
			target: &types.Build{
				Event:       &pullEvent,
				Branch:      &branchDev,
				HeadRef:     &branchPatch,
				EventAction: &actionSync,
			},
			current: &types.Build{
				Event:       &pullEvent,
				Branch:      &branchDev,
				HeadRef:     &branchDev,
				EventAction: &actionSync,
			},
			want: false,
		},
		{
			name: "Cancelable Merge Group",
			target: &types.Build{
				Event:  &mergeEvent,
				Branch: &branchMergeGroupA,
			},
			current: &types.Build{
				Event:  &mergeEvent,
				Branch: &branchMergeGroupB,
			},
			want: true,
		},
		{
			name: "Pull Ineligible Action",
			target: &types.Build{
				Event:       &pullEvent,
				Branch:      &branchDev,
				HeadRef:     &branchPatch,
				EventAction: &actionEdited,
			},
			current: &types.Build{
				Event:       &pullEvent,
				Branch:      &branchDev,
				HeadRef:     &branchDev,
				EventAction: &actionSync,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCancelable(tt.target, tt.current, "merge/"); got != tt.want {
				t.Errorf("test %s: isCancelable() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_ShouldAutoCancel(t *testing.T) {
	// setup types
	pushEvent := constants.EventPush
	pullEvent := constants.EventPull
	mergeEvent := constants.EventMergeGroup
	tagEvent := constants.EventTag

	branchDev := "dev"
	branchPatch := "patch-1"
	branchMergeGroup := "merge/123abc"

	statusPendingApproval := constants.StatusPendingApproval

	actionOpened := constants.ActionOpened
	actionSync := constants.ActionSynchronize

	tests := []struct {
		name   string
		opts   *pipeline.CancelOptions
		build  *types.Build
		branch string
		want   bool
	}{
		{
			name: "Wrong Event",
			opts: &pipeline.CancelOptions{
				Running:       true,
				Pending:       true,
				DefaultBranch: true,
			},
			build: &types.Build{
				Event:  &tagEvent,
				Branch: &branchPatch,
			},
			branch: branchDev,
			want:   false,
		},
		{
			name: "Auto Cancel Disabled",
			opts: &pipeline.CancelOptions{
				Running:       false,
				Pending:       false,
				DefaultBranch: false,
			},
			build: &types.Build{
				Event:  &pushEvent,
				Branch: &branchPatch,
			},
			branch: branchDev,
			want:   false,
		},
		{
			name: "Eligible Push",
			opts: &pipeline.CancelOptions{
				Running:       true,
				Pending:       true,
				DefaultBranch: false,
			},
			build: &types.Build{
				Event:  &pushEvent,
				Branch: &branchPatch,
			},
			branch: branchDev,
			want:   true,
		},
		{
			name: "Eligible Push - Default Branch",
			opts: &pipeline.CancelOptions{
				Running:       true,
				Pending:       true,
				DefaultBranch: true,
			},
			build: &types.Build{
				Event:  &pushEvent,
				Branch: &branchDev,
			},
			branch: branchDev,
			want:   true,
		},
		{
			name: "Push Mismatch - Default Branch",
			opts: &pipeline.CancelOptions{
				Running:       true,
				Pending:       true,
				DefaultBranch: false,
			},
			build: &types.Build{
				Event:  &pushEvent,
				Branch: &branchDev,
			},
			branch: branchDev,
			want:   false,
		},
		{
			name: "Eligible Pull",
			opts: &pipeline.CancelOptions{
				Running:       true,
				Pending:       true,
				DefaultBranch: false,
			},
			build: &types.Build{
				Event:       &pullEvent,
				Branch:      &branchDev,
				EventAction: &actionSync,
			},
			branch: branchDev,
			want:   true,
		},
		{
			name: "Pull Mismatch - Action",
			opts: &pipeline.CancelOptions{
				Running:       true,
				Pending:       true,
				DefaultBranch: false,
			},
			build: &types.Build{
				Event:       &pullEvent,
				Branch:      &branchDev,
				EventAction: &actionOpened,
			},
			branch: branchDev,
			want:   false,
		},
		{
			name: "Pending Approval Build",
			opts: &pipeline.CancelOptions{
				Running:       false,
				Pending:       false,
				DefaultBranch: false,
			},
			build: &types.Build{
				Event:       &pullEvent,
				Branch:      &branchDev,
				EventAction: &actionOpened,
				Status:      &statusPendingApproval,
			},
			branch: branchDev,
			want:   true,
		},
		{
			name: "Merge Group build",
			opts: &pipeline.CancelOptions{
				Running:       false,
				Pending:       false,
				DefaultBranch: false,
			},
			build: &types.Build{
				Event:       &mergeEvent,
				Branch:      &branchMergeGroup,
				EventAction: &actionOpened,
				Status:      &statusPendingApproval,
			},
			branch: branchDev,
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShouldAutoCancel(tt.opts, tt.build, tt.branch); got != tt.want {
				t.Errorf("test %s: ShouldAutoCancel() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}
