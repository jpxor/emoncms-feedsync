package utils

import (
	"math"
	"testing"
)

func TestParseUnixTimestamp(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{"Valid timestamp", "1234567890", 1234567890, false},
		{"Zero", "0", 0, false},
		{"Large number", "9223372036854775807", 9223372036854775807, false},
		{"Invalid character", "123a456", 0, true},
		{"Empty string", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUnixTimestamp(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUnixTimestamp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseUnixTimestamp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFloatLowPrecision(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    float32
		wantErr bool
	}{
		{"Positive integer", "123", 123.0, false},
		{"Negative integer", "-456", -456.0, false},
		{"Positive float", "78.9", 78.9, false},
		{"Negative float", "-12.34", -12.34, false},
		{"Three decimal places", "0.123", 0.123, false},
		{"More than three decimal places", "1.234567", 1.234, false},
		{"Zero", "0", 0, false},
		{"Invalid character", "1.23a", 0, true},
		{"Multiple dots", "1.2.3", 0, true},
		{"Empty string", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFloatLowPrecision(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFloatLowPrecision() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if math.Abs(float64(got-tt.want)) > 1e-6 {
				t.Errorf("ParseFloatLowPrecision() = %v, want %v", got, tt.want)
			}
		})
	}
}
