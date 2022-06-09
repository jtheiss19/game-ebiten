package ecs

import (
	"reflect"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

type system interface {
	Update(dt float32)
	Draw(screen *ebiten.Image)
	GetRequiredComponents() []reflect.Type
	AddEntity(comps map[reflect.Type]Component)
	RemoveEntity(id ID)
	Initilizer()
}

type RequiredComponents interface{}

func Fill(thingToFill reflect.Value, compsToFillWith map[reflect.Type]Component) {
	v := thingToFill.Elem()

	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)
		for _, comp := range compsToFillWith {
			if reflect.TypeOf(comp) == f.Type() {
				logrus.Trace("system field match, setting")
				f.Set(reflect.ValueOf(comp))
			}
		}
	}
}
