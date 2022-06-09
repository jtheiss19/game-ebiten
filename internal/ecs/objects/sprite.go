package objects

import (
	"game/internal/ecs"
	"game/internal/ecs/components"

	"github.com/go-gl/mathgl/mgl32"
)

func NewSprite(world *ecs.World, position mgl32.Vec2, isNetworked bool) ecs.ID {
	// CUBE
	compTransform := &components.Transform2D{
		BaseComponent: &ecs.BaseComponent{},
		WorldPosition: position,
		WorldRotation: 0,
		WorldScale:    [2]float32{1, 1},
	}
	compSprite := &components.Sprite{
		BaseComponent: &ecs.BaseComponent{},
		Image:         "/some/file/path",
		SpriteSheetX:  1,
		SpriteSheetY:  1,
	}
	comps := []ecs.Component{compSprite, compTransform}
	if isNetworked {
		network := &components.Network{
			BaseComponent: &ecs.BaseComponent{},
		}
		comps = append(comps, network)
	}
	return world.AddEntity(comps)
}
