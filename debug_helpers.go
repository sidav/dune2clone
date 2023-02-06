package main

import (
	"fmt"
	"time"
)

func debugWrite(msg string) {
	if config.DebugOutput {
		fmt.Println(msg)
	}
}

func debugWritef(msg string, args ...interface{}) {
	if config.DebugOutput {
		fmt.Printf(msg, args...)
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
