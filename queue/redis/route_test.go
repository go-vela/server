// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package redis

import (
	"strings"
	"testing"

	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/pipeline"
)

func TestRedis_Client_Route(t *testing.T) {
	// setup
	client, _ := NewTest("vela")
	tests := []struct {
		want   string
		worker pipeline.Worker
	}{

		//  pipeline with not worker passed
		{
			want:   constants.DefaultRoute,
			worker: pipeline.Worker{},
		},
		{
			want:   "vela",
			worker: pipeline.Worker{},
		},
		{
			want:   "16cpu8gb",
			worker: pipeline.Worker{Flavor: "16cpu8gb"},
		},
		{
			want:   "16cpu8gb:gcp",
			worker: pipeline.Worker{Flavor: "16cpu8gb", Platform: "gcp"},
		},
		{
			want:   "gcp",
			worker: pipeline.Worker{Platform: "gcp"},
		},
	}

	// run
	for _, test := range tests {
		got, err := client.Route(&test.worker)

		if err != nil {
			t.Errorf("Route returned err: %v", err)
		}

		if !strings.EqualFold(got, test.want) {
			t.Errorf("Route is %v, want %v", got, test.want)
		}
	}
}
