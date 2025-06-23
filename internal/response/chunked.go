package response

import (
	"fmt"
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
	n, err := w.Write([]byte("0\r\n\r\n"))
	if err != nil {
		return 0, err
	}
	return n, nil
}
