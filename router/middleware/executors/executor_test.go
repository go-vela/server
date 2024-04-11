// SPDX-License-Identifier: Apache-2.0

package executors

import (
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"

	api "github.com/go-vela/server/api/types"
)

func TestExecutors_Retrieve(t *testing.T) {
	// setup types
	eID := int64(1)
	e := api.Executor{ID: &eID}
	want := []api.Executor{e}

	// setup context
	gin.SetMode(gin.TestMode)

	context, _ := gin.CreateTestContext(nil)
	ToContext(context, want)

	// run test
	got := Retrieve(context)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Retrieve is %v, want %v", got, want)
	}
}
