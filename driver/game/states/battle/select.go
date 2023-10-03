package battle

import (
	"fmt"
	"found-relics/driver/game"
	"found-relics/pkg/gfxutil"
	"found-relics/pkg/rpg/combat"
	"found-relics/pkg/rpg/combat/testdata"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type SelectState struct {
	Options []Option
	txt     *text.Text
}

type Option struct {
	Name    string
	Factory func() *combat.Battle
}

var keys = []pixelgl.Button{
	pixelgl.Key1,
	pixelgl.Key2,
	pixelgl.Key3,
	pixelgl.Key4,
	pixelgl.Key5,
	pixelgl.Key6,
	pixelgl.Key7,
	pixelgl.Key8,
	pixelgl.Key9,
}

func NewExampleBattles() *SelectState {
	return &SelectState{
		Options: []Option{
			{
				Name:    "Simple 2v1",
				Factory: testdata.Battle2v1,
			},
		},
	}
}

func (s *SelectState) Init(game *game.Game, win *pixelgl.Window) {
	atlas, err := gfxutil.LoadTTF("res/fonts/alagard-16.ttf", 32)
	if err != nil {
		panic(err)
	}
	s.txt = text.New(pixel.ZV, atlas)
}

func (s *SelectState) Tick(game *game.Game, win *pixelgl.Window, dt float64) {
	win.Clear(colornames.Black)
	s.txt.Clear()
	s.txt.WriteString("Select a battle:\n")
	for idx, opt := range s.Options {
		s.txt.WriteString(fmt.Sprintf(" %d) %s", idx+1, opt.Name))
		if win.JustPressed(keys[idx]) {
			game.SwapState(&State{Battle: opt.Factory()})
		}
	}
	s.txt.Draw(win, pixel.IM.Moved(pixel.V(20, win.Bounds().H()-s.txt.Bounds().H())))

	if win.JustPressed(pixelgl.KeyEscape) {
		win.SetClosed(true)
	}
}
