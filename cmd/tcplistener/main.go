package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close() // nolint: errcheck
		defer close(out)

		var line bytes.Buffer
		buf := make([]byte, 8)
		for {
			n, err := f.Read(buf)

			if n > 0 {
				data := buf[:n]
				if bytes.ContainsRune(data, '\n') {
					lines := bytes.Split(data, []byte{'\n'})

					for i, l := range lines {
						_, err := line.Write(l)
						if err != nil {
							log.Fatal(err)
						}

						// The last segment is not a complete line probs
						if i == len(lines)-1 {
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

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:42069")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close() // nolint: errcheck

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("New connection accepted!")

		lines := getLinesChannel(conn)
		for line := range lines {
			fmt.Printf("read: %s\n", line)
		}
	}
}
