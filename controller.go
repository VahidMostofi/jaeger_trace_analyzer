package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Controller to provide an interface to other services use this
type Controller struct {
}

// Handle handles each http request
func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("err in reading request body")
			w.WriteHeader(400)
			return
		}

		f, err := NewFetcherInput(b)
		if err != nil {
			log.Println("err in parsing FetcherInput:", err)
			b, err = json.Marshal(err)
			if err != nil {
				panic(err)
			}
			w.Write(b)
			w.WriteHeader(500)
			return
		}

		res, err := GatherTraceInfo(f, &FixedBookstoreTraceParser{})
		if err != nil {
			log.Println(err)

			w.WriteHeader(500)
			return
		}

		b, err = json.Marshal(res)
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		w.WriteHeader(200)
		return
	}
}
