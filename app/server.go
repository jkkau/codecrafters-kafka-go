package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

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
	mustNot(err, "Failed to bind to port 9092")

	conn, err := l.Accept()
	mustNot(err, "Accept error")

	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	mustNot(err, "Read msg error")

	// parse received data
	request_api_key := binary.BigEndian.Uint16(buf[4:])
	request_api_version := binary.BigEndian.Uint16(buf[6:])
	corelation_id := binary.BigEndian.Uint32(buf[8:])

	// construct response
	// https://kafka.apache.org/protocol.html#The_Messages_ApiVersions
	// https://forum.codecrafters.io/t/question-about-handle-apiversions-requests-stage/1743/5
	responseBuf := make([]byte, 30)
	binary.BigEndian.PutUint32(responseBuf[0:4], 30-4) // message size = sizeof(Header)+sizeof(Body)
	binary.BigEndian.PutUint32(responseBuf[4:8], corelation_id)

	error_code := uint16(35) // code 35 means not supported
	versions, ok := supportedApiVersions[request_api_key]
	if ok {
		for _, v := range versions {
			if v == request_api_version {
				error_code = 0
			}
		}
	}

	binary.BigEndian.PutUint16(responseBuf[8:], error_code)
	responseBuf[10] = 3                              // (number+1) of api_keys
	binary.BigEndian.PutUint16(responseBuf[11:], 18) // api_key APIVersions
	binary.BigEndian.PutUint16(responseBuf[13:], 0)  // min_version
	binary.BigEndian.PutUint16(responseBuf[15:], 4)  // max_version
	responseBuf[17] = 0                              // _tagged_fields
	binary.BigEndian.PutUint16(responseBuf[18:], 75) // api_key DescribeTopicPartitions
	binary.BigEndian.PutUint16(responseBuf[20:], 0)  // min_version
	binary.BigEndian.PutUint16(responseBuf[22:], 0)  // max_version
	responseBuf[24] = 0                              // _tagged_fields
	binary.BigEndian.PutUint32(responseBuf[25:], 0)  // throttle_time_ms
	responseBuf[29] = 0                              // _tagged_fields

	// send message
	_, err = conn.Write(responseBuf)
	mustNot(err, "Write response error")
}
