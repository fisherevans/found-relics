package assets

import (
	"bytes"
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"log"
)

var (
	//go:embed images/temp_bg.png
	img_tempBg []byte
)

var Images = ImageAssets{}

type ImageAssets struct {
	TempBattleBackground *ebiten.Image
}

func loadImages() {
	fmt.Println("loading images...")
	Images.TempBattleBackground = loadImage(img_tempBg)
	fmt.Println("loaded all images")
}

func loadImage(data []byte) *ebiten.Image {
	reader := bytes.NewReader(img_tempBg)
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}
