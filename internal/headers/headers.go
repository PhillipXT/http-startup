package headers

import (
    "bytes"
    "fmt"
    "log"
    "slices"
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

    // Don't trim space in the key because we need to check for invalid
    // spaces, i.e. "Host :", so we need the key to be "Host "
    key := strings.ToLower(line[:colon])
    value := strings.TrimSpace(line[colon + 1:])

    log.Printf("Header line: %s\n", line)
    log.Printf("Parsed header: [%s] == [%s]\n", key, value)

    validChars := []byte{'!','#','$','%','&','\'','*','+','-','.','^','_','`','|','~'}

    for _, char := range []byte(key) {
        if (char <= 'a' || char >= 'z') && (char <= 'A' || char >= 'Z') && (char <= '0' || char >= '9') && !slices.Contains(validChars, char) {
            return 0, false, fmt.Errorf("Header key contains illegal characters: %s", line)
        }
    }

    h[key] = value

    return index + 2, false, nil
}

func NewHeaders() Headers {
    return make(map[string]string)
}
