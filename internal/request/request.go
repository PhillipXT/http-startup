package request

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/PhillipXT/http-startup/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       requestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateDone
)

const crlf = "\r\n"
const bufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {

	buffer := make([]byte, bufferSize)

	readToIndex := 0

	request := &Request{
		Headers: headers.NewHeaders(),
		state:   requestStateInitialized,
	}

	for request.state != requestStateDone {
		if readToIndex >= len(buffer) {
			newB := make([]byte, len(buffer)*2)
			copy(newB, buffer)
			buffer = newB
		}

		bytesRead, err := reader.Read(buffer[readToIndex:])
		if err != nil {
			if err == io.EOF {
				if request.state != requestStateDone {
					return nil, fmt.Errorf("incomplete request:  state (%d) bytes (%d)", request.state, bytesRead)
				}
			}
			return nil, err
		}

		readToIndex += bytesRead

		bytesParsed, err := request.parse(buffer[:readToIndex])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[bytesParsed:])
		readToIndex -= bytesParsed
	}

	return request, nil
}

func parseRequestLine(data []byte) (*RequestLine, int, error) {

	index := bytes.Index(data, []byte(crlf))
	if index == -1 {
		return nil, 0, nil
	}

	line := string(data[:index])

	request_line, err := requestLineFromString(line)
	if err != nil {
		return nil, 0, err
	}

	return request_line, index + 2, nil
}

func requestLineFromString(line string) (*RequestLine, error) {

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

	request_line := RequestLine{
		HttpVersion:   version,
		RequestTarget: target,
		Method:        method,
	}

	return &request_line, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	// Loop through and make sure we process all complete headers.  If
	// we don't do that here, the main loop could get an EOF while we
	// still have unprocessed data in our buffer.
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		totalBytesParsed += n
		if n == 0 {
			break
		}
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		requestLine, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = requestStateParsingHeaders
		return n, nil
	case requestStateParsingHeaders:
		log.Printf("Parsing data: [%s]", string(data))
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		} else if done {
			r.state = requestStateDone
		}
		return n, nil
	case requestStateDone:
		return 0, fmt.Errorf("error: trying to read but operation is complete")
	default:
		return 0, fmt.Errorf("unknown state")
	}
}
