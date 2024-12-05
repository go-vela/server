// SPDX-License-Identifier: Apache-2.0

package schema

import (
	"testing"
)

func TestSchema_NewPipelineSchema(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "basic schema generation",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPipelineSchema()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPipelineSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("NewPipelineSchema() returned nil schema without error")
			}
			if !tt.wantErr {
				if got.Title != "Vela Pipeline Configuration" {
					t.Errorf("NewPipelineSchema() title = %v, want %v", got.Title, "Vela Pipeline Configuration")
				}
				if got.AdditionalProperties != nil {
					t.Error("NewPipelineSchema() AdditionalProperties should be nil")
				}
			}
		})
	}
}
