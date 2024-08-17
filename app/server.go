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

type Config struct {
	directory string
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
	var response strings.Builder
	buildStatusLine(&response, http.StatusOK)
	buildDelineator(&response)

	fmt.Fprint(conn, response.String())
}

func handleNotFound(conn net.Conn) {
	var response strings.Builder
	buildStatusLine(&response, http.StatusNotFound)
	buildDelineator(&response)

	fmt.Fprint(conn, response.String())
}

func handleEcho(conn net.Conn, path string) {
	echoArg := strings.TrimPrefix(path, "/echo/")

	var response strings.Builder
	buildStatusLine(&response, http.StatusOK)
	buildPlainTextHeaders(&response, echoArg)
	response.WriteString(echoArg)

	fmt.Fprint(conn, response.String())
}

func handleUserAgent(conn net.Conn, headers map[string]string) {
	userAgent := headers["user-agent"]

	var response strings.Builder
	buildStatusLine(&response, http.StatusOK)
	buildPlainTextHeaders(&response, userAgent)
	response.WriteString(userAgent)

	fmt.Fprint(conn, response.String())
}

func handleFiles(conn net.Conn, path string) {
	fileName := strings.TrimPrefix(path, "/files/")
	content, err := os.ReadFile(filepath.Join(config.directory, fileName))
	if err != nil && os.IsNotExist(err) {
		handleNotFound(conn)
		return
	}

	var response strings.Builder
	buildStatusLine(&response, http.StatusOK)
	buildOctetStreamHeaders(&response, content)
	response.WriteString(string(content))

	fmt.Fprint(conn, response.String())
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

func buildPlainTextHeaders(builder *strings.Builder, content string) {
	headers := map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": strconv.Itoa(len(content)),
	}

	for k, v := range headers {
		builder.WriteString(fmt.Sprintf("%s: %s", k, v))
		buildDelineator(builder)
	}

	buildDelineator(builder)
}

func buildOctetStreamHeaders(builder *strings.Builder, content []byte) {
	headers := map[string]string{
		"Content-Type":   "application/octet-stream",
		"Content-Length": strconv.Itoa(len(content)),
	}

	for k, v := range headers {
		builder.WriteString(fmt.Sprintf("%s: %s", k, v))
		buildDelineator(builder)
	}

	buildDelineator(builder)
}
