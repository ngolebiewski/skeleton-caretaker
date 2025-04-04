// controls/controls.go
package controls

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/ngolebiewski/gogamejammin/types"
)

type Controls struct {
	LastTapTime int
}

func (c *Controls) HandleInput(g *types.Game) bool {
	moving := false
	// Keyboard Controls
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.X += g.Speed
		g.IsMovingLeft = false
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.X -= g.Speed
		g.IsMovingLeft = true
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.Y -= g.Speed
		moving = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.Y += g.Speed
		moving = true
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.IsAttacking = true
		g.AttackCounter = 0
	}

	// Touch Controls
	if g.EnableTouch {
		touches := ebiten.TouchIDs()
		width, height := ebiten.ScreenSizeInFullscreen()
		if len(touches) > 0 {
			x, y := ebiten.TouchPosition(touches[0])
			if x < width/3 {
				g.X -= g.Speed
				g.IsMovingLeft = true
				moving = true
			} else if x > width*2/3 {
				g.X += g.Speed
				g.IsMovingLeft = false
				moving = true
			}
			if y < height/3 {
				g.Y -= g.Speed
				moving = true
			} else if y > height*2/3 {
				g.Y += g.Speed
				moving = true
			}
		}
		// Double-tap attack implementation.
		if len(touches) == 1 {
			if inpututil.TouchPressDuration(touches[0]) == 1 {
				// Check if a tap happened very recently.
				if c.LastTapTime > 0 && (g.Count-c.LastTapTime) < 15 { // 15 frames for double tap. Adjust as needed.
					g.IsAttacking = true
					g.AttackCounter = 0
					c.LastTapTime = 0 //reset tap time.
				} else {
					c.LastTapTime = g.Count //set the time of the first tap.
				}
			}
		}
	}
	return moving
}
