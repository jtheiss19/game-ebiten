package main

import (
	"fmt"
	"game/internal/ecs"
	"game/internal/ecs/systems"
	"game/internal/wrapper"
	_ "image/png"
	"log"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/sirupsen/logrus"
)

const (
	windowName = "Test Game"
	fpsCap     = 60
)

var (
	windowSize = mgl32.Vec2{1080, 720}
)

type Game struct {
	gameWorld *ecs.World
}

func (g *Game) Update() error {
	g.gameWorld.Update(0)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	msg := fmt.Sprintf(
		"TPS: %0.2f\nFPS: %0.2f",
		ebiten.CurrentTPS(),
		ebiten.CurrentFPS(),
	)

	g.gameWorld.Draw(screen)

	ebitenutil.DebugPrint(screen, msg)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	logrus.Info("starting up...")

	logrus.Info("creating base systems...")
	world := ecs.NewWorld()

	newRenderer := systems.NewEbitenSpriteRenderer()
	world.AddSystem(newRenderer)

	newPlayerController := systems.NewPlayerController()
	world.AddSystem(newPlayerController)

	newMultiplayerSystem := systems.NewMultiplayerSystem(false)
	world.AddSystem(newMultiplayerSystem)

	newInputHandler := systems.NewInputHandler()
	world.AddSystem(newInputHandler)

	logrus.Info("creating base objects...")
	// objects.NewSprite(world, mgl32.Vec2{0, 0})
	// objects.NewSprite(world, mgl32.Vec2{2, 0})
	// objects.NewPlayer(world, false)

	logrus.Info("Configuring Window...")

	wrapper.ConfigureWindowSize(windowSize)
	wrapper.ConfigureWindowName(windowName)

	newGame := &Game{
		gameWorld: world,
	}
	logrus.Info("Running game...")

	err := ebiten.RunGame(newGame)
	if err != nil {
		log.Fatal(err)
	}
}
