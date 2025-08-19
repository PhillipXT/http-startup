package headers

import (
    "fmt"
    "bytes"
    "strings"
)

type Headers map[string]string

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

    // If no CRLF is found, we don't have enough data to process yet
    index := bytes.Index(data, []byte(crlf))
    if index == -1 {
        return 0, false, nil
    }

    // If CRLF is found at the beginning, we are at the end of the headers
    if index == 0 {
        return 2, true, nil
    }

    line := strings.TrimSpace(string(data[:index]))

    colon := strings.Index(line, ":")
    if colon == -1 {
        return 0, false, fmt.Errorf("Malformed header: %s", line)
    }

    key := strings.TrimSpace(line[:colon])
    value := strings.TrimSpace(line[colon + 1:])

    if len(key) != colon {
        return 0, false, fmt.Errorf("Space found inside header key: %s", line)
    }

    h[key] = value

    return index + 2, false, nil
}

func NewHeaders() Headers {
    return make(map[string]string)
}
