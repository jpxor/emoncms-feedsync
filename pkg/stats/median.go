package stats

import "sort"

func MedianValue(vals []float32) float32 {
	if len(vals) == 0 {
		return 0
	}
	// Sort the slice
	sortedVals := make([]float32, len(vals))
	copy(sortedVals, vals)
	sort.Slice(sortedVals, func(i, j int) bool {
		return sortedVals[i] < sortedVals[j]
	})
	// Find the median
	mid := len(sortedVals) / 2
	if len(sortedVals)%2 == 0 {
		return (sortedVals[mid-1] + sortedVals[mid]) / 2
	}
	return sortedVals[mid]
}
