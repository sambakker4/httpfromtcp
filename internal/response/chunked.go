package response

import (
	"fmt"

	"github.com/sambakker4/httpfromtcp/internal/headers"
)

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	hex := fmt.Sprintf("%X", len(p))
	n, err := w.Write([]byte(fmt.Sprintf("%s\r\n%s\r\n", hex, string(p))))
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	n, err := w.Write([]byte("0\r\n"))
	if err != nil {
		return 0, err
	}
	return n, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	for key, val := range h {
		_, err := w.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, val)))
		if err != nil {
			return err
		}
	}
	w.Write([]byte("\r\n"))
	return nil
}
