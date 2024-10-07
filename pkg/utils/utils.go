package utils

import "fmt"

func ParseUnixTimestamp(s string) (int64, error) {
	var result int64
	if len(s) == 0 {
		return 0, fmt.Errorf("invalid timestamp: empty")
	}
	for i := 0; i < len(s); i++ {
		digit := s[i] - '0'
		if digit > 9 {
			return 0, fmt.Errorf("invalid character in timestamp: %c", s[i])
		}
		result = result*10 + int64(digit)
	}
	return result, nil
}

// parses up to 3 decimal places: 0.123
func ParseFloatLowPrecision(s string) (float32, error) {
	var intPart, fracPart int32
	var fracDiv float32 = 1
	var negative bool
	var i int

	if len(s) == 0 {
		return 0, fmt.Errorf("invalid float: empty string")
	}

	if s[0] == '-' {
		negative = true
		i++
	}
	for ; i < len(s) && s[i] != '.'; i++ {
		digit := s[i] - '0'
		if digit > 9 {
			return 0, fmt.Errorf("invalid character in float: %c", s[i])
		}
		intPart = intPart*10 + int32(digit)
	}
	if i < len(s) && s[i] == '.' {
		i++
		for j := 0; j < 3 && i < len(s); j++ {
			digit := s[i] - '0'
			if digit > 9 {
				return 0, fmt.Errorf("invalid character in float: %c", s[i])
			}
			fracPart = fracPart*10 + int32(digit)
			fracDiv *= 10
			i++
		}
	}
	result := float32(intPart) + float32(fracPart)/fracDiv
	if negative {
		result = -result
	}
	return result, nil
}
