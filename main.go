package main

import (
	"fmt"
	"net/http"
	"strings"
)

type Controller struct {
	Aggregator *Aggregator
}

func (c *Controller) Handler(w http.ResponseWriter, r *http.Request) {
	ip := GetIP(r)
	if !(strings.HasPrefix(ip, "50.99.77.228") || strings.HasPrefix(ip, "[::1]")) {
		c.Logger.Printf("request came from: %s , rejected\n", ip)
		w.WriteHeader(403)
		return
	}
	if r.URL.EscapedPath() == "/start"{
		c := &Controller{
			Aggregator : Aggregator{
				AggregatedResult: make([]map[string]interface{}, 0),
				Inputs:           make(chan Input),
				Results:          make(chan map[string]interface{}),
			}
		}
	}
}

func main() {
	port := flag.String("port number", "8657", "port to listent to")
	fmt.Println("server started and listening to port", *port)
	http.ListenAndServe(":"+*port, http.HandlerFunc(c.Handler))
	// input := []Input{{"12bf2fda8ac13441", "getbook"}, {"d320d5903dcba868", "editbook"}, {"bd4c08790c556854", "login"}}
	// a.Start(input, 12)
}
