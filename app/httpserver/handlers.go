package httpserver

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/config"
)

var validEncodings = map[string]struct{}{
	"gzip": struct{}{},
}

func handleEcho(request *Request) {
	echoArg := strings.TrimPrefix(request.path, "/echo/")
	response := NewResponse(StatusOK, echoArg)
	cs, shouldEncode := request.headers["accept-encoding"]
	if shouldEncode {
		_, validEncoding := validEncodings[cs]
		if validEncoding {
			response.headers["Content-Encoding"] = cs
		}
	}
	response.Send(request)
}

func handleGetFiles(request *Request) {
	fileName := strings.TrimPrefix(request.path, "/files/")
	content, err := os.ReadFile(filepath.Join(config.Directory, fileName))
	if err != nil && os.IsNotExist(err) {
		handleNotFound(request)
		return
	}
	response := NewResponse(StatusOK, string(content))
	response.headers["Content-Type"] = "application/octet-stream"
	response.Send(request)
}

func handleNotFound(request *Request) {
	response := NewResponse(StatusNotFound, "")
	response.Send(request)
}

func handlePostFiles(request *Request) {
	fileName := strings.TrimPrefix(request.path, "/files/")
	filePath := filepath.Join(config.Directory, fileName)
	os.WriteFile(filePath, request.body, 0644)
	response := NewResponse(StatusCreated, "")
	response.Send(request)

}

func handleRoot(request *Request) {
	response := NewResponse(StatusOK, "")
	response.Send(request)
}

func handleUserAgent(request *Request) {
	userAgent := request.headers["user-agent"]
	response := NewResponse(StatusOK, userAgent)
	response.Send(request)
}
