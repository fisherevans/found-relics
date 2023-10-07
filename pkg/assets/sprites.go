package assets

import (
	"fmt"
	"github.com/faiface/pixel"
	"image"
	"os"
)

func LoadSprites() *Sprites {
	s := &Sprites{}
	for bin := 1; bin <= 8; bin++ {
		s.BattleIcons = append(s.BattleIcons, battleIconSheet(bin))
	}
	return s
}

type Sprites struct {
	BattleIcons []*SpriteSheet
}

type SpriteSheet struct {
	Location   string
	SpriteSize int

	Rows    int
	Columns int

	Picture pixel.Picture
	Sprites []*pixel.Sprite
}

type SpriteReference func(r Resources) *pixel.Sprite

func BattleIconReference(sheet, row, column int) SpriteReference {
	return func(r Resources) *pixel.Sprite {
		return r.Sprites.BattleIcons[sheet-1].GetSprite(row, column)
	}
}

func (ss *SpriteSheet) Load() {
	pic, err := LoadPicture(ss.Location)
	if err != nil {
		panic(fmt.Errorf("failed to load sprite sheet from %s: %w", ss.Location, err))
	}
	ss.Picture = pic
	ss.Columns = int(ss.Picture.Bounds().W()) / ss.SpriteSize
	ss.Rows = int(ss.Picture.Bounds().H()) / ss.SpriteSize
	for r := 1; r <= ss.Rows; r++ {
		for c := 1; c <= ss.Columns; c++ {
			ss.Sprites = append(ss.Sprites, GetSprite(pic, ss.SpriteSize, r, c))
		}
	}
	//fmt.Printf("loaded %d row x %d col @%d sprite sheet from %s\n", ss.Rows, ss.Columns, ss.SpriteSize, ss.Location)
}

func (ss *SpriteSheet) GetSprite(row, column int) *pixel.Sprite {
	id := ss.Columns*(row-1) + (column - 1)
	//fmt.Printf("r%d, c%d > id%d\n", row, column, id)
	return ss.Sprites[id]
}

func LoadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	//d, _ := os.Getwd()
	//fmt.Println(d)
	if err != nil {
		return nil, fmt.Errorf("failed to access file from os: %w", err)
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		s, _ := file.Stat()
		return nil, fmt.Errorf("failed to decode image file %s (%db): %w", file.Name(), s.Size(), err)
	}
	return pixel.PictureDataFromImage(img), nil
}

func GetSprite(sheet pixel.Picture, size, row, col int) *pixel.Sprite {
	bx := sheet.Bounds().Min.X + float64((col-1)*size)
	by := sheet.Bounds().Max.Y - float64(row*size)
	//fmt.Printf("size: %d, r:%d, c:%d, bx: %f, by:%f -- min:%s\n", size, row, col, bx, by, sheet.Bounds().Min.Text())
	bounds := pixel.R(bx, by, bx+float64(size), by+float64(size))
	return pixel.NewSprite(sheet, bounds)
}

func battleIconSheet(number int) *SpriteSheet {
	ss := &SpriteSheet{
		Location:   fmt.Sprintf("res/sprites/battle_icons/sheet_%03d.png", number),
		SpriteSize: 50,
	}
	ss.Load()
	return ss
}
