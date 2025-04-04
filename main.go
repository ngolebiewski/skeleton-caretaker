package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"github.com/ngolebiewski/gogamejammin/controls"
	"github.com/ngolebiewski/gogamejammin/types"
)

//go:embed assets/skeleton_enemy.png
var skeletonEnemyPNG []byte
var skeletonImage *ebiten.Image

const (
	frameWidth       = 64
	frameHeight      = 64
	walkFrameCount   = 12
	attackFrameCount = 13
	idleFrameCount   = 4
	speed            = 2
	scaleFactor      = 2
)

type MyGame struct {
	types.Game // Embed types.Game
}

func (g *MyGame) Update() error {
	g.Count++

	width, _ := ebiten.ScreenSizeInFullscreen()
	if width < 600 {
		g.EnableTouch = true
	} else {
		g.EnableTouch = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF11) || inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.Fullscreen = !g.Fullscreen
		ebiten.SetFullscreen(g.Fullscreen)
	}

	controlHandler := controls.Controls{}
	moving := controlHandler.HandleInput(&g.Game) // Pass the embedded Game

	if g.IsAttacking {
		g.AttackCounter++
		if g.AttackCounter >= attackFrameCount*5 {
			g.IsAttacking = false
		}
	}

	if !moving {
		g.IdleCounter++
	} else {
		g.IdleCounter = 0
	}

	return nil
}

func (g *MyGame) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	var sx, sy int
	var frameIndex int

	if g.IsAttacking {
		frameIndex = (g.AttackCounter / 5) % attackFrameCount
		sy = 0
	} else if g.IdleCounter > 0 {
		frameIndex = (g.IdleCounter / 10) % idleFrameCount
		sy = 3 * frameHeight
	} else {
		frameIndex = (g.Count / 5) % walkFrameCount
		sy = 2 * frameHeight
	}

	sx = frameIndex * frameWidth

	if g.IsMovingLeft {
		op.GeoM.Scale(-scaleFactor, scaleFactor)
		op.GeoM.Translate(float64(frameWidth*scaleFactor), 0)
	} else {
		op.GeoM.Scale(scaleFactor, scaleFactor)
	}

	op.GeoM.Translate(g.X, g.Y)

	frame := skeletonImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image)
	screen.DrawImage(frame, op)
}

func (g *MyGame) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	img, _, err := image.Decode(bytes.NewReader(skeletonEnemyPNG))
	if err != nil {
		log.Fatal(err)
	}
	skeletonImage = ebiten.NewImageFromImage(img)

	game := &MyGame{types.Game{X: 100, Y: 100, Fullscreen: false, Speed: speed}}

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
