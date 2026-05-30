package response

import (
	"bytes"
	"errors"
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
	err1 := WriteHeaders(&w.Buf, headers)
	_, err2 := w.Buf.Write([]byte("\r\n"))
	return errors.Join(err1, err2)
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.Buf.Write(p)
	return n, err
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	n1, err1 := w.Buf.Write([]byte(strconv.Itoa(len(p))))
	n2, err2 := w.Buf.Write([]byte("\r\n"))
	n3, err3 := w.Buf.Write(p)
	n4, err4 := w.Buf.Write([]byte("\r\n"))
	return n1 + n2 + n3 + n4, errors.Join(err1, err2, err3, err4)
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
