package main

import (
	"dune2clone/fibrandom"
	"flag"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	config yamlConfig
	rnd    fibrandom.FibRandom
)

func main() {
	config.initFromFileOrCreate()
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

	rl.InitWindow(int32(WINDOW_W), int32(WINDOW_H), "DAS IST KEIN DUNE 2!")
	rl.SetTargetFPS(int32(config.TargetFPS))
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.SetExitKey(rl.KeyEscape)

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
