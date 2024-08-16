package main

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"os"
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

	conn, err := listener.Accept()
	defer listener.Close()

	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	requestLine, _ := bufio.NewReader(conn).ReadString('\n')
	path := strings.Fields(requestLine)[1]

	routePath(conn, path)

}

func routePath(conn net.Conn, path string) {
	if path == "/" {
		handleRoot(conn)
	} else if strings.HasPrefix(path, "/echo/") {
		handleEcho(conn, path)
	} else {
		handleNotFound(conn)
	}
}

func handleRoot(conn net.Conn) {
	fmt.Fprintf(
		conn,
		"%s %d %s\r\n\r\n",
		httpVersion,
		http.StatusOK,
		statusCodeToReasonPhrase[http.StatusOK],
	)
}

func handleNotFound(conn net.Conn) {
	fmt.Fprintf(
		conn,
		"%s %d %s\r\n\r\n",
		httpVersion,
		http.StatusNotFound,
		statusCodeToReasonPhrase[http.StatusNotFound],
	)
}

func handleEcho(conn net.Conn, path string) {
	echoArg := strings.TrimPrefix(path, "/echo/")
	fmt.Fprintf(
		conn,
		"%s %d %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		httpVersion,
		http.StatusOK,
		statusCodeToReasonPhrase[http.StatusOK],
		len(echoArg),
		echoArg,
	)
}
