package httpserver

import (
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/codecrafters-io/http-server-starter-go/app/config"
)

func handleEcho(conn net.Conn, path string) {
	echoArg := strings.TrimPrefix(path, "/echo/")
	response := getResponse(http.StatusOK, echoArg)
	respond(conn, getResponseString(response))
}

func handleGetFiles(conn net.Conn, path string) {
	fileName := strings.TrimPrefix(path, "/files/")
	content, err := os.ReadFile(filepath.Join(config.Directory, fileName))
	if err != nil && os.IsNotExist(err) {
		handleNotFound(conn)
		return
	}
	response := getResponse(http.StatusOK, string(content))
	response.headers["Content-Type"] = "application/octet-stream"
	respond(conn, getResponseString(response))
}

func handlePostFiles(conn net.Conn, path string, data []byte) {
	fileName := strings.TrimPrefix(path, "/files/")
	filePath := filepath.Join(config.Directory, fileName)
	os.WriteFile(filePath, data, 0644)
	response := getResponse(http.StatusCreated, "")
	respond(conn, getResponseString(response))

}

func handleRoot(conn net.Conn) {
	response := getResponse(http.StatusOK, "")
	respond(conn, getResponseString(response))
}

func handleUserAgent(conn net.Conn, headers map[string]string) {
	userAgent := headers["user-agent"]
	response := getResponse(http.StatusOK, userAgent)
	respond(conn, getResponseString(response))
}
