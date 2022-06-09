package systems

import (
	"game/internal/ecs"
	"game/internal/ecs/components"
	"reflect"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

// -------------------- SPECIAL COMPONENTS -------------------------------------------//
type PlayerControllerRequiredComponents struct {
	Entities []*PlayerControllerComponents
}

type PlayerControllerComponents struct {
	PlayerTransform *components.Transform2D
	PlayerInput     *components.Input
	PlayerData      *components.Player
}

// -------------------- Main Component -------------------------------------------//

type PlayerControllerSystem struct {
	trackedEntities *PlayerControllerRequiredComponents
}

func NewPlayerController() *PlayerControllerSystem {
	return &PlayerControllerSystem{
		trackedEntities: &PlayerControllerRequiredComponents{},
	}
}

// -------------------- Custom Functionality -------------------------------------------//

func (pc *PlayerControllerSystem) Update(dt float32) {

	// Handle Input Events
	for _, entity := range pc.trackedEntities.Entities {

		input := entity.PlayerInput

		// Quit
		if input.IsKeyPressed(ebiten.KeyEscape) {
			logrus.Exit(0)
		}

		// Keyboard
		forwardVec := mgl32.Vec2{0, -1}
		rightVec := mgl32.Vec2{1, 0}

		moveDirection := mgl32.Vec2{}
		if input.IsKeyPressed(ebiten.KeyW) {
			moveDirection = moveDirection.Add(forwardVec)
		}
		if input.IsKeyPressed(ebiten.KeyS) {
			moveDirection = moveDirection.Add(forwardVec.Mul(-1))
		}
		if input.IsKeyPressed(ebiten.KeyD) {
			moveDirection = moveDirection.Add(rightVec)
		}
		if input.IsKeyPressed(ebiten.KeyA) {
			moveDirection = moveDirection.Add(rightVec.Mul(-1))
		}

		playerPosition := entity.PlayerTransform
		speed := 0.1
		playerPosition.AddWorldPosition(moveDirection.Mul(float32(speed)))
	}
}

func (pc *PlayerControllerSystem) Draw(screen *ebiten.Image) {
}

// -------------------- BoilerPlate Code -------------------------------------------//

func (pc *PlayerControllerSystem) GetRequiredComponents() []reflect.Type {
	reqComponentsStruct := &PlayerControllerRequiredComponents{}

	v := reflect.ValueOf(reqComponentsStruct).Elem()

	returnTyppc := []reflect.Type{}
	for j := 0; j < v.NumField(); j++ {
		reqField := v.Field(j)
		switch reqField.Type().Kind() {
		case reflect.Slice:
			returnTyppc = append(returnTyppc, reqField.Type().Elem())
		case reflect.Ptr:
			returnTyppc = append(returnTyppc, reqField.Elem().Type())
		default:
			logrus.Error("no field match found")
		}
	}

	return returnTyppc
}

func (pc *PlayerControllerSystem) AddEntity(comps map[reflect.Type]ecs.Component) {
	logrus.Trace("adding entity to Player Controller System")

	for _, reqComp := range pc.GetRequiredComponents() {
		if ecs.SatisfySystemRequirements(comps, reqComp) {

			f := reflect.ValueOf(pc.trackedEntities).Elem()
			for j := 0; j < f.NumField(); j++ {
				reqField := f.Field(j)
				reqFieldType := reqField.Type().Elem()
				if reqFieldType == reqComp {
					newReqFieldEntry := reflect.New(reqFieldType.Elem())
					ecs.Fill(newReqFieldEntry, comps)

					reqFieldElem := reqField

					logrus.Debug("Setting Player Controller System entity element")
					reqFieldElem.Set(reflect.Append(reqFieldElem, newReqFieldEntry))
				}
			}
		}
	}

}

func (pc *PlayerControllerSystem) RemoveEntity(id ecs.ID) {
	// Called when a entity needs removed
}

func (pc *PlayerControllerSystem) Initilizer() {
	// some code thats ran on world join
}
