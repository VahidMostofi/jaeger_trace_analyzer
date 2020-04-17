package main

// GatherTraceInfo ...
func GatherTraceInfo(f *FetcherInput, traceAggregator TraceAggregator) (map[string][]map[string]interface{}, error) {
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

		}
		detailsMap[trace.TraceType] = append(detailsMap[trace.TraceType], details)
	}

	return detailsMap, nil
}
