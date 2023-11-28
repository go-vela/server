// SPDX-License-Identifier: Apache-2.0

package build

import (
	"testing"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
)

func Test_isCancelable(t *testing.T) {
	// setup types
	pushEvent := constants.EventPush
	pullEvent := constants.EventPull
	tagEvent := constants.EventTag

	branchDev := "dev"
	branchPatch := "patch-1"

	actionOpened := constants.ActionOpened
	actionSync := constants.ActionSynchronize
	actionEdited := constants.ActionEdited

	tests := []struct {
		name    string
		target  *library.Build
		current *library.Build
		want    bool
	}{
		{
			name: "Wrong Event",
			target: &library.Build{
				Event:  &tagEvent,
				Branch: &branchDev,
			},
			current: &library.Build{
				Event:  &pushEvent,
				Branch: &branchDev,
			},
			want: false,
		},
		{
			name: "Cancelable Push",
			target: &library.Build{
				Event:  &pushEvent,
				Branch: &branchDev,
			},
			current: &library.Build{
				Event:  &pushEvent,
				Branch: &branchDev,
			},
			want: true,
		},
		{
			name: "Push Branch Mismatch",
			target: &library.Build{
				Event:  &pushEvent,
				Branch: &branchDev,
			},
			current: &library.Build{
				Event:  &pushEvent,
				Branch: &branchPatch,
			},
			want: false,
		},
		{
			name: "Event Mismatch",
			target: &library.Build{
				Event:  &pushEvent,
				Branch: &branchDev,
			},
			current: &library.Build{
				Event:   &pullEvent,
				Branch:  &branchDev,
				HeadRef: &branchPatch,
			},
			want: false,
		},
		{
			name: "Cancelable Pull",
			target: &library.Build{
				Event:       &pullEvent,
				Branch:      &branchDev,
				HeadRef:     &branchPatch,
				EventAction: &actionOpened,
			},
			current: &library.Build{
				Event:       &pullEvent,
				Branch:      &branchDev,
				HeadRef:     &branchPatch,
				EventAction: &actionSync,
			},
			want: true,
		},
		{
			name: "Pull Head Ref Mismatch",
			target: &library.Build{
				Event:       &pullEvent,
				Branch:      &branchDev,
				HeadRef:     &branchPatch,
				EventAction: &actionSync,
			},
			current: &library.Build{
				Event:       &pullEvent,
				Branch:      &branchDev,
				HeadRef:     &branchDev,
				EventAction: &actionSync,
			},
			want: false,
		},
		{
			name: "Pull Ineligible Action",
			target: &library.Build{
				Event:       &pullEvent,
				Branch:      &branchDev,
				HeadRef:     &branchPatch,
				EventAction: &actionEdited,
			},
			current: &library.Build{
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
			if got := isCancelable(tt.target, tt.current); got != tt.want {
				t.Errorf("test %s: isCancelable() = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func Test_ShouldAutoCancel(t *testing.T) {
	// setup types
	pushEvent := constants.EventPush
	pullEvent := constants.EventPull
	tagEvent := constants.EventTag

	branchDev := "dev"
	branchPatch := "patch-1"

	actionOpened := constants.ActionOpened
	actionSync := constants.ActionSynchronize

	tests := []struct {
		name   string
		opts   *pipeline.CancelOptions
		build  *library.Build
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
			build: &library.Build{
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
			build: &library.Build{
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
			build: &library.Build{
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
			build: &library.Build{
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
			build: &library.Build{
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
			build: &library.Build{
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
			build: &library.Build{
				Event:       &pullEvent,
				Branch:      &branchDev,
				EventAction: &actionOpened,
			},
			branch: branchDev,
			want:   false,
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
