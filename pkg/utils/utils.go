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

// encodes up to 3 decimal places: 0.123
func AppendFloat(b []byte, f float32) []byte {
	if f < 0 {
		b = append(b, '-')
		f = -f
	}
	intPart := uint64(f)
	b = AppendUInt(b, intPart)

	fracPart := uint64((f - float32(intPart)) * 1000)
	if fracPart > 0 {
		b = append(b, '.')

		if fracPart < 10 {
			b = append(b, '0', '0')
		} else if fracPart < 100 {
			b = append(b, '0')
		}
		b = AppendUInt(b, fracPart)
		for len(b) > 0 && b[len(b)-1] == '0' {
			b = b[:len(b)-1]
		}
		if len(b) > 0 && b[len(b)-1] == '.' {
			b = b[:len(b)-1]
		}
	}
	return b
}

func AppendUInt(b []byte, v uint64) []byte {
	if v == 0 {
		return append(b, '0')
	}
	var temp [20]byte
	pos := len(temp)

	for v > 0 {
		pos--
		temp[pos] = byte(v%10) + '0'
		v /= 10
	}
	return append(b, temp[pos:]...)

}
