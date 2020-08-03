package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"

	"github.com/paddlesteamer/gcrproxy/internal/utils"
)

var proxyAddr = "http://127.0.0.1:8080"

func handleConnection(conn net.Conn) {
	defer conn.Close()

	connID := make([]byte, 8)
	rand.Read(connID)

	for {
		req, err := utils.ReadSome(conn)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("error reading proxy request: %v\n", err)
				return
			}

			_, err := http.Post(proxyAddr, "application/octet-stream", bytes.NewBuffer(append(connID, 0x1)))
			if err != nil {
				fmt.Printf("error while sending connectionclose request: %v\n", err)
			}

			return
		}

		resp, err := http.Post(proxyAddr, "application/octet-stream", bytes.NewBuffer(append(append(connID, 0x0), req...)))
		if err != nil {
			fmt.Printf("error while forwarding request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		buf := make([]byte, 1024)
		for {
			n, err := resp.Body.Read(buf)
			if err != nil {
				if err == io.EOF {
					if _, err := conn.Write(buf[:n]); err != nil {
						fmt.Printf("error forwarding proxy response: %v\n", err)
					}

					break
				}

				fmt.Printf("error reading proxy response: %v\n", err)
				return
			}

			if _, err := conn.Write(buf[:n]); err != nil {
				fmt.Printf("error forwarding proxy response: %v\n", err)
				return
			}

		}

	}
}

func main() {
	if os.Getenv("PROXY") != "" {
		proxyAddr = os.Getenv("PROXY")
	}

	serv, err := net.Listen("tcp", "127.0.0.1:1080")
	if err != nil {
		fmt.Printf("couldn't start server: %v\n", err)
		return
	}

	for {
		conn, err := serv.Accept()
		if err != nil {
			fmt.Printf("error while accepting connection: %v\n", err)
			continue
		}

		go handleConnection(conn)
	}

}
