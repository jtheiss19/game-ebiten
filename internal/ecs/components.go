package ecs

import (
	"reflect"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ID string

type Component interface {
	GetComponentID() ID
	SetComponentID(id ID)
	GenerateComponentID() ID
}

type BaseComponent struct {
	EntityID ID
}

func (bc *BaseComponent) GetComponentID() ID {
	return bc.EntityID
}

func (bc *BaseComponent) SetComponentID(id ID) {
	bc.EntityID = id
}

func (bc *BaseComponent) GenerateComponentID() ID {
	id := ID(uuid.New().String())
	bc.EntityID = id
	return id
}

func ObjectAssign(target interface{}, object interface{}) {
	// object atributes values in target atributes values
	// using pattern matching (https://golang.org/pkg/reflect/#Value.FieldByName)
	// https://stackoverflow.com/questions/35590190/how-to-use-the-spread-operator-in-golang
	t := reflect.ValueOf(target).Elem()
	o := reflect.ValueOf(object).Elem()

	logrus.Infof("got %v and %v", t, o)
	for i := 0; i < o.NumField(); i++ {
		for j := 0; j < t.NumField(); j++ {
			tName := t.Field(j).Type()
			oName := o.Field(i).Type()
			logrus.Infof("names %v and %v", tName, oName)
			if tName == oName {
				t.Field(j).Set(o.Field(i))
			}
		}
	}
}

func ToStructPtr(obj interface{}) interface{} {
	vp := reflect.New(reflect.TypeOf(obj))
	vp.Elem().Set(reflect.ValueOf(obj))
	return vp.Interface()
}
