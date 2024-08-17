package httpserver

import (
	"net"
	"net/http"
	"strings"
)

func routeRequest(conn net.Conn, request *request) {
	method := extractMethod(request.line)

	switch method {
	case "GET":
		routeGetRequest(conn, request)
	case "POST":
		routePostRequest(conn, request)
	default:
		handleNotFound(conn)
	}

}

func handleNotFound(conn net.Conn) {
	response := getResponse(http.StatusNotFound, "")
	respond(conn, getResponseString(response))
}

func routeGetRequest(conn net.Conn, request *request) {
	path := extractPath(request.line)

	switch {
	case path == "/":
		handleRoot(conn)
	case strings.HasPrefix(path, "/echo/"):
		handleEcho(conn, path)
	case strings.HasPrefix(path, "/user-agent"):
		handleUserAgent(conn, request.headers)
	case strings.HasPrefix(path, "/files/"):
		handleGetFiles(conn, path)
	default:
		handleNotFound(conn)
	}
}

func routePostRequest(conn net.Conn, request *request) {
	path := extractPath(request.line)

	switch {
	case strings.HasPrefix(path, "/files/"):
		handlePostFiles(conn, path, request.body)
	default:
		handleNotFound(conn)
	}
}
