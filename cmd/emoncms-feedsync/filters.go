package main

import (
	"fmt"
	"jpxor/emoncms/feedsync/pkg/stats"
	"jpxor/emoncms/feedsync/pkg/utils"
	"math"
	"strings"
)

type Filter func(data string) (string, error)

type FilterMap struct {
	filtermap map[string]Filter
}

func NewFilterMap(fconfigs []FilterConfig) (*FilterMap, error) {
	fm := FilterMap{
		filtermap: make(map[string]Filter),
	}
	for _, fconfig := range fconfigs {
		switch fconfig.Name {

		case "minmax":
			if len(fconfig.Args) != 2 {
				return nil, fmt.Errorf("minmax requires two args: [min, max]")
			}
			min := fconfig.Args[0]
			max := fconfig.Args[1]
			for _, feedName := range fconfig.Feeds {
				fm.Add(feedName, MinMaxFilter(min, max))
			}

		case "median":
			if len(fconfig.Args) != 2 {
				return nil, fmt.Errorf("median requires two args: [window size, max offset from median]")
			}
			window := int(fconfig.Args[0])
			maxoff := fconfig.Args[1]
			for _, feedName := range fconfig.Feeds {
				fm.Add(feedName, MedianFilter(window, maxoff))
			}

		default:
			return nil, fmt.Errorf("unknown filter name: %s", fconfig.Name)
		}
	}
	return &fm, nil
}

func (m FilterMap) Apply(name, data string) (string, error) {
	if filter, ok := m.filtermap[name]; ok {
		return filter(data)
	}
	return data, nil
}

func (m FilterMap) Add(name string, filter Filter) {
	m.filtermap[name] = filter
}

type DataPoint struct {
	Timestamp int64
	Value     float32
}

func ParseDataStr(datastr string) ([]DataPoint, error) {
	capacity := (1 + strings.Count(datastr, ",")) / 2
	if capacity == 0 {
		return []DataPoint{}, nil
	}
	dataPoints := make([]DataPoint, 0, capacity)

	// Remove outer brackets
	datastr = datastr[2 : len(datastr)-1]

	for len(datastr) > 0 {
		idx := strings.IndexByte(datastr, ',')
		if idx == -1 {
			break
		}
		timestamp, err := utils.ParseUnixTimestamp(datastr[:idx])
		if err != nil {
			return nil, err
		}
		datastr = datastr[idx+1:]
		idx = strings.IndexByte(datastr, ']')
		if idx == -1 {
			break
		}
		if datastr[:idx] != "null" {
			value, err := utils.ParseFloatLowPrecision(datastr[:idx])
			if err != nil {
				return nil, err
			}
			dataPoints = append(dataPoints, DataPoint{Timestamp: timestamp, Value: float32(value)})
		}
		if len(datastr) > idx+3 {
			datastr = datastr[idx+3:]
		} else {
			break
		}
	}
	return dataPoints, nil
}

func EncodeDataStr(data []DataPoint) string {
	if len(data) == 0 {
		return "[]"
	}
	// Estimate size: 4 bytes for "[[]]", 3 bytes for each separator "],[",
	// and approximately 20 bytes per DataPoint
	estimatedSize := 4 + (len(data)-1)*3 + len(data)*20
	buf := make([]byte, 0, estimatedSize)

	buf = append(buf, '[', '[')
	for i, dp := range data {
		if i > 0 {
			buf = append(buf, ']', ',', '[')
		}
		buf = utils.AppendUInt(buf, uint64(dp.Timestamp))
		buf = append(buf, ',')
		buf = utils.AppendFloat(buf, dp.Value)
	}
	buf = append(buf, ']', ']')
	return string(buf)
}

func extractValues(window []DataPoint) []float32 {
	values := make([]float32, len(window))
	for i, dp := range window {
		values[i] = dp.Value
	}
	return values
}

func MinMaxFilter(min, max float32) Filter {
	test := func(val float32) bool {
		return val >= min && val <= max
	}
	return func(datastr string) (string, error) {
		data, err := ParseDataStr(datastr)
		if err != nil {
			return datastr, err
		}
		var filteredData []DataPoint
		for _, datapoint := range data {
			if test(datapoint.Value) {
				filteredData = append(filteredData, datapoint)
			}
		}
		return EncodeDataStr(filteredData), nil
	}
}

func MedianFilter(windowSize int, maxdiff float32) Filter {
	test := func(median, val float32) bool {
		return math.Abs(float64(median-val)) <= float64(maxdiff)
	}
	return func(datastr string) (string, error) {
		data, err := ParseDataStr(datastr)
		if err != nil {
			return datastr, err
		}
		var filteredData []DataPoint
		for i, datapoint := range data {
			left := max(0, i-windowSize/2)
			right := min(i+windowSize/2, len(data))
			window := data[left:right]

			median := stats.MedianValue(extractValues(window))
			if test(median, datapoint.Value) {
				filteredData = append(filteredData, datapoint)
			}
		}
		return EncodeDataStr(filteredData), nil
	}
}
