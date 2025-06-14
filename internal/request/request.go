package request

import (
	"errors"
	"io"
	"strings"
	"unicode"
)

type Request struct {
	RequestLine RequestLine
	Enum          int
}

const initialized = 1
const done = 2
const bufferSize = 8

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	req := Request{Enum: initialized}

	for req.Enum != done {
		if len(buf) == cap(buf) {
			newBuffer := make([]byte, len(buf) * 2, len(buf) * 2)
			copy(newBuffer, buf)
			buf = newBuffer
		}

		n, err := reader.Read(buf[readToIndex:])
		if errors.Is(err, io.EOF) {
			req.Enum = done
			break
		}

		if err != nil {
			return &Request{}, err
		}

		readToIndex += n

		n, err = req.parse(buf[:readToIndex])
		if err != nil {
			return &Request{}, err
		}

		newSlice := make([]byte, len(buf[n:]), cap(buf))
		copy(newSlice, buf[n:])
		buf = newSlice
		readToIndex -= n
	}

	return &req, nil
}

func parseRequestLine(s []byte) (*RequestLine, int, error) {
	if !strings.Contains(string(s), "\r\n") {
		return nil, 0, nil
	}

	line := strings.Split(string(s), "\r\n")[0]

	if len(strings.Split(line, " ")) != 3 {
		return nil, 0, errors.New("invalid parts of request line")
	}

	parts := strings.Split(line, " ")

	method := parts[0]
	if method == "" {
		return nil, 0, errors.New("method is not all uppercase letters")
	}

	for _, char := range method {
		if !unicode.IsUpper(char) {
			return nil, 0, errors.New("method is not all uppercase letters")
		}
	}

	target := parts[1]
	if !strings.Contains(target, "/") {
		return nil, 0, errors.New("invalid target")
	}
	httpVersion := parts[2]

	if httpVersion != "HTTP/1.1" {
		return nil, 0, errors.New("no support for versions other than HTTP/1.1")
	}

	version, _ := strings.CutPrefix(httpVersion, "HTTP/")

	return &RequestLine{
		HttpVersion:   version,
		RequestTarget: target,
		Method:        method,
	}, len(line) + len("\r\n"), nil
}

func (r *Request) parse(data []byte) (int, error) {
	if r.Enum == initialized {
		newRequestLine, n, err := parseRequestLine(data)

		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *newRequestLine
		r.Enum = done
		return len(data), nil
	}

	if r.Enum == done {
		return 0, errors.New("error: trying to read data from a done state")
	}

	return 0, errors.New("error: unknown state")
}
