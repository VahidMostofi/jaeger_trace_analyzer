package main

import "github.com/montanaflynn/stats"

// Stats stores stat information for each endpoint
type Stats struct {
	Average      float64 `json:average"`
	Percentile90 float64 `json:percentile90"`
	Percentile95 float64 `json:percentile95"`
	Percentile99 float64 `json:percentile99"`
}

// ComputeStats ...
func ComputeStats(traceInfo []map[string]interface{}) map[string]*Stats {
	temp := make(map[string][]float64)

	for key := range traceInfo[0] {
		if key != "startTime" && key != "duration" {
			temp[key] = make([]float64, len(traceInfo))
		}
	}

	for i := 0; i < len(traceInfo); i++ {
		for key := range temp {
			temp[key][i] = traceInfo[i][key].(float64)
		}
	}

	res := make(map[string]*Stats)
	for key := range temp {
		avg, _ := stats.Mean(temp[key])
		p90, _ := stats.Percentile(temp[key], 90)
		p95, _ := stats.Percentile(temp[key], 95)
		p99, _ := stats.Percentile(temp[key], 99)
		stats := &Stats{
			Average:      avg,
			Percentile90: p90,
			Percentile95: p95,
			Percentile99: p99,
		}
		res[key] = stats
	}
	return res
}
