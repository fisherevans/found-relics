package battle

import (
	"fmt"
	"found-relics/driver/foundrelics/assets"
	"found-relics/driver/foundrelics/drawutil"
	"found-relics/driver/foundrelics/state"
	"found-relics/pkg/rpg/combat"
	"found-relics/pkg/style"
	"github.com/faiface/pixel"
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"image"
)

type HudRenderer struct {
	faceName      font.Face
	faceBarNumber font.Face

	padding int

	barBorderWidth      int
	barHealthWidthScale float64

	time float64
}

func NewHudRenderer() *HudRenderer {
	cr := &HudRenderer{
		faceName:      assets.Fonts.TitleSmall,
		faceBarNumber: assets.Fonts.TextSmall.Bold,

		barBorderWidth:      6,
		barHealthWidthScale: 0.2,

		padding: 12,
	}
	return cr
}

func (r *HudRenderer) flashSelected(a, b pixel.RGBA) pixel.RGBA {
	return style.Flash(a, b, r.time, 3)
}

func (r *HudRenderer) Update(game state.Game, battle *combat.Battle, dt float64, elapsed combat.Time) {
	r.time += dt
}

type HudRenderCtx struct {
	cfg    *HudRenderer
	target *ebiten.Image
	bounds image.Rectangle
}

func (r *HudRenderer) Draw(x, y int, game state.Game, battle *combat.Battle, target *ebiten.Image) {
	hud := &HudRenderCtx{
		cfg:    r,
		target: target,
		bounds: image.Rectangle{
			Min: image.Point{X: x, Y: y},
			Max: image.Point{X: x, Y: y},
		},
	}

	for i, c := range battle.PlayerTeam {
		if i != 0 {
			hud.PadY()
			hud.PadY()
			hud.PadY()
		}
		hud.RenderCharacterName(c)
		hud.PadY()
		hud.RenderHealth(c)
	}
}

func (h *HudRenderCtx) RenderCharacterName(c *combat.BattleCharacter) {
	b := drawutil.NewTextDrawer(h.cfg.faceName, h.bounds.Min.X, h.bounds.Max.Y, h.cfg.flashSelected(style.ColorHighlightBright, style.ColorHighlightDark)).
		Shadowed(4, 4, colornames.Black).
		Draw(c.Details.Name, h.target)
	h.bounds = h.bounds.Union(b)
}

func (h *HudRenderCtx) RenderHealth(c *combat.BattleCharacter) {
	bp := h.cfg.barBorderWidth
	innerWidth := int(float64(c.Details.MaxLife) * h.cfg.barHealthWidthScale)
	innerHeight := h.cfg.faceBarNumber.Metrics().CapHeight.Ceil()
	outerBarWidth := innerWidth + bp*4

	hpTxt := fmt.Sprintf("%d / %d HP", c.LagLife, c.Details.MaxLife)
	//hpDrawer := drawutil.NewTextDrawer(h.cfg.faceBarNumber, outerBorder.Max.X+h.cfg.padding, h.bounds.Max.Y+bp*2, style.ColorBright1)
	hpDrawer := drawutil.NewTextDrawer(h.cfg.faceBarNumber, 0, 0, style.ColorBright1)
	hpBounds := hpDrawer.BoundsOf(hpTxt)

	barImg := gg.NewContext(outerBarWidth+h.cfg.padding+hpBounds.Dx()+h.cfg.padding, innerHeight+bp*4)
	hpDrawer.Move(h.bounds.Min.X+outerBarWidth+h.cfg.padding, h.bounds.Max.Y+bp*2)

	barImg.SetColor(style.ColorDark1.Scaled(0.5))
	barImg.DrawRoundedRectangle(float64(outerBarWidth-bp-h.cfg.padding), 0, float64(barImg.Width()-(outerBarWidth-bp-h.cfg.padding)), float64(barImg.Height()), float64(bp*2))
	barImg.Fill()

	barImg.SetColor(style.ColorDark1)
	barImg.DrawRoundedRectangle(0, 0, float64(outerBarWidth), float64(barImg.Height()), float64(bp*2))
	barImg.Fill()

	barImg.SetColor(style.ColorBright1)
	barImg.DrawRoundedRectangle(float64(bp), float64(bp), float64(outerBarWidth-bp*2), float64(barImg.Height()-bp*2), float64(bp))
	barImg.Fill()

	barImg.SetColor(style.ColorHealthBg)
	barImg.DrawRoundedRectangle(float64(bp*2), float64(bp*2), float64(outerBarWidth-bp*4), float64(barImg.Height()-bp*4), float64(bp))
	barImg.Fill()

	barImgOpt := &ebiten.DrawImageOptions{}
	barImgOpt.GeoM.Translate(float64(h.bounds.Min.X), float64(h.bounds.Max.Y))
	h.target.DrawImage(ebiten.NewImageFromImage(barImg.Image()), barImgOpt)

	outerBorder := image.Rect(h.bounds.Min.X, h.bounds.Max.Y, h.bounds.Min.X+innerWidth+h.cfg.barBorderWidth*4, h.bounds.Max.Y+innerHeight+bp*4)
	//drawutil.DrawRect(outerBorder, style.ColorGray, h.target)

	innerBorder := outerBorder.Inset(bp)
	//drawutil.DrawRect(innerBorder, style.ColorBright1, h.target)

	bar := innerBorder.Inset(bp)
	//drawutil.DrawRect(bar, style.ColorHealthBg, h.target)

	barFgDx := int(float64(bar.Dx()) * (float64(c.Life) / float64(c.Details.MaxLife)))
	bar.Max.X -= bar.Dx() - barFgDx
	drawutil.DrawRect(bar, style.ColorHealth, h.target)

	barLagDx := int((float64(c.LagLife) - float64(c.Life)) / float64(c.Details.MaxLife) * float64(innerWidth))
	bar.Min.X = bar.Max.X + barLagDx
	drawutil.DrawRect(bar, style.ColorHealthLag, h.target)

	h.bounds = h.bounds.Union(hpDrawer.Draw(hpTxt, h.target))

	h.bounds = h.bounds.Union(outerBorder)
}

func (h *HudRenderCtx) PadY() {
	h.bounds.Max.Y += h.cfg.padding
}
