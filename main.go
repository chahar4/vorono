package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	d "github.com/uaraven/delaunay/delaunay"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	Dots = 100
)

var inputDots int

func main() {
	inputDotsS := os.Args[1]
	inputDots, _ = strconv.Atoi(inputDotsS)
	fmt.Println(inputDots)

	screenWidth := int32(800)
	screenHeight := int32(450)
	dots := initDots(inputDots)
	p := dotsToPoints(dots)
	delaunay := d.InitDelaunay(p)
	d.Delaunay

	for _, x := range dots {
		fmt.Print(x, " ")
	}
	fmt.Println("")

	rl.InitWindow(screenWidth, screenHeight, "raylib [shapes] example - cubic-bezier lines")

	// start := rl.NewVector2(0, 0)
	// end := rl.NewVector2(float32(screenWidth), float32(screenHeight))
	rl.SetTargetFPS(60)

	// Create offscreen render target
	target := rl.LoadRenderTexture(int32(screenWidth), int32(screenHeight))
	defer rl.UnloadRenderTexture(target)

	dirty := true // Flag to redraw texture when dots change

	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyEnter) {
			dots = initDots(inputDots)
			dirty = true
		}

		// Only redraw circles when dots change
		if dirty {
			rl.BeginTextureMode(target)
			rl.ClearBackground(rl.RayWhite)
			for i := 0; i < len(dots); i += 2 {
				rl.DrawCircle(int32(dots[i]), int32(dots[i+1]), 3, rl.Black)
			}
			rl.EndTextureMode()
			dirty = false
		}

		// Per-frame: just draw one textured quad
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// NOTE: RenderTexture's texture has flipped Y axis in raylib
		rl.DrawTextureRec(
			target.Texture,
			rl.NewRectangle(0, 0, float32(screenWidth), -float32(screenHeight)),
			rl.NewVector2(0, 0),
			rl.White,
		)

		rl.DrawText(fmt.Sprintf("Dots: %d | FPS: %.0f", len(dots)/2, rl.GetFPS()), 10, 10, 20, rl.DarkGray)
		rl.EndDrawing()
	}
}

func initDots(len int) []int {
	dots := make([]int, 0)
	for range len {
		x := rand.Intn(800)
		y := rand.Intn(450)

		dots = append(dots, x)
		dots = append(dots, y)
	}
	return dots
}

func dotsToPoints(dots []int) []d.Point {
	points := make([]d.Point, inputDots)
	for i := 0; i < len(dots); i += 2 {
		points = append(points, d.Point{
			X: float64(dots[i]),
			Y: float64(dots[i+1])})
	}
	return points
}
