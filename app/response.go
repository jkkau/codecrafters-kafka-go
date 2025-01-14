package main

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
