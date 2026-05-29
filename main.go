package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		panic(err)
	}

	var line bytes.Buffer
	buf := make([]byte, 8)
	for {
		n, err := f.Read(buf)

		if n > 0 {
			data := buf[:n]
			if bytes.ContainsRune(data, '\n') {
				lines := bytes.Split(data, []byte{'\n'})
				if len(lines) != 2 {
					panic("expected only one \\n")
				}

				_, err := line.Write(lines[0])
				if err != nil {
					panic(err)
				}

				fmt.Printf("read: %s\n", line.String())

				line.Reset()

				_, err = line.Write(lines[1])
				if err != nil {
					panic(err)
				}
			} else {
				_, err := line.Write(data)
				if err != nil {
					panic(err)
				}
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
	}
}
