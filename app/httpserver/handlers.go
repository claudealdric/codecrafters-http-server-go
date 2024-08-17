package httpserver

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/config"
)

func handleEcho(request *request) {
	echoArg := strings.TrimPrefix(request.path, "/echo/")
	response := NewResponse(http.StatusOK, echoArg)
	response.Send(request)
}

func handleGetFiles(request *request) {
	fileName := strings.TrimPrefix(request.path, "/files/")
	content, err := os.ReadFile(filepath.Join(config.Directory, fileName))
	if err != nil && os.IsNotExist(err) {
		handleNotFound(request)
		return
	}
	response := NewResponse(http.StatusOK, string(content))
	response.headers["Content-Type"] = "application/octet-stream"
	response.Send(request)
}

func handleNotFound(request *request) {
	response := NewResponse(http.StatusNotFound, "")
	response.Send(request)
}

func handlePostFiles(request *request) {
	fileName := strings.TrimPrefix(request.path, "/files/")
	filePath := filepath.Join(config.Directory, fileName)
	os.WriteFile(filePath, request.body, 0644)
	response := NewResponse(http.StatusCreated, "")
	response.Send(request)

}

func handleRoot(request *request) {
	response := NewResponse(http.StatusOK, "")
	response.Send(request)
}

func handleUserAgent(request *request) {
	userAgent := request.headers["user-agent"]
	response := NewResponse(http.StatusOK, userAgent)
	response.Send(request)
}
