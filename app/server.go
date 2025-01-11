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

func initSupportedApiVersions() (supportedApiVersions map[uint16][]uint16) {
	supportedApiVersions = make(map[uint16][]uint16)
	supportedApiVersions[ApiVersions] = []uint16{4}	

	return supportedApiVersions
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	supportedApiVersions := initSupportedApiVersions()

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
	/*
	* Response body for apiKey 18:
	* ApiVersions Response (Version: 0) => error_code [api_versions] 
  	* error_code => INT16
    * api_versions => api_key min_version max_version 
    *   api_key => INT16
    *   min_version => INT16
    *   max_version => INT16
	*
	*/
	responseBuf := make([]byte, 16)
	binary.BigEndian.PutUint32(responseBuf[0:4], 16) // message size
	binary.BigEndian.PutUint32(responseBuf[4:8], corelation_id)

	error_code := uint16(35)  // code 35 means not supported
	versions, ok := supportedApiVersions[request_api_key]
	if ok {
		for _, v := range versions {
			if v == request_api_version {
				error_code = 0
			}
		}
	}

	binary.BigEndian.PutUint16(responseBuf[8:], error_code)
	binary.BigEndian.PutUint16(responseBuf[10:], 18)   // api_key
	binary.BigEndian.PutUint16(responseBuf[12:], 4)    // min_version
	binary.BigEndian.PutUint16(responseBuf[14:], 4)    // max_version

	// send message
	_, err = conn.Write(responseBuf)
	mustNot(err, "Write response error")
}
