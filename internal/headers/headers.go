package headers

import (
	"bytes"
	"errors"
)

var ErrInvalidField = errors.New("invalid header field")

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (int, bool, error) {
	if !bytes.Contains(data, []byte("\r\n")) {
		return 0, false, nil
	}

	parsed := 0
	prevEmpty := false
	for field := range bytes.SplitSeq(data, []byte("\r\n")) {
		if len(field) == 0 {
			if prevEmpty {
				return parsed, true, nil
			}
			prevEmpty = true
			continue
		}

		key, value, err := parseHeader(field)
		if err != nil {
			return parsed, false, err
		}

		// We add 2 because of '\r\n'
		parsed += len(field) + 2

		h[key] = value
	}

	return parsed, false, nil
}

func parseHeader(fieldLine []byte) (string, string, error) {
	fieldSplit := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(fieldSplit) != 2 {
		return "", "", ErrInvalidField
	}

	key := fieldSplit[0]
	if bytes.HasPrefix(key, []byte(" ")) || bytes.HasSuffix(key, []byte(" ")) {
		return "", "", ErrInvalidField
	}
	value := bytes.TrimSpace(fieldSplit[1])

	return string(key), string(value), nil
}

func bytesToString(bytes [][]byte) []string {
	result := make([]string, 0, len(bytes))
	for _, b := range bytes {
		result = append(result, string(b))
	}
	return result
}
