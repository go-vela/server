// SPDX-License-Identifier: Apache-2.0

package pipeline

import (
	"testing"
)

func TestPipeline_Deployment_Empty(t *testing.T) {
	// setup tests
	tests := []struct {
		deployment *Deployment
		want       bool
	}{
		{
			deployment: &Deployment{Targets: []string{"foo"}},
			want:       false,
		},
		{
			deployment: &Deployment{Parameters: ParameterMap{"foo": new(Parameter)}},
			want:       false,
		},
		{
			deployment: new(Deployment),
			want:       true,
		},
		{
			want: true,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.deployment.Empty()

		if got != test.want {
			t.Errorf("Empty is %v, want %t", got, test.want)
		}
	}
}

func TestPipeline_Deployment_Validate(t *testing.T) {
	// setup types
	fullDeployConfig := &Deployment{
		Targets: []string{"east", "west", "north", "south"},
		Parameters: ParameterMap{
			"alpha": {
				Description: "foo",
				Type:        "string",
				Required:    true,
				Options:     []string{"bar", "baz"},
			},
			"beta": {
				Description: "bar",
				Type:        "string",
				Required:    false,
			},
			"gamma": {
				Description: "baz",
				Type:        "integer",
				Required:    true,
				Min:         -2,
				Max:         2,
			},
			"delta": {
				Description: "qux",
				Type:        "boolean",
				Required:    false,
			},
			"epsilon": {
				Description: "quux",
				Type:        "number",
			},
		},
	}

	// setup tests
	tests := []struct {
		inputTarget  string
		inputParams  map[string]string
		deployConfig *Deployment
		wantErr      bool
	}{
		{ // nil deployment config
			inputTarget: "north",
			inputParams: map[string]string{
				"alpha": "foo",
				"beta":  "bar",
			},
			wantErr: false,
		},
		{ // empty deployment config
			inputTarget: "north",
			inputParams: map[string]string{
				"alpha": "foo",
				"beta":  "bar",
			},
			deployConfig: new(Deployment),
			wantErr:      false,
		},
		{ // correct target and required params
			inputTarget: "west",
			inputParams: map[string]string{
				"alpha": "bar",
				"gamma": "1",
			},
			deployConfig: fullDeployConfig,
			wantErr:      false,
		},
		{ // correct target and wrong integer type for param gamma
			inputTarget: "east",
			inputParams: map[string]string{
				"alpha": "bar",
				"beta":  "test1",
				"gamma": "string",
			},
			deployConfig: fullDeployConfig,
			wantErr:      true,
		},
		{ // correct target and wrong boolean type for param delta
			inputTarget: "south",
			inputParams: map[string]string{
				"alpha": "bar",
				"beta":  "test2",
				"gamma": "2",
				"delta": "not-bool",
			},
			deployConfig: fullDeployConfig,
			wantErr:      true,
		},
		{ // correct target and wrong option for param alpha
			inputTarget: "south",
			inputParams: map[string]string{
				"alpha": "bazzy",
				"beta":  "test2",
				"gamma": "2",
				"delta": "true",
			},
			deployConfig: fullDeployConfig,
			wantErr:      true,
		},
		{ // correct target and gamma value over max
			inputTarget: "north",
			inputParams: map[string]string{
				"alpha": "bar",
				"beta":  "bar",
				"gamma": "3",
			},
			deployConfig: fullDeployConfig,
			wantErr:      true,
		},
		{ // correct target and gamma value under min
			inputTarget: "north",
			inputParams: map[string]string{
				"alpha": "baz",
				"beta":  "bar",
				"gamma": "-3",
			},
			deployConfig: fullDeployConfig,
			wantErr:      true,
		},
		{ // correct target and some number provided for epsilon param
			inputTarget: "north",
			inputParams: map[string]string{
				"alpha":   "bar",
				"beta":    "bar",
				"gamma":   "1",
				"delta":   "true",
				"epsilon": "42",
			},
			deployConfig: fullDeployConfig,
			wantErr:      false,
		},
	}

	// run tests
	for _, test := range tests {
		err := test.deployConfig.Validate(test.inputTarget, test.inputParams)

		if err != nil && !test.wantErr {
			t.Errorf("Deployment.Validate returned err: %v", err)
		}

		if err == nil && test.wantErr {
			t.Errorf("Deployment.Validate did not return err")
		}
	}
}
