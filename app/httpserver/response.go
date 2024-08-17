package httpserver

import (
	"fmt"
	"net"
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

type response struct {
	statusCode int
	headers    map[string]string
	body       string
}

func getResponse(statusCode int, body string) *response {
	headers := make(map[string]string)

	if body != "" {
		headers["Content-Type"] = "text/plain"
		headers["Content-Length"] = strconv.Itoa(len(body))
	}

	return &response{
		statusCode: statusCode,
		headers:    headers,
		body:       body,
	}
}

func getResponseString(response *response) string {
	var s strings.Builder

	buildStatusLine(&s, response.statusCode)

	if response.headers != nil {
		buildHeaders(&s, response)
	}

	if response.body != "" {
		s.WriteString(response.body)
	}

	return s.String()
}

func respond(c net.Conn, s string) {
	fmt.Fprint(c, s)
}

func buildDelineator(builder *strings.Builder) {
	builder.WriteString("\r\n")
}

func buildHeaders(builder *strings.Builder, response *response) {
	response.headers["Content-Length"] = strconv.Itoa(len(response.body))

	if response.headers["Content-Type"] == "application/octet-stream" {
		buildOctetStreamHeaders(builder, response)
	} else {
		buildPlainTextHeaders(builder, response)
	}
}

func buildOctetStreamHeaders(builder *strings.Builder, response *response) {
	for k, v := range response.headers {
		builder.WriteString(fmt.Sprintf("%s: %s", k, v))
		buildDelineator(builder)
	}

	buildDelineator(builder)
}

func buildPlainTextHeaders(builder *strings.Builder, response *response) {
	for k, v := range response.headers {
		builder.WriteString(fmt.Sprintf("%s: %s", k, v))
		buildDelineator(builder)
	}

	buildDelineator(builder)
}

func buildStatusLine(builder *strings.Builder, statusCode int) {
	builder.WriteString(fmt.Sprintf(
		"%s %d %s",
		httpVersion,
		statusCode,
		statusCodeToReasonPhrase[statusCode],
	))
	buildDelineator(builder)
}
