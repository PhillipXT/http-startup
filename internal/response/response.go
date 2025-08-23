package response

import (
	"fmt"
	"io"

	"github.com/PhillipXT/http-startup/internal/headers"
)

type StatusCode int

const (
	StatusCodeOK                  StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

var statusCodes = map[StatusCode]string{
	200: "OK",
	400: "Bad Request",
	500: "Internal Server Error",
}

func getStatusLine(statusCode StatusCode) []byte {

	reasonPhrase := statusCodes[statusCode]
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase))
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	_, err := w.Write(getStatusLine(statusCode))
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	header := headers.NewHeaders()
	header["Content-Type"] = "text/plain"
	header["Content-Length"] = fmt.Sprintf("%d", contentLen)
	header["Connection"] = "close"
	return header
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		fmt.Printf("%s: %s\n", key, value)
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
