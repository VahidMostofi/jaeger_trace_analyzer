curl --location --request POST 'http://localhost:8657' \
--header 'Content-Type: application/json' \
--header 'Content-Type: text/plain' \
--data-raw '{
	"url": "http://server204:16686",
	"start": 1587598229149861,
	"end": 0,
	"limit": 1000000, 
	"lookback" : "",
	"maxDuration" : 0,
	"minDuration" : 0,
	"service" : "gateway"
}'
