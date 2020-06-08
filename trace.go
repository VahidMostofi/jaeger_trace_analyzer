package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"

	"net/http"
)

// Trace struct contains a group of spans
type Trace struct {
	Spans     []*Span          `json:"spans"`
	SpansMap  map[string]*Span //map from operationName to Span object
	TraceType string           `json:"-"`
	HasRoot   bool             `json:"-"`
	TraceID   string           `json:"traceID"`
}

// Span struct contains information about each span
type Span struct {
	StartTime     float64 `json:"startTime"`
	Duration      float64 `json:"duration"`
	OperationName string  `json:"operationName"`
	SpanID        string  `json:"spanID"`
	TraceID       string  `json:"traceID"`
	IsRoot        bool    `json:"-"`
}

// FetcherInput specifies the properties of the get request to the Jaeger
type FetcherInput struct {
	URL         string `json:"url"` // example of a valid url : http://server204:16686
	End         uint64 `json:"end"`
	Start       uint64 `json:"start"`
	Limit       int    `json:"limit"`
	Lookback    string `json:"lookback"`
	MaxDuration uint64 `json:"maxDuration"`
	MinDuration uint64 `json:"minDuration"`
	Service     string `json:"service"`
}

// TraceFetcher the interface which is in charge of getting all traces based on the input from Jaeger
type TraceFetcher interface {
	FetchTraces(input *FetcherInput) ([]*Trace, error)
}

// NewFetcherInput creates a new FetcherInput object from byte array object
func NewFetcherInput(b []byte) (*FetcherInput, error) {
	f := &FetcherInput{}
	err := json.Unmarshal(b, f)
	if err != nil {
		return nil, fmt.Errorf("unable to parse FetcherInput json: %w", err)
	}
	var validURL = regexp.MustCompile(`^https?://`)
	if !validURL.MatchString(f.URL) {
		return nil, fmt.Errorf("make sure the url contains http:// or https://")
	}
	return f, nil
}

// SimpleTraceFetcher is the implementation of TraceFetcher
type SimpleTraceFetcher struct{}

// FetchTraces fetches traces based on the input
func (t *SimpleTraceFetcher) FetchTraces(input *FetcherInput) ([]*Trace, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", input.URL+"/api/traces", nil)
	if err != nil {
		return nil, fmt.Errorf("error in creating the request to fetch: %w", err)
	}

	q := req.URL.Query()
	if input.End != 0 {
		q.Add("end", strconv.FormatUint(input.End, 10))
	}
	if input.Start != 0 {
		q.Add("start", strconv.FormatUint(input.Start, 10))
	}
	if len(input.Service) > 0 {
		q.Add("service", input.Service)
	}
	if input.Limit > 0 {
		q.Add("limit", strconv.Itoa(input.Limit))
	}
	if len(input.Lookback) > 0 {
		q.Add("lookback", input.Lookback)
	}
	if input.MaxDuration > 0 {
		q.Add("maxDuration", strconv.FormatUint(input.MaxDuration, 10))
	}
	if input.MinDuration > 0 {
		q.Add("minDuration", strconv.FormatUint(input.MinDuration, 10))
	}

	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error in executing the request to fetch: %w", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	results := &struct {
		Data []*Trace `json:"data"`
	}{}

	json.Unmarshal(b, results)

	// create spansMap from spans
	for _, trace := range results.Data {
		trace.SpansMap = make(map[string]*Span)
		for _, span := range trace.Spans {
			trace.SpansMap[span.OperationName] = span
			if span.SpanID == span.TraceID {
				span.IsRoot = true
				trace.TraceType = span.OperationName
				trace.HasRoot = true
			}
		}
	}
	return results.Data, nil
}
