package assets

import (
	_ "embed"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"log"
)

var (
	//go:embed fonts/bitmap/alagard-16.otf
	font_Alagard16 []byte

	//go:embed fonts/bitmap/runescape-16.ttf
	font_Runescape16 []byte

	//go:embed fonts/bitmap/pixelade-13.otf
	font_Pixelade13 []byte

	//go:embed fonts/bitmap/nokia-8.otf
	font_Nokia8 []byte
)

var Fonts = _FontAssets{}

type _FontAssets struct {
	TitleLarge font.Face
	TitleSmall font.Face

	TextRegular font.Face
	TextSmall   font.Face
	TextTiny    font.Face
}

func loadFonts() {
	fmt.Println("loading fonts...")
	Fonts.TitleLarge = loadFont(font_Alagard16, 32)
	Fonts.TitleSmall = loadFont(font_Alagard16, 16)
	Fonts.TextRegular = loadFont(font_Runescape16, 16)
	Fonts.TextSmall = loadFont(font_Pixelade13, 13)
	Fonts.TextTiny = loadFont(font_Nokia8, 8)
	fmt.Println("loaded all fonts")
}

func loadFont(data []byte, size int) font.Face {
	tt, err := opentype.Parse(data)
	if err != nil {
		log.Fatal(err)
	}
	ff, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    float64(size),
		DPI:     72,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	return ff
}
