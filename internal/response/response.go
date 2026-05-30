package response

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/diego-velez/http-from-scratch-course/internal/headers"
)

type Writer struct {
	Buf bytes.Buffer
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	return WriteStatusLine(&w.Buf, statusCode)
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	return WriteHeaders(&w.Buf, headers)
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	_, _ = w.Buf.Write([]byte("\r\n"))
	_, _ = w.Buf.Write(p)
	return 0, nil
}

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

var statusCodeText = map[StatusCode]string{
	StatusOK:                  "OK",
	StatusBadRequest:          "Bad Request",
	StatusInternalServerError: "Internal Server Error"}

func (sc StatusCode) text() string {
	return statusCodeText[sc]
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	out := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, statusCode.text())
	_, _ = w.Write([]byte(out))
	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	var buf bytes.Buffer
	for key, value := range headers {
		field := fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := buf.Write([]byte(field))
		if err != nil {
			return err
		}
	}
	_, err := w.Write(buf.Bytes())
	return err
}
