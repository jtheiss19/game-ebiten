package components

import (
	"game/internal/ecs"
)

type Player struct {
	*ecs.BaseComponent

	PlayerID int
}
