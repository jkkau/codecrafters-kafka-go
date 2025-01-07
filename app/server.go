package main

import (
	"fmt"
	"net"
	"os"
	"encoding/binary"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func mustNot(err error, msg string) {
	if err != nil {
		fmt.Println(msg)
		os.Exit(1)
	}
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:9092")
	mustNot(err, "Failed to bind to port 9092");

	conn, err := l.Accept()
	mustNot(err, "Accept error")

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	mustNot(err, "Read msg error");

	responseBuf := make([]byte, 8)
	binary.BigEndian.PutUint32(responseBuf[0:4], 0)
	binary.BigEndian.PutUint32(responseBuf[4:8], 7)

	_, err = conn.Write(responseBuf)
	mustNot(err, "Write response error")
}
