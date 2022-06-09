package components

import (
	"game/internal/ecs"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Input struct {
	*ecs.BaseComponent

	Keys []ebiten.Key
}

func (in *Input) IsKeyPressed(keyToCheck ebiten.Key) bool {
	for _, key := range in.Keys {
		if key == keyToCheck {
			return true
		}
	}
	return false
}

func (in *Input) RefreshKeyboardState() {
	in.Keys = inpututil.AppendPressedKeys(in.Keys[:0])
}
