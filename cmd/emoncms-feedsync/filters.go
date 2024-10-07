package main

import "encoding/json"

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

func ParseDataStr(datastr string) ([][]float64, error) {
	var parsed [][]float64
	err := json.Unmarshal([]byte(datastr), &parsed)
	if err != nil {
		return nil, err
	}
	return parsed, nil
}

func EncodeDataStr(data [][]float64) (string, error) {
	encoded, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}
