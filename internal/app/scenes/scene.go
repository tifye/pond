package scenes

import "github.com/hajimehoshi/ebiten/v2"

type Scene interface {
	Update() Scene
	Draw(target *ebiten.Image)
}
