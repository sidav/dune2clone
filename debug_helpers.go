package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

func debugWrite(msg string) {
	if config.DebugOutput {
		log.Println(msg)
	}
}

func debugWritef(msg string, args ...interface{}) {
	if config.DebugOutput {
		log.Printf(msg, args...)
	}
}

type debugTimeInfo struct {
	logicName           string
	duration            time.Duration
	maxRecordedDuration time.Duration
	criticalDuration    time.Duration

	// for mean calculation
	totalChanges           int
	meanDurAccumulator     time.Duration
	calculatedMeanDuration int64
}

func raylibTraceLogFn(x int, str string) {
	debugWritef(str)
}

func (dti *debugTimeInfo) setNewValue(dur time.Duration) {
	if dur > dti.maxRecordedDuration {
		dti.maxRecordedDuration = dur
	}
	dti.duration = dur
	dti.meanDurAccumulator += dur

	dti.totalChanges++
	dti.calculatedMeanDuration = int64(dti.meanDurAccumulator)
	dti.calculatedMeanDuration /= int64(dti.totalChanges)
	const changesToResetMean = 100
	if dti.totalChanges == changesToResetMean {
		dti.meanDurAccumulator = 0
		dti.totalChanges = 0
	}
}

func recoverPanicToFile() {
	if x := recover(); x != nil {
		fo, err := os.Create("crash_report.log")
		if err != nil {
			panic(err)
		}
		fo.Write([]byte(fmt.Sprintf("Panic: %v", x)))

		if err := fo.Close(); err != nil {
			panic(err)
		}

		// Panic again for a crash
		panic(x)
	}
}
