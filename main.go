package main

import (
	"dune2clone/fibrandom"
	"flag"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var rnd fibrandom.FibRandom

func main() {
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
	rl.SetTargetFPS(RENDERER_DESIRED_FPS)
	rl.SetWindowState(rl.FlagWindowResizable)
	rl.SetExitKey(rl.KeyEscape)

	rnd.InitDefault()
	loadResources()

	//for i := 0; i <= 360; i+=10 {
	//	debugWritef("%ddeg is %d sector\n", i, degreeToRotationFrameNumber(i))
	//}
	game := game{}
	game.startGame()


	rl.CloseWindow()
}
