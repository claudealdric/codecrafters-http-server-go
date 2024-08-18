package main

import (
	"fmt"
	"log"
	"net"

	"github.com/codecrafters-io/http-server-starter-go/app/config"
	"github.com/codecrafters-io/http-server-starter-go/app/httpserver"
)

const port = 4221

func main() {
	fmt.Println("Logs from your program will appear here!")

	config.ParseFlags()

	listener, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalf("Failed to bind to port %d\n", port)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Error accepting connection:", err)
		}

		go httpserver.ProcessRequest(conn)
	}
}
