package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/gen2brain/raylib-go/raylib"
	"github.com/pzsz/voronoi"
)

const (
	SCREEN_WIDTH  = 736
	SCREEN_HEIGHT = 736
)

var inputDots int
var dots []float64
var screenWidth int32
var screenHeight int32
var gown image.Image
var luminanceMap []float64

func main() {
	screenWidth = int32(SCREEN_WIDTH)
	screenHeight = int32(SCREEN_HEIGHT)
	rl.InitWindow(screenWidth, screenHeight, "voronoi")
	setup()
	rl.SetTargetFPS(60)
	draw(screenWidth, screenHeight)

}

var sites []voronoi.Vertex

func setup() {
	inputDotsS := os.Args[1]
	inputDots, _ = strconv.Atoi(inputDotsS)
	img, err := getImageFromFilePath("./gown.jpg")
	if err != nil {
		log.Fatal("Error getting image from file", err)
	}
	gown = img
	dots = initDots(inputDots)
	sites = dotsToSites(dots)
	sites = sanitize(sites)

	luminanceMap = make([]float64, screenWidth*screenHeight)
	for y := 0; y < int(screenHeight); y++ {
		for x := 0; x < int(screenWidth); x++ {
			lum := luminancePercent(gown.At(x, y))
			luminanceMap[y*int(screenWidth)+x] = 1.0 - (lum / 100.0)
		}
	}
}

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func draw(screenWidth, screenHeight int32) {
	for !rl.WindowShouldClose() {
		if rl.IsKeyPressed(rl.KeyEnter) {
			dots = initDots(inputDots)
			sites = dotsToSites(dots)
			sites = sanitize(sites)

		}

		bbox := voronoi.NewBBox(0, float64(screenWidth), 0, float64(screenHeight))
		diagram := voronoi.ComputeDiagram(sites, bbox, true)

		type Accumulator struct {
			SumX, SumY, SumW float64
		}
		weightedCentroids := make([]voronoi.Vertex, len(sites))

		for i, cell := range diagram.Cells {
			poly := make([]voronoi.Vertex, 0, len(cell.Halfedges))
			minX, maxX := math.Inf(1), math.Inf(-1)
			minY, maxY := math.Inf(1), math.Inf(-1)

			for _, he := range cell.Halfedges {
				p := he.GetStartpoint()
				poly = append(poly, p)
				if p.X < minX {
					minX = p.X
				}
				if p.X > maxX {
					maxX = p.X
				}
				if p.Y < minY {
					minY = p.Y
				}
				if p.Y > maxY {
					maxY = p.Y
				}
			}

			geoCentroid := centroid(poly) // Fallback

			startX := int(math.Max(0, math.Floor(minX)))
			endX := int(math.Min(float64(screenWidth-1), math.Ceil(maxX)))
			startY := int(math.Max(0, math.Floor(minY)))
			endY := int(math.Min(float64(screenHeight-1), math.Ceil(maxY)))

			var sumX, sumY, sumW float64

			for y := startY; y <= endY; y++ {
				rowOffset := y * int(screenWidth)
				for x := startX; x <= endX; x++ {
					// Standard point-in-polygon check
					if isPointInPolygon(float64(x), float64(y), poly) {
						// Instant array lookup instead of slow gown.At()
						w := luminanceMap[rowOffset+x]
						sumX += float64(x) * w
						sumY += float64(y) * w
						sumW += w
					}
				}
			}

			if sumW > 0.0001 {
				weightedCentroids[i].X = sumX / sumW
				weightedCentroids[i].Y = sumY / sumW
			} else {
				weightedCentroids[i] = geoCentroid
			}
		}

		// movement
		dt := rl.GetFrameTime()
		t := float32(1.0 - math.Exp(-3.0*float64(dt)))

		for j := range sites {
			pos := rl.Vector2Lerp(
				rl.NewVector2(float32(sites[j].X), float32(sites[j].Y)),
				rl.NewVector2(float32(weightedCentroids[j].X), float32(weightedCentroids[j].Y)),
				t,
			)
			sites[j].X = float64(pos.X)
			sites[j].Y = float64(pos.Y)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		for _, cell := range diagram.Cells {
			for _, he := range cell.Halfedges {
				p := he.GetStartpoint()
				p2 := he.GetEndpoint()
				rl.DrawLine(int32(p.X), int32(p.Y), int32(p2.X), int32(p2.Y), rl.Black)
			}
		}
		for _, point := range sites {
			rl.DrawCircle(int32(point.X), int32(point.Y), 3, rl.Black)
		}

		rl.DrawText(fmt.Sprintf("Dots: %d | FPS: %d", len(sites), rl.GetFPS()), 10, 10, 20, rl.DarkGray)
		rl.EndDrawing()
	}
}

func initDots(len int) []float64 {
	dots := make([]float64, 0)
	for i := 0; i < len; i++ {
		x := rand.Intn(SCREEN_WIDTH)
		y := rand.Intn(SCREEN_HEIGHT)

		if float64(rand.Intn(100)) > luminancePercent(gown.At(x, y)) {
			dots = append(dots, float64(x))
			dots = append(dots, float64(y))
		} else {
			i--
		}

	}
	return dots
}

func isPointInPolygon(x, y float64, poly []voronoi.Vertex) bool {
	inside := false
	j := len(poly) - 1
	for i := 0; i < len(poly); i++ {
		if ((poly[i].Y > y) != (poly[j].Y > y)) &&
			(x < (poly[j].X-poly[i].X)*(y-poly[i].Y)/(poly[j].Y-poly[i].Y)+poly[i].X) {
			inside = !inside
		}
		j = i
	}
	return inside
}

func luminancePercent(c color.Color) float64 {
	r, g, b, _ := c.RGBA()

	redPercent := float64(r) / 65535 * 100
	greenPercent := float64(g) / 65535 * 100
	bluePercent := float64(b) / 65535 * 100

	return redPercent*0.2126 + greenPercent*0.7152 + bluePercent*0.0722
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

func sanitize(points []voronoi.Vertex) []voronoi.Vertex {
	clean := make([]voronoi.Vertex, 0, len(points))
	const minDist = 0.5

	for _, p := range points {
		p.X += (rand.Float64() - 0.5) * 0.01
		p.Y += (rand.Float64() - 0.5) * 0.01

		tooClose := false
		for _, c := range clean {
			dx := p.X - c.X
			dy := p.Y - c.Y
			if dx*dx+dy*dy < minDist*minDist {
				tooClose = true
				break
			}
		}

		if !tooClose {
			clean = append(clean, p)
		}
	}
	return clean
}
