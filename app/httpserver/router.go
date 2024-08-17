package httpserver

import (
	"strings"
)

func RouteRequest(request *Request) {
	switch request.method {
	case "GET":
		routeGetRequest(request)
	case "POST":
		routePostRequest(request)
	default:
		handleNotFound(request)
	}

}

func routeGetRequest(request *Request) {
	switch {
	case request.path == "/":
		handleRoot(request)
	case strings.HasPrefix(request.path, "/echo/"):
		handleEcho(request)
	case strings.HasPrefix(request.path, "/user-agent"):
		handleUserAgent(request)
	case strings.HasPrefix(request.path, "/files/"):
		handleGetFiles(request)
	default:
		handleNotFound(request)
	}
}

func routePostRequest(request *Request) {
	switch {
	case strings.HasPrefix(request.path, "/files/"):
		handlePostFiles(request)
	default:
		handleNotFound(request)
	}
}
