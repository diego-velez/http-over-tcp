package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ErrBadRequestLine = errors.New("invalid request line")

func RequestFromReader(reader io.Reader) (*Request, error) {
	r := &Request{}
	parsedRequestLine := false

	lines := getLinesChannel(reader)
	for l := range lines {
		if !parsedRequestLine {
			parsedRequestLine = true
			requestLine, err := parseRequestLine(l)
			if err != nil {
				return nil, fmt.Errorf("%w: %w", ErrBadRequestLine, err)
			}
			r.RequestLine = requestLine
		}
	}

	return r, nil
}

func parseRequestLine(line string) (RequestLine, error) {
	requestLineSplit := strings.Split(line, " ")
	if len(requestLineSplit) != 3 {
		return RequestLine{}, fmt.Errorf("expected 3 content for the request line, but got %d", len(requestLineSplit))
	}

	protocolVersion := requestLineSplit[2]
	protocolVersionSplit := strings.Split(protocolVersion, "/")
	if len(protocolVersionSplit) != 2 {
		return RequestLine{}, fmt.Errorf("invalid protocol version")
	}

	return RequestLine{
		Method:        requestLineSplit[0],
		RequestTarget: requestLineSplit[1],
		HttpVersion:   protocolVersionSplit[1]}, nil
}

func getLinesChannel(f io.Reader) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer close(out)

		var line bytes.Buffer
		buf := make([]byte, 8)
		for {
			n, err := f.Read(buf)

			if n > 0 {
				data := buf[:n]
				if bytes.ContainsRune(data, '\n') {
					lines := bytes.Split(data, []byte{'\r', '\n'})

					for i, l := range lines {
						_, err := line.Write(l)
						if err != nil {
							log.Fatal(err)
						}

						// The last segment is not a complete line probs
						if i == len(lines)-1 {
							// This should detect \r\n dividing the header and the body
							if len(l) == 0 {
								return
							}
							break
						}

						out <- line.String()

						line.Reset()
					}
				} else {
					_, err := line.Write(data)
					if err != nil {
						log.Fatal(err)
					}
				}
			}

			if err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				log.Fatal(err)
			}
		}
	}()

	return out
}
