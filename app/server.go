package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

func main() {
	fmt.Println("start!")

	l, err := net.Listen("tcp", "0.0.0.0:9092")
	mustNot(err, "Failed to bind to port 9092")

	conn, err := l.Accept()
	mustNot(err, "Accept error")

	var req Request
	err = binary.Read(conn, binary.BigEndian, &req)
	mustNot(err, "read conn error")

	responseBuf := BuildResponse(&req)

	// send message
	writtenSize := 0
	for writtenSize < len(responseBuf) {
		len, err := conn.Write(responseBuf)
		mustNot(err, "Write response error")
		writtenSize += len
	}
	fmt.Printf("send %d bytes response success.\n", writtenSize)
}
