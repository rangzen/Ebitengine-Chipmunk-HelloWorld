package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/jakecoffman/cp"
	"golang.org/x/image/colornames"

	"log"
)

// See the original at https://chipmunk-physics.net/release/ChipmunkLatest-Docs/#Intro-HelloChipmunk
// Values are changed due to screen size.

const (
	title              = "Hello Chipmunk (World)"
	simulateMaxSeconds = 6
	screenWidth        = 800
	screenHeight       = 600
)

var (
	ball = ebiten.NewImage(5, 5)
)

func init() {
	ball.Fill(color.White)
}

func main() {
	log.Println(title)
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(title)
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}

type Game struct {
	space    *cp.Space
	ballBody *cp.Body
	time     float64
}

func NewGame() *Game {
	// Create an empty space.
	gravity := cp.Vector{Y: 100}
	space := cp.NewSpace()
	space.SetGravity(gravity)

	// Add a static line segment shape for the ground.
	// We'll make it slightly tilted so the ball will roll off.
	// We attach it to a static body to tell Chipmunk it shouldn't be movable.
	ground := cp.NewSegment(
		space.StaticBody,
		cp.Vector{},
		cp.Vector{X: screenWidth, Y: screenHeight},
		0,
	)
	ground.SetFriction(1)
	space.AddShape(ground)

	// Now let's make a ball that falls onto the line and rolls off.
	// First we need to make a cpBody to hold the physical properties of the object.
	// These include the mass, position, velocity, angle, etc. of the object.
	// Then we attach collision shapes to the cpBody to give it a size and shape.

	var radius float64 = 5
	var mass float64 = 1

	// The moment of inertia is like mass for rotation
	// Use the cp.MomentFor*() functions to help you approximate it.
	moment := cp.MomentForCircle(mass, 0, radius, cp.Vector{})

	// The Space.Add*() functions return the thing that you are adding.
	// It's convenient to create and add an object in one line.
	ballBody := space.AddBody(cp.NewBody(mass, moment))
	ballBody.SetPosition(cp.Vector{X: screenWidth / 2, Y: screenHeight / 4})

	// Now we create the collision shape for the ball.
	// You can create multiple collision shapes that point to the same body.
	// They will all be attached to the body and move around to follow it.
	ballShape := space.AddShape(cp.NewCircle(ballBody, radius, cp.Vector{}))
	ballShape.SetFriction(0.7)

	return &Game{
		space:    space,
		ballBody: ballBody,
	}
}

func (g *Game) Update() error {
	// Now that it's all set up, we simulate all the objects in the space by
	// stepping forward through time in small increments called steps.
	// It is *highly* recommended to use a fixed size time step.
	timeStep := 1.0 / float64(ebiten.MaxTPS())
	g.time += timeStep
	if g.time < simulateMaxSeconds {
		g.space.Step(timeStep)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(colornames.Black)

	// Ground
	ebitenutil.DrawLine(screen, 0, 0, screenWidth, screenHeight, color.White)

	// Ball
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(200.0/255.0, 200.0/255.0, 200.0/255.0, 1)
	op.GeoM.Translate(g.ballBody.Position().X, g.ballBody.Position().Y)
	screen.DrawImage(ball, op)

	if g.time < simulateMaxSeconds {
		pos := g.ballBody.Position()
		vel := g.ballBody.Velocity()
		ebitenutil.DebugPrint(
			screen,
			fmt.Sprintf(
				"Time is %5.2f. ballBody is at (%5.2f, %5.2f). It's velocity is (%5.2f, %5.2f)",
				g.time, pos.X, pos.Y, vel.X, vel.Y,
			))
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}
