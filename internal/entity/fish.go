package entity

import (
	"bytes"
	_ "embed"
	"image/color"
	"image/jpeg"
	"math"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/tifye/pond/colors"
	"github.com/tifye/pond/pkg/mathutil"
)

//go:embed shaders/texture_shader.go
var textureShader []byte

//go:embed assets/koi-skin.jpg
var koiSkin []byte

type Fish struct {
	Position mathutil.Point

	spine BoneChain

	texture       *ebiten.Image
	textureShader *ebiten.Shader

	mesh mesh

	curvature float64

	pectoralLeftFin  *fin
	pectoralRightFin *fin

	ventralLeftFin  *fin
	ventralRightFin *fin

	barbelLeft  mathutil.Point
	barbelRight mathutil.Point

	tail *CaudralFin
}

type mesh struct {
	vertices []ebiten.Vertex
	indicies []uint16
}

type FishConfig struct {
	SpawnPoint       mathutil.Point
	BoneChainUpdater BoneChainUpdater
}

func NewFish(cfg FishConfig) *Fish {
	bones := make([]mathutil.Point, 16)
	boneLengths := make([]float64, len(bones)-1)
	bones[0] = cfg.SpawnPoint

	for i := 1; i < len(bones); i++ {
		bones[i] = cfg.SpawnPoint
		boneLengths[i-1] = 16 * (float64(len(bones)-i) / float64(len(bones)))
	}

	ts, err := ebiten.NewShader(textureShader)
	if err != nil {
		panic(err)
	}

	skinImage, err := jpeg.Decode(bytes.NewReader(koiSkin))
	if err != nil {
		panic(err)
	}
	texture := ebiten.NewImageFromImage(skinImage)

	spine := BoneChain{
		Joints:  bones,
		Lengths: boneLengths,
		Updater: BoneChainFABRIKUpdater{
			MinAngle: math.Pi * 0.85,
		},
	}
	if cfg.BoneChainUpdater != nil {
		spine.Updater = cfg.BoneChainUpdater
	}

	f := &Fish{
		spine: spine,

		texture:       texture,
		textureShader: ts,

		pectoralLeftFin:  newFin(35, 25, true),
		pectoralRightFin: newFin(35, 25, false),

		ventralLeftFin:  newFin(25, 20, true),
		ventralRightFin: newFin(25, 20, false),

		tail: NewTail(),
	}

	f.BuildMesh()

	return f
}

func (f *Fish) Update() {
	f.spine.Position = f.Position
	f.spine.Update()

	{
		v := f.spine.Joints[0].Subtract(f.spine.Joints[1])

		f.barbelLeft = v.Rotate(-math.Pi * 0.3).
			MultiplyScalar(f.spine.Lengths[0] * 1.1).
			Add(f.spine.Joints[0])

		f.barbelRight = v.Rotate(math.Pi * 0.3).
			MultiplyScalar(f.spine.Lengths[0] * 1.1).
			Add(f.spine.Joints[0])
	}

	f.curvature = f.spine.Curvature()

	f.tail.Position = f.spine.Joints[len(f.spine.Joints)-1]
	f.tail.Update()

	f.updateMesh()
}

func (f *Fish) Draw(target *ebiten.Image) {

	f.drawPectoralFins(target)
	f.drawVentralFins(target)
	f.drawBarbels(target)

	target.DrawTrianglesShader(
		f.mesh.vertices,
		f.mesh.indicies,
		f.textureShader,
		&ebiten.DrawTrianglesShaderOptions{
			Images: [4]*ebiten.Image{
				f.texture,
			},
		},
	)

	f.tail.Draw(target)

}

func (f *Fish) drawBarbels(target *ebiten.Image) {
	length := 8.0
	thicknes := 2
	col := colors.White

	stokeOpts := &vector.StrokeOptions{
		LineCap:  vector.LineCapRound,
		LineJoin: vector.LineJoinBevel,
		Width:    float32(thicknes),
	}

	drawOpts := &vector.DrawPathOptions{}
	drawOpts.ColorScale.ScaleWithColor(col)

	v := f.barbelLeft.Subtract(f.spine.Joints[0]).
		Rotate(-math.Pi * .3).
		Normalize().
		MultiplyScalar(length).
		Add(f.barbelLeft)

	p := vector.Path{}
	p.MoveTo(float32(f.barbelLeft.X), float32(f.barbelLeft.Y))
	p.LineTo(float32(v.X), float32(v.Y))
	vector.StrokePath(target, &p, stokeOpts, drawOpts)

	v = f.barbelRight.Subtract(f.spine.Joints[0]).
		Rotate(math.Pi * .3).
		Normalize().
		MultiplyScalar(length).
		Add(f.barbelRight)

	p = vector.Path{}
	p.MoveTo(float32(f.barbelRight.X), float32(f.barbelRight.Y))
	p.LineTo(float32(v.X), float32(v.Y))
	vector.StrokePath(target, &p, stokeOpts, drawOpts)
}

func (f *Fish) drawPectoralFins(target *ebiten.Image) {
	anchorPoint := f.spine.Joints[1]
	angle := f.spine.Joints[1].AngleBetween(f.spine.Joints[2])

	f.pectoralLeftFin.Rotation = angle + math.Pi*0.35 - f.curvature*0.25
	f.pectoralLeftFin.Position = anchorPoint
	f.pectoralLeftFin.draw(target)
	f.pectoralLeftFin.draw(target)

	f.pectoralRightFin.Rotation = angle + math.Pi*0.65 - f.curvature*0.25
	f.pectoralRightFin.Position = anchorPoint
	f.pectoralRightFin.draw(target)
	f.pectoralRightFin.draw(target)
}

func (f *Fish) drawVentralFins(target *ebiten.Image) {
	anchorPoint := f.spine.Joints[4]
	angle := f.spine.Joints[4].AngleBetween(f.spine.Joints[6])

	f.ventralLeftFin.Rotation = angle + math.Pi*0.25 - f.curvature*0.15
	f.ventralLeftFin.Position = anchorPoint
	f.ventralLeftFin.draw(target)
	f.ventralLeftFin.draw(target)

	f.ventralRightFin.Rotation = angle + math.Pi*0.75 - f.curvature*0.15
	f.ventralRightFin.Position = anchorPoint
	f.ventralRightFin.draw(target)
	f.ventralRightFin.draw(target)
}

func debugCircle(t *ebiten.Image, p mathutil.Point, c color.Color) {
	vector.FillCircle(
		t,
		float32(p.X),
		float32(p.Y),
		2,
		c,
		false,
	)
}

func (f *Fish) updateMesh() {
	leftBodyStartIdx := 0
	rightBodyStartIdx := len(f.spine.Joints)
	headStartIdx := rightBodyStartIdx + len(f.spine.Joints)

	for i := 0; i < len(f.spine.Joints)-1; i++ {
		curBone := f.spine.Joints[i+1]
		nextBone := f.spine.Joints[i]
		segment := curBone.Subtract(nextBone)

		l := segment.RotateCounterClockwise().Add(nextBone)
		r := segment.RotateClockwise().Add(nextBone)

		f.mesh.vertices[leftBodyStartIdx+i].DstX = float32(l.X)
		f.mesh.vertices[leftBodyStartIdx+i].DstY = float32(l.Y)
		f.mesh.vertices[rightBodyStartIdx+i].DstX = float32(r.X)
		f.mesh.vertices[rightBodyStartIdx+i].DstY = float32(r.Y)
	}

	lastBone := f.spine.Joints[len(f.spine.Joints)-1]
	last2Bone := f.spine.Joints[len(f.spine.Joints)-2]
	segment := lastBone.Subtract(last2Bone)

	l := segment.RotateCounterClockwise().Add(last2Bone)
	r := segment.RotateClockwise().Add(last2Bone)

	f.mesh.vertices[leftBodyStartIdx+len(f.spine.Joints)-1].DstX = float32(l.X)
	f.mesh.vertices[leftBodyStartIdx+len(f.spine.Joints)-1].DstY = float32(l.Y)
	f.mesh.vertices[rightBodyStartIdx+len(f.spine.Joints)-1].DstX = float32(r.X)
	f.mesh.vertices[rightBodyStartIdx+len(f.spine.Joints)-1].DstY = float32(r.Y)

	// Head index
	v := f.spine.Joints[0].Subtract(f.spine.Joints[1]).Normalize()

	vl := v.Rotate(-math.Pi * 0.2).MultiplyScalar(f.spine.Lengths[0]).Add(f.spine.Joints[0])
	f.mesh.vertices[headStartIdx].DstX = float32(vl.X)
	f.mesh.vertices[headStartIdx].DstY = float32(vl.Y)

	vr := v.Rotate(math.Pi * 0.2).MultiplyScalar(f.spine.Lengths[0]).Add(f.spine.Joints[0])
	f.mesh.vertices[headStartIdx+1].DstX = float32(vr.X)
	f.mesh.vertices[headStartIdx+1].DstY = float32(vr.Y)
}

func (f *Fish) BuildMesh() {
	nHeadVertices := 2

	vertices := make([]ebiten.Vertex, len(f.spine.Joints)*2+nHeadVertices)
	indicies := make([]uint16, (len(vertices)-2)*3)

	leftBodyStartIdx := 0
	rightBodyStartIdx := len(f.spine.Joints)
	headStartIdx := rightBodyStartIdx + len(f.spine.Joints)

	leftBodyBoneIdx := func(boneIdx int) uint16 {
		return uint16(leftBodyStartIdx + boneIdx)
	}

	rightBodyBoneIdx := func(boneIdx int) uint16 {
		return uint16(rightBodyStartIdx + boneIdx)
	}

	i := 0
	for boneIdx := 0; boneIdx < len(f.spine.Joints)-1; boneIdx++ {
		indicies[i+0] = leftBodyBoneIdx(boneIdx)
		indicies[i+1] = rightBodyBoneIdx(boneIdx)
		indicies[i+2] = rightBodyBoneIdx(boneIdx + 1)
		i += 3

		indicies[i+0] = leftBodyBoneIdx(boneIdx)
		indicies[i+1] = rightBodyBoneIdx(boneIdx + 1)
		indicies[i+2] = leftBodyBoneIdx(boneIdx + 1)
		i += 3
	}

	// Head indices
	indicies[i+0] = leftBodyBoneIdx(0)
	indicies[i+1] = uint16(headStartIdx)
	indicies[i+2] = rightBodyBoneIdx(0)
	i += 3

	indicies[i+0] = uint16(headStartIdx)
	indicies[i+1] = uint16(headStartIdx + 1)
	indicies[i+2] = rightBodyBoneIdx(0)
	i += 3

	// Set UVs
	{
		for i := range f.spine.Joints {
			vertices[leftBodyStartIdx+i].Custom0 = 0
			vertices[leftBodyStartIdx+i].Custom1 = float32(i) / (float32(len(f.spine.Joints) - 1))
		}

		for i := range f.spine.Joints {
			vertices[rightBodyStartIdx+i].Custom0 = 1
			vertices[rightBodyStartIdx+i].Custom1 = float32(i) / (float32(len(f.spine.Joints) - 1))
		}

		vertices[headStartIdx].Custom0 = 0
		vertices[headStartIdx].Custom1 = 0
		vertices[headStartIdx+1].Custom0 = 1
		vertices[headStartIdx+1].Custom1 = 1
	}

	for i := range vertices {
		vertices[i].ColorA = 1
	}

	f.mesh = mesh{
		vertices: vertices,
		indicies: indicies,
	}
}

//go:embed shaders/fin_shader.go
var finShaderSrc []byte
var finShader *ebiten.Shader

func init() {
	var err error
	finShader, err = ebiten.NewShader(finShaderSrc)
	if err != nil {
		panic(err)
	}
}

type fin struct {
	Position mathutil.Point

	Rotation float64

	size   mathutil.Size
	scaleX float64
	img    *ebiten.Image
}

func newFin(width, height int, flipped bool) *fin {
	scaleX := 1.0
	if !flipped {
		scaleX = -1.0
	}
	return &fin{
		size:   mathutil.NewPoint(float64(width), float64(height)),
		img:    ebiten.NewImage(width, height),
		scaleX: scaleX,
	}
}

func (f *fin) draw(target *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	opts.GeoM.Scale(f.scaleX, 1.0)
	opts.GeoM.Translate(0.0, f.size.Y*-0.5)
	opts.GeoM.Rotate(f.Rotation)
	opts.GeoM.Translate(f.Position.X, f.Position.Y)

	f.img.DrawRectShader(
		f.img.Bounds().Dx(),
		f.img.Bounds().Dy(),
		finShader,
		nil,
	)

	target.DrawImage(f.img, opts)
}

type CaudralFin struct {
	Position mathutil.Point

	upperLobe          BoneChain
	lowerLobe          BoneChain
	lowerLobeReflected []mathutil.Point

	lobeConnection mathutil.Point

	mesh    mesh
	texture *ebiten.Image
}

func NewTail() *CaudralFin {
	topN := 10
	topL := 40.0
	topBones := make([]mathutil.Point, topN)
	topBoneLengths := make([]float64, len(topBones)-1)

	for i := 1; i < len(topBones); i++ {
		topBoneLengths[i-1] = topL / float64(topN)
	}

	botN := 10
	botL := topL * 1.4
	botBones := make([]mathutil.Point, botN)
	botBoneLengths := make([]float64, len(botBones)-1)

	for i := 1; i < len(botBones); i++ {
		botBoneLengths[i-1] = botL / float64(botN)
	}

	texture := ebiten.NewImage(1, 1)
	texture.Set(0, 0, colors.Neutral300)

	vertices := make([]ebiten.Vertex, len(botBones)+len(topBones)+1) // +1 for the lobe connection
	indicies := make([]uint16, (len(vertices)-1)*3)

	lobeConnectionIdx := 0
	topVertsOffset := 1
	botVertsOffset := topVertsOffset + len(topBones)

	j := 0
	for boneIdx := 0; boneIdx < len(topBones)-1; boneIdx++ {
		indicies[j+0] = uint16(topVertsOffset + boneIdx)
		indicies[j+1] = uint16(topVertsOffset + boneIdx + 1)
		indicies[j+2] = uint16(lobeConnectionIdx)
		j += 3
	}

	for boneIdx := 0; boneIdx < len(botBones)-1; boneIdx++ {
		indicies[j+0] = uint16(botVertsOffset + boneIdx + 1)
		indicies[j+1] = uint16(botVertsOffset + boneIdx)
		indicies[j+2] = uint16(lobeConnectionIdx)
		j += 3
	}

	for i := range vertices {
		vertices[i].ColorR = 1
		vertices[i].ColorG = 1
		vertices[i].ColorB = 1
		vertices[i].ColorA = 1
	}

	return &CaudralFin{
		upperLobe: BoneChain{
			Joints:  topBones,
			Lengths: topBoneLengths,
			Updater: BoneChainFABRIKUpdater{
				MinAngle: math.Pi,
			},
		},

		lowerLobe: BoneChain{
			Joints:  botBones,
			Lengths: botBoneLengths,
			Updater: BoneChainFABRIKUpdater{
				MinAngle: math.Pi * 0.9,
			},
		},

		lowerLobeReflected: slices.Clone(botBones),

		mesh: mesh{
			vertices: vertices,
			indicies: indicies,
		},

		texture: texture,
	}
}

func (t *CaudralFin) Update() {
	t.upperLobe.Position = t.Position
	t.upperLobe.Update()

	t.lowerLobe.Position = t.Position
	t.lowerLobe.Update()

	anchor := t.upperLobe.Joints[0]
	nv := t.upperLobe.Joints[1].
		Subtract(anchor).
		Normalize().
		Rotate(math.Pi * 0.5)

	t.lowerLobeReflected[0] = t.lowerLobe.Joints[0]
	t.lowerLobeReflected[1] = t.lowerLobe.Joints[1]
	for i, p := range t.upperLobe.Joints[2:] {
		local := p.Subtract(anchor)
		dot := nv.Dot(local)

		offset := nv.MultiplyScalar(2 * dot)
		reflected := local.Subtract(offset).Add(anchor)

		t.lowerLobeReflected[i+2] = reflected
	}

	a := t.upperLobe.Joints[len(t.upperLobe.Joints)-1]
	b := t.lowerLobeReflected[len(t.lowerLobeReflected)-1]
	t.lobeConnection = a.Subtract(b).
		MultiplyScalar(0.5).
		Add(b).
		Subtract(t.Position).
		MultiplyScalar(0.8).
		Add(t.Position)

	t.updateMesh()
}

func (t *CaudralFin) updateMesh() {
	lobeConnectionIdx := 0
	topVertsOffset := 1
	botVertsOffset := topVertsOffset + len(t.upperLobe.Joints)

	t.mesh.vertices[lobeConnectionIdx].DstX = float32(t.lobeConnection.X)
	t.mesh.vertices[lobeConnectionIdx].DstY = float32(t.lobeConnection.Y)

	for i, p := range t.upperLobe.Joints {
		t.mesh.vertices[topVertsOffset+i].DstX = float32(p.X)
		t.mesh.vertices[topVertsOffset+i].DstY = float32(p.Y)
	}

	for i, p := range t.lowerLobe.Joints {
		t.mesh.vertices[botVertsOffset+i].DstX = float32(p.X)
		t.mesh.vertices[botVertsOffset+i].DstY = float32(p.Y)
	}
}

func (t *CaudralFin) Draw(target *ebiten.Image) {
	var path vector.Path

	path.MoveTo(float32(t.Position.X), float32(t.Position.Y))
	for _, p := range t.upperLobe.Joints[:len(t.upperLobe.Joints)-2] {
		path.LineTo(float32(p.X), float32(p.Y))
	}

	path.MoveTo(float32(t.Position.X), float32(t.Position.Y))
	for _, p := range t.lowerLobe.Joints[:len(t.lowerLobe.Joints)-1] {
		path.LineTo(float32(p.X), float32(p.Y))
	}

	vector.StrokePath(
		target,
		&path,
		&vector.StrokeOptions{
			Width:    6,
			LineCap:  vector.LineCapRound,
			LineJoin: vector.LineJoinRound,
		},
		nil,
	)

	// debugCircle(target, t.lobeConnection, colors.White)

	target.DrawTriangles(
		t.mesh.vertices,
		t.mesh.indicies,
		t.texture,
		nil,
	)
}
