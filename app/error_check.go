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
