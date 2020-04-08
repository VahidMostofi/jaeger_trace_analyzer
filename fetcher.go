package main

import (
	"io/ioutil"
	"net/http"
	"sync"
)

type Input struct {
	Id          string
	RequestType string
}

type Fetcher struct {
	TraceParser          TraceParser
	BookstoreTraceParser BookstoreTraceParser
}

func (f *Fetcher) fetchTrace(id, request_type string) map[string]interface{} {
	resp, err := http.Get("http://localhost:16686/api/traces/" + id)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	spans, _ := f.TraceParser.ParseTrace(data)
	res := make(map[string]interface{})
	if request_type == "login" {
		res, _ = f.BookstoreTraceParser.ParseLoginTrace(spans)
	} else if request_type == "getbook" {
		res, _ = f.BookstoreTraceParser.ParseGetBookTrace(spans)
	} else if request_type == "editbook" {
		res, _ = f.BookstoreTraceParser.ParseEditBookTrace(spans)
	}
	res["id"] = id
	res["request_type"] = request_type
	return res
}

func (f *Fetcher) StartFetching(results chan map[string]interface{}, inputs chan Input) {
	for i := range inputs {
		results <- f.fetchTrace(i.Id, i.RequestType)
	}
}

type Aggregator struct {
	AggregatedResult []map[string]interface{}
	Inputs           chan Input
	Results          chan map[string]interface{}
}

func (a *Aggregator) Start(inputs []Input, NumWorkers int) {
	for i := 0; i < NumWorkers; i++ {
		f := &Fetcher{
			TraceParser:          &SimpleTraceParser{},
			BookstoreTraceParser: &SimpleBookstoreTraceParser{},
		}
		go f.StartFetching(a.Results, a.Inputs)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for r := range a.Results {
			a.AggregatedResult = append(a.AggregatedResult, r)
			if len(a.AggregatedResult) == len(inputs) {
				break
			}
		}
		wg.Done()
	}()
	for _, i := range inputs {
		a.Inputs <- i
	}
	wg.Wait()
	return
}
