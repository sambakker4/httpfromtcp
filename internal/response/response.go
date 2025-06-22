package response

import (
	"errors"
	"io"
	"strconv"

	"github.com/sambakker4/httpfromtcp/internal/headers"
)

type StatusCode int
type WriterState int

type Writer struct {
	StatusCode StatusCode
	state      WriterState
	Writer     io.Writer
}

const (
	writerStateStatusLine WriterState = iota
	writerStateHeaders
	writerStateBody
)

const (
	Success             StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func (w *Writer) Write(b []byte) (int, error) {
	n, err := w.Writer.Write(b)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writerStateStatusLine {
		return errors.New("error: writing request in the wrong order")
	}
	w.state = writerStateHeaders

	switch statusCode {
	case Success:
		_, err := w.Write([]byte("HTTP/1.1 200 OK\r\n"))
		return err
	case BadRequest:
		_, err := w.Write([]byte("HTTP/1.1 400 Bad Request\r\n"))
		return err
	case InternalServerError:
		_, err := w.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n"))
		return err
	default:
		return errors.New("error: unknown status code")
	}
}
func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != writerStateHeaders {
		return errors.New("error: writing request in the wrong order")
	}
	w.state = writerStateBody

	for key, val := range headers {
		_, err := w.Write([]byte(key + ": " + val + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.Write(p)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["content-length"] = strconv.Itoa(contentLen)
	headers["connection"] = "close"
	headers["content-type"] = "text/plain"

	return headers
}
