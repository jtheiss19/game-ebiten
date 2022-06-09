package objects

import (
	"game/internal/ecs"
	"game/internal/ecs/components"

	"github.com/hajimehoshi/ebiten/v2"
)

func NewPlayer(world *ecs.World, isNetworked bool, playerID int) ecs.ID {
	compCameraTransform := &components.Transform2D{
		BaseComponent: &ecs.BaseComponent{},
		WorldPosition: [2]float32{0, 0},
		WorldRotation: 0,
		WorldScale:    [2]float32{1, 1},
	}
	compCamera := &components.Camera{
		BaseComponent: &ecs.BaseComponent{},
		IsActive:      true,
	}
	compInput := &components.Input{
		BaseComponent: &ecs.BaseComponent{},
		Keys:          []ebiten.Key{},
	}
	compSprite := &components.Sprite{
		BaseComponent: &ecs.BaseComponent{},
		Image:         "/some/file/path",
		SpriteSheetX:  2,
		SpriteSheetY:  2,
	}
	player := &components.Player{
		BaseComponent: &ecs.BaseComponent{},
		PlayerID:      playerID,
	}
	comps := []ecs.Component{compCamera, compCameraTransform, compInput, compSprite, player}
	if isNetworked {
		network := &components.Network{
			BaseComponent: &ecs.BaseComponent{},
		}
		comps = append(comps, network)
	}
	return world.AddEntity(comps)
}
