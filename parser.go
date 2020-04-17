package main

// TraceAggregator has ParseTrace which parses one trace to a map
type TraceAggregator interface {
	ParseTrace(trace *Trace) (map[string]interface{}, error)
}
