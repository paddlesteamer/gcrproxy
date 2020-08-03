package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/elazarl/goproxy"
	"github.com/paddlesteamer/gcrproxy/internal/utils"
)

const proxyAddr = "127.0.0.1:2080"

var connPool = map[string]net.Conn{}

func handler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		fmt.Printf("couldn't get content-length header: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	buf := make([]byte, length)

	if _, err := io.ReadFull(r.Body, buf); err != nil {
		fmt.Printf("couldn't read post body: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	connID := strconv.FormatUint(binary.LittleEndian.Uint64(buf[:8]), 10)
	termFlag := buf[8] == 0x1
	req := buf[9:]

	conn, found := connPool[connID]
	if !found {
		conn, err = net.Dial("tcp", proxyAddr)
		if err != nil {
			fmt.Printf("couldn't connect to local proxy: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		connPool[connID] = conn
	}

	if termFlag {
		conn.Close()

		delete(connPool, connID)
		return
	}

	if _, err := conn.Write(req); err != nil {
		fmt.Printf("couldn't forward post body to local proxy: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := utils.ReadSome(conn)
	if err != nil {
		if err == io.EOF {
			_, err = w.Write(resp)
			if err != nil {
				fmt.Printf("couldn't write response to tunnel: %v\n", err)
				w.WriteHeader(http.StatusInternalServerError)
			}

			return
		}

		fmt.Printf("couldn't read response from tunnel: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		fmt.Printf("couldn't write response to tunnel: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func serveProxy(proxy *goproxy.ProxyHttpServer) {
	log.Fatal(http.ListenAndServe(proxyAddr, proxy))
}

func main() {
	port := "8080"

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	proxy := goproxy.NewProxyHttpServer()

	go serveProxy(proxy)

	web := http.NewServeMux()
	web.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":"+port, web))
}
