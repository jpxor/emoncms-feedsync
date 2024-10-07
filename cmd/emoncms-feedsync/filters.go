package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Filter func(data string) (string, error)

type FilterMap struct {
	filtermap map[string]Filter
}

func NewFilterMap() FilterMap {
	return FilterMap{
		filtermap: make(map[string]Filter),
	}
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
	dataPoints := make([]DataPoint, 0, capacity)

	// Remove outer brackets
	datastr = datastr[2 : len(datastr)-2]

	for len(datastr) > 0 {
		idx := strings.IndexByte(datastr, ',')
		if idx == -1 {
			break
		}
		timestamp, err := strconv.ParseInt(datastr[:idx], 10, 64)
		if err != nil {
			return nil, err
		}
		datastr = datastr[idx+1:]
		idx = strings.IndexByte(datastr, ']')
		if idx == -1 {
			break
		}
		value, err := strconv.ParseFloat(datastr[:idx], 32)
		if err != nil {
			return nil, err
		}
		dataPoints = append(dataPoints, DataPoint{Timestamp: timestamp, Value: float32(value)})

		if len(datastr) > idx+3 {
			datastr = datastr[idx+3:]
		} else {
			break
		}
	}
	return dataPoints, nil
}

func EncodeDataStr(data []DataPoint) (string, error) {
	var sb strings.Builder

	// Estimate size: 4 bytes for "[[]]", 3 bytes for each separator "],[",
	// and approximately 20 bytes per DataPoint (assuming 10 digits for timestamp and 8 for float)
	estimatedSize := 4 + (len(data)-1)*3 + len(data)*20
	sb.Grow(estimatedSize)

	sb.WriteString("[[")

	for i, dp := range data {
		if i > 0 {
			sb.WriteString("],[")
		}
		sb.WriteString(fmt.Sprintf("%d,%g", dp.Timestamp, dp.Value))
	}

	sb.WriteString("]]")
	return sb.String(), nil

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
		encoded, err := EncodeDataStr(filteredData)
		if err != nil {
			return datastr, err
		}
		return encoded, nil
	}
}
