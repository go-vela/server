// SPDX-License-Identifier: Apache-2.0

//go:build ignore

package main

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/go-vela/server/schema"
)

func main() {
	js, err := schema.NewPipelineSchema()
	if err != nil {
		logrus.Fatal("schema generation failed:", err)
	}

	// output json
	j, err := json.MarshalIndent(js, "", "  ")
	if err != nil {
		logrus.Fatal(err)
	}

	fmt.Printf("%s\n", j)
}
