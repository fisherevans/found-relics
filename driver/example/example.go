package example

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"image"
	_ "image/png"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"time"
)

var fragmentShader = `
#version 330 core

in vec2 vTexCoords;
out vec4 fragColor;

uniform sampler2D uTexture;
uniform vec4 uTexBounds;

// custom uniforms
uniform float uSpeed;
uniform float uTime;

void main() {
    vec2 t = vTexCoords / uTexBounds.zw;
	vec3 influence = texture(uTexture, t).rgb;

    if (influence.r + influence.g + influence.b > 0.3) {
		t.y += cos(t.x * 40.0 + (uTime * uSpeed))*0.005;
		t.x += cos(t.y * 40.0 + (uTime * uSpeed))*0.01;
	}

    vec3 col = texture(uTexture, t).rgb;
	fragColor = vec4(col * vec3(0.6, 0.6, 1.2),texture(uTexture, t).a);
}
`

func Run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Combat Prototype",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	//win.SetSmooth(true)

	face, err := loadTTF("res/fonts/alagard-16.ttf", 16*4)
	atlas := text.NewAtlas(face, text.ASCII)
	txt := text.New(pixel.V(0, 0), atlas)
	//txt.LineHeight = atlas.LineHeight() * 1.5

	sheet, err := loadPicture("trees.png")
	if err != nil {
		panic(err)
	}
	sheetBatch := pixel.NewBatch(&pixel.TrianglesData{}, sheet)

	var trees []*pixel.Sprite
	for x := 0; x < 3; x++ {
		for y := 0; y < 3; y++ {
			trees = append(trees, tileSprite(sheet, 32, x, y))
		}
	}

	type Planted struct {
		sprite    *pixel.Sprite
		transform pixel.Matrix
	}
	var planted []Planted

	var uTime, uSpeed float32 = 0.0, 5.0
	win.Canvas().SetUniform("uTime", &uTime)
	win.Canvas().SetUniform("uSpeed", &uSpeed)
	win.Canvas().SetFragmentShader(fragmentShader)

	cameraPos := pixel.ZV
	cameraSpeed := 500.0
	cameraZoom := 4.0
	cameraZoomSpeed := 1.2

	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	targetFps := time.Tick(time.Second / 120)
	start := time.Now()
	last := start
	lastDraw := time.Now()
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		uTime = float32(time.Since(start).Seconds())
		if win.Pressed(pixelgl.KeyRight) {
			uSpeed += 0.1
		}
		if win.Pressed(pixelgl.KeyLeft) {
			uSpeed -= 0.1
		}

		if win.Pressed(pixelgl.KeyA) {
			cameraPos.X -= cameraSpeed * dt
		}
		if win.Pressed(pixelgl.KeyD) {
			cameraPos.X += cameraSpeed * dt
		}
		if win.Pressed(pixelgl.KeyS) {
			cameraPos.Y -= cameraSpeed * dt
		}
		if win.Pressed(pixelgl.KeyW) {
			cameraPos.Y += cameraSpeed * dt
		}

		cameraZoom *= math.Pow(cameraZoomSpeed, win.MouseScroll().Y)
		cameraZoom = math.Max(cameraZoom, 0.2)
		cameraZoom = math.Min(cameraZoom, 4)

		camera := pixel.IM.
			Scaled(cameraPos, cameraZoom).
			Moved(win.Bounds().Center().Sub(cameraPos))
		win.SetMatrix(camera)

		if win.Pressed(pixelgl.MouseButtonLeft) && time.Since(lastDraw).Seconds() > 0.05 {
			lastDraw = last
			planted = append(planted, Planted{
				sprite: trees[rand.Intn(len(trees))],
				transform: pixel.IM.
					Scaled(pixel.ZV, rand.Float64()*4+1).
					Moved(camera.Unproject(win.MousePosition())),
			})
		}

		win.Clear(colornames.Forestgreen)

		sheetBatch.Clear()
		for _, tree := range planted {
			tree.sprite.Draw(sheetBatch, tree.transform)
		}
		sheetBatch.Draw(win)

		//txt.WriteString(win.Typed())
		typeShadowed(txt, win.Typed())
		if win.JustPressed(pixelgl.KeyEnter) || win.Repeated(pixelgl.KeyEnter) {
			typeShadowed(txt, ".\n")
			//txt.WriteString(".\n")
		}
		txt.Draw(win, pixel.IM)

		win.Update()

		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d | Trees: %d | Zoom: %0.1f", cfg.Title, frames, len(planted), cameraZoom))
			frames = 0
		default:
		}

		<-targetFps
	}
}

func typeShadowed(txt *text.Text, v string) {
	shadow := pixel.V(4.0, -4.0)

	origD := txt.Dot
	txt.Color = colornames.Black
	txt.Dot = origD.Add(shadow)
	txt.WriteString(v)

	txt.Dot = origD
	txt.Color = colornames.Whitesmoke
	txt.WriteString(v)
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func tileSprite(sheet pixel.Picture, size, x, y int) *pixel.Sprite {
	bx := sheet.Bounds().Min.X + float64(x*size)
	by := sheet.Bounds().Min.Y + float64(y*size)
	bounds := pixel.R(bx, by, bx+float64(size), by+float64(size))
	return pixel.NewSprite(sheet, bounds)
}

func loadTTF(path string, size float64) (font.Face, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	font, err := truetype.Parse(bytes)
	if err != nil {
		return nil, err
	}

	return truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	}), nil
}
