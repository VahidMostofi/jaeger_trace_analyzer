package main

import (
	"encoding/json"
	"strings"
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
		if strings.HasPrefix(operationName, "post: http://auth:3007/auth/login") {
			operationName = "auth_service_login"
		} else if strings.HasPrefix(operationName, "get: http://books:3009/books/") {
			if len(operationName) == len("get: http://books:3009/books/") {
				operationName = "book_service_list"
			} else if len(operationName) > len("get: http://books:3009/books/") {
				operationName = "book_service_getone"
			}
		} else if strings.HasPrefix(operationName, "put: http://books:3009/books/") {
			operationName = "book_service_edit"
		}
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
	gateway2auth := spans["login"].StartTime - spans["auth_service_login"].StartTime
	inAuth := spans["login"].Duration
	inAuthDB := spans["queryDB"].Duration
	inAuthGenJWT := spans["generateAuthToken"].Duration
	auth2gateway := spans["auth_login"].Duration - spans["auth_service_login"].Duration - (spans["auth_service_login"].StartTime - spans["auth_login"].StartTime)
	info := make(map[string]interface{})
	info["total"] = total
	info["gateway2auth"] = gateway2auth
	info["inAuth"] = inAuth
	info["inAuthDB"] = inAuthDB
	info["inAuthGenJWT"] = inAuthGenJWT
	info["auth2gateway"] = auth2gateway
	return info, nil
}

func (b *SimpleBookstoreTraceParser) ParseGetBookTrace(spans map[string]Span) (map[string]interface{}, error) {
	total := spans["get_book"].Duration
	authenticate := spans["authenticate"].Duration
	gateway2books := spans["getone"].StartTime - spans["book_service_getone"].StartTime
	inBooks := spans["getone"].Duration
	inBooksDB := spans["DB"].Duration
	books2gateway := spans["get_book"].Duration - spans["book_service_getone"].Duration - (spans["book_service_getone"].StartTime - spans["get_book"].StartTime)
	info := make(map[string]interface{})
	info["total"] = total
	info["authenticate"] = authenticate
	info["gateway2books"] = gateway2books
	info["inBooks"] = inBooks
	info["inBooksDB"] = inBooksDB
	info["books2gateway"] = books2gateway
	return info, nil
}

func (b *SimpleBookstoreTraceParser) ParseEditBookTrace(spans map[string]Span) (map[string]interface{}, error) {
	total := spans["update_book"].Duration
	authenticate := spans["authenticate"].Duration
	gateway2books := spans["update"].StartTime - spans["book_service_edit"].StartTime
	inBooks := spans["update"].Duration
	inBooksDB := spans["DB"].Duration
	books2gateway := spans["update_book"].Duration - spans["book_service_edit"].Duration - (spans["book_service_edit"].StartTime - spans["update_book"].StartTime)
	info := make(map[string]interface{})
	info["total"] = total
	info["authenticate"] = authenticate
	info["gateway2books"] = gateway2books
	info["inBooks"] = inBooks
	info["inBooksDB"] = inBooksDB
	info["books2gateway"] = books2gateway
	return info, nil
}
