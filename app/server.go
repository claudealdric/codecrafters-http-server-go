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

func main() {
	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Printf("Failed to bind to port %d\n", port)
		os.Exit(1)
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	statusCodeToReasonPhrase := map[int]string{
		http.StatusOK:       "OK",
		http.StatusNotFound: "Not Found",
	}

	reader := bufio.NewReader(conn)
	requestLine, _ := reader.ReadString('\n')
	target := strings.Fields(requestLine)[1]

	switch target {
	case "/":
		fmt.Fprintf(
			conn,
			"%s %d %s\r\n\r\n",
			httpVersion,
			http.StatusOK,
			statusCodeToReasonPhrase[http.StatusOK],
		)
	default:
		fmt.Fprintf(
			conn,
			"%s %d %s\r\n\r\n",
			httpVersion,
			http.StatusNotFound,
			statusCodeToReasonPhrase[http.StatusNotFound],
		)
	}

}
