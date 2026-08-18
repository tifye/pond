package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tifye/pond/internal/app/scenes"
	"github.com/tifye/pond/internal/entity"
	"github.com/tifye/pond/internal/test"
)

func main() {
	app := test.NewTestApp(test.TestAppConfig{})
	app.Run(&scenes.BookstrapScene{
		Next: &scene{
			sprite: entity.NewDragonFlySprite(200, 200),
		},
	})
}

type scene struct {
	sprite *ebiten.Image
}

func (s *scene) Update() scenes.Scene {
	return nil
}

func (s *scene) Draw(target *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	target.DrawImage(s.sprite, opts)
}
