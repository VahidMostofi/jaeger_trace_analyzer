package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {
	port := flag.String("port number", "8657", "port to listent to")
	fmt.Println("server started and listening to port", *port)
	c := &Controller{}
	http.ListenAndServe(":"+*port, c)
}
