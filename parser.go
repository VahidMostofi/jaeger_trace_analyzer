package main

import (
	"encoding/json"
)

type Span struct {
	StartTime float64 `json:"startTime"`
	Duration  float64 `json:"duration"`
}

type TraceParser interface {
	ParseTrace(data []byte) (map[string]Span, error)
}

type SimpleTraceParser struct{}

func (s *SimpleTraceParser) ParseTrace(data []byte) (map[string]Span, error) {
	r := &struct {
		Data []map[string]interface{} `json:"data"`
	}{}
	json.Unmarshal(data, r)

	spans := map[string]Span{}
	for _, v := range r.Data[0]["spans"].([]interface{}) {
		v2 := v.(map[string]interface{})
		s := Span{
			StartTime: v2["startTime"].(float64),
			Duration:  v2["duration"].(float64),
		}
		operationName := v2["operationName"].(string)
		spans[operationName] = s
	}
	return spans, nil
}

type BookstoreTraceParser interface {
	ParseLoginTrace(spans map[string]Span) (map[string]interface{}, error)
	ParseGetBookTrace(spans map[string]Span) (map[string]interface{}, error)
	ParseEditBookTrace(spans map[string]Span) (map[string]interface{}, error)
}

type SimpleBookstoreTraceParser struct{}

func (b *SimpleBookstoreTraceParser) ParseLoginTrace(spans map[string]Span) (map[string]interface{}, error) {
	total := spans["auth_login"].Duration
	// gateway2auth := spans["login"].StartTime - spans["auth_service_login"].StartTime
	getService := spans["auth_req_login"].Duration
	connectToService := spans["auth_connect"].Duration
	inAuth := spans["login"].Duration
	inAuthDB := spans["queryDB"].Duration
	inAuthGenJWT := spans["generateAuthToken"].Duration
	waitTime := spans["login"].StartTime - (spans["auth_connect"].StartTime + spans["auth_connect"].Duration)
	// auth2gateway := spans["auth_login"].Duration - spans["auth_service_login"].Duration - (spans["auth_service_login"].StartTime - spans["auth_login"].StartTime)
	info := make(map[string]interface{})
	info["total"] = total
	info["getService"] = getService
	info["inAuth"] = inAuth
	info["inAuthDB"] = inAuthDB
	info["inAuthGenJWT"] = inAuthGenJWT
	info["connectToService"] = connectToService
	info["waitTime"] = waitTime
	return info, nil
}

func (b *SimpleBookstoreTraceParser) ParseGetBookTrace(spans map[string]Span) (map[string]interface{}, error) {
	total := spans["get_book"].Duration
	authenticate := spans["authenticate"].Duration
	getService := spans["books_get_book"].Duration
	// gateway2books := spans["getone"].StartTime - spans["book_service_getone"].StartTime
	connectToService := spans["books_connect"].Duration
	inBooks := spans["getone"].Duration
	inBooksDB := spans["DB"].Duration
	// books2gateway := spans["get_book"].Duration - spans["book_service_getone"].Duration - (spans["book_service_getone"].StartTime - spans["get_book"].StartTime)
	waitTime := spans["getone"].StartTime - (spans["books_connect"].Duration + spans["books_connect"].StartTime)
	info := make(map[string]interface{})
	info["total"] = total
	info["authenticate"] = authenticate
	info["connectToService"] = connectToService
	info["inBooks"] = inBooks
	info["inBooksDB"] = inBooksDB
	info["waitTime"] = waitTime
	info["getService"] = getService
	return info, nil
}

func (b *SimpleBookstoreTraceParser) ParseEditBookTrace(spans map[string]Span) (map[string]interface{}, error) {
	total := spans["update_book"].Duration
	authenticate := spans["authenticate"].Duration
	getService := spans["books_edit_book"].Duration
	// gateway2books := spans["update"].StartTime - spans["book_service_edit"].StartTime
	connectToService := spans["books_connect"].Duration
	inBooks := spans["update"].Duration
	inBooksDB := spans["DB"].Duration
	// books2gateway := spans["update_book"].Duration - spans["book_service_edit"].Duration - (spans["book_service_edit"].StartTime - spans["update_book"].StartTime)
	waitTime := spans["update"].StartTime - (spans["books_connect"].Duration + spans["books_connect"].StartTime)
	info := make(map[string]interface{})
	info["total"] = total
	info["authenticate"] = authenticate
	info["getService"] = getService
	info["inBooks"] = inBooks
	info["inBooksDB"] = inBooksDB
	info["connectToService"] = connectToService
	info["waitTime"] = waitTime
	return info, nil
}
