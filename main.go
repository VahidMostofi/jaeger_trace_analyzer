package main

import (
	"fmt"
	"net/http"
	"strings"
	"flag"
	"encoding/json"
	"io/ioutil"
)

type Controller struct {
	Aggregator *Aggregator
}

type RequsetInput struct{
	Data []Input `json:"data"`
}



// GetIP gets a requests IP address by reading off the forwarded-for
// header (for proxies) and falls back to use the remote address.
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
func (c *Controller) Handler(w http.ResponseWriter, r *http.Request) {
	ip := GetIP(r)
	if !(strings.HasPrefix(ip, "50.99.77.228") || strings.HasPrefix(ip, "[::1]")) {
		fmt.Printf("request came from: %s , rejected\n", ip)
		w.WriteHeader(403)
		return
	}
	if r.URL.EscapedPath() == "/start"{
		if c.Aggregator == nil{
			c.Aggregator = &Aggregator{
				AggregatedResult: make([]map[string]interface{}, 0),
				Inputs:           make(chan Input),
				Results:          make(chan map[string]interface{}),
			}
			
			data := RequsetInput{}
			b, e := ioutil.ReadAll(r.Body)
			if e!=nil{
				panic(e)
			}
			err := json.Unmarshal(b, &data)
			if err != nil{
				panic(err)
			}
			input := data.Data
			fmt.Println("stared")
			c.Aggregator.Start(input, 5)
			b, e = json.Marshal(c.Aggregator.AggregatedResult)
			if e != nil{
				panic(e)
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(b)
			c.Aggregator = nil
			fmt.Println("done")
		}else{
			w.WriteHeader(400)
		}		
	} else if r.URL.EscapedPath() == "/ready"{
		w.WriteHeader(200)
		return
	}
}

func main() {
	port := flag.String("port number", "8657", "port to listent to")
	fmt.Println("server started and listening to port", *port)
	c := &Controller{}
	http.ListenAndServe(":"+*port, http.HandlerFunc(c.Handler))
}
