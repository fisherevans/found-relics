package assets

import (
	"bytes"
	"embed"
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	_ "image/png"
	"log"
)

var (
	//go:embed images/temp_bg.png
	img_tempBg []byte

	//go:embed images/move_icons/move_icons_sheet_*.png
	img_moveIcons embed.FS
)

func loadImages() {
	fmt.Println("loading images...")
	Images.TempBattleBackground = loadImage(img_tempBg)
	Images.MoveIconSheets = loadSpriteList(50, &img_moveIcons)
	fmt.Println("loaded all images")
}

var Images = _ImageAssets{}

type _ImageAssets struct {
	TempBattleBackground *ebiten.Image
	MoveIconSheets       []*_SpriteSheet
}

type _SpriteSheet struct {
	Image      *ebiten.Image
	SpriteSize int
}

func (ss *_SpriteSheet) Sprite(row, col int) *ebiten.Image {
	x := (col - 1) * ss.SpriteSize
	y := (row - 1) * ss.SpriteSize
	return ss.Image.SubImage(image.Rect(x, y, x+ss.SpriteSize, y+ss.SpriteSize)).(*ebiten.Image)
}

func loadSpriteList(spriteSize int, efs *embed.FS) []*_SpriteSheet {
	var sheets []*_SpriteSheet
	for _, img := range loadImageList(efs) {
		sheets = append(sheets, &_SpriteSheet{
			Image:      img,
			SpriteSize: spriteSize,
		})
	}
	return sheets
}

func loadImageList(efs *embed.FS) []*ebiten.Image {
	datas, err := getAllFileBytes(efs)
	if err != nil {
		log.Fatal("failed to load embedded images", err)
	}
	var images []*ebiten.Image
	for _, data := range datas {
		images = append(images, loadImage(data))
	}
	return images
}

func loadImage(data []byte) *ebiten.Image {
	reader := bytes.NewReader(data)
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(img)
}
