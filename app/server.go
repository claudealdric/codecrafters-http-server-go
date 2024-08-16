package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
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

func main() {
	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Printf("Failed to bind to port %d\n", port)
		os.Exit(1)
	}
	defer listener.Close()

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	processRequest(conn)
}

func processRequest(conn net.Conn) {
	defer conn.Close()

	requestLine, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		return
	}

	path := extractPath(requestLine)
	routeRequest(conn, path)

}

func extractPath(requestLine string) string {
	return strings.Fields(requestLine)[1]
}

func routeRequest(conn net.Conn, path string) {
	switch {
	case path == "/":
		handleRoot(conn)
	case strings.HasPrefix(path, "/echo/"):
		handleEcho(conn, path)
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

func buildStatusLine(b *strings.Builder, statusCode int) {
	b.WriteString(fmt.Sprintf(
		"%s %d %s",
		httpVersion,
		statusCode,
		statusCodeToReasonPhrase[statusCode],
	))
	buildDelineator(b)
}

func buildDelineator(b *strings.Builder) {
	b.WriteString("\r\n")
}

func buildPlainTextHeaders(b *strings.Builder, content string) {
	headers := map[string]string{
		"Content-Type":   "text/plain",
		"Content-Length": strconv.Itoa(len(content)),
	}

	for k, v := range headers {
		b.WriteString(fmt.Sprintf("%s: %s", k, v))
		buildDelineator(b)
	}

	buildDelineator(b)
}
