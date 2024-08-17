package httpserver

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const httpVersion = "HTTP/1.1"

var statusCodeToReasonPhrase = map[int]string{
	http.StatusOK:       "OK",
	http.StatusNotFound: "Not Found",
	http.StatusCreated:  "Created",
}

type Response struct {
	statusCode int
	headers    map[string]string
	body       string
	builder    *strings.Builder
}

func NewResponse(statusCode int, body string) *Response {
	headers := make(map[string]string)
	var builder strings.Builder

	if body != "" {
		headers["Content-Type"] = "text/plain"
		headers["Content-Length"] = strconv.Itoa(len(body))
	}

	return &Response{
		statusCode: statusCode,
		headers:    headers,
		body:       body,
		builder:    &builder,
	}
}

func (r *Response) Send(request *request) {
	fmt.Fprint(request.conn, r.String())
}

func (r *Response) String() string {
	r.buildStatusLine()

	if r.headers != nil {
		r.buildHeaders()
	}

	if r.body != "" {
		r.builder.WriteString(r.body)
	}

	return r.builder.String()
}

func (r *Response) buildDelineator() {
	r.builder.WriteString("\r\n")
}

func (r *Response) buildHeaders() {
	r.headers["Content-Length"] = strconv.Itoa(len(r.body))

	if r.headers["Content-Type"] == "application/octet-stream" {
		r.buildOctetStreamHeaders()
	} else {
		r.buildPlainTextHeaders()
	}
}

func (r *Response) buildOctetStreamHeaders() {
	for k, v := range r.headers {
		r.builder.WriteString(fmt.Sprintf("%s: %s", k, v))
		r.buildDelineator()
	}

	r.buildDelineator()
}

func (r *Response) buildPlainTextHeaders() {
	for k, v := range r.headers {
		r.builder.WriteString(fmt.Sprintf("%s: %s", k, v))
		r.buildDelineator()
	}

	r.buildDelineator()
}

func (r *Response) buildStatusLine() {
	r.builder.WriteString(fmt.Sprintf(
		"%s %d %s",
		httpVersion,
		r.statusCode,
		statusCodeToReasonPhrase[r.statusCode],
	))
	r.buildDelineator()
}
