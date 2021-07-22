package api

import (
	"github.com/go-vela/types/pipeline"
	"testing"
)

func Test_skipEmptyBuild(t *testing.T) {
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
			if got := skipEmptyBuild(tt.args.p); got != tt.want {
				t.Errorf("skipEmptyBuild() = %v, want %v", got, tt.want)
			}
		})
	}
}
