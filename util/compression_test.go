// SPDX-License-Identifier: Apache-2.0

package util

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/constants"
)

func TestDatabase_CompressAndDecompress(t *testing.T) {
	// setup tests
	tests := []struct {
		name    string
		failure bool
		level   int
		dataSet [][]byte
	}{
		{
			name:    "compression level -1",
			failure: false,
			level:   constants.CompressionNegOne,
			dataSet: [][]byte{[]byte("foo"), complexYAML(), complexJSON()},
		},
		{
			name:    "compression level 0",
			failure: false,
			level:   constants.CompressionZero,
			dataSet: [][]byte{[]byte("foo"), complexYAML(), complexJSON()},
		},
		{
			name:    "compression level 1",
			failure: false,
			level:   constants.CompressionOne,
			dataSet: [][]byte{[]byte("foo"), complexYAML(), complexJSON()},
		},
		{
			name:    "compression level 2",
			failure: false,
			level:   constants.CompressionTwo,
			dataSet: [][]byte{[]byte("foo"), complexYAML(), complexJSON()},
		},
		{
			name:    "compression level 3",
			failure: false,
			level:   constants.CompressionThree,
			dataSet: [][]byte{[]byte("foo"), complexYAML(), complexJSON()},
		},
		{
			name:    "compression level 4",
			failure: false,
			level:   constants.CompressionFour,
			dataSet: [][]byte{[]byte("foo"), complexYAML(), complexJSON()},
		},
		{
			name:    "compression level 5",
			failure: false,
			level:   constants.CompressionFive,
			dataSet: [][]byte{[]byte("foo"), complexYAML(), complexJSON()},
		},
		{
			name:    "compression level 6",
			failure: false,
			level:   constants.CompressionSix,
			dataSet: [][]byte{[]byte("foo"), complexYAML(), complexJSON()},
		},
		{
			name:    "compression level 7",
			failure: false,
			level:   constants.CompressionSeven,
			dataSet: [][]byte{[]byte("foo"), complexYAML(), complexJSON()},
		},
		{
			name:    "compression level 8",
			failure: false,
			level:   constants.CompressionEight,
			dataSet: [][]byte{[]byte("foo"), complexYAML(), complexJSON()},
		},
		{
			name:    "compression level 9",
			failure: false,
			level:   constants.CompressionNine,
			dataSet: [][]byte{[]byte("foo"), complexYAML(), complexJSON()},
		},
	}

	// run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, data := range test.dataSet {
				got, err := Compress(test.level, data)

				if test.failure {
					if err == nil {
						t.Errorf("compress for %s should have returned err", test.name)
					}

					return
				}

				if err != nil {
					t.Errorf("compress for %s returned err: %v", test.name, err)
				}

				decompressed, err := Decompress(got)
				if err != nil {
					t.Errorf("decompress for %s returned err: %v", test.name, err)
				}

				if !reflect.DeepEqual(decompressed, data) {
					t.Errorf("compress for %s is %v, want %v", test.name, string(decompressed), string(data))
				}
			}
		})
	}
}

func complexJSON() []byte {
	return []byte(`
{ 
    "id": "step_name",
    "version": "1",
    "metadata":{
        "clone":true,
        "environment":["steps","services","secrets"]},
    "worker":{},
    "steps":[
        {
            "id":"step_github_octocat_1_init",
            "directory":"/vela/src/github.com/github/octocat",
            "environment": {"BUILD_AUTHOR":"Octocat"}
        }
    ]
}
`)
}

func complexYAML() []byte {
	return []byte(`
version: 1

worker:
  flavor: large
  platform: docker

services:
  - name: redis
    image: redis

steps:
  - name: ping
    image: redis
    commands:
      - redis-cli -h redis ping
`)
}
