package wrapper

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

var (
	windowSize   = mgl32.Vec2{}
	windowCenter = mgl32.Vec2{}

	windowName string
)

func ConfigureWindowSize(newWindowSize mgl32.Vec2) {
	logrus.Info("Configuring Window...")
	windowSize = newWindowSize
	windowCenter = windowSize.Mul(0.5)
	offSet = windowCenter.Sub(desiredDimensions.Mul(0.5))

	w, h := newWindowSize.Elem()
	ebiten.SetWindowSize(int(w), int(h))

}

func ConfigureWindowName(newWindowName string) {
	windowName = newWindowName
	ebiten.SetWindowTitle(newWindowName)
}
