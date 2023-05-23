// Copyright (c) 2023 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package api

import (
	"testing"

	"github.com/go-vela/types/pipeline"
)

func Test_SkipEmptyBuild(t *testing.T) {
	type args struct {
		p *pipeline.Build
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{"only init stage", args{p: &pipeline.Build{Stages: []*pipeline.Stage{
			{
				Name: "init",
			},
		}}}, "skipping build since only init stage found"},
		{"init and clone stages", args{p: &pipeline.Build{Stages: []*pipeline.Stage{
			{
				Name: "init",
			},
			{
				Name: "clone",
			},
		}}}, "skipping build since only init and clone stages found"},
		{"three stages", args{p: &pipeline.Build{Stages: []*pipeline.Stage{
			{
				Name: "init",
			},
			{
				Name: "clone",
			},
			{
				Name: "foo",
			},
		}}}, ""},
		{"only init step", args{p: &pipeline.Build{Steps: []*pipeline.Container{
			{
				Name: "init",
			},
		}}}, "skipping build since only init step found"},
		{"init and clone steps", args{p: &pipeline.Build{Steps: []*pipeline.Container{
			{
				Name: "init",
			},
			{
				Name: "clone",
			},
		}}}, "skipping build since only init and clone steps found"},
		{"three steps", args{p: &pipeline.Build{Steps: []*pipeline.Container{
			{
				Name: "init",
			},
			{
				Name: "clone",
			},
			{
				Name: "foo",
			},
		}}}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SkipEmptyBuild(tt.args.p); got != tt.want {
				t.Errorf("SkipEmptyBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}
