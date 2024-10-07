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

func TestAppendFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    float32
		initial  []byte
		expected []byte
	}{
		{"Positive integer", 123, []byte(""), []byte("123")},
		{"Negative integer", -456, []byte(""), []byte("-456")},
		{"Positive float", 78.9, []byte(""), []byte("78.9")},
		{"Negative float", -12.34, []byte(""), []byte("-12.34")},
		{"Three decimal places", 0.123, []byte(""), []byte("0.123")},
		{"More than three decimal places", 1.234567, []byte(""), []byte("1.234")},
		{"Zero", 0, []byte(""), []byte("0")},
		{"Append to existing", 5.67, []byte("Value: "), []byte("Value: 5.67")},
		{"Very small positive", 0.001, []byte(""), []byte("0.001")},
		{"Very small negative", -0.001, []byte(""), []byte("-0.001")},
		{"Large number", 123456.789, []byte(""), []byte("123456.789")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AppendFloat(tt.initial, tt.input)
			if string(result) != string(tt.expected) {
				t.Errorf("AppendFloat() = %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestAppendUInt(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		initial  []byte
		expected []byte
	}{
		{"Zero", 0, []byte(""), []byte("0")},
		{"Positive small", 123, []byte(""), []byte("123")},
		{"Positive large", 9876543210, []byte(""), []byte("9876543210")},
		{"Append to existing", 42, []byte("Number: "), []byte("Number: 42")},
		{"Max int64", 9223372036854775807, []byte(""), []byte("9223372036854775807")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AppendUInt(tt.initial, tt.input)
			if string(result) != string(tt.expected) {
				t.Errorf("AppendUInt() = %s, want %s", result, tt.expected)
			}
		})
	}
}
