package scenes

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tifye/pond/internal/entity"
	"github.com/tifye/pond/internal/entity/lilypad"
	"github.com/tifye/pond/pkg/mathutil"
)

type loadResult struct {
	lilypads *lilypad.Lilypads
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
	p.WaitFor = time.Second * 3

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
	lilypads := lilypad.NewUsingCirclePacking(lilypad.Config{
		MinRadius: 5,
	})

	time.Sleep(p.WaitFor)
	p.readych <- loadResult{
		lilypads: lilypads,
	}
}

func (p *LoadingScene) Update() Scene {
	select {
	case result := <-p.readych:
		pondScene := &PondScene{
			ScreenWidth:  int(p.ScreenWidth),
			ScreenHeight: int(p.ScreenHeight),
			lilypads:     result.lilypads,
		}
		pondScene.Initialize([]*entity.Fish{p.yin, p.yang})
		return pondScene
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
