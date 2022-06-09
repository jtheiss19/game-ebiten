package ecs

import (
	"reflect"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/sirupsen/logrus"
)

var (
	activeWorld *World
)

type World struct {
	systems      []system
	entityLookup map[ID]map[reflect.Type]Component
}

func NewWorld() *World {
	activeWorld = &World{
		systems:      []system{},
		entityLookup: map[ID]map[reflect.Type]Component{},
	}
	return activeWorld
}

func GetActiveWorld() *World {
	return activeWorld
}

func (wrld *World) Update(dt float32) {
	for _, system := range wrld.systems {
		system.Update(dt)
	}
}

func (wrld *World) Draw(screen *ebiten.Image) {
	for _, system := range wrld.systems {
		system.Draw(screen)
	}
}

func (wrld *World) AddComponent(comp Component) {
	logrus.Trace("adding component")
	id := comp.GetComponentID()

	if comps, ok := wrld.entityLookup[id]; ok {
		comps[reflect.TypeOf(comp)] = comp
		wrld.entityLookup[id] = comps
		wrld.checkCompsForNewSystemMatch(comps)
	} else {
		compsMap := map[reflect.Type]Component{}
		compsMap[reflect.TypeOf(comp)] = comp
		wrld.entityLookup[id] = compsMap
		wrld.checkCompsForNewSystemMatch(compsMap)
	}

}

func (wrld *World) RemoveCompoent(comp Component) {
	logrus.Info("removing component")
	id := comp.GetComponentID()

	if comps, ok := wrld.entityLookup[id]; ok {
		newComps := map[reflect.Type]Component{}
		for _, compItem := range comps {
			if reflect.TypeOf(compItem) != reflect.TypeOf(comp) {
				newComps[reflect.TypeOf(comp)] = comp
			}
		}
		wrld.entityLookup[id] = newComps

		for _, system := range wrld.systems {
			for _, reqComp := range system.GetRequiredComponents() {
				if !SatisfySystemRequirements(newComps, reqComp) {
					system.RemoveEntity(id)
				}
			}
		}
	}
}

func (wrld *World) AddEntity(comps []Component) ID {
	id := ID(uuid.New().String())
	logrus.Info("Creating and adding new entity: ", id)
	for _, comp := range comps {
		comp.SetComponentID(id)
	}

	compsMap := map[reflect.Type]Component{}
	for _, comp := range comps {
		compsMap[reflect.TypeOf(comp)] = comp
	}

	wrld.entityLookup[id] = compsMap

	wrld.checkCompsForNewSystemMatch(compsMap)

	return id
}

func (wrld *World) GetEntity(entityID ID) map[reflect.Type]Component {
	if compsMap, ok := wrld.entityLookup[entityID]; !ok {
		return map[reflect.Type]Component{}
	} else {
		// Create the target map
		targetMap := map[reflect.Type]Component{}

		// Copy from the original map to the target map
		for key, value := range compsMap {
			targetMap[key] = value
		}
		return targetMap
	}
}

func (wrld *World) RemoveEntity(id ID) {
	logrus.Info("removing entity")
	for _, system := range wrld.systems {
		system.RemoveEntity(id)
	}
}

func (wrld *World) AddSystem(system system) {
	logrus.Info("initing system")
	system.Initilizer()
	logrus.Info("adding system")
	wrld.systems = append(wrld.systems, system)
}

func (wrld *World) checkCompsForNewSystemMatch(comps map[reflect.Type]Component) {
	logrus.Trace("checking for system matches")
	for _, system := range wrld.systems {
		requirementComponents := system.GetRequiredComponents()
		for _, requirement := range requirementComponents {
			logrus.Debugf("Evaluating requirements for %v for comps %v", requirement, comps)
			if SatisfySystemRequirements(comps, requirement) {
				logrus.Debug("found system match")
				system.AddEntity(comps)
			}
		}
	}
}

func SatisfySystemRequirements(comps map[reflect.Type]Component, reqComps reflect.Type) bool {
	if reqComps.Kind() == reflect.Ptr {
		reqComps = reqComps.Elem()
	}

	for j := 0; j < reqComps.NumField(); j++ {
		reqCompFieldType := reqComps.Field(j).Type
		if !func() bool { // return if there is a compoenent that doesn't have a match
			if _, ok := comps[reqCompFieldType]; ok {
				return true
			} else {
				logrus.Debugf("The Given Comp of Type %v does not exist in comps %v", reqCompFieldType, comps)
			}
			return false
		}() {
			return false // Stop execution here as we found a type that is requeid but doesn't exist
		} else {
			continue // Keep checking each required type to make sure it exists
		}
	}
	return true // congrats all the required types have been found in the structure.
}
