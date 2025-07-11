package headers

import (
	"errors"
	"regexp"
	"strings"
)

type Headers map[string]string

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	if !strings.Contains(string(data), "\r\n") {
		return 0, false, nil
	}

	if strings.Index(string(data), "\r\n") == 0 {
		return len("\r\n"), true, nil
	}

	line := strings.Split(string(data), "\r\n")

	str := strings.TrimSpace(line[0])
	header := strings.SplitN(str, ":", 2)
	if len(header) != 2 {
		return 0, false, errors.New("error: header must be in <key>: <value>, format")
	}

	key, value := header[0], header[1]
	if key[len(key)-1] == ' ' {
		return 0, false, errors.New("error: header format must be <key>: <value>")
	}

	re := regexp.MustCompile(`^[A-Za-z0-9!#$%&'*+\-.\^_` + "`" + `|~]+$`)
	if !re.MatchString(key) || len(key) == 0 {
		return 0, false, errors.New("error: key contains invalid characters")
	}

	n = len(line[0]) + len("\r\n")
	err = nil
	done = false
	value = strings.TrimSpace(value)

	if h.Get(key) != "" {
		h[strings.ToLower(key)] = h.Get(key) + ", " + value
		return
	}

	h[strings.ToLower(key)] = value
	return
}

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Get(s string) string {
	val, ok := h[strings.ToLower(s)]
	if !ok {
		return ""
	}
	return val
}

func (h Headers) Set(key string, val string) {
	h[strings.ToLower(key)] = strings.ToLower(val)
}
