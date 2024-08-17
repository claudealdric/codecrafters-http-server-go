package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	port        = 4221
	httpVersion = "HTTP/1.1"
)

var statusCodeToReasonPhrase = map[int]string{
	http.StatusOK:       "OK",
	http.StatusNotFound: "Not Found",
}
var config Config

func main() {
	fmt.Println("Logs from your program will appear here!")

	parseFlags()

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Printf("Failed to bind to port %d\n", port)
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			os.Exit(1)
		}

		go processRequest(conn)
	}
}

func parseFlags() {
	flag.StringVar(
		&config.directory,
		"directory",
		"",
		"Specify the directory where files are stored",
	)
	flag.Parse()
}

func processRequest(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request:", err.Error())
		return
	}

	headers, err := getHeaders(reader)
	if err != nil {
		fmt.Println("Error reading headers:", err.Error())
		return
	}

	path := extractPath(requestLine)
	routeRequest(conn, path, headers)
}

func getHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if line == "\r\n" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.ToLower(strings.TrimSpace(parts[0]))
			value := (strings.TrimSpace(parts[1]))
			headers[key] = value
		}
	}

	return headers, nil
}

func extractPath(requestLine string) string {
	return strings.Fields(requestLine)[1]
}

func routeRequest(conn net.Conn, path string, headers map[string]string) {
	switch {
	case path == "/":
		handleRoot(conn)
	case strings.HasPrefix(path, "/echo/"):
		handleEcho(conn, path)
	case strings.HasPrefix(path, "/user-agent"):
		handleUserAgent(conn, headers)
	case strings.HasPrefix(path, "/files/"):
		handleFiles(conn, path)
	default:
		handleNotFound(conn)
	}
}

func handleRoot(conn net.Conn) {
	response := getResponse(http.StatusOK, "")
	respond(conn, getResponseString(response))
}

func handleNotFound(conn net.Conn) {
	response := getResponse(http.StatusNotFound, "")
	respond(conn, getResponseString(response))
}

func handleEcho(conn net.Conn, path string) {
	echoArg := strings.TrimPrefix(path, "/echo/")
	response := getResponse(http.StatusOK, echoArg)
	respond(conn, getResponseString(response))
}

func handleUserAgent(conn net.Conn, headers map[string]string) {
	userAgent := headers["user-agent"]
	response := getResponse(http.StatusOK, userAgent)
	respond(conn, getResponseString(response))
}

func handleFiles(conn net.Conn, path string) {
	fileName := strings.TrimPrefix(path, "/files/")
	content, err := os.ReadFile(filepath.Join(config.directory, fileName))
	if err != nil && os.IsNotExist(err) {
		handleNotFound(conn)
		return
	}
	response := getResponse(http.StatusOK, string(content))
	response.headers["Content-Type"] = "application/octet-stream"
	respond(conn, getResponseString(response))
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

func buildDelineator(builder *strings.Builder) {
	builder.WriteString("\r\n")
}

func buildHeaders(builder *strings.Builder, response *Response) {
	response.headers["Content-Length"] = strconv.Itoa(len(response.body))

	if response.headers["Content-Type"] == "application/octet-stream" {
		buildOctetStreamHeaders(builder, response)
	} else {
		buildPlainTextHeaders(builder, response)
	}
}

func buildPlainTextHeaders(builder *strings.Builder, response *Response) {
	for k, v := range response.headers {
		builder.WriteString(fmt.Sprintf("%s: %s", k, v))
		buildDelineator(builder)
	}

	buildDelineator(builder)
}

func buildOctetStreamHeaders(builder *strings.Builder, response *Response) {
	for k, v := range response.headers {
		builder.WriteString(fmt.Sprintf("%s: %s", k, v))
		buildDelineator(builder)
	}

	buildDelineator(builder)
}

func getResponseString(response *Response) string {
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

func getResponse(statusCode int, body string) *Response {
	headers := make(map[string]string)

	if body != "" {
		headers["Content-Type"] = "text/plain"
		headers["Content-Length"] = strconv.Itoa(len(body))
	}

	return &Response{
		statusCode: statusCode,
		headers:    headers,
		body:       body,
	}
}

func respond(c net.Conn, s string) {
	fmt.Fprint(c, s)
}

type Config struct {
	directory string
}

type Response struct {
	statusCode int
	headers    map[string]string
	body       string
}
