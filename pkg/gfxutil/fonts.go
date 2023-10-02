package gfxutil

import (
	"github.com/faiface/pixel/text"
	"github.com/golang/freetype/truetype"
	"io/ioutil"
	"os"
)

func LoadTTF(path string, size float64) (*text.Atlas, error) {
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

	face := truetype.NewFace(font, &truetype.Options{
		Size:              size,
		GlyphCacheEntries: 1,
	})

	return text.NewAtlas(face, text.ASCII), nil
}
