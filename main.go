package main

import (
	"net/http"
)

func main() {
	// f := &FetcherInput{
	// 	URL:         "http://server204:16686",
	// 	End:         0,
	// 	Start:       1587099027978000,
	// 	Limit:       10,
	// 	Lookback:    "",
	// 	MaxDuration: 0,
	// 	MinDuration: 0,
	// 	Service:     "gateway",
	// }

	// res, err := GatherTraceInfo(f, &FixedBookstoreTraceParser{})
	// if err != nil {
	// 	panic(err)
	// }

	controller := &Controller{}
	http.ListenAndServe(":9890", controller)

}
