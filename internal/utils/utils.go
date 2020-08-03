package utils

import (
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
)

func ReadSome(conn net.Conn) ([]byte, error) {
	timeoutCount := 0

	buf := make([]byte, 1024)
	data := []byte{}
	for {
		conn.SetReadDeadline(time.Now().Add(time.Millisecond * 10))

		n, err := conn.Read(buf)
		if err != nil {
			if err == io.EOF {
				data = append(data, buf[:n]...)

				return data, err
			} else if err, ok := err.(net.Error); ok && err.Timeout() {
				data = append(data, buf[:n]...)

				if timeoutCount > 5 {
					return data, nil
				}

				timeoutCount++
				continue
			}

			return nil, errors.Wrap(err, "error while reading from connection")
		}

		timeoutCount = 0
		data = append(data, buf[:n]...)
	}

}
