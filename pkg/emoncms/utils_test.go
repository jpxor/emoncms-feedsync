package emoncms

import (
	"testing"
)

func TestSplitFeedDataString(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantLeft  string
		wantRight string
	}{
		{
			name:      "Empty string",
			input:     "",
			wantLeft:  "",
			wantRight: "",
		},
		{
			name:      "Single pair",
			input:     "[[1.0,2.0]]",
			wantLeft:  "[[1.0,2.0]]",
			wantRight: "",
		},
		{
			name:      "Two pairs",
			input:     "[[1.0,2.0],[3.0,4.0]]",
			wantLeft:  "[[1.0,2.0]]",
			wantRight: "[[3.0,4.0]]",
		},
		{
			name:      "Evennumber of pairs",
			input:     "[[1.0,2.0],[3.0,4.0],[5.0,6.0],[7.0,8.0]]",
			wantLeft:  "[[1.0,2.0],[3.0,4.0]]",
			wantRight: "[[5.0,6.0],[7.0,8.0]]",
		},
		{
			name:      "Odd number of pairs",
			input:     "[[1.0,2.0],[3.0,4.0],[5.0,6.0]]",
			wantLeft:  "[[1.0,2.0]]",
			wantRight: "[[3.0,4.0],[5.0,6.0]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLeft, gotRight := splitFeedDataString(tt.input)
			if gotLeft != tt.wantLeft {
				t.Errorf("splitFeedDataString() gotLeft = %v, want %v", gotLeft, tt.wantLeft)
			}
			if gotRight != tt.wantRight {
				t.Errorf("splitFeedDataString() gotRight = %v, want %v", gotRight, tt.wantRight)
			}
		})
	}
}
