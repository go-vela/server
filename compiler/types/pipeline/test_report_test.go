package pipeline

import "testing"

func TestPipeline_TestReport_Empty(t *testing.T) {
	// setup tests
	tests := []struct {
		report *TestReport
		want   bool
	}{
		{
			report: &TestReport{Results: []string{"foo"}},
			want:   false,
		},
		{
			report: new(TestReport),
			want:   true,
		},
	}

	// run tests
	for _, test := range tests {
		got := test.report.Empty()

		if got != test.want {
			t.Errorf("Empty is %v, want %t", got, test.want)
		}
	}
}
