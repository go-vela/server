// SPDX-License-Identifier: Apache-2.0

package redis

import (
	"strings"
	"testing"

	"github.com/go-vela/server/compiler/types/pipeline"
	"github.com/go-vela/server/constants"
)

func TestRedis_Client_Route(t *testing.T) {
	// setup
	client, _ := NewTest(_signingPrivateKey, _signingPublicKey, "vela", "16cpu8gb", "16cpu8gb:gcp", "gcp")
	tests := []struct {
		success bool
		want    string
		worker  pipeline.Worker
	}{

		//  pipeline with not worker passed
		{
			success: true,
			want:    constants.DefaultRoute,
			worker:  pipeline.Worker{},
		},
		{
			success: true,
			want:    "vela",
			worker:  pipeline.Worker{},
		},
		{
			success: true,
			want:    "16cpu8gb",
			worker:  pipeline.Worker{Flavor: "16cpu8gb"},
		},
		{
			success: true,
			want:    "16cpu8gb:gcp",
			worker:  pipeline.Worker{Flavor: "16cpu8gb", Platform: "gcp"},
		},
		{
			success: true,
			want:    "gcp",
			worker:  pipeline.Worker{Platform: "gcp"},
		},
		{
			success: false,
			want:    "",
			worker:  pipeline.Worker{Flavor: "bad", Platform: "route"},
		},
		{
			success: false,
			want:    "",
			worker:  pipeline.Worker{Flavor: "bad"},
		},
	}

	// run
	for _, test := range tests {
		got, err := client.Route(&test.worker)

		if test.success && err != nil {
			t.Errorf("Route returned err: %v", err)
		}

		if !test.success && err == nil {
			t.Errorf("Route returned %s, want err", got)
		}

		if !strings.EqualFold(got, test.want) {
			t.Errorf("Route is %v, want %v", got, test.want)
		}
	}
}
