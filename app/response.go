package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type Response struct {
	Size            int32
	CorelationId    int32
	ErrorCode       int16
	NumberOfApiKeys byte
	ApiKey          int16
	MinVersion      int16
	MaxVersion      int16
	TaggedFields1   byte
	ThrottleTimeMs  int32
	TaggedFields2   byte
}

func (r *Response) Send(conn net.Conn) (int, error) {
	buf := new(bytes.Buffer)
	// send message
	err := binary.Write(buf, binary.BigEndian, *r)
	mustNot(err, "Write binary error")

	totalWritten := 0
	responseBuf := buf.Bytes()
	for totalWritten < len(responseBuf) {
		n, err := conn.Write(responseBuf[totalWritten:])
		if err != nil {
			return totalWritten, err
		}
		totalWritten += n
	}

	fmt.Printf("send total message: %d bytes\n", totalWritten)
	return totalWritten, nil
}
