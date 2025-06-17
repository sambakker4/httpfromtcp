package request

import (
	"errors"
	"io"
	"strings"
	"unicode"

	"github.com/sambakker4/httpfromtcp/internal/headers"
)

type Request struct {
	RequestLine RequestLine
	Headers      headers.Headers
	state       int
}

const (
	requestStateInitialized = iota
	requestStateParsingHeaders
	requestStateDone          
)
const bufferSize = 8

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf := make([]byte, bufferSize, bufferSize)
	readToIndex := 0

	req := Request{state: requestStateInitialized, Headers: headers.NewHeaders()}

	for req.state != requestStateDone {
		if len(buf) == cap(buf) {
			newBuffer := make([]byte, len(buf)*2, len(buf)*2)
			copy(newBuffer, buf)
			buf = newBuffer
		}

		n, err := reader.Read(buf[readToIndex:])
		if errors.Is(err, io.EOF) {
			req.state = requestStateDone
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
	totalBytesParsed := 0
	for r.state != requestStateDone {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return totalBytesParsed, nil
		}

		totalBytesParsed += n
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case requestStateInitialized:
		newRequestLine, n, err := parseRequestLine(data)

		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, nil
		}

		r.RequestLine = *newRequestLine
		r.state = requestStateParsingHeaders
		return n, nil

	case requestStateParsingHeaders:
		n, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		
		if n == 0 {
			return 0, nil
		}
		
		if done {
			r.state = requestStateDone
		}
		return n, nil

	case requestStateDone:
		return 0, errors.New("error: trying to read data from a requestStateDone state")

	default:
		return 0, errors.New("error: unknown state")
	}
}
