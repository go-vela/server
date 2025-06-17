// SPDX-License-Identifier: Apache-2.0

package types

import (
	"testing"
)

func TestTypes_StorageInfo_Getters(t *testing.T) {
	// setup tests
	tests := []struct {
		sI   *StorageInfo
		want *StorageInfo
	}{
		{
			sI:   testStorageInfo(),
			want: testStorageInfo(),
		},
		{
			sI:   new(StorageInfo),
			want: new(StorageInfo),
		},
	}

	// run tests
	for _, test := range tests {
		if test.sI.GetAccessKey() != test.want.GetAccessKey() {
			t.Errorf("GetAccessKey is %v, want %v", test.sI.GetAccessKey(), test.want.GetAccessKey())
		}

		if test.sI.GetSecretKey() != test.want.GetSecretKey() {
			t.Errorf("GetSecretKey is %v, want %v", test.sI.GetSecretKey(), test.want.GetSecretKey())
		}

		if test.sI.GetStorageAddress() != test.want.GetStorageAddress() {
			t.Errorf("GetStorageAddress is %v, want %v", test.sI.GetStorageAddress(), test.want.GetStorageAddress())
		}
		if test.sI.GetStorageBucket() != test.want.GetStorageBucket() {
			t.Errorf("GetStorageBucket is %v, want %v", test.sI.GetStorageBucket(), test.want.GetStorageBucket())
		}
	}
}

func TestTypes_StorageInfo_Setters(t *testing.T) {
	// setup types
	var sI *StorageInfo

	// setup tests
	tests := []struct {
		sI   *StorageInfo
		want *StorageInfo
	}{
		{
			sI:   testStorageInfo(),
			want: testStorageInfo(),
		},
		{
			sI:   sI,
			want: new(StorageInfo),
		},
	}

	// run tests
	for _, test := range tests {
		test.sI.SetAccessKey(test.want.GetAccessKey())
		test.sI.SetSecretKey(test.want.GetSecretKey())
		test.sI.SetStorageAddress(test.want.GetStorageAddress())
		test.sI.SetStorageBucket(test.want.GetStorageBucket())

		if test.sI.GetAccessKey() != test.want.GetAccessKey() {
			t.Errorf("GetAccessKey is %v, want %v", test.sI.GetAccessKey(), test.want.GetAccessKey())
		}

		if test.sI.GetSecretKey() != test.want.GetSecretKey() {
			t.Errorf("GetSecretKey is %v, want %v", test.sI.GetSecretKey(), test.want.GetSecretKey())
		}

		if test.sI.GetStorageAddress() != test.want.GetStorageAddress() {
			t.Errorf("GetStorageAddress is %v, want %v", test.sI.GetStorageAddress(), test.want.GetStorageAddress())
		}
	}
}

// testStorageInfo is a test helper function to register a StorageInfo
// type with all fields set to a fake value.
func testStorageInfo() *StorageInfo {
	sI := new(StorageInfo)
	sI.SetAccessKey("fakeAccessKey")
	sI.SetSecretKey("fakeSecretKey")
	sI.SetStorageAddress("http://localhost:8080")

	return sI
}
