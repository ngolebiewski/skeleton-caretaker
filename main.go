package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

type Game struct {
	count         int
	x, y          float64
	isMovingLeft  bool
	isAttacking   bool
	attackCounter int
	idleCounter   int
	enableTouch   bool
	fullscreen    bool
}

func (g *Game) Update() error {
	g.count++

	width, _ := ebiten.ScreenSizeInFullscreen()
	if width < 600 {
		g.enableTouch = true
	} else {
		g.enableTouch = false
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF11) || inpututil.IsKeyJustPressed(ebiten.KeyF) {
		g.fullscreen = !g.fullscreen
		ebiten.SetFullscreen(g.fullscreen)
	}

	moving := false
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.x += speed
		g.isMovingLeft = false
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.x -= speed
		g.isMovingLeft = true
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.y -= speed
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.y += speed
		moving = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.isAttacking = true
		g.attackCounter = 0
	}

	if g.enableTouch {
		touches := ebiten.TouchIDs()
		width, height := ebiten.ScreenSizeInFullscreen()
		if len(touches) > 0 {
			x, y := ebiten.TouchPosition(touches[0])
			if x < width/3 {
				g.x -= speed
				g.isMovingLeft = true
				moving = true
			} else if x > width*2/3 {
				g.x += speed
				g.isMovingLeft = false
				moving = true
			}
			if y < height/3 {
				g.y -= speed
				moving = true
			} else if y > height*2/3 {
				g.y += speed
				moving = true
			}
		}
		if len(touches) > 1 || (len(touches) == 1 && inpututil.TouchPressDuration(touches[0]) == 1) {
			g.isAttacking = true
			g.attackCounter = 0
		}
	}

	if g.isAttacking {
		g.attackCounter++
		if g.attackCounter >= attackFrameCount*5 {
			g.isAttacking = false
		}
	}

	if !moving {
		g.idleCounter++
	} else {
		g.idleCounter = 0
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	var sx, sy int
	var frameIndex int

	if g.isAttacking {
		frameIndex = (g.attackCounter / 5) % attackFrameCount
		sy = 0
	} else if g.idleCounter > 0 {
		frameIndex = (g.idleCounter / 10) % idleFrameCount
		sy = 3 * frameHeight
	} else {
		frameIndex = (g.count / 5) % walkFrameCount
		sy = 2 * frameHeight
	}

	sx = frameIndex * frameWidth

	if g.isMovingLeft {
		op.GeoM.Scale(-scaleFactor, scaleFactor)
		op.GeoM.Translate(float64(frameWidth*scaleFactor), 0)
	} else {
		op.GeoM.Scale(scaleFactor, scaleFactor)
	}

	op.GeoM.Translate(g.x, g.y)

	frame := skeletonImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image)
	screen.DrawImage(frame, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	img, _, err := image.Decode(bytes.NewReader(skeletonEnemyPNG))
	if err != nil {
		log.Fatal(err)
	}
	skeletonImage = ebiten.NewImageFromImage(img)

	game := &Game{x: 100, y: 100, fullscreen: false}
	// ebiten.SetFullscreen(true) // Enable fullscreen mode at the start

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
