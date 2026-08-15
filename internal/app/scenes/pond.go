package scenes

import (
	_ "embed"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tifye/pond/internal/entity"
	"github.com/tifye/pond/internal/entity/lilypad"
	"github.com/tifye/pond/pkg/agent"
	"github.com/tifye/pond/pkg/mathutil"
)

//go:embed pixel.kage
var pixelShaderSrc []byte
var pixelShader *ebiten.Shader

//go:embed crt.kage
var crtShaderSrc []byte
var crtShader *ebiten.Shader

//go:embed waves.kage
var wavesShaderSrc []byte
var wavesShader *ebiten.Shader

//go:embed vingette.kage
var vingetteShaderSrc []byte
var vingetteShader *ebiten.Shader

//go:embed shadow.kage
var shadowShaderSrc []byte
var shadowShader *ebiten.Shader

func init() {
	if s, err := ebiten.NewShader(shadowShaderSrc); err != nil {
		panic(err)
	} else {
		shadowShader = s
	}

	if s, err := ebiten.NewShader(pixelShaderSrc); err != nil {
		panic(err)
	} else {
		pixelShader = s
	}

	if s, err := ebiten.NewShader(vingetteShaderSrc); err != nil {
		panic(err)
	} else {
		vingetteShader = s
	}

	if s, err := ebiten.NewShader(wavesShaderSrc); err != nil {
		panic(err)
	} else {
		wavesShader = s
	}

	if s, err := ebiten.NewShader(crtShaderSrc); err != nil {
		panic(err)
	} else {
		crtShader = s
	}
}

type PondScene struct {
	ScreenWidth  int
	ScreenHeight int

	buffer *ebiten.Image

	agents *agent.Agents
	fish   []*entity.Fish

	shadowLayer      *ebiten.Image
	shadowShaderOpts *ebiten.DrawRectShaderOptions

	lilypads   *lilypad.Lilypads
	flowers    *lilypad.Flowers
	background *ebiten.Image
}

func (p *PondScene) Initialize(defaultFish []*entity.Fish) {
	p.initFish(defaultFish)
	p.initLilypads()

	p.shadowLayer = ebiten.NewImage(p.ScreenWidth, p.ScreenHeight)
	p.shadowShaderOpts = &ebiten.DrawRectShaderOptions{
		Images: [4]*ebiten.Image{
			p.shadowLayer,
		},
		Uniforms: map[string]any{
			"Offset": []float32{10.0, -15.0},
		},
	}

	p.buffer = ebiten.NewImage(p.ScreenWidth, p.ScreenHeight)
}

func (p *PondScene) initFish(defaultFish []*entity.Fish) {
	numFish := max(2, len(defaultFish))

	p.agents = agent.NewAgents(uint(numFish))
	p.agents.AddBehaviour(agent.NewWander(p.agents.Num(), 1, 50, 50, math.Pi))
	p.agents.AddBehaviour(agent.Boundry(float64(p.ScreenWidth), float64(p.ScreenHeight), 200, 0.05))

	p.fish = make([]*entity.Fish, 0, numFish)

	for i := range defaultFish {
		p.fish = append(p.fish, defaultFish[i])
		p.agents.Set(uint(i),
			defaultFish[i].Position,
			mathutil.NewPoint(0, 0),
			mathutil.NewPoint(0, 0),
		)

	}

	for len(p.fish) < numFish {
		p.fish = append(p.fish, entity.NewFish(entity.FishConfig{
			SpawnPoint: mathutil.Point{
				X: float64(p.ScreenWidth) * 0.5,
				Y: float64(p.ScreenHeight) * 0.5,
			},
			BoneChainUpdater: entity.BoneChainFABRIKUpdater{
				MinAngle: math.Pi * 0.85,
			},
		}))
	}
}

func (p *PondScene) initLilypads() {
	if p.lilypads != nil {
		return
	}

	r := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

	perlin := mathutil.NewPerlin(r)
	renderScale := 2.0
	noiseWidth := int(float64(p.ScreenWidth) / renderScale)
	noiseHeight := int(float64(p.ScreenHeight) / renderScale)

	noiesMap := perlin.Map2D(noiseWidth, noiseHeight, 0.006)
	p.lilypads = lilypad.NewUsingNoiseThreshold(noiesMap, noiseWidth, noiseHeight, 2, 0.5, lilypad.DefaultFromNoiseConfig)

	// noise := ebiten.NewImage(noiseWidth, noiseHeight)
	// noise.WritePixels(PixelBufferFromNoise(noiesMap, noiseWidth, noiseHeight, 2.0))
}

func (p *PondScene) Update() Scene {
	elapsed := float64(ebiten.Tick()) / float64(ebiten.TPS())
	delta := 1 / float64(ebiten.TPS())

	p.agents.Update(elapsed, delta)

	for i := range p.agents.Num() {
		p.fish[i].Position = p.agents.Position(uint(i))
		p.fish[i].Update()
	}

	return nil
}

func (p *PondScene) Draw(target *ebiten.Image) {
	p.shadowLayer.Clear()

	for _, f := range p.fish {
		f.Draw(p.shadowLayer)
	}

	p.lilypads.Draw(p.shadowLayer)
	p.flowers.Draw(p.shadowLayer)

	target.DrawImage(p.background, nil)

	target.DrawRectShader(
		target.Bounds().Dx(),
		target.Bounds().Dy(),
		shadowShader,
		p.shadowShaderOpts,
	)

	target.DrawImage(p.shadowLayer, nil)

	AddPostProccessing(target, p.buffer)
}

func AddPostProccessing(target, buffer *ebiten.Image) {
	buffer.Clear()
	// target.Clear()

	buffer.DrawRectShader(
		target.Bounds().Dx(),
		target.Bounds().Dy(),
		pixelShader,
		&ebiten.DrawRectShaderOptions{
			Images: [4]*ebiten.Image{
				target,
			},
			Uniforms: map[string]any{
				"Resolution": []float32{float32(target.Bounds().Dx()), float32(target.Bounds().Dy())},
			},
		},
	)
	target.DrawImage(buffer, nil)

	buffer.DrawRectShader(
		target.Bounds().Dx(),
		target.Bounds().Dy(),
		vingetteShader,
		&ebiten.DrawRectShaderOptions{
			Images: [4]*ebiten.Image{
				target,
			},
			Uniforms: map[string]any{
				"Resolution": []float32{float32(target.Bounds().Dx()), float32(target.Bounds().Dy())},
			},
		},
	)
	target.DrawImage(buffer, nil)

	buffer.DrawRectShader(
		target.Bounds().Dx(),
		target.Bounds().Dy(),
		wavesShader,
		&ebiten.DrawRectShaderOptions{
			Images: [4]*ebiten.Image{
				target,
			},
			Uniforms: map[string]any{
				"Time":       float64(ebiten.Tick()) / float64(ebiten.TPS()),
				"Resolution": []float32{float32(target.Bounds().Dx()), float32(target.Bounds().Dy())},
			},
		},
	)

	target.DrawImage(buffer, nil)
	buffer.DrawRectShader(
		target.Bounds().Dx(),
		target.Bounds().Dy(),
		crtShader,
		&ebiten.DrawRectShaderOptions{
			Images: [4]*ebiten.Image{
				target,
			},
			Uniforms: map[string]any{
				"Time":       float64(ebiten.Tick()) / float64(ebiten.TPS()),
				"Resolution": []float32{float32(target.Bounds().Dx()), float32(target.Bounds().Dy())},
			},
		},
	)
	target.DrawImage(buffer, nil)
}

func PixelBufferFromNoise(noise []float64, width, height int) []byte {
	buff := make([]byte, 4*width*height)
	idx := 0

	for y := range height {
		for x := range width {
			noiseValue := noise[y*width+x]
			pixelValue := byte(noiseValue * 255)
			buff[idx+0] = pixelValue
			buff[idx+1] = pixelValue
			buff[idx+2] = pixelValue
			buff[idx+3] = 255
			idx += 4
		}
	}

	return buff
}
