package main

import (
	"fmt"
	"time"
)

func generateFileName() string {
	return fmt.Sprintf("%d", generateRandomNumber())
}

func generateRandomNumber() int64 {
	return time.Now().UnixNano()
}
