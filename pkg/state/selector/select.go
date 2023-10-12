package selector

import (
	"fmt"
	"found-relics/assets"
	"found-relics/pkg/rpg/combat"
	"found-relics/pkg/state"
	"found-relics/pkg/testdata"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	"os"
)

type SelectState struct {
	Selected int
	Options  []Option
}

type Option struct {
	Name    string
	Factory func() *combat.Battle
}

func NewExampleBattles() *SelectState {
	return &SelectState{
		Options: []Option{
			{
				Name:    "Simple 2v1",
				Factory: testdata.Battle2v1,
			},
			{
				Name:    "The same 2v1",
				Factory: testdata.Battle2v1,
			},
		},
	}
}

func (s *SelectState) Update(game state.Game, dt float64) {
	if game.Controller().JustPressed(state.InputUp) {
		s.Selected = (len(s.Options) + s.Selected - 1) % len(s.Options)
	}
	if game.Controller().JustPressed(state.InputDown) {
		s.Selected = (s.Selected + 1) % len(s.Options)
	}
	if game.Controller().JustPressed(state.InputAliasSelect) {
		game.EnterBattle(s.Options[s.Selected].Factory())
	}
	if game.Controller().JustPressed(state.InputMenu) {
		os.Exit(0)
	}
}

func (s *SelectState) Draw(game state.Game, screen *ebiten.Image) {
	x := float64(screen.Bounds().Dx()) * 0.1
	y := float64(assets.Fonts.TitleLarge.Metrics().Height.Round())
	text.DrawWithOptions(screen, "Select a Battle!", assets.Fonts.TitleLarge, drawTextOpt(x, y, colornames.White))
	y += float64(assets.Fonts.TextRegular.Metrics().Height.Round())
	for idx, opt := range s.Options {
		y += float64(assets.Fonts.TextRegular.Metrics().Height.Round())
		color := colornames.Gray
		if idx == s.Selected {
			color = colornames.Blue
		}
		text.DrawWithOptions(screen, fmt.Sprintf(" %d - %s", idx+1, opt.Name), assets.Fonts.TextRegular, drawTextOpt(x, y, color))
	}
}

func drawTextOpt(x, y float64, clr color.Color) *ebiten.DrawImageOptions {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(math.Round(x), math.Round(y))
	opt.ColorScale.ScaleWithColor(clr)
	return opt
}
