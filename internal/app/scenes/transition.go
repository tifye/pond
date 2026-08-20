package scenes

import (
	_ "embed"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed circle_transition.kage
var circleTransitionShaderSrc []byte
var circleTransitionShader *ebiten.Shader

func init() {
	if s, err := ebiten.NewShader(circleTransitionShaderSrc); err != nil {
		panic(err)
	} else {
		circleTransitionShader = s
	}
}

type TransitionScene struct {
	Before Scene
	After  Scene

	before    *ebiten.Image
	after     *ebiten.Image
	auxillery *ebiten.Image

	ScreenWidth  int
	ScreenHeight int

	step float64
	t    float64
	dur  time.Duration
}

func (t *TransitionScene) Initialize() {
	t.dur = time.Second * 50
	t.step = 1 / (t.dur.Seconds() * float64(ebiten.TPS()))

	t.before = ebiten.NewImage(int(t.ScreenWidth), int(t.ScreenHeight))
	t.after = ebiten.NewImage(int(t.ScreenWidth), int(t.ScreenHeight))
	t.auxillery = ebiten.NewImage(int(t.ScreenWidth), int(t.ScreenHeight))
}

func (t *TransitionScene) Update() Scene {
	t.t += t.step

	// We compare strictly larger than so that we always animate
	// the case where t = 1.0.
	if t.t > 1 {
		return t.After
	}

	return t.After.Update()
}

func (t *TransitionScene) Draw(target *ebiten.Image) {
	t.before.Clear()
	t.Before.Draw(t.before)

	t.after.Clear()
	t.After.Draw(t.after)

	t.auxillery.Clear()
	t.auxillery.DrawRectShader(
		t.ScreenWidth,
		t.ScreenHeight,
		circleTransitionShader,
		&ebiten.DrawRectShaderOptions{
			Images: [4]*ebiten.Image{
				t.before,
				t.after,
			},
			Uniforms: map[string]any{
				"Origin":    []float64{float64(t.ScreenWidth) * 0.5, float64(t.ScreenHeight) * 0.5},
				"MinRadius": 100.0,
				"MaxRadius": float64(t.ScreenWidth)*0.5 + float64(t.ScreenHeight)*0.5,
				"T":         t.t,
			},
		},
	)

	opts := &ebiten.DrawImageOptions{}
	opts.Blend = ebiten.BlendSourceIn
	t.auxillery.DrawImage(t.after, opts)

	opts.Blend = ebiten.BlendDestinationOver
	t.auxillery.DrawImage(t.before, opts)

	target.DrawImage(t.auxillery, nil)
}
