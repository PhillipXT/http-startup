package request

import (
    "bytes"
    "fmt"
    "io"
    "strings"
)

type Request struct {
    RequestLine RequestLine
}

type RequestLine struct {
    HttpVersion string
    RequestTarget string
    Method string
}

const (
    crlf = "\r\n"
)

func RequestFromReader (reader io.Reader) (*Request, error) {

    data, err := io.ReadAll(reader)
    if err != nil {
        return nil, err
    }

    index := bytes.Index(data, []byte(crlf))
    if index == -1 {
        return nil, fmt.Errorf("Missing CRLF in request line")
    }

    line := string(data[:index])

    request_line, err := ParseRequestLine(line)
    if err != nil {
        return nil, err
    }

    request := Request {
        RequestLine: *request_line,
    }

    return &request, nil
}

func ParseRequestLine(line string) (*RequestLine, error) {

    parts := strings.Split(line, " ")

    if len(parts) != 3 {
        return nil, fmt.Errorf("Malformed request line: %si", line)
    }

    method := parts[0]
    if method != strings.ToUpper(method) {
        return nil, fmt.Errorf("Invalid method: %s", method)
    }

    target := parts[1]

    version_parts := strings.Split(parts[2], "/")
    if len(version_parts) != 2 {
        return nil, fmt.Errorf("Malformed version: %s", parts[2])
    }

    http_part := version_parts[0]
    if http_part != "HTTP" {
        return nil, fmt.Errorf("Unsupported protocol: %s", http_part)
    }

    version := version_parts[1]
    if version != "1.1" {
        return nil, fmt.Errorf("Unsupported version: %s", version_parts[1])
    }

    request_line := RequestLine {
        HttpVersion: version,
        RequestTarget: target,
        Method: method,
    }

    return &request_line, nil
}
