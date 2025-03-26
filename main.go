package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed assets/skeleton_enemy.png
var skeletonEnemyPNG []byte
var skeletonImage *ebiten.Image

const (
	frameWidth       = 64 // Each frame is 64x64
	frameHeight      = 64
	walkFrameCount   = 12 // 12 frames for walk animation in 3rd row
	attackFrameCount = 13 // 13 frames for attack animation in 1st row
	idleFrameCount   = 4  // 4 frames for idle animation in 4th row
	speed            = 2  // Movement speed of the character
	scaleFactor      = 2  // Scale factor to draw the character twice as large
)

// const (
// 	screenWidth  = 320
// 	screenHeight = 240
// )

type Game struct {
	count         int
	x, y          float64 // Character's position
	isMovingLeft  bool    // Flag to track if character is moving left
	isAttacking   bool    // Flag to track if spacebar is pressed for attack animation
	attackCounter int     // Counter to track the frames of the attack animation
	idleCounter   int     // Counter to track the frames of the idle animation
}

// Update is called every frame (60 FPS by default)
func (g *Game) Update() error {
	g.count++

	// Check if the character is not attacking to handle movement
	if !g.isAttacking {
		moving := false // Track if the character is moving

		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			g.x += speed
			g.isMovingLeft = false // Moving right, no flip
			moving = true
		}
		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			g.x -= speed
			g.isMovingLeft = true // Moving left, enable flip
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

		// Start attack animation if spacebar is pressed
		if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.isAttacking = true
			g.attackCounter = 0 // Reset attack animation counter
		}

		// If the character is not moving, play the idle animation
		if !moving {
			g.idleCounter++
		} else {
			g.idleCounter = 0 // Reset idle counter when moving
		}
	}

	// If attacking, increment the attack counter
	if g.isAttacking {
		g.attackCounter++
		// Once the attack animation has played all frames, stop attacking
		if g.attackCounter >= attackFrameCount*5 { // Assuming 5 frames per animation cycle
			g.isAttacking = false
		}
	}

	return nil
}

// Draw is called every frame to draw the game elements on the screen
func (g *Game) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	var sx, sy int
	var frameIndex int

	if g.isAttacking {
		// Attack animation (first row)
		frameIndex = (g.attackCounter / 5) % attackFrameCount
		sy = 0 // First row for attack animation
	} else if g.idleCounter > 0 {
		// Idle animation (fourth row)
		frameIndex = (g.idleCounter / 10) % idleFrameCount // Slower idle animation
		sy = 3 * frameHeight                               // Fourth row for idle animation
	} else {
		// Walk animation (third row)
		frameIndex = (g.count / 5) % walkFrameCount
		sy = 2 * frameHeight // Third row for walk animation
	}

	sx = frameIndex * frameWidth

	// Flip horizontally if moving left
	if g.isMovingLeft {
		op.GeoM.Scale(-scaleFactor, scaleFactor)              // Flip on X-axis and scale by 2
		op.GeoM.Translate(float64(frameWidth*scaleFactor), 0) // Translate to correct position after flip
	} else {
		op.GeoM.Scale(scaleFactor, scaleFactor) // Scale both X and Y by 2
	}

	// Set up translation for where to draw the sprite (after scaling)
	op.GeoM.Translate(-float64(frameWidth*scaleFactor)/2, -float64(frameHeight*scaleFactor)/2)
	op.GeoM.Translate(g.x+float64(screen.Bounds().Dx()/2), g.y+float64(screen.Bounds().Dy()/2)) // Translate to character's position

	// Draw the current frame from the sprite sheet
	screen.DrawImage(skeletonImage.SubImage(image.Rect(sx, sy, sx+frameWidth, sy+frameHeight)).(*ebiten.Image), op)
}

// Layout sets the logical screen size
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func main() {
	// Load skeletonImage.png from the assets directory
	// file, err := os.Open("assets/skeleton_enemy.png")
	// if err != nil {
	// 	log.Fatalf("Failed to open image: %v", err)
	// }
	// defer file.Close()

	// img, _, err := image.Decode(file)
	// if err != nil {
	// 	log.Fatalf("Failed to decode image: %v", err)
	// }
	// skeletonImage = ebiten.NewImageFromImage(img)

	// Load skeletonImage.png from embedded data
	img, _, err := image.Decode(bytes.NewReader(skeletonEnemyPNG))
	if err != nil {
		log.Fatalf("Failed to decode embedded image: %v", err)
	}
	// Convert standard image.Image to an ebiten.Image
	skeletonImage = ebiten.NewImageFromImage(img)

	// Set initial character position to the center of the screen
	game := &Game{
		x: 0,
		y: 0,
	}
	// ebiten.SetWindowSize(screenWidth, screenHeight)
	// Enable fullscreen

	ebiten.SetFullscreen(true)
	ebiten.SetWindowTitle("Skeleton Animation Demo Fullscreen")

	// Run the game
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
