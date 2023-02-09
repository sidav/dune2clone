package main

import (
	"dune2clone/fibrandom"
	"flag"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"log"
	"os"
	"time"
)

var (
	config yamlConfig
	rnd    fibrandom.FibRandom
)

func main() {
	defer recoverPanicToFile()
	config.initFromFileOrCreate()

	if config.LogToFile {
		f, err := os.OpenFile(fmt.Sprintf("debug_output%s.log", time.Now().Format("2006_01_02_15_04_05")),
			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			panic(fmt.Sprintf("Error opening file: %v", err))
		}
		defer f.Close()
		log.SetOutput(f)
	}

	// geometry.SetDegreesInCircleAmount(100)

	runSanity := flag.Bool("sanity", false, "Perform static data sanity")
	runBalanceTester := flag.Bool("test-balance", false, "Perform balance check")
	flag.Parse()
	nonDefaultFlagSet := false
	if *runSanity {
		performAllDataSanityChecks()
		nonDefaultFlagSet = true
	}
	if *runBalanceTester {
		bt := balanceTester{}
		bt.testCombatBalance()
		nonDefaultFlagSet = true
	}
	if nonDefaultFlagSet {
		return
	}

	rl.SetTraceLogCallback(raylibTraceLogFn)
	rl.InitWindow(int32(WINDOW_W), int32(WINDOW_H), "DAS IST KEIN DUNE 2!")
	rl.SetTargetFPS(int32(config.TargetFPS))
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.SetExitKey(rl.KeyF12)

	rnd.InitDefault()
	rendererInstance := &renderer{}
	loadResources(rendererInstance)

	//for i := 0; i <= 360; i+=10 {
	//	debugWritef("%ddeg is %d sector\n", i, degreeToRotationFrameNumber(i))
	//}
	game := game{}
	game.render = rendererInstance
	game.startGame()

	rl.CloseWindow()
}
