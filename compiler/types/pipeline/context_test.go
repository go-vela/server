// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"context"
	"testing"
)

func TestPipeline_BuildFromContext(t *testing.T) {
	// setup types
	b := &Build{ID: "1"}

	// setup tests
	tests := []struct {
		ctx  context.Context
		want *Build
	}{
		{
			ctx:  context.WithValue(context.Background(), buildKey, b),
			want: b,
		},
		{
			ctx:  context.Background(),
			want: nil,
		},
		{
			ctx:  context.WithValue(context.Background(), buildKey, "foo"),
			want: nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := BuildFromContext(test.ctx)

		if got != test.want {
			t.Errorf("BuildFromContext is %v, want %v", got, test.want)
		}
	}
}

func TestPipeline_BuildWithContext(t *testing.T) {
	// setup types
	want := &Build{ID: "1"}

	// setup context
	ctx := BuildWithContext(context.Background(), want)

	// run test
	got := ctx.Value(buildKey)

	if got != want {
		t.Errorf("BuildWithContext is %v, want %v", got, want)
	}
}

func TestPipeline_SecretFromContext(t *testing.T) {
	// setup types
	s := &Secret{Name: "foo"}

	// setup tests
	tests := []struct {
		ctx  context.Context
		want *Secret
	}{
		{
			ctx:  context.WithValue(context.Background(), secretKey, s),
			want: s,
		},
		{
			ctx:  context.Background(),
			want: nil,
		},
		{
			ctx:  context.WithValue(context.Background(), secretKey, "foo"),
			want: nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := SecretFromContext(test.ctx)

		if got != test.want {
			t.Errorf("SecretFromContext is %v, want %v", got, test.want)
		}
	}
}

func TestPipeline_SecretWithContext(t *testing.T) {
	// setup types
	want := &Secret{Name: "foo"}

	// setup context
	ctx := SecretWithContext(context.Background(), want)

	// run test
	got := ctx.Value(secretKey)

	if got != want {
		t.Errorf("SecretWithContext is %v, want %v", got, want)
	}
}

func TestPipeline_StageFromContext(t *testing.T) {
	// setup types
	s := &Stage{Name: "foo"}

	// setup tests
	tests := []struct {
		ctx  context.Context
		want *Stage
	}{
		{
			ctx:  context.WithValue(context.Background(), stageKey, s),
			want: s,
		},
		{
			ctx:  context.Background(),
			want: nil,
		},
		{
			ctx:  context.WithValue(context.Background(), stageKey, "foo"),
			want: nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := StageFromContext(test.ctx)

		if got != test.want {
			t.Errorf("StageFromContext is %v, want %v", got, test.want)
		}
	}
}

func TestPipeline_StageWithContext(t *testing.T) {
	// setup types
	want := &Stage{Name: "foo"}

	// setup context
	ctx := StageWithContext(context.Background(), want)

	// run test
	got := ctx.Value(stageKey)

	if got != want {
		t.Errorf("StageWithContext is %v, want %v", got, want)
	}
}

func TestPipeline_ContainerFromContext(t *testing.T) {
	// setup types
	c := &Container{Name: "foo"}

	// setup tests
	tests := []struct {
		ctx  context.Context
		want *Container
	}{
		{
			ctx:  context.WithValue(context.Background(), containerKey, c),
			want: c,
		},
		{
			ctx:  context.Background(),
			want: nil,
		},
		{
			ctx:  context.WithValue(context.Background(), containerKey, "foo"),
			want: nil,
		},
	}

	// run tests
	for _, test := range tests {
		got := ContainerFromContext(test.ctx)

		if got != test.want {
			t.Errorf("ContainerFromContext is %v, want %v", got, test.want)
		}
	}
}

func TestPipeline_ContainerWithContext(t *testing.T) {
	// setup types
	want := &Container{ID: "1"}

	// setup context
	ctx := ContainerWithContext(context.Background(), want)

	// run test
	got := ctx.Value(containerKey)

	if got != want {
		t.Errorf("ContainerWithContext is %v, want %v", got, want)
	}
}
