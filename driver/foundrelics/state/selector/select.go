package selector

import (
	"fmt"
	"found-relics/driver/foundrelics/assets"
	"found-relics/driver/foundrelics/state"
	"found-relics/pkg/rpg/combat"
	"found-relics/pkg/rpg/combat/testdata"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	if inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		s.Selected = (len(s.Options) + s.Selected - 1) % len(s.Options)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) || inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		s.Selected = (s.Selected + 1) % len(s.Options)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		game.EnterBattle(s.Options[s.Selected].Factory())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}
}

func (s *SelectState) Draw(game state.Game, screen *ebiten.Image) {
	x := float64(screen.Bounds().Dx()) * 0.1
	y := float64(assets.Fonts.TitleOutlineSmall.Metrics().Height.Round())
	text.DrawWithOptions(screen, "Select a Battle!", assets.Fonts.TitleOutlineSmall, drawTextOpt(x, y, colornames.White))
	y += float64(assets.Fonts.TextMedium.Bold.Metrics().Height.Round())
	for idx, opt := range s.Options {
		y += float64(assets.Fonts.TextMedium.Bold.Metrics().Height.Round())
		color := colornames.Gray
		if idx == s.Selected {
			color = colornames.Blue
		}
		text.DrawWithOptions(screen, fmt.Sprintf(" %d - %s", idx+1, opt.Name), assets.Fonts.TextMedium.Bold, drawTextOpt(x, y, color))
	}

}

func drawTextOpt(x, y float64, clr color.Color) *ebiten.DrawImageOptions {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(math.Round(x), math.Round(y))
	opt.ColorScale.ScaleWithColor(clr)
	return opt
}
