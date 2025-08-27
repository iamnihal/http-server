package main

import (
	"fmt"
	"strings"
)

type HTTPRequestSchema struct {
	RequestLine
	Headers map[string]string
	Body    string
}

type RequestLine struct {
	Method      string
	Path        string
	HTTPVersion string
}

const (
	GET     = "GET"
	HEAD    = "HEAD"
	POST    = "POST"
	PUT     = "PUT"
	DELETE  = "DELETE"
	CONNECT = "CONNECT"
	OPTIONS = "OPTOINS"
	TRACE   = "TRACE"
)

var allowedHTTPMethods = map[string]string{
	"GET":     GET,
	"HEAD":    HEAD,
	"POST":    POST,
	"PUT":     PUT,
	"DELETE":  DELETE,
	"CONNECT": CONNECT,
	"OPTIONS": OPTIONS,
	"TRACE":   TRACE,
}

func parseRequestLine(rl string) RequestLine {
	requestLine := strings.Split(rl, " ")
	httpMethod := requestLine[0]
	path := requestLine[1]
	httpVersion := requestLine[2]

	if _, ok := allowedHTTPMethods[httpMethod]; !ok {
		fmt.Println("Invalid method")
	}

	if (httpVersion != "HTTP/1.0") && (httpVersion != "HTTP/1.1") && (httpVersion != "HTTP/2") && (httpVersion != "HTTP/3") {
		fmt.Println("Invalid HTTP version!")
	}

	if len(strings.TrimSpace(path)) == 0 {
		fmt.Println("Invalid path")
	}

	// fmt.Printf("HTTP Method: %s\nPath: %s\nHTTP Version: %s\n", httpMethod, path, httpVersion)

	return RequestLine{
		Method:      httpMethod,
		Path:        path,
		HTTPVersion: httpVersion,
	}
}

func parseHeaders(h []string) map[string]string {

	headers := make(map[string]string)

	for _, val := range h {
		header := strings.Split(strings.TrimSpace(val), ":")
		headers[header[0]] = header[1]
		headers[header[0]] = header[1]
	}

	return headers
}

func parseBody(b []string) {
	// for key, val := range b {
	// 	fmt.Println(key, val)
	// }
}

func parseHTTPRequest(r string) HTTPRequestSchema {
	message := strings.Split(strings.TrimSpace(r), "\n")
	httpRequestLine := parseRequestLine(message[0])
	var i int
	var v string
	for i, v = range message[1:] {
		if len(v) == 0 {
			parseBody(message[i+1:])
			break
		}
	}

	hostHeaders := parseHeaders(message[1 : i+1])

	return HTTPRequestSchema{
		RequestLine: httpRequestLine,
		Headers:     hostHeaders,
	}

}

func make_http_request(r HTTPRequestSchema) {
	fmt.Println(r)
}

func main() {
	httpRawRequest := `GET /to/some/path HTTP/1.1
	Host: nihalchoudhary.in
	Origin: nihalchoudhary.in
	Cookie: session_token

	username=admin&password=admin
	`
	httpRequest := parseHTTPRequest(httpRawRequest)

	fmt.Println(httpRequest.RequestLine.Method)
	fmt.Println(httpRequest.RequestLine.Path)
	fmt.Println(httpRequest.RequestLine.HTTPVersion)
	fmt.Println(httpRequest.Headers)
	make_http_request(httpRequest)
}
