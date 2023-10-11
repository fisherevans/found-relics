package selector

import (
	"fmt"
	"found-relics/assets"
	"found-relics/pkg/drawutil"
	"found-relics/pkg/rpg/combat"
	"found-relics/pkg/state"
	"found-relics/pkg/style"
	"found-relics/pkg/testdata"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

	clr := color.RGBA{100, 100, 100, 100}

	ebitenutil.DrawRect(screen, 700, 100, 100, 110, style.ColorHealth)
	b1 := drawutil.NewTextDrawer(assets.Fonts.TextSmall.Regular, 700, 100, style.ColorDark1).
		Bounded(100, 110, drawutil.AlignCenter, drawutil.AlignMiddle).
		WithAlignModeHorizontal(drawutil.HAlignLine).
		WithAlignModeVertical(drawutil.VAlignLineHeight).
		Draw("ppp.\nppppgqy.\nqy.", screen)
	ebitenutil.DrawRect(screen, float64(b1.Min.X), float64(b1.Min.Y), float64(b1.Dx()), float64(b1.Dy()), clr)

	ebitenutil.DrawRect(screen, 700, 300, 100, 110, style.ColorHealth)
	b2 := drawutil.NewTextDrawer(assets.Fonts.TextSmall.Regular, 700, 300, style.ColorDark1).
		Bounded(100, 110, drawutil.AlignCenter, drawutil.AlignMiddle).
		Draw("ppppp.\npy.", screen)
	ebitenutil.DrawRect(screen, float64(b2.Min.X), float64(b2.Min.Y), float64(b2.Dx()), float64(b2.Dy()), clr)

	ebitenutil.DrawRect(screen, 700, 500, 100, 110, style.ColorHealth)
	b3 := drawutil.NewTextDrawer(assets.Fonts.TextSmall.Regular, 700, 500, style.ColorDark1).
		Bounded(100, 110, drawutil.AlignLeft, drawutil.AlignTop).
		WithAlignModeVertical(drawutil.VAlignLineHeight).
		Draw("Ppp.", screen)
	ebitenutil.DrawRect(screen, float64(b3.Min.X), float64(b3.Min.Y), float64(b3.Dx()), float64(b3.Dy()), clr)

}

func drawTextOpt(x, y float64, clr color.Color) *ebiten.DrawImageOptions {
	opt := &ebiten.DrawImageOptions{}
	opt.GeoM.Translate(math.Round(x), math.Round(y))
	opt.ColorScale.ScaleWithColor(clr)
	return opt
}
