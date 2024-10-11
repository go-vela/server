// SPDX-License-Identifier: Apache-2.0

package yaml

import (
	"reflect"
	"testing"

	"github.com/go-vela/server/compiler/types/pipeline"
)

func TestYaml_Worker_ToPipeline(t *testing.T) {
	// setup tests
	tests := []struct {
		worker *Worker
		want   *pipeline.Worker
	}{
		{
			worker: &Worker{
				Flavor:   "8cpu16gb",
				Platform: "gcp",
			},
			want: &pipeline.Worker{
				Flavor:   "8cpu16gb",
				Platform: "gcp",
			},
		},
	}

	// run tests
	for _, test := range tests {
		got := test.worker.ToPipeline()

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ToPipeline is %v, want %v", got, test.want)
		}
	}
}
