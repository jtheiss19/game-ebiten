package components

import (
	"game/internal/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type Transform2D struct {
	*ecs.BaseComponent

	WorldPosition mgl32.Vec2
	WorldRotation float64
	WorldScale    mgl32.Vec2
}

func (tf *Transform2D) AddWorldPosition(positionToAdd mgl32.Vec2) {
	tf.WorldPosition = tf.WorldPosition.Add(positionToAdd)
}

func (tf *Transform2D) AddWorldRotation(angle float64) {
	tf.WorldRotation = tf.WorldRotation + angle
}
