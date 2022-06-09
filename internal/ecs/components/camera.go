package components

import (
	"game/internal/ecs"

	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	*ecs.BaseComponent

	LookDirection mgl32.Vec3
	IsActive      bool
}
