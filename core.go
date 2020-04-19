package main

// TracesInfo ...
type TracesInfo struct {
	Info  map[string][]map[string]interface{} `json:"info"`
	Stats map[string]map[string]*Stats        `json:"stats"`
}

// GatherTraceInfo ...
func GatherTraceInfo(f *FetcherInput, traceAggregator TraceAggregator) (*TracesInfo, error) {
	traceFetcher := &SimpleTraceFetcher{}
	traces, err := traceFetcher.FetchTraces(f)
	if err != nil {
		panic(err)
	}

	detailsMap := make(map[string][]map[string]interface{})

	for _, trace := range traces {
		if _, ok := detailsMap[trace.TraceType]; !ok {
			detailsMap[trace.TraceType] = make([]map[string]interface{}, 0)
		}
		details, err := traceAggregator.ParseTrace(trace)
		if err != nil {

		} else {
			detailsMap[trace.TraceType] = append(detailsMap[trace.TraceType], details)
		}
	}

	t := &TracesInfo{}
	t.Info = detailsMap
	t.Stats = make(map[string]map[string]*Stats)
	for key, traceInfo := range t.Info {
		t.Stats[key] = ComputeStats(traceInfo)
	}
	return t, nil
}
