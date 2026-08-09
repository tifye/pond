package app

import (
	_ "embed"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tifye/pond/internal/entity"
	"github.com/tifye/pond/pkg/agent"
	"github.com/tifye/pond/pkg/mathutil"
)

//go:embed pixel_shader.go
var pixelShader []byte

//go:embed shadow_shader.go
var shadowShaderData []byte
var shadowShader *ebiten.Shader

func init() {
	if s, err := ebiten.NewShader(shadowShaderData); err != nil {
		panic(err)
	} else {
		shadowShader = s
	}
}

type App struct {
	dt      float64
	elapsed float64

	debugColor color.Color

	mousePos                     mathutil.Point
	screenCenterX, screenCenterY float64

	agents *agent.Agents

	pixelShader *ebiten.Shader

	fish      []*entity.Fish
	fishLayer *ebiten.Image

	offscreen *ebiten.Image

	lilypads *entity.LilyPads
	patches  *entity.LilypadPatches
}

func NewApp() *App {
	w, h := ebiten.Monitor().Size()
	// Transpacy doesnt work when setting window to exact size
	// of monitor. My guess is that it thinks that nothing needs to render
	// behind it and that changes the properties of the window
	ebiten.SetWindowSize(w-1, h-1)
	// ebiten.SetFullscreen(true)
	ebiten.SetWindowMousePassthrough(true)
	ebiten.SetWindowDecorated(false)
	ebiten.SetCursorMode(ebiten.CursorModeHidden)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowTitle("Pond")

	screenCenterX := float64(w / 2)
	screenCenterY := float64(h / 2)

	ps, err := ebiten.NewShader(pixelShader)
	if err != nil {
		panic(err)
	}

	fish := make([]*entity.Fish, 2)
	for i := range fish {
		fish[i] = entity.NewFish(entity.FishConfig{
			SpawnPoint: mathutil.Point{
				X: screenCenterX,
				Y: screenCenterY,
			},
			BoneChainUpdater: entity.BoneChainFABRIKUpdater{
				MinAngle: math.Pi * 0.85,
			},
		})
	}

	agents := agent.NewAgents(uint(len(fish)))
	agents.AddBehaviour(agent.NewWander(agents.Num(), 1, 50, 50, math.Pi))
	agents.AddBehaviour(agent.Boundry(float64(w), float64(h), 200, 0.05))

	lpMinHeight := 25.0
	lpMaxHeight := 50.0
	nlp := 100
	lpp := make([]mathutil.Point, 0, nlp)
	lps := make([]int, 0, nlp)
	lpr := make([]float64, 0, nlp)
	for range nlp {
		lpp = append(lpp, mathutil.Point{
			X: rand.Float64() * float64(w),
			Y: rand.Float64() * float64(h),
		})

		rn := rand.Float64()
		lps = append(lps, int(lpMinHeight+(lpMaxHeight-lpMinHeight)*rn))
		lpr = append(lpr, rn*math.Pi*2)
	}

	r := rand.New(rand.NewPCG(0, 0))
	patches := entity.GenerateLilypadPatches(r, w, h, entity.LilypadPatchesOptions{
		MinAmount: 2,
		MaxAmount: 5,
		MinSize:   200,
		MaxSize:   500,
	})

	return &App{
		debugColor: color.RGBA{R: 125, G: 200, B: 85, A: 255},

		screenCenterX: screenCenterX,
		screenCenterY: screenCenterY,
		dt:            1.0 / 60.0, // default for ebiten. It runs on fixed update rate

		fish:   fish,
		agents: agents,

		pixelShader: ps,
		offscreen:   ebiten.NewImage(w, h),
		fishLayer:   ebiten.NewImage(w, h),

		lilypads: entity.NewLilyPads(lpp, lps, lpr),
		patches:  patches,
	}
}

func (a *App) Run() error {
	runOpts := ebiten.RunGameOptions{
		ScreenTransparent: true,
	}
	return ebiten.RunGameWithOptions(a, &runOpts)
}

func (a *App) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	a.elapsed += a.dt

	mx, my := ebiten.CursorPosition()
	a.mousePos.X = float64(mx)
	a.mousePos.Y = float64(my)

	a.agents.Update(a.elapsed, a.dt)

	for i := range a.agents.Num() {
		a.fish[i].Position = a.agents.Position(uint(i))
		a.fish[i].Update()
	}

	// a.fish[len(a.fish)-1].Position = a.mousePos
	// a.fish[len(a.fish)-1].Update()

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	mx, my := ebiten.CursorPosition()

	opts := &ebiten.DrawRectShaderOptions{}
	opts.Uniforms = map[string]any{
		"Time":   float32(ebiten.Tick()) / float32(ebiten.TPS()),
		"Cursor": []float32{float32(mx), float32(my)},
	}

	a.fishLayer.Clear()

	for _, f := range a.fish {
		f.Draw(a.fishLayer)
	}

	a.lilypads.Draw(a.fishLayer)

	// for i, p := range a.patches.Positions {
	// 	s := a.patches.Sizes[i]
	// 	debugCircle(screen, p, colors.Rose600, 5)

	// 	vector.FillRect(
	// 		screen,
	// 		float32(p.X),
	// 		float32(p.Y),
	// 		float32(s.X),
	// 		float32(s.Y),
	// 		colors.WithAlpha(colors.Violet600, 100),
	// 		false,
	// 	)

	// }

	screen.DrawRectShader(
		screen.Bounds().Dx(),
		screen.Bounds().Dy(),
		shadowShader,
		&ebiten.DrawRectShaderOptions{
			Images: [4]*ebiten.Image{
				a.fishLayer,
			},
			Uniforms: map[string]any{
				"Offset": []float32{10.0, -15.0},
			},
		},
	)

	screen.DrawImage(a.fishLayer, nil)

	a.offscreen.DrawRectShader(
		screen.Bounds().Dx(),
		screen.Bounds().Dy(),
		a.pixelShader,
		&ebiten.DrawRectShaderOptions{
			Images: [4]*ebiten.Image{
				screen,
			},
			Uniforms: map[string]any{
				"Resolution": []float32{float32(screen.Bounds().Dx()), float32(screen.Bounds().Dy())},
			},
		},
	)

	screen.Clear()
	screen.DrawImage(a.offscreen, nil)

	a.offscreen.Clear()
}

func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.Monitor().Size()
}

func debugCircle(t *ebiten.Image, p mathutil.Point, c color.Color, s float32) {
	vector.FillCircle(
		t,
		float32(p.X),
		float32(p.Y),
		s,
		c,
		false,
	)
}
