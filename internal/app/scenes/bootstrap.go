package scenes

import "github.com/hajimehoshi/ebiten/v2"

type BookstrapScene struct {
	Next Scene
}

// Draw implements [Scene].
func (b *BookstrapScene) Draw(target *ebiten.Image) {
	panic("unimplemented")
}

type initializer interface {
	Initialize()
}

// Update implements [Scene].
func (b *BookstrapScene) Update() Scene {
	if g, ok := b.Next.(initializer); ok {
		g.Initialize()
	}

	return b.Next
}

var _ Scene = (*BookstrapScene)(nil)
