package battle

import (
	"fmt"
	"found-relics/assets"
	"found-relics/pkg/drawutil"
	"found-relics/pkg/rpg/combat"
	"found-relics/pkg/rpg/combat/moves"
	"found-relics/pkg/state"
	"found-relics/pkg/style"
	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"
	"image"
	"math"
	"strings"
)

type HudRenderer struct {
	padding int

	barBorderWidth      int
	barHealthWidthScale float64

	queueBoxBorder float64

	moveWidthScale float64
	moveMargin     float64
	moveBorder     float64
	movePadding    float64
	moveTextHeight float64

	currentMoveInnerWidth float64

	beatDividerWidth float64

	moveSlotSize   float64
	moveSlotBorder float64

	time float64
}

func NewHudRenderer() *HudRenderer {
	cr := &HudRenderer{
		padding: 4,

		barBorderWidth:        2,
		barHealthWidthScale:   0.05,
		currentMoveInnerWidth: 80,

		moveWidthScale: 0.02,
		moveMargin:     3,
		moveBorder:     2,
		movePadding:    1,

		queueBoxBorder: 2,

		moveSlotSize:   25,
		moveSlotBorder: 2,

		beatDividerWidth: 1,
	}
	cr.moveTextHeight = float64(assets.Fonts.TextSmall.Metrics().CapHeight.Ceil()) + cr.movePadding*2
	return cr
}

func (r *HudRenderer) flashSelected(a, b colorful.Color) colorful.Color {
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

func (r *HudRenderer) Draw(x, y int, game state.Game, battle *combat.Battle, target *ebiten.Image, selectedChar int) {
	hud := &HudRenderCtx{
		cfg:    r,
		target: target,
		bounds: image.Rectangle{
			Min: image.Point{X: x, Y: y},
			Max: image.Point{X: x + 1, Y: y + 1},
		},
	}

	for i, c := range battle.PlayerTeam {
		selected := i == selectedChar
		if i != 0 {
			hud.PadY()
			hud.PadY()
		}
		hud.RenderCharacterName(c, selected)
		hud.PadY()
		hud.RenderHealth(c, selected)
		hud.PadY()
		hud.RenderMoveQueue(c, selected)
		hud.PadY()
		hud.RenderMoveSlots(c, selected)
	}
}

func (h *HudRenderCtx) RenderCharacterName(c *combat.BattleCharacter, isSelected bool) {
	clr := style.ColorBright2
	if isSelected {
		clr = h.cfg.flashSelected(style.ColorHighlightBright, style.ColorHighlightDark)
	}
	b := drawutil.NewTextDrawer(assets.Fonts.TitleLarge, h.bounds.Min.X+h.cfg.padding, h.bounds.Max.Y, clr).
		Shadowed(1, 1, colornames.Black).
		//Move(0, assets.Fonts.TitleLarge.Metrics().Height.Round()).
		Draw(c.Details.Name, h.target)
	h.bounds = h.bounds.Union(b)
}

func (h *HudRenderCtx) RenderHealth(c *combat.BattleCharacter, isSelected bool) {
	bp := h.cfg.barBorderWidth
	innerWidth := int(float64(c.Details.MaxLife) * h.cfg.barHealthWidthScale)
	innerHeight := int(float64(assets.Fonts.TextTiny.Metrics().CapHeight.Ceil()) * 1)
	outerBarWidth := innerWidth + bp*4

	hpTxt := strings.ReplaceAll(fmt.Sprintf("%d", c.LagLife), "0", "O")
	maxHpTxt := strings.ReplaceAll(fmt.Sprintf("%d", c.Details.MaxLife), "0", "O")
	outOfText := strings.ReplaceAll(fmt.Sprintf("/%d HP", c.Details.MaxLife), "0", "O")
	hpDrawer := drawutil.NewTextDrawer(assets.Fonts.TextTiny, 0, 0, style.ColorBright1)
	hpDrawer.Shadowed(1, 1, style.ColorDark1)
	hpBounds := hpDrawer.BoundsOf(hpTxt)
	maxHpBounds := hpDrawer.BoundsOf(maxHpTxt)

	hpDrawer.Face = assets.Fonts.TextTiny
	outOfBounds := hpDrawer.BoundsOf(outOfText)
	hpDrawer.Face = assets.Fonts.TextTiny

	hpDrawer.Move(h.bounds.Min.X+outerBarWidth+h.cfg.padding, h.bounds.Max.Y+bp*2)

	barImg := gg.NewContext(outerBarWidth+h.cfg.padding+maxHpBounds.Dx()+h.cfg.padding/4+outOfBounds.Dx()+h.cfg.padding, innerHeight+bp*4)

	barImg.SetColor(style.Transparent(style.ColorDark1, 0.5))

	barImg.DrawRoundedRectangle(float64(outerBarWidth-bp-h.cfg.padding), 0, float64(barImg.Width()-(outerBarWidth-bp-h.cfg.padding)), float64(barImg.Height()), float64(bp*2))
	barImg.Fill()

	barImg.SetColor(style.ColorDark1)
	barImg.DrawRoundedRectangle(0, 0, float64(outerBarWidth), float64(barImg.Height()), float64(bp*2))
	barImg.Fill()

	barImg.SetColor(style.ColorBright1)
	barImg.DrawRoundedRectangle(float64(bp), float64(bp), float64(outerBarWidth-bp*2), float64(barImg.Height()-bp*2), float64(bp))
	barImg.Fill()

	barImg.SetColor(style.ColorHealthBg)
	barImg.DrawRectangle(float64(bp*2), float64(bp*2), float64(innerWidth), float64(innerHeight))
	barImg.Fill()

	healthWidth := int(float64(innerWidth) * (float64(c.Life) / float64(c.Details.MaxLife)))
	barImg.SetColor(style.ColorHealth)
	barImg.DrawRectangle(float64(bp*2), float64(bp*2), float64(healthWidth), float64(innerHeight))
	barImg.Fill()

	lagWidth := int((float64(c.LagLife) - float64(c.Life)) / float64(c.Details.MaxLife) * float64(innerWidth))
	if lagWidth > 0 {
		barImg.SetColor(style.ColorHealthLagDown)
	} else {
		barImg.SetColor(style.ColorHealthLagUp)
	}
	barImg.DrawRectangle(float64(bp*2+healthWidth), float64(bp*2), float64(lagWidth), float64(innerHeight))
	barImg.Fill()

	barImgOpt := &ebiten.DrawImageOptions{}
	barImgOpt.GeoM.Translate(float64(h.bounds.Min.X), float64(h.bounds.Max.Y))
	h.target.DrawImage(ebiten.NewImageFromImage(barImg.Image()), barImgOpt)

	h.bounds = h.bounds.Union(image.Rect(h.bounds.Min.X, h.bounds.Max.Y, h.bounds.Min.X+barImg.Width(), h.bounds.Max.Y+barImg.Height()))

	hpDrawer.Move(maxHpBounds.Dx()-hpBounds.Dx(), 0)
	hpDrawer.Draw(hpTxt, h.target)

	hpDrawer.Move(hpBounds.Dx()+h.cfg.padding/4, (assets.Fonts.TextTiny.Metrics().CapHeight.Ceil()-assets.Fonts.TextTiny.Metrics().CapHeight.Ceil())/2)
	hpDrawer.Face = assets.Fonts.TextTiny
	hpDrawer.AlignHorizontal = drawutil.AlignLeft
	hpDrawer.Color(style.ColorBright3)
	hpDrawer.Draw(outOfText, h.target)
}

func (h *HudRenderCtx) PadY() {
	h.bounds.Max.Y += h.cfg.padding
}

func (h *HudRenderCtx) RenderMoveQueue(char *combat.BattleCharacter, isSelected bool) {
	pct, transitionPct := 0.0, 0.0
	if len(char.MoveQueue) > 0 {
		pct = math.Min(float64(char.MoveQueue[0].ElapsedTime)/400.0, 1.0)
		transitionPct = drawutil.InterpolateSin(pct)
	}

	canvasH := h.cfg.queueBoxBorder*2 + h.cfg.moveMargin*2 + h.cfg.moveBorder*2 + h.cfg.movePadding*2 + h.cfg.moveTextHeight
	canvasW := h.cfg.queueBoxBorder*2 + h.cfg.currentMoveInnerWidth + h.cfg.moveWidthScale*float64(char.Details.MaxMoveQueueDepth.ToCombatTime()) + h.cfg.moveMargin

	// make canvas
	canvas := gg.NewContext(int(canvasW), int(canvasH))

	// draw background box
	drawutil.NewStrokedRectangle(0, 0, canvasW, canvasH).Stroked(h.cfg.queueBoxBorder, drawutil.StrokeInside).Draw(canvas)
	canvas.SetColor(style.ColorDark2)
	canvas.FillPreserve()
	canvas.SetColor(style.ColorDark3)
	if isSelected {
		canvas.SetColor(h.cfg.flashSelected(style.ColorHighlightBright, style.ColorHighlightDark))
	}
	canvas.Stroke()

	// draw beat lines
	for b := combat.Beats(1); b < char.Details.MaxMoveQueueDepth; b++ {
		bx := float64(b.ToCombatTime())*h.cfg.moveWidthScale + h.cfg.currentMoveInnerWidth + h.cfg.queueBoxBorder
		bh := (canvasH - h.cfg.queueBoxBorder*2) * 0.66
		by := (canvasH - bh) / 2.0
		canvas.DrawRectangle(bx, by, 0, bh)
		canvas.SetColor(style.ColorDark3)
		canvas.SetLineWidth(h.cfg.beatDividerWidth)
		canvas.Stroke()
	}

	// draw current move box
	var currentColor colorful.Color
	if len(char.MoveQueue) == 0 {
		currentColor = h.cfg.flashSelected(style.ColorDark2, style.ColorDark3)
	} else {
		currentColor = style.Flash(style.ColorPurpleDark, style.ColorPurpleBright, h.cfg.time, 1)
	}
	canvas.DrawRectangle(h.cfg.queueBoxBorder, h.cfg.queueBoxBorder, h.cfg.currentMoveInnerWidth, canvasH-h.cfg.queueBoxBorder*2)
	canvas.SetColor(currentColor)
	canvas.Fill()

	// write no current move warning
	if len(char.MoveQueue) == 0 {
		clr := h.cfg.flashSelected(style.ColorGray, style.ColorGray)
		clr = style.ColorBright2
		drawutil.NewTextDrawer(assets.Fonts.TextSmall, int(h.cfg.queueBoxBorder), int(h.cfg.queueBoxBorder), clr).
			Bounded(int(h.cfg.currentMoveInnerWidth), int(float64(canvas.Height())-h.cfg.queueBoxBorder*2), drawutil.AlignCenter, drawutil.AlignMiddle).
			DrawToCanvas("No Move!", canvas)
	}

	// draw queued moves
	if len(char.MoveQueue) > 0 {
		x := h.cfg.currentMoveInnerWidth + h.cfg.queueBoxBorder
		// print upcoming
		for id, m := range char.MoveQueue {
			boxH := h.cfg.moveTextHeight + h.cfg.moveBorder*2 + h.cfg.movePadding*2

			totalW := float64(m.Move.Duration.ToCombatTime()) * h.cfg.moveWidthScale
			moveBoxW := totalW - h.cfg.moveMargin
			fillColor, borderColor := style.ColorDark3, style.ColorBright2
			nameColor := style.ColorBright2
			textDx := 0.0
			borderThickness := h.cfg.moveBorder
			boxDx := 0.0
			// current move shrinks, aligned left
			if id == 0 {
				totalW = float64(m.Move.Duration.ToCombatTime()-m.ElapsedTime) * h.cfg.moveWidthScale
				moveBoxW = totalW - h.cfg.moveMargin
				fillColor = style.Flash(fillColor, currentColor, transitionPct, 0.5)
				borderColor = style.Flash(borderColor, currentColor, transitionPct, 0.5)
				borderThickness = borderThickness + (boxH-borderThickness*2)*transitionPct
				textDx = -1.0 * transitionPct * h.cfg.currentMoveInnerWidth
				nameColor = style.Flash(nameColor, style.ColorBright1, transitionPct, 0.5)
			}
			if totalW <= 0 {
				fmt.Printf("wtf %f\n", totalW)
			}

			boxX := math.Round(x) + h.cfg.moveMargin + boxDx
			boxY := h.cfg.queueBoxBorder + h.cfg.moveMargin
			boxW := math.Max(moveBoxW-boxDx, 0)

			//borderThickness = math.Min(borderThickness, boxW/2.0)

			drawutil.NewStrokedRectangle(boxX, boxY, boxW, boxH).
				Stroked(borderThickness, drawutil.StrokeInside).
				Draw(canvas)

			canvas.SetColor(fillColor)
			canvas.FillPreserve()
			canvas.SetColor(borderColor)
			canvas.Stroke()

			td := drawutil.NewTextDrawer(assets.Fonts.TextSmall, int(boxX+h.cfg.movePadding*2+h.cfg.moveBorder), int(boxY+h.cfg.movePadding), nameColor).
				Move(int(textDx), 0).
				Bounded(int(boxW-h.cfg.movePadding*4-h.cfg.moveBorder*2), int(boxH-h.cfg.movePadding*2), drawutil.AlignLeft, drawutil.AlignMiddle)
			if id == 0 && transitionPct >= 0.5 {
				td.Shadowed(1, 1, style.ColorDark1)
			}
			td.DrawToCanvas(string(m.Move.Name), canvas)
			x += totalW
		}
	}

	barImgOpt := &ebiten.DrawImageOptions{}
	barImgOpt.GeoM.Translate(float64(h.bounds.Min.X), float64(h.bounds.Max.Y))
	h.target.DrawImage(ebiten.NewImageFromImage(canvas.Image()), barImgOpt)
	h.bounds = h.bounds.Union(drawutil.NewSizedRect(h.bounds.Min.X, h.bounds.Max.Y, canvas.Width(), canvas.Height()))

}

func (h *HudRenderCtx) RenderMoveSlots(char *combat.BattleCharacter, selected bool) {
	canvas := gg.NewContext(int(h.cfg.moveSlotSize*4+h.cfg.moveSlotBorder*16), int(h.cfg.moveSlotSize+h.cfg.moveSlotBorder*2))
	textDrawer := drawutil.NewTextDrawer(assets.Fonts.TextTiny, h.bounds.Min.X+int(h.cfg.moveSlotBorder*2), h.bounds.Max.Y+int(h.cfg.moveSlotSize+h.cfg.moveSlotBorder), style.ColorBright1).
		Bounded(int(h.cfg.moveSlotSize), assets.Fonts.TextTiny.Metrics().Height.Round(), drawutil.AlignCenter, drawutil.AlignBottom)

	remainingQueue := char.Details.MaxMoveQueueDepth.ToCombatTime() - char.MoveQueueTimeDepth
	for i, moveId := range char.Details.Moves.AsSlice() {
		if moveId == moves.None {
			continue
		}
		move := moves.Get(moveId)
		x := float64(h.cfg.moveSlotSize+h.cfg.moveSlotBorder*4)*float64(i) + h.cfg.moveSlotBorder*2
		y := h.cfg.moveSlotBorder
		opt := &ebiten.DrawImageOptions{}
		opt.GeoM.Scale(
			h.cfg.moveSlotSize/float64(move.Sprite.Bounds().Dx()),
			h.cfg.moveSlotSize/float64(move.Sprite.Bounds().Dy()))
		opt.GeoM.Translate(float64(h.bounds.Min.X)+x, float64(h.bounds.Max.Y)+y)
		h.target.DrawImage(move.Sprite, opt)
		drawutil.NewStrokedRectangle(x, y, h.cfg.moveSlotSize, h.cfg.moveSlotSize).
			Stroked(h.cfg.moveSlotBorder, drawutil.StrokeOutside).
			Draw(canvas)

		if move.Duration.ToCombatTime() > remainingQueue {
			canvas.SetColor(style.Transparent(style.ColorGray, 0.5))
			canvas.FillPreserve()
			canvas.SetColor(style.ColorDark3)
		} else {
			canvas.SetColor(style.ColorBright2)
		}
		canvas.StrokePreserve()
		canvas.ClearPath()

		textDrawer.Draw(fmt.Sprintf("%d", i+1), h.target)
		textDrawer.Move(int(h.cfg.moveSlotSize+h.cfg.moveSlotBorder*4), 0)
	}
	barImgOpt := &ebiten.DrawImageOptions{}
	barImgOpt.GeoM.Translate(float64(h.bounds.Min.X), float64(h.bounds.Max.Y))
	h.target.DrawImage(ebiten.NewImageFromImage(canvas.Image()), barImgOpt)
	h.bounds = h.bounds.Union(drawutil.NewSizedRect(h.bounds.Min.X, h.bounds.Max.Y, canvas.Width(), canvas.Height()+assets.Fonts.TextTiny.Metrics().Height.Round()+int(h.cfg.moveSlotBorder)))
}
