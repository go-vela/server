// SPDX-License-Identifier: Apache-2.0

package settings

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypes_Queue_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		queue *Queue
		want  *Queue
	}{
		{
			queue: testQueueSettings(),
			want:  testQueueSettings(),
		},
		{
			queue: new(Queue),
			want:  new(Queue),
		},
	}

	// run tests
	for _, test := range tests {
		if !reflect.DeepEqual(test.queue.GetRoutes(), test.want.GetRoutes()) {
			t.Errorf("GetRoutes is %v, want %v", test.queue.GetRoutes(), test.want.GetRoutes())
		}
	}
}

func TestTypes_Queue_Setters(t *testing.T) {
	// setup types
	var qs *Queue

	// setup tests
	tests := []struct {
		queue *Queue
		want  *Queue
	}{
		{
			queue: testQueueSettings(),
			want:  testQueueSettings(),
		},
		{
			queue: qs,
			want:  new(Queue),
		},
	}

	// run tests
	for _, test := range tests {
		test.queue.SetRoutes(test.want.GetRoutes())

		if !reflect.DeepEqual(test.queue.GetRoutes(), test.want.GetRoutes()) {
			t.Errorf("SetRoutes is %v, want %v", test.queue.GetRoutes(), test.want.GetRoutes())
		}
	}
}

func TestTypes_Queue_String(t *testing.T) {
	// setup types
	qs := testQueueSettings()

	want := fmt.Sprintf(`{
  Routes: %s,
}`,
		qs.GetRoutes(),
	)

	// run test
	got := qs.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testQueueSettings is a test helper function to create a Queue
// type with all fields set to a fake value.
func testQueueSettings() *Queue {
	qs := new(Queue)

	qs.SetRoutes([]string{"vela"})

	return qs
}
