// SPDX-License-Identifier: Apache-2.0

package github

import (
	"context"
	"reflect"
	"testing"
)

func TestGitHub_MergeQueuePrefix(t *testing.T) {
	// setup types
	want := mergeQueueBranchPrefix

	_service, err := New(
		context.Background(),
		WithAddress("https://github.com/"),
		WithClientID("foo"),
		WithClientSecret("bar"),
		WithServerAddress("https://vela-server.example.com"),
		WithStatusContext("continuous-integration/vela"),
		WithWebUIAddress("https://vela.example.com"),
	)
	if err != nil {
		t.Errorf("unable to create scm service: %v", err)
	}

	// run test
	got := _service.MergeQueueBranchPrefix()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("MergeQueueBranchPrefix is %v, want %v", got, want)
	}
}
