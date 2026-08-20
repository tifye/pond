package scenes

import (
	_ "embed"
	"math"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tifye/pond/internal/entity"
	"github.com/tifye/pond/internal/entity/lilypad"
	"github.com/tifye/pond/pkg/mathutil"
)

//go:embed gradient.kage
var gradientShaderSrc []byte
var gradientShader *ebiten.Shader

func init() {
	if gs, err := ebiten.NewShader(gradientShaderSrc); err != nil {
		panic(err)
	} else {
		gradientShader = gs
	}
}

type loadResult struct {
	lilypads   *lilypad.Lilypads
	background *ebiten.Image
	flowers    *lilypad.Flowers
}

type LoadingScene struct {
	ScreenWidth  float64
	ScreenHeight float64

	WaitFor time.Duration
	waitSec int64

	yin         *entity.Fish
	yang        *entity.Fish
	followPoint mathutil.Point

	readych chan loadResult
}

func (p *LoadingScene) Initialize() {
	p.WaitFor = time.Microsecond * 20

	p.yin = entity.NewFish(entity.FishConfig{
		SpawnPoint: mathutil.NewPoint(p.ScreenWidth*0.5, p.ScreenHeight*0.5),
		BoneChainUpdater: entity.BoneChainFABRIKUpdater{
			MinAngle: math.Pi * 0.85,
		},
	})
	p.yang = entity.NewFish(entity.FishConfig{
		SpawnPoint: mathutil.NewPoint(p.ScreenWidth*0.5, p.ScreenHeight*0.5),
		BoneChainUpdater: entity.BoneChainFABRIKUpdater{
			MinAngle: math.Pi * 0.85,
		},
	})

	p.readych = make(chan loadResult, 1)
	go p.load()
}

func (p *LoadingScene) load() {
	r := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

	perlin := mathutil.NewPerlin(r)
	renderScale := 2.0
	noiseWidth := int(float64(p.ScreenWidth) / renderScale)
	noiseHeight := int(float64(p.ScreenHeight) / renderScale)

	noiesMap := perlin.FBMMap2D(noiseWidth, noiseHeight, 0.006, 4, 0.35, 3)
	lilypads := lilypad.NewUsingNoiseThreshold(noiesMap, noiseWidth, noiseHeight, 2, 0.5, lilypad.DefaultFromNoiseConfig)

	flowers := lilypad.NewFlowers(int(p.ScreenWidth), int(p.ScreenHeight), 0.3, lilypads.Positions)

	pixels := PixelBufferFromNoise(noiesMap, noiseWidth, noiseHeight)
	temp := ebiten.NewImage(noiseWidth, noiseHeight)
	temp.WritePixels(pixels)

	noise := ebiten.NewImage(int(p.ScreenWidth), int(p.ScreenHeight))

	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Scale(renderScale, renderScale)
	noise.DrawImage(temp, opts)

	background := ebiten.NewImage(int(p.ScreenWidth), int(p.ScreenHeight))
	background.DrawRectShader(
		int(p.ScreenWidth),
		int(p.ScreenHeight),
		gradientShader,
		&ebiten.DrawRectShaderOptions{
			Images: [4]*ebiten.Image{
				noise,
			},
		},
	)

	time.Sleep(p.WaitFor)

	p.readych <- loadResult{
		lilypads:   lilypads,
		flowers:    flowers,
		background: background,
	}
}

func (p *LoadingScene) Update() Scene {
	select {
	case result := <-p.readych:
		pondScene := &PondScene{
			ScreenWidth:    int(p.ScreenWidth),
			ScreenHeight:   int(p.ScreenHeight),
			lilypads:       result.lilypads,
			flowers:        result.flowers,
			background:     result.background,
			backgroundBuff: ebiten.NewImageFromImage(result.background),
		}
		pondScene.Initialize([]*entity.Fish{p.yin, p.yang})

		return &TransitionScene{
			Before:       p,
			After:        pondScene,
			ScreenWidth:  int(p.ScreenWidth),
			ScreenHeight: int(p.ScreenHeight),
		}

	default:
	}

	amplitude := 175.0
	phase := float64(ebiten.Tick()) / float64(ebiten.TPS()) * 2.0

	// Yin
	{
		x := math.Cos(phase)*amplitude + p.ScreenWidth*0.5
		y := math.Sin(phase)*amplitude + p.ScreenHeight*0.5

		p.followPoint.X = x
		p.followPoint.Y = y

		p.yin.Position = p.followPoint
		p.yin.Update()
	}

	// Yang
	{
		x := math.Cos(phase+math.Pi)*amplitude + p.ScreenWidth*0.5
		y := math.Sin(phase+math.Pi)*amplitude + p.ScreenHeight*0.5

		p.followPoint.X = x
		p.followPoint.Y = y

		p.yang.Position = p.followPoint
		p.yang.Update()
	}

	return nil
}

func (p *LoadingScene) Draw(target *ebiten.Image) {
	p.yin.Draw(target)
	p.yang.Draw(target)
}
