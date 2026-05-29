package main

import (
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

	buf := make([]byte, 8)
	for {
		n, err := f.Read(buf)

		if n > 0 {
			fmt.Printf("read (%d bytes): %s\n", n, string(buf[:n]))
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			panic(err)
		}
	}
}
