package httpserver

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"strconv"
	"strings"
)

const (
	StatusOK       = uint(200)
	StatusCreated  = uint(201)
	StatusNotFound = uint(404)
	httpVersion    = "HTTP/1.1"
)

var statusCodeToReasonPhrase = map[uint]string{
	StatusOK:       "OK",
	StatusNotFound: "Not Found",
	StatusCreated:  "Created",
}

type Response struct {
	statusCode uint
	headers    map[string]string
	body       string
	builder    *strings.Builder
}

func NewResponse(statusCode uint, body string) *Response {
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

func (r *Response) Send(request *Request) {
	fmt.Fprint(request.conn, r.String())
}

func (r *Response) String() string {
	r.buildStatusLine()

	if r.headers["Content-Encoding"] == "gzip" {
		r.body = gzipCompress(r.body)
	}

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

func gzipCompress(input string) string {
	var compressedData bytes.Buffer
	writer := gzip.NewWriter(&compressedData)
	_, err := writer.Write([]byte(input))
	if err != nil {
		return ""
	}
	err = writer.Close()
	if err != nil {
		return ""
	}
	return compressedData.String()
}
