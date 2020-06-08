Jaeger trace analyzer

I developed this tool to use it in my research at University of Calgary. It is supposed to gather, aggreagte, process and analyse the Trace Data from a Jaeger server. At the moment it only works with a specific application (I will make it general as soon as I can).
## How to run
### Without Docker
```
go run *.go
```
### With Docker
* Bulding the Docker image:
```
sudo docker build . -t trace_analyser:latest
```
* If the jaeger server is located on the local host, start the container using
```
sudo docker run --rm --network host --name trace_analyser -d trace_analyser:latest
```
* otherwise use:
```
sudo docker run --rm -p 8657:8657 --name trace_analyser -d trace_analyser:latest
```
## How to use
You need to send http post requests to port which the server is listening to.
In this example, the Jaeger Trace Analyzer is running on localhost, port 8657.
XXXX is the ip for the Jaeger server
PPP is the Jaeger UI port
```
curl --location --request POST 'http://localhost:8657' \
--header 'Content-Type: application/json' \
--header 'Content-Type: text/plain' \
--data-raw '{
	"url": "http://XXXX:PPP",
	"start": 1587598229149861,
	"end": 0,
	"limit": 1000000, 
	"lookback" : "",
	"maxDuration" : 0,
	"minDuration" : 0,
	"service" : "gateway"
}'
```
## ToDo
- [ ] Better Readme
- [ ] Configure port number using environment variables and command line
- [ ] Make it more general! Now it only works with Bookstore application
- [ ] Configuration and possiblity on how to get the data, instead of having one static approach
- [ ] Use RPC instead of http request
- [ ] Better project
