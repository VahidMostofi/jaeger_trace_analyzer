package main

import (
	"fmt"
	"log"
)

// FixedBookstoreTraceParser parses one trace of bookstore, this will be deprecated soon
type FixedBookstoreTraceParser struct{}

// ParseTrace is the implemeation of ParseTrace from TraceParser interface
func (b *FixedBookstoreTraceParser) ParseTrace(trace *Trace) (map[string]interface{}, error) {

	var parseInfo map[string]interface{}
	var err error
	if len(trace.SpansMap) < 3 {
		return nil, fmt.Errorf("Wrong request!")
	}
	switch trace.TraceType {
	case "auth_login":
		parseInfo, err = b.parseLoginTrace(trace.SpansMap)
		break
	case "get_book":
		parseInfo, err = b.parseGetBookTrace(trace.SpansMap)
		break
	case "update_book":
		parseInfo, err = b.parseEditBookTrace(trace.SpansMap)
		break
	}
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return parseInfo, nil
}

// parseLoginTrace ...
func (b *FixedBookstoreTraceParser) parseLoginTrace(spans map[string]*Span) (map[string]interface{}, error) {
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

// parseGetBookTrace ...
func (b *FixedBookstoreTraceParser) parseGetBookTrace(spans map[string]*Span) (map[string]interface{}, error) {
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

// parseEditBookTrace ...
func (b *FixedBookstoreTraceParser) parseEditBookTrace(spans map[string]*Span) (map[string]interface{}, error) {
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
