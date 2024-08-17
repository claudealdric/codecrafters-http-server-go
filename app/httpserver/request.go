package httpserver

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

type request struct {
	conn    net.Conn
	method  string
	path    string
	headers map[string]string
	body    []byte
}

func ProcessRequest(conn net.Conn) {
	defer conn.Close()

	request, err := newRequest(conn)
	if err != nil {
		fmt.Println("Error parsing the request:", err.Error())
		return
	}

	routeRequest(request)
}

func newRequest(conn net.Conn) (*request, error) {
	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	headers, err := getHeaders(reader)
	if err != nil {
		return nil, err
	}

	body, err := getRequestBody(reader, headers)
	if err != nil {
		return nil, err
	}

	return &request{
		conn:    conn,
		method:  extractMethod(requestLine),
		path:    extractPath(requestLine),
		headers: headers,
		body:    body,
	}, nil
}

func getHeaders(reader *bufio.Reader) (map[string]string, error) {
	headers := make(map[string]string)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if line == "\r\n" {
			break
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key := strings.ToLower(strings.TrimSpace(parts[0]))
			value := (strings.TrimSpace(parts[1]))
			headers[key] = value
		}
	}

	return headers, nil
}

func getRequestBody(
	reader *bufio.Reader,
	headers map[string]string,
) ([]byte, error) {
	contentLength := headers["content-length"]
	if contentLength == "" {
		return nil, nil
	}

	length, err := strconv.Atoi(contentLength)
	if err != nil {
		return nil, err
	}

	body := make([]byte, length)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func extractMethod(requestLine string) string {
	return strings.Fields(requestLine)[0]
}

func extractPath(requestLine string) string {
	return strings.Fields(requestLine)[1]
}
