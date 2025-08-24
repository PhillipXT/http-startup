package response

import (
	"errors"
	"fmt"
	"io"

	"github.com/PhillipXT/http-startup/internal/headers"
)

type writerState int

const (
	writerStateInitialized writerState = 0
	writerStateHeaders     writerState = 1
	writerStateBody        writerState = 2
	writerStateChunking    writerState = 3
	writerStateComplete    writerState = 4
)

type Writer struct {
	writer io.Writer
	state  writerState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
		state:  writerStateInitialized,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if r := checkState(w.state, writerStateInitialized); r != "" {
		return errors.New(r)
	}
	_, err := w.writer.Write(getStatusLine(statusCode))
	if err == nil {
		w.state = writerStateHeaders
	}
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if r := checkState(w.state, writerStateHeaders); r != "" {
		return errors.New(r)
	}

	chunking := false
	if headers["Transfer-Encoding"] == "chunked" {
		chunking = true
	}

	for key, value := range headers {
		fmt.Printf("Header: %s: %s\n", key, value)
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	if err == nil {
		if chunking {
			w.state = writerStateChunking
		} else {
			w.state = writerStateBody
		}
	}
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if r := checkState(w.state, writerStateBody); r != "" {
		return 0, errors.New(r)
	}
	n, err := w.writer.Write(p)
	if err == nil {
		w.state = writerStateComplete
	}
	return n, nil
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if r := checkState(w.state, writerStateChunking); r != "" {
		return 0, errors.New(r)
	}
	t := []byte(fmt.Sprintf("%x\r\n%s\r\n", len(p), string(p)))
	return w.writer.Write(t)
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if r := checkState(w.state, writerStateChunking); r != "" {
		return 0, errors.New(r)
	}
	t := []byte("0\r\n")
	w.state = writerStateComplete
	return w.writer.Write(t)
}

func (w *Writer) WriteTrailers(trailers headers.Headers) error {
	if r := checkState(w.state, writerStateComplete); r != "" {
		return errors.New(r)
	}

	for key, value := range trailers {
		fmt.Printf("Trailer: %s: %s\n", key, value)
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return err
		}
	}
	fmt.Println("Finished writing trailers")
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}

func checkState(current, requested writerState) string {
	switch current {
	case requested:
		return ""
	case writerStateInitialized:
		return "Must set the StatusLine first."
	case writerStateHeaders:
		return "Expecting headers to be set next."
	case writerStateBody:
		return "Expecting body to be set next."
	case writerStateChunking:
		return "Expecting chunked body to be set next."
	case writerStateComplete:
		return "Response is already complete."
	default:
		return "Unknown state"
	}
}
