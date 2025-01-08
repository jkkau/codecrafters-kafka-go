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

// Api keys definition
const (
	ApiVersions uint16 = 18
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	supportedApiVersions := make(map[uint16][]uint16)
	supportedApiVersions[ApiVersions] = []uint16{4}

	l, err := net.Listen("tcp", "0.0.0.0:9092")
	mustNot(err, "Failed to bind to port 9092");

	conn, err := l.Accept()
	mustNot(err, "Accept error")

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	mustNot(err, "Read msg error");

	// parse received data
	request_api_key := binary.BigEndian.Uint16(buf[4:])
	request_api_version := binary.BigEndian.Uint16(buf[6:])
	corelation_id := binary.BigEndian.Uint32(buf[8:])

	// construct response
	responseBuf := make([]byte, 10)
	binary.BigEndian.PutUint32(responseBuf[0:4], 0)
	binary.BigEndian.PutUint32(responseBuf[4:8], corelation_id)

	error_code := uint16(35)  // not supported
	versions, ok := supportedApiVersions[request_api_key]
	if ok {
		for _, v := range versions {
			if v == request_api_version {
				error_code = 0
			}
		}
	}
	binary.BigEndian.PutUint16(responseBuf[8:], error_code)

	// send message
	_, err = conn.Write(responseBuf)
	mustNot(err, "Write response error")
}
