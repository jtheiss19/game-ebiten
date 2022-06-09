package main

import (
	"game/internal/ecs"
	"game/internal/ecs/systems"
	"game/internal/wrapper"
	_ "image/png"
	"math"
	"time"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/sirupsen/logrus"
)

const (
	windowName = "Server"
)

var (
	windowSize = mgl32.Vec2{100, 100}
	fpsCap     = 60
)

func main() {
	logrus.Info("starting up...")

	logrus.Info("creating base systems...")
	world := ecs.NewWorld()

	newMultiplayerSystem := systems.NewMultiplayerSystem(true)
	world.AddSystem(newMultiplayerSystem)

	newPlayerController := systems.NewPlayerController()
	world.AddSystem(newPlayerController)

	// logrus.Info("creating base objects...")
	// for i := 0; i < 18; i++ {
	// 	objects.NewSprite(world, mgl32.Vec2{float32(i%6 + 1), float32(i/6 + 1)}, true)
	// }

	logrus.Info("Configuring Window...")
	wrapper.ConfigureWindowSize(windowSize)
	wrapper.ConfigureWindowName(windowName)

	logrus.Info("Running game...")
	mainLoop := func() {
		world.Update(0)
	}
	loopOnTimer(fpsCap, mainLoop)
}

func loopOnTimer(TPS int, doSomethingFunc func()) {

	waitTime := time.Second / time.Duration(TPS)

	timeStart := time.Now()
	lastTickDuration := time.Since(timeStart)
	lastTickCount := 1
	for range time.Tick(waitTime) {
		timeStart = time.Now()

		logrus.Debug("took %v to render %v frames with a max frame time of %v\n", lastTickDuration, lastTickCount, waitTime)

		cumulativeWaitTime := waitTime * time.Duration(lastTickCount)

		if lastTickDuration > cumulativeWaitTime {
			logrus.Warnf(
				"took to long (%v > %v) to proccess %v frames from last measurement.",
				lastTickDuration,
				cumulativeWaitTime,
				lastTickCount,
			)
			if lastTickCount > 5 {
				logrus.Error(
					"can't keep up, reducing next call to 1 tick.",
				)
				lastTickCount = 1
				doSomethingFunc()
			} else {
				totalElapsedTicks := lastTickDuration.Seconds() / waitTime.Seconds()
				logrus.Warnf(
					"requires %v extra back to back ticks to catch up",
					math.Ceil(totalElapsedTicks)-1,
				)
				lastTickCount = int(math.Ceil(totalElapsedTicks))
				for i := 0; i < lastTickCount; i++ {
					doSomethingFunc()
				}
			}

		} else {
			lastTickCount = 1
			doSomethingFunc()
		}

		lastTickDuration = time.Since(timeStart)
	}
}
