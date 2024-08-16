package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	reader := bufio.NewReader(conn)
	requestLine, _ := reader.ReadString('\n')
	target := strings.Fields(requestLine)[1]

	switch target {
	case "/":
		fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n\r\n")
	default:
		fmt.Fprint(conn, "HTTP/1.1 404 Not Found\r\n\r\n")
	}

}
