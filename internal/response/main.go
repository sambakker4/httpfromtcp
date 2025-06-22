package response

import (
	"errors"
	"io"
	"strconv"

	"github.com/sambakker4/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	Success             StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
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

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["content-length"] = strconv.Itoa(contentLen)
	headers["connection"] = "close"
	headers["content-type"] = "text/plain"

	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, val := range headers {
		_, err := w.Write([]byte(key + ": " + val + "\r\n"))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}
