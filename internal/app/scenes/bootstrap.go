package scenes

import "github.com/hajimehoshi/ebiten/v2"

type BookstrapScene struct {
	Next Scene
}

// Draw implements [Scene].
func (b *BookstrapScene) Draw(target *ebiten.Image) {
	panic("unimplemented")
}

type Initializer interface {
	Initialize()
}

// Update implements [Scene].
func (b *BookstrapScene) Update() Scene {
	return b.Next
}

var _ Scene = (*BookstrapScene)(nil)
