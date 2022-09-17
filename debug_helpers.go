package main

import (
	"fmt"
	"time"
)

func debugWrite(msg string) {
	if DEBUG_OUTPUT {
		fmt.Println(msg)
	}
}

func debugWritef(msg string, args ...interface{}) {
	if DEBUG_OUTPUT {
		fmt.Printf(msg, args...)
	}
}

type debugTimeInfo struct {
	logicName        string
	duration         time.Duration
	criticalDuration time.Duration
}
