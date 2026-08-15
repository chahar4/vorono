# 🌌 Voronoi Stippling & Relaxation in Go

A fun, real-time interactive Voronoi stippling visualizer written in Go! 

This project uses **Lloyd's Relaxation Algorithm** to distribute points across an image based on its luminance (darkness). The points gracefully animate and settle into a beautiful stippling pattern, creating more dots in darker areas and fewer dots in lighter areas.

![alt text](https://github.com/chahar4/vorono/blob/master/image/output.png)

## ✨ Features

- **Real-time Animation:** Smooth interpolation of points using `raylib-go`.
- **Image-Driven Density:** Uses an underlying image (`gown.jpg`) to calculate weighted centroids. Darker pixels exert a stronger "gravitational pull" on the Voronoi sites.
- **Optimized Rendering:** Pre-calculates a luminance map and uses bounding-box point-in-polygon checks for high performance.
- **Interactive:** Press `Enter` to instantly randomize and restart the simulation!

## 🧠 How It Works

1. **Luminance Mapping:** The app loads an image and converts it into a 1D array of luminance values (0.0 to 1.0) for ultra-fast pixel lookups.
2. **Initial Placement:** Random dots are generated. Dots falling on lighter pixels are rejected, naturally biasing the initial distribution toward darker areas of the image.
3. **Voronoi Diagram:** The `pzsz/voronoi` library calculates the Voronoi cells for the current sites.
4. **Weighted Centroids:** For each cell, the algorithm calculates the "Center of Mass" (weighted by the luminance map). 
5. **Relaxation (Lloyd's Algorithm):** The sites are smoothly interpolated (`Lerp`) toward their new weighted centroids. This loop runs every frame, creating a mesmerizing animation as the points snap into their optimal stippling positions.

## 🛠️ Prerequisites

This project uses [raylib-go](https://github.com/gen2brain/raylib-go), which requires the **Raylib C library** to be installed on your system.

### Linux (Ubuntu/Debian)
```bash
sudo apt-get install libraylib-dev
```
### macOS
```bash
brew install raylib
```
### Windows
Follow the raylib-go [Windows setup guide](https://github.com/gen2brain/raylib-go#windows) (usually involves installing via MSYS2 or using the pre-compiled binaries).

## 🚀 Installation & Usage

1. **Clone the repository:**
   ```bash
   git clone https://github.com/yourusername/voronoi-stippling.git
   cd voronoi-stippling
   ```

2. **Add an image:**
   Place a `.jpg` image in the root directory and name it `gown.jpg`. 
   *(Note: The screen size is hardcoded to `736x736`. For best results, use a square image of this resolution, or the app will scale/sample it!)*

3. **Install dependencies:**
   ```bash
   go mod init voronoi-stippling
   go get github.com/gen2brain/raylib-go/raylib
   go get github.com/pzsz/voronoi
   go mod tidy
   ```

4. **Run the project:**
   Pass the desired number of dots as a command-line argument (e.g., `2000`):
   ```bash
   go run main.go 2000
   ```

## 🎮 Controls

| Key | Action |
| :--- | :--- |
| `Enter` | Randomize dots and restart the relaxation animation. |
| `Esc` / Close Window | Exit the application. |

## 📦 Dependencies

- [raylib-go](https://github.com/gen2brain/raylib-go) - Go bindings for raylib (rendering and window management).
- [pzsz/voronoi](https://github.com/pzsz/voronoi) - Go implementation of Fortune's algorithm for Voronoi diagrams.
- Standard Go `image` and `image/jpeg` packages for image decoding.

## 📜 License

This project is open-source and available under the [MIT License](LICENSE). Feel free to use, modify, and distribute!

---
*Made with ❤️ and Go*
