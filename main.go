package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
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
		headers[strings.ToLower(header[0])] = header[1]
	}

	return headers
}

func parseBody(b []string) string {
	var body strings.Builder
	for _, val := range b {
		body.WriteString(strings.TrimSpace(val))
	}
	return body.String()
}

func parseHTTPRequest(r string) HTTPRequestSchema {
	message := strings.Split(strings.TrimSpace(r), "\n")
	httpRequestLine := parseRequestLine(message[0])
	var i int
	var v string
	var httpBody string
	for i, v = range message[1:] {
		if len(v) == 0 {
			httpBody = parseBody(message[i+2:])
			break
		}
	}

	httpHeaders := parseHeaders(message[1 : i+1])

	return HTTPRequestSchema{
		RequestLine: httpRequestLine,
		Headers:     httpHeaders,
		Body:        httpBody,
	}
}

func make_http_request(r HTTPRequestSchema) {

	client := &http.Client{
		Timeout: time.Second * 60,
	}

	host, ok := r.Headers["host"]

	var url string
	if ok {
		url = "https://" + strings.TrimSpace(host)
	} else {
		fmt.Println("Host header doesn't exist!")
	}

	var req *http.Request
	var err error

	if (r.Method == "POST") || (r.Method == "PUT") {
		req, err = http.NewRequest(r.Method, url+r.Path, strings.NewReader(r.Body))

		if err != nil {
			fmt.Println("Error while forming request!")
			return
		}
	} else {
		req, err = http.NewRequest(r.Method, url+r.Path, nil)

		if err != nil {
			fmt.Println("Error while forming request!")
			return
		}
	}

	for k, v := range r.Headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error")
		return
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	fmt.Println(string(bodyBytes))

}

func main() {
	httpRawRequest := `GET /posts HTTP/1.1
	Host: nihalchoudhary.in
	Origin: nihalchoudhary.in
	Cookie: session_token

	username=admin&password=admin
	`
	httpRequest := parseHTTPRequest(httpRawRequest)
	make_http_request(httpRequest)
}
