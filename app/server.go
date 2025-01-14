package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
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
	ApiVersions int16 = 18
)

func initSupportedApiVersions() (supportedApiVersions map[int16][]int16) {
	supportedApiVersions = make(map[int16][]int16)
	supportedApiVersions[ApiVersions] = []int16{4}

	return supportedApiVersions
}

func main() {
	fmt.Println("start work.")
	// supportedApiVersions := initSupportedApiVersions()
	_ = initSupportedApiVersions()

	l, err := net.Listen("tcp", "0.0.0.0:9092")
	mustNot(err, "Failed to bind to port 9092")

	conn, err := l.Accept()
	mustNot(err, "Accept error")

	var req Request
	err = binary.Read(conn, binary.BigEndian, &req)
	mustNot(err, "read conn error")
	fmt.Printf("receive req: %v\n", req)

	// construct response
	// https://kafka.apache.org/protocol.html#The_Messages_ApiVersions
	// https://forum.codecrafters.io/t/question-about-handle-apiversions-requests-stage/1743/5
	// error_code := int16(35) // code 35 means not supported
	// versions, ok := supportedApiVersions[req.RequestApiKey]
	// if ok {
	// 	for _, v := range versions {
	// 		if v == req.RequestApiVersion {
	// 			error_code = 0
	// 		}
	// 	}
	// }
	error_code := int16(0)
	rsp := Response{
		Size:            23 - 4,
		CorelationId:    req.CorelationId,
		ErrorCode:       error_code,
		NumberOfApiKeys: 2,
		ApiKey:          18,
		MinVersion:      4,
		MaxVersion:      4,
		TaggedFields1:   0,
		ThrottleTimeMs:  0,
		TaggedFields2:   0,
	}
	fmt.Printf("construct rsp: %v\n", rsp)

	// send message
	err = binary.Write(conn, binary.BigEndian, rsp)
	mustNot(err, "Write response error")
}
