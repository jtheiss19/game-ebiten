package wrapper

import (
	"image"
	"log"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

var (
	desiredDimensions = mgl32.Vec2{100, 100}
	offSet            = mgl32.Vec2{}
)

type Texture struct {
	image *ebiten.Image
	cols  int
	rows  int
}

type Transform struct {
	Position mgl32.Vec2
	Scale    mgl32.Vec2
	Rotation float64
}

func NewTexture(filepath string, sheetSize ...mgl32.Vec2) *Texture {
	dimensions := mgl32.Vec2{1, 1}
	if len(sheetSize) > 0 {
		dimensions = sheetSize[0]
	}

	infile, err := os.Open(filepath)
	if err != nil {
		logrus.Error(err)
	}
	defer infile.Close()

	img, _, err := image.Decode(infile)
	if err != nil {
		log.Fatal(err)
	}

	origEbitenImage := ebiten.NewImageFromImage(img)

	w, h := origEbitenImage.Size()
	ebitenImage := ebiten.NewImage(w, h)

	op := &ebiten.DrawImageOptions{}

	sX, sY := desiredDimensions.X()/float32(w)*dimensions.X(), desiredDimensions.Y()/float32(h)*dimensions.Y()
	op.GeoM.Scale(float64(sX), float64(sY))

	ebitenImage.DrawImage(origEbitenImage, op)

	return &Texture{
		image: ebitenImage,
	}
}

func (tex *Texture) Draw(screen *ebiten.Image, drawWhere Transform, camera Transform) {
	position := drawWhere.Position
	scale := drawWhere.Scale
	rotation := drawWhere.Rotation
	drawOptions := getImageDrawOptions(position, scale, rotation)
	screen.DrawImage(tex.image, &drawOptions)
}

func (tex *Texture) DrawItem(screen *ebiten.Image, drawWhere Transform, camera Transform, sheetPositionX, sheetPositionY int) {
	position := drawWhere.Position.Sub(camera.Position)
	scale := mgl32.Vec2{drawWhere.Scale.X() * camera.Scale.X(), drawWhere.Scale.Y() * camera.Scale.Y()}
	rotation := drawWhere.Rotation - camera.Rotation

	drawOptions := getImageDrawOptions(position, scale, rotation)

	sheetPositionX = sheetPositionX - 1
	sheetPositionY = sheetPositionY - 1

	subImage := tex.image.SubImage(
		image.Rect(
			sheetPositionX*int(desiredDimensions.X()),
			sheetPositionY*int(desiredDimensions.Y()),
			(sheetPositionX+1)*int(desiredDimensions.X()),
			(sheetPositionY+1)*int(desiredDimensions.Y())),
	).(*ebiten.Image)

	screen.DrawImage(subImage, &drawOptions)
}

func getImageDrawOptions(position, scale mgl32.Vec2, rotation float64) ebiten.DrawImageOptions {
	op := ebiten.DrawImageOptions{}
	op.GeoM.Reset()
	sX, sY := scale.Elem()
	op.GeoM.Scale(float64(sX), float64(sY))

	position = position.Mul(100).Add(offSet)
	pX, pY := position.Elem()
	op.GeoM.Translate(float64(pX), float64(pY))
	return op
}
