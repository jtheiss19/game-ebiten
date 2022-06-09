package systems

import (
	"game/internal/ecs"
	"game/internal/ecs/components"
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

// -------------------- SPECIAL COMPONENTS -------------------------------------------//
type InputHandlerRequiredComponents struct {
	Entities []*InputHandlerComponents
}

type InputHandlerComponents struct {
	PlayerInput *components.Input
}

// -------------------- Main Component -------------------------------------------//

type InputHandlerSystem struct {
	trackedEntities *InputHandlerRequiredComponents
}

func NewInputHandler() *InputHandlerSystem {
	return &InputHandlerSystem{
		trackedEntities: &InputHandlerRequiredComponents{},
	}
}

// -------------------- Custom Functionality -------------------------------------------//

func (pc *InputHandlerSystem) Update(dt float32) {

	// Handle Input Updates
	for _, entity := range pc.trackedEntities.Entities {
		input := entity.PlayerInput
		input.RefreshKeyboardState()
	}
}

func (pc *InputHandlerSystem) Draw(screen *ebiten.Image) {
}

// -------------------- BoilerPlate Code -------------------------------------------//

func (pc *InputHandlerSystem) GetRequiredComponents() []reflect.Type {
	reqComponentsStruct := &InputHandlerRequiredComponents{}

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

func (pc *InputHandlerSystem) AddEntity(comps map[reflect.Type]ecs.Component) {
	logrus.Trace("adding entity to Input Handler System")

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

func (pc *InputHandlerSystem) RemoveEntity(id ecs.ID) {
	// Called when a entity needs removed
}

func (pc *InputHandlerSystem) Initilizer() {
	// some code thats ran on world join
}
