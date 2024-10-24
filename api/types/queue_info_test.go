// SPDX-License-Identifier: Apache-2.0

package types

import (
	"testing"
)

func TestTypes_QueueInfo_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		qR   *QueueInfo
		want *QueueInfo
	}{
		{
			qR:   testQueueInfo(),
			want: testQueueInfo(),
		},
		{
			qR:   new(QueueInfo),
			want: new(QueueInfo),
		},
	}

	// run tests
	for _, test := range tests {
		if test.qR.GetQueueAddress() != test.want.GetQueueAddress() {
			t.Errorf("GetQueueAddress is %v, want %v", test.qR.GetQueueAddress(), test.want.GetQueueAddress())
		}

		if test.qR.GetPublicKey() != test.want.GetPublicKey() {
			t.Errorf("GetPublicKey is %v, want %v", test.qR.GetPublicKey(), test.want.GetPublicKey())
		}
	}
}

func TestTypes_QueueInfo_Setters(t *testing.T) {
	// setup types
	var w *QueueInfo

	// setup tests
	tests := []struct {
		qR   *QueueInfo
		want *QueueInfo
	}{
		{
			qR:   testQueueInfo(),
			want: testQueueInfo(),
		},
		{
			qR:   w,
			want: new(QueueInfo),
		},
	}

	// run tests
	for _, test := range tests {
		test.qR.SetQueueAddress(test.want.GetQueueAddress())
		test.qR.SetPublicKey(test.want.GetPublicKey())

		if test.qR.GetQueueAddress() != test.want.GetQueueAddress() {
			t.Errorf("GetQueueAddress is %v, want %v", test.qR.GetQueueAddress(), test.want.GetQueueAddress())
		}

		if test.qR.GetPublicKey() != test.want.GetPublicKey() {
			t.Errorf("GetPublicKey is %v, want %v", test.qR.GetPublicKey(), test.want.GetPublicKey())
		}
	}
}

// testQueueInfo is a test helper function to register a QueueInfo
// type with all fields set to a fake value.
func testQueueInfo() *QueueInfo {
	w := new(QueueInfo)
	w.SetQueueAddress("http://localhost:8080")
	w.SetPublicKey("CuS+EQAzofbk3tVFS3bt5f2tIb4YiJJC4nVMFQYQElg=")

	return w
}
