package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"

	d "github.com/uaraven/delaunay/delaunay"

	"github.com/gen2brain/raylib-go/raylib"
	"github.com/pzsz/voronoi"
)

const (
	Dots = 100
)

var inputDots int
var dots []float64
var actualPoints []voronoi.Vertex
var triangles []d.Triangle

func main() {
	screenWidth := int32(800)
	screenHeight := int32(450)
	rl.InitWindow(screenWidth, screenHeight, "voronoi")
	setup()
	rl.SetTargetFPS(60)
	draw(screenWidth, screenHeight)

}

func setup() {
	inputDotsS := os.Args[1]
	inputDots, _ = strconv.Atoi(inputDotsS)

	dots = initDots(inputDots)

	p := dotsToPoints(dots)
	actualPoints = dotsToPoints(dots)
	delaunay := d.InitDelaunay(p)
	triangles = delaunay.Triangulate()
}

func draw(screenWidth, screenHeight int32) {
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
			for _, point := range actualPoints {
				rl.DrawCircle(int32(point.X), int32(point.Y), 3, rl.Black)
				rl.DrawText(point.String(), int32(point.X)+5, int32(point.Y), 12, rl.DarkBlue)
			}
			for _, triangle := range triangles {
				rl.DrawLine(int32(triangle.P1.X), int32(triangle.P1.Y), int32(triangle.P2.X), int32(triangle.P2.Y), rl.Green)
				rl.DrawLine(int32(triangle.P2.X), int32(triangle.P2.Y), int32(triangle.P3.X), int32(triangle.P3.Y), rl.Green)
				rl.DrawLine(int32(triangle.P3.X), int32(triangle.P3.Y), int32(triangle.P1.X), int32(triangle.P1.Y), rl.Green)
				rl.DrawCircleLines(int32(triangle.Circumcenter.X), int32(triangle.Circumcenter.Y), float32(triangle.Circumradius), rl.Gray)
				rl.DrawCircle(int32(triangle.Circumcenter.X), int32(triangle.Circumcenter.Y), 2, rl.Red)
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

		rl.DrawText(fmt.Sprintf("Dots: %d | FPS: %d", len(dots)/2, rl.GetFPS()), 10, 10, 20, rl.DarkGray)
		rl.EndDrawing()
	}
}

func initDots(len int) []float64 {
	dots := make([]float64, 0)
	for range len {
		x := rand.Intn(800)
		y := rand.Intn(450)

		dots = append(dots, float64(x))
		dots = append(dots, float64(y))
	}
	return dots
}

func dotsToPoints(dots []float64) []voronoi.Vertex {
	sites := make([]voronoi.Vertex, 0)
	for i := 0; i < len(dots); i += 2 {
		sites = append(sites, voronoi.Vertex{
			X: dots[i],
			Y: dots[i+1],
		})
	}

	return sites
}
