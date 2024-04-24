// SPDX-License-Identifier: Apache-2.0

package types

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTypes_QueueBuild_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		buildQueue *QueueBuild
		want       *QueueBuild
	}{
		{
			buildQueue: testQueueBuild(),
			want:       testQueueBuild(),
		},
		{
			buildQueue: new(QueueBuild),
			want:       new(QueueBuild),
		},
	}

	// run tests
	for _, test := range tests {
		if test.buildQueue.GetNumber() != test.want.GetNumber() {
			t.Errorf("GetNumber is %v, want %v", test.buildQueue.GetNumber(), test.want.GetNumber())
		}

		if test.buildQueue.GetStatus() != test.want.GetStatus() {
			t.Errorf("GetStatus is %v, want %v", test.buildQueue.GetStatus(), test.want.GetStatus())
		}

		if test.buildQueue.GetCreated() != test.want.GetCreated() {
			t.Errorf("GetCreated is %v, want %v", test.buildQueue.GetCreated(), test.want.GetCreated())
		}

		if test.buildQueue.GetFullName() != test.want.GetFullName() {
			t.Errorf("GetFullName is %v, want %v", test.buildQueue.GetFullName(), test.want.GetFullName())
		}
	}
}

func TestTypes_QueueBuild_Setters(t *testing.T) {
	// setup types
	var b *QueueBuild

	// setup tests
	tests := []struct {
		buildQueue *QueueBuild
		want       *QueueBuild
	}{
		{
			buildQueue: testQueueBuild(),
			want:       testQueueBuild(),
		},
		{
			buildQueue: b,
			want:       new(QueueBuild),
		},
	}

	// run tests
	for _, test := range tests {
		test.buildQueue.SetNumber(test.want.GetNumber())
		test.buildQueue.SetStatus(test.want.GetStatus())
		test.buildQueue.SetCreated(test.want.GetCreated())
		test.buildQueue.SetFullName(test.want.GetFullName())

		if test.buildQueue.GetNumber() != test.want.GetNumber() {
			t.Errorf("SetNumber is %v, want %v", test.buildQueue.GetNumber(), test.want.GetNumber())
		}

		if test.buildQueue.GetStatus() != test.want.GetStatus() {
			t.Errorf("SetStatus is %v, want %v", test.buildQueue.GetStatus(), test.want.GetStatus())
		}

		if test.buildQueue.GetCreated() != test.want.GetCreated() {
			t.Errorf("SetCreated is %v, want %v", test.buildQueue.GetCreated(), test.want.GetCreated())
		}

		if test.buildQueue.GetFullName() != test.want.GetFullName() {
			t.Errorf("SetFullName is %v, want %v", test.buildQueue.GetFullName(), test.want.GetFullName())
		}
	}
}

func TestTypes_QueueBuild_String(t *testing.T) {
	// setup types
	b := testQueueBuild()

	want := fmt.Sprintf(`{
  Created: %d,
  FullName: %s,
  Number: %d,
  Status: %s,
}`,
		b.GetCreated(),
		b.GetFullName(),
		b.GetNumber(),
		b.GetStatus(),
	)

	// run test
	got := b.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

// testQueueBuild is a test helper function to create a QueueBuild
// type with all fields set to a fake value.
func testQueueBuild() *QueueBuild {
	b := new(QueueBuild)

	b.SetNumber(1)
	b.SetStatus("running")
	b.SetCreated(1563474076)
	b.SetFullName("github/octocat")

	return b
}
