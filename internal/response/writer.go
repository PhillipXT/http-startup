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
	writerStateComplete    writerState = 3
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
	for key, value := range headers {
		fmt.Printf("%s: %s\n", key, value)
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	if err == nil {
		w.state = writerStateBody
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
	case writerStateComplete:
		return "Response is already complete."
	default:
		return "Unknown state"
	}
}
