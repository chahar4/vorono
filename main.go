package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/gen2brain/raylib-go/raylib"
	"github.com/pzsz/voronoi"
)

const (
	Dots = 100
)

var inputDots int
var dots []float64
var screenWidth int32
var screenHeight int32
var diagram *voronoi.Diagram
var centroidedSites []voronoi.Vertex

func main() {
	screenWidth = int32(800)
	screenHeight = int32(450)
	rl.InitWindow(screenWidth, screenHeight, "voronoi")
	setup()
	rl.SetTargetFPS(60)
	draw(screenWidth, screenHeight)

}

var sites []voronoi.Vertex

func setup() {
	inputDotsS := os.Args[1]
	inputDots, _ = strconv.Atoi(inputDotsS)
	dots = initDots(inputDots)
	sites = dotsToSites(dots)
}

func draw(screenWidth, screenHeight int32) {
	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyEnter) {
			dots = initDots(inputDots)
			sites = dotsToSites(dots)
		}

		bbox := voronoi.NewBBox(0, float64(screenWidth), 0, float64(screenHeight))
		diagram := voronoi.ComputeDiagram(sites, bbox, true)

		centroidedSites := make([]voronoi.Vertex, 0, len(sites))
		for _, cell := range diagram.Cells {
			poly := make([]voronoi.Vertex, 0, len(cell.Halfedges))
			for _, he := range cell.Halfedges {
				poly = append(poly, he.GetStartpoint())
			}
			centroidedSites = append(centroidedSites, centroid(poly))
		}

		dt := rl.GetFrameTime()
		t := float32(1.0 - math.Exp(-1.0*float64(dt)))
		for j := range sites {
			pos := rl.Vector2Lerp(
				rl.NewVector2(float32(sites[j].X), float32(sites[j].Y)),
				rl.NewVector2(float32(centroidedSites[j].X), float32(centroidedSites[j].Y)),
				t,
			)
			sites[j].X = float64(pos.X)
			sites[j].Y = float64(pos.Y)
		}

		// 5. Draw DIRECTLY to the screen
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Draw Voronoi lines
		for _, cell := range diagram.Cells {
			for _, he := range cell.Halfedges {
				p := he.GetStartpoint()
				p2 := he.GetEndpoint()
				rl.DrawLine(int32(p.X), int32(p.Y), int32(p2.X), int32(p2.Y), rl.Red)
			}
		}

		// Draw actual points (Black)
		for _, point := range sites {
			rl.DrawCircle(int32(point.X), int32(point.Y), 3, rl.Black)
		}

		// Draw target centroids (Green)
		for _, point := range centroidedSites {
			rl.DrawCircle(int32(point.X), int32(point.Y), 3, rl.Green)
		}

		rl.DrawText(fmt.Sprintf("Dots: %d | FPS: %d", len(sites), rl.GetFPS()), 10, 10, 20, rl.DarkGray)
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

func dotsToSites(dots []float64) []voronoi.Vertex {
	sites := make([]voronoi.Vertex, 0)
	for i := 0; i < len(dots); i += 2 {
		sites = append(sites, voronoi.Vertex{
			X: dots[i],
			Y: dots[i+1],
		})
	}

	return sites
}

func centroid(poly []voronoi.Vertex) voronoi.Vertex {
	var a, cx, cy float64
	n := len(poly)
	for i := 0; i < n; i++ {
		p, q := poly[i], poly[(i+1)%n]
		cross := p.X*q.Y - q.X*p.Y
		a += cross
		cx += (p.X + q.X) * cross
		cy += (p.Y + q.Y) * cross
	}
	a *= 0.5
	if math.Abs(a) < 1e-12 {
		return poly[0]
	}
	return voronoi.Vertex{X: cx / (6 * a), Y: cy / (6 * a)}
}
