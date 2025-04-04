// types.go
package types

type Game struct {
	Count         int
	X, Y          float64
	IsMovingLeft  bool
	IsAttacking   bool
	AttackCounter int
	IdleCounter   int
	EnableTouch   bool
	Fullscreen    bool
	Speed         float64
}
