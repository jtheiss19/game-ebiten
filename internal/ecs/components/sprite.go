package components

import (
	"game/internal/ecs"
)

type Sprite struct {
	*ecs.BaseComponent

	Image string

	SpriteSheetX, SpriteSheetY int
}
