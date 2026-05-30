package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/diego-velez/http-from-scratch-course/internal/headers"
)

var ErrBadRequestLine = errors.New("invalid request line")

type RequestState int

const (
	StateStartLine RequestState = iota
	StateHeader
	StateBody
	StateDone
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers

	buf   bytes.Buffer
	state RequestState
}

func NewRequest() *Request {
	return &Request{Headers: headers.NewHeaders(), state: StateStartLine}
}

func (r *Request) parse(data []byte) (int, error) {
	if r.state == StateDone {
		return 0, nil
	}

	_, err := r.buf.Write(data)
	if err != nil {
		return 0, err
	}

	bufBytes := r.buf.Bytes()
	if !bytes.Contains(bufBytes, []byte("\r\n")) {
		return 0, nil
	}

	// Account for the 2 bytes of '\r\n'
	parsed := 2
	lines := bytes.Split(bufBytes, []byte("\r\n"))
	for i, l := range lines {
		// We assume that the last line is incomplete so we do not parse it
		if i == len(lines)-1 {
			// We only want to keep in the buffer unparsed bytes
			r.buf.Reset()
			_, err := r.buf.Write(l)
			if err != nil {
				return 0, err
			}
			break
		}

		if r.state == StateDone {
			break
		}

		switch r.state {
		case StateStartLine:
			r.state = StateHeader
			requestLine, n, err := parseRequestLine(l)
			parsed += n
			if err != nil {
				return parsed, err
			}
			r.RequestLine = requestLine
		case StateHeader:
			n, done, err := r.Headers.Parse(l)
			parsed += n
			if err != nil {
				return parsed, err
			}
			if done {
				r.state = StateDone
			}
		default:
			log.Fatal("unknown parser state")
		}
	}

	return parsed, nil
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := NewRequest()

	for r.state != StateDone {
		buf := make([]byte, 8)
		n, err := reader.Read(buf)

		if n > 0 {
			_, err := r.parse(buf[:n])
			if err != nil {
				return nil, err
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
	}

	return r, nil
}

func parseRequestLine(line []byte) (RequestLine, int, error) {
	requestLineSplit := bytes.Split(line, []byte{' '})
	if len(requestLineSplit) != 3 {
		return RequestLine{}, 0, fmt.Errorf("expected 3 content for the request line, but got %d", len(requestLineSplit))
	}

	protocolVersion := requestLineSplit[2]
	protocolVersionSplit := bytes.Split(protocolVersion, []byte{'/'})
	if len(protocolVersionSplit) != 2 {
		return RequestLine{}, 0, fmt.Errorf("invalid protocol version")
	}

	return RequestLine{
		Method:        string(requestLineSplit[0]),
		RequestTarget: string(requestLineSplit[1]),
		HttpVersion:   string(protocolVersionSplit[1])}, len(line), nil
}

func bytesToString(bytes [][]byte) []string {
	result := make([]string, 0, len(bytes))
	for _, b := range bytes {
		result = append(result, string(b))
	}
	return result
}
