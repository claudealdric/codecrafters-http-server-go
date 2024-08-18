package httpserver

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/config"
)

var supportedEncodings = map[string]struct{}{
	"gzip": {},
}

func handleEcho(request *Request) {
	echoArg := strings.TrimPrefix(request.path, "/echo/")
	response := NewResponse(StatusOK, echoArg)
	encoding := getSupportedEncoding(request)
	if encoding != "" {
		response.headers["Content-Encoding"] = encoding
	}
	response.Send(request)
}

func getSupportedEncoding(request *Request) string {
	encondingsString := request.headers["accept-encoding"]
	if encondingsString == "" {
		return ""
	}
	encodingsSlice := strings.Split(encondingsString, ", ")
	for _, encoding := range encodingsSlice {
		if _, isSupported := supportedEncodings[encoding]; isSupported {
			return encoding
		}
	}
	return ""
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
	err := os.WriteFile(filePath, request.body, 0644)
	if err != nil {
		fmt.Printf("error writing to file %q\n: %v", filePath, err)
		return
	}
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
