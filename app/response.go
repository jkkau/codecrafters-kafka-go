package main

import (
	"encoding/binary"
)

const (
	MAX_BUFFER_SIZE = 1024
)

// construct response
// https://kafka.apache.org/protocol.html#The_Messages_ApiVersions
// https://forum.codecrafters.io/t/question-about-handle-apiversions-requests-stage/1743/5
func BuildResponse(req *Request) []byte {
	responseBuf := make([]byte, MAX_BUFFER_SIZE)
	offset := 4

	binary.BigEndian.PutUint32(responseBuf[offset:], uint32(req.CorelationId))
	offset += 4

	error_code := buildErrorCode(req)
	binary.BigEndian.PutUint16(responseBuf[offset:], error_code)
	offset += 2

	offset = buildApiKeys(responseBuf, offset)

	binary.BigEndian.PutUint32(responseBuf[offset:], 0) // throttle_time_ms
	offset += 4
	responseBuf[offset] = 0 // _tagged_fields
	offset += 1
	// Now offset is the current length of response message
	binary.BigEndian.PutUint32(responseBuf[0:4], uint32(offset-4)) // message size field = sizeof(header)+sizeof(Body)

	return responseBuf[:offset]
}

func buildErrorCode(req *Request) uint16 {
	error_code := uint16(35) // code 35 means not supported
	versions, ok := SupportedApiVersions[req.RequestApiKey]
	if ok {
		for _, v := range versions {
			if v == req.RequestApiVersion {
				error_code = 0
			}
		}
	}
	return error_code
}

func buildApiKeys(responseBuf []byte, offset int) int {
	responseBuf[offset] = byte(len(SupportedApiVersions) + 1) // (number+1) of api_keys
	offset += 1
	for apiKey, minMaxVersion := range SupportedApiVersions {
		binary.BigEndian.PutUint16(responseBuf[offset:], uint16(apiKey)) // api_key
		offset += 2
		binary.BigEndian.PutUint16(responseBuf[offset:], uint16(minMaxVersion[0]))  // min_version
		offset += 2
		binary.BigEndian.PutUint16(responseBuf[offset:], uint16(minMaxVersion[1]))  // max_version
		offset += 2
		responseBuf[offset] = 0                              // _tagged_fields
		offset+=1
	}
	return offset
}
