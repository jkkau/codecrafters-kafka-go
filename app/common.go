package main

import (
	"fmt"
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
	ApiVersions             int16 = 18
	DescribeTopicPartitions int16 = 75
)

var SupportedApiVersions map[int16][2]int16

func init() {
	SupportedApiVersions = make(map[int16][2]int16)
	SupportedApiVersions[ApiVersions] = [2]int16{0, 4}             // minVersion: 0, maxVersion: 4
	SupportedApiVersions[DescribeTopicPartitions] = [2]int16{0, 0} // minVersion: 0, maxVersion: 0
}
