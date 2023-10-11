package assets

import (
	_ "embed"
	"fmt"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"log"
)

var (
	//go:embed fonts/RaiderCrusader-D9XE.otf
	font_RaiderCrusader []byte

	//go:embed fonts/RaiderCrusader3D-ABd6.otf
	font_RaiderCrusader3D []byte

	//go:embed fonts/SimplyMono-Bold.ttf
	font_SimplyMono []byte

	//go:embed fonts/MintSansExtraLight.otf
	font_MintLight []byte

	//go:embed fonts/MintSansRegular.otf
	font_MintRegular []byte

	//go:embed fonts/MintSansBold.otf
	font_MintBold []byte
)

var Fonts = _FontAssets{}

type _FontAssets struct {
	TitleLarge        font.Face
	TitleOutlineLarge font.Face

	TitleSmall        font.Face
	TitleOutlineSmall font.Face

	TitleTiny        font.Face
	TitleOutlineTiny font.Face

	TextLarge  StyledFont
	TextMedium StyledFont
	TextSmall  StyledFont

	NumbersMicro        font.Face
	NumbersSmall        font.Face
	NumbersOutlineSmall font.Face

	NumbersMedium        font.Face
	NumbersOutlineMedium font.Face
}

type StyledFont struct {
	Regular font.Face
	Bold    font.Face
	Thin    font.Face
}

func loadFonts() {
	fmt.Println("loading fonts...")
	Fonts.TitleLarge = loadFont(font_RaiderCrusader, 64)
	Fonts.TitleOutlineLarge = loadFont(font_RaiderCrusader3D, 64)

	Fonts.TitleSmall = loadFont(font_RaiderCrusader, 32)
	Fonts.TitleOutlineSmall = loadFont(font_RaiderCrusader3D, 32)

	Fonts.TitleTiny = loadFont(font_RaiderCrusader, 24)
	Fonts.TitleOutlineTiny = loadFont(font_RaiderCrusader3D, 24)

	Fonts.TextLarge = StyledFont{
		Regular: loadFont(font_MintRegular, 24),
		Bold:    loadFont(font_MintBold, 24),
		Thin:    loadFont(font_MintLight, 24),
	}

	Fonts.TextMedium = StyledFont{
		Regular: loadFont(font_MintRegular, 18),
		Bold:    loadFont(font_MintBold, 18),
		Thin:    loadFont(font_MintLight, 18),
	}

	Fonts.TextSmall = StyledFont{
		Regular: loadFont(font_MintRegular, 12),
		Bold:    loadFont(font_MintBold, 12),
		Thin:    loadFont(font_MintLight, 12),
	}

	Fonts.NumbersMicro = loadFont(font_SimplyMono, 8)
	Fonts.NumbersSmall = loadFont(font_SimplyMono, 12)
	Fonts.NumbersMedium = loadFont(font_SimplyMono, 24)

	fmt.Println("loaded all fonts")
}

func loadFont(data []byte, size float64) font.Face {
	tt, err := opentype.Parse(data)
	if err != nil {
		log.Fatal(err)
	}
	ff, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    size,
		DPI:     150,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}
	return ff
}
