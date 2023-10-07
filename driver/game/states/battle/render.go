package battle

import (
	"fmt"
	"found-relics/driver/game"
	"found-relics/pkg/gfxutil"
	"found-relics/pkg/rpg/combat"
	"found-relics/pkg/style"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/text"
	"math"
)

type CharacterRenderer struct {
	game *game.Game

	imd       *imdraw.IMDraw
	txtSmall  *text.Text
	txtMedium *text.Text
	txtLarge  *text.Text

	padding float64

	barBorderThickness float64
	healthWidthScale   float64
	healthHeight       float64

	moveWidthScale        float64
	moveHeight            float64
	moveMargin            float64
	moveBorder            float64
	movePadding           float64
	currentMoveInnerWidth float64

	queueBoxPadding float64
	queueBoxBorder  float64

	beatDividerWidth float64

	slotSize   float64
	slotMargin float64

	textShadow pixel.Vec

	time float64
}

func NewCharacterRenderer(game *game.Game, atlasSmall, atlasMedium, atlasLarge *text.Atlas) *CharacterRenderer {
	cr := &CharacterRenderer{
		game: game,

		imd:       imdraw.New(nil),
		txtSmall:  text.New(pixel.ZV, atlasSmall),
		txtMedium: text.New(pixel.ZV, atlasMedium),
		txtLarge:  text.New(pixel.ZV, atlasLarge),
		padding:   4.0,

		barBorderThickness: 4.0,
		healthWidthScale:   0.05,

		moveWidthScale: 0.05,
		moveMargin:     6,
		moveBorder:     2,
		movePadding:    4,

		queueBoxPadding: 4,
		queueBoxBorder:  4,

		beatDividerWidth: 1,

		slotSize:   50,
		slotMargin: 4,

		textShadow: pixel.V(2, -2),
	}

	cr.healthHeight = cr.txtSmall.BoundsOf("0").H() - cr.barBorderThickness

	largeMoveBounds := cr.txtSmall.BoundsOf("Some Long Move")
	cr.currentMoveInnerWidth = largeMoveBounds.W()
	cr.moveHeight = largeMoveBounds.H() + cr.movePadding*2 + cr.moveBorder*2

	return cr
}

type RenderCtx struct {
	target pixel.Target
	matrix pixel.Matrix
	maxY   float64
	maxX   float64
}

func (ctx *RenderCtx) TrackMaxX(x float64) {
	ctx.maxX = math.Max(x, ctx.maxX)
}

func (cr *CharacterRenderer) Update(dt float64) {
	cr.time += dt
}

func (cr *CharacterRenderer) Render(target pixel.Target, pos pixel.Vec, chars []*combat.BattleCharacter, selectedCid int) pixel.Vec {
	cr.imd.Clear()
	cr.txtLarge.Clear()
	cr.txtSmall.Clear()

	ctx := &RenderCtx{
		target: target,
		matrix: pixel.IM.Moved(pos),
	}

	for cid := len(chars) - 1; cid >= 0; cid -= 1 {
		char := chars[cid]
		isSelected := cid == selectedCid

		cr.renderSlots(ctx, char)

		cr.renderMoveQueue(ctx, char, isSelected)
		ctx.maxY += cr.padding

		cr.renderBar(ctx, char.Life, char.LagLife, char.Details.MaxLife, "HP", cr.healthWidthScale, style.ColorHealth, style.ColorHealthLag, style.ColorHealthBg)
		ctx.maxY += cr.padding * 2

		cr.renderName(ctx, char.Details.Name, isSelected)
	}

	cr.imd.SetMatrix(ctx.matrix)
	cr.imd.Draw(target)
	cr.txtLarge.Draw(target, ctx.matrix)
	cr.txtSmall.Draw(target, ctx.matrix)

	return pixel.V(ctx.maxX, ctx.maxY)
}

func (cr *CharacterRenderer) renderName(ctx *RenderCtx, name string, selected bool) {
	color := style.ColorBright2
	if selected {
		color = cr.flashSelected(style.ColorHighlightBright, style.ColorHighlightDark)
	}
	cr.txtLarge.Dot = pixel.V(cr.padding, ctx.maxY)
	cr.textShadowed(cr.txtLarge, name, color)
	ctx.maxY += cr.txtLarge.BoundsOf(name).H()
	ctx.TrackMaxX(cr.txtLarge.Dot.X)
}

func (cr *CharacterRenderer) renderBar(ctx *RenderCtx, current, lag, max int, label string, widthScale float64, colorFg, colorLag, colorBg pixel.RGBA) {
	v := func(x, y float64) pixel.Vec {
		return pixel.V(x+cr.barBorderThickness, y+cr.barBorderThickness+ctx.maxY)
	}
	barWidth := float64(max) * widthScale
	currentWidth := float64(current) / float64(max) * barWidth
	lagWidth := float64(lag) / float64(max) * barWidth

	cr.imd.Color = style.ColorDark1
	cr.imd.EndShape = imdraw.RoundEndShape
	cr.imd.Push(v(0, cr.healthHeight),
		v(0, 0),
		v(barWidth, 0),
		v(barWidth, cr.healthHeight))
	cr.imd.Polygon(cr.barBorderThickness * 2)

	cr.imd.Color = style.ColorBright1
	cr.imd.EndShape = imdraw.RoundEndShape
	cr.imd.Push(v(0, cr.healthHeight),
		v(0, 0),
		v(barWidth, 0),
		v(barWidth, cr.healthHeight))
	cr.imd.Polygon(cr.barBorderThickness * 1)

	cr.imd.EndShape = imdraw.NoEndShape
	cr.imd.Color = colorBg
	cr.imd.Push(v(0, cr.healthHeight/2), v(barWidth, cr.healthHeight/2))
	cr.imd.Line(cr.healthHeight)

	cr.imd.Color = colorFg
	cr.imd.Push(v(0, cr.healthHeight/2), v(currentWidth, cr.healthHeight/2))
	cr.imd.Line(cr.healthHeight)

	if lag >= 0 {
		cr.imd.Color = colorLag
		cr.imd.Push(v(currentWidth, cr.healthHeight/2), v(lagWidth, cr.healthHeight/2))
		cr.imd.Line(cr.healthHeight)
	}

	cr.txtSmall.Dot = v(barWidth+cr.barBorderThickness*2+cr.padding, 0)
	cr.textShadowed(cr.txtSmall, fmt.Sprintf("%d", current), style.ColorBright1)
	cr.textShadowed(cr.txtSmall, fmt.Sprintf("/%d %s", max, label), style.ColorBright3)

	ctx.maxY += cr.healthHeight + cr.barBorderThickness*2
	ctx.TrackMaxX(cr.txtSmall.Dot.X)
}

func (cr *CharacterRenderer) renderMoveQueue(ctx *RenderCtx, char *combat.BattleCharacter, selected bool) {
	v := func(x, y float64) pixel.Vec {
		return pixel.V(x, y+ctx.maxY)
	}
	pct, transitionPct := 0.0, 0.0
	if len(char.MoveQueue) > 0 {
		pct = math.Min(float64(char.MoveQueue[0].ElapsedTime)/400.0, 1.0)
		transitionPct = gfxutil.InterpolateSin(pct)
	}

	queueBox := gfxutil.Box(
		v(cr.queueBoxBorder, cr.queueBoxBorder),
		cr.currentMoveInnerWidth+cr.moveWidthScale*float64(char.Details.MaxMoveQueueDepth.ToCombatTime())+cr.queueBoxPadding,
		cr.moveHeight+cr.queueBoxPadding*2)
	queueBorderColor := style.ColorDark3
	if selected {
		queueBorderColor = cr.flashSelected(style.ColorHighlightBright, style.ColorHighlightDark)
	}
	gfxutil.RectShape{
		Bounds:         queueBox,
		Fill:           true,
		FillColor:      style.ColorDark2,
		StrokePosition: gfxutil.StrokeOuter,
		StrokeWidth:    cr.queueBoxBorder,
		StrokeColor:    queueBorderColor,
		StrokeShape:    imdraw.RoundEndShape,
	}.Draw(cr.imd)
	margin := cr.queueBoxPadding + cr.queueBoxBorder

	for b := combat.Beats(0); b < char.Details.MaxMoveQueueDepth; b++ {
		bx := float64(b.ToCombatTime())*cr.moveWidthScale + cr.currentMoveInnerWidth + margin
		if b > 0 {
			bx -= cr.moveMargin / 2.0
		}
		gfxutil.RectShape{
			Bounds: gfxutil.Box(
				v(bx, cr.queueBoxBorder+cr.beatDividerWidth),
				0, queueBox.H()-cr.queueBoxBorder/2.0),
			StrokePosition: gfxutil.StrokeOuter,
			StrokeWidth:    cr.beatDividerWidth,
			StrokeColor:    style.ColorDark3,
			StrokeShape:    imdraw.NoEndShape,
		}.Draw(cr.imd)
	}

	currentColor := style.Flash(style.ColorPurpleDark.Scaled(0.8), style.ColorPurpleBright, transitionPct, 0.5)
	if len(char.MoveQueue) == 0 {
		currentColor = cr.flashSelected(style.ColorDark3.Scaled(1.2), style.ColorDark3)
	}
	currentBox := gfxutil.Box(v(margin, margin), cr.currentMoveInnerWidth, cr.moveHeight)
	gfxutil.RectShape{
		Bounds:    currentBox,
		Fill:      true,
		FillColor: currentColor,
	}.Draw(cr.imd)
	if len(char.MoveQueue) == 0 {
		gfxutil.Text{
			String: "No Move!",
			Color:  cr.flashSelected(style.ColorGray.Scaled(1.2), style.ColorGray),
			Bounds: currentBox,
			VAlign: gfxutil.Middle,
			HAlign: gfxutil.Center,
		}.Draw(cr.txtSmall, pixel.V(0, 2))
	}

	if len(char.MoveQueue) > 0 {
		x := cr.currentMoveInnerWidth
		// print upcoming
		for id, m := range char.MoveQueue {
			totalW := float64(m.Move.Duration.ToCombatTime()) * cr.moveWidthScale
			moveBoxW := totalW - cr.moveMargin
			fillColor, borderColor := style.ColorDark3, style.ColorBright2
			nameColor := style.ColorBright2
			textDx := 0.0
			borderThickness := cr.moveBorder
			if id == 0 {
				// box
				totalW = float64(m.Move.Duration.ToCombatTime()-m.ElapsedTime) * cr.moveWidthScale
				moveBoxW = totalW
				fillColor = style.Flash(fillColor, style.ColorPurpleDark, transitionPct, 0.5)
				borderColor = style.Flash(borderColor, style.ColorPurpleBright, transitionPct, 0.5)
				borderThickness = borderThickness + borderThickness*transitionPct*2.0
				// drawtext.go
				textDx = -1.0 * transitionPct * cr.currentMoveInnerWidth
				nameColor = style.Flash(nameColor, style.ColorBright1, transitionPct, 0.5)
			}
			if totalW <= 0 {
				fmt.Printf("wtf %f\n", totalW)
			}
			box := gfxutil.Box(v(math.Round(x)+margin, margin), moveBoxW, cr.moveHeight)
			gfxutil.RectShape{
				Bounds:         box,
				Fill:           true,
				FillColor:      fillColor,
				StrokeWidth:    borderThickness,
				StrokeColor:    borderColor,
				StrokeShape:    imdraw.RoundEndShape,
				StrokePosition: gfxutil.StrokeInner,
			}.Draw(cr.imd)
			gfxutil.Text{
				String: m.Move.Name,
				Color:  nameColor,
				Bounds: box,
				VAlign: gfxutil.Middle,
				HAlign: gfxutil.Left,
			}.Draw(cr.txtSmall, pixel.V(cr.movePadding+cr.moveBorder+textDx, 2))
			x += totalW
		}
	}
	ctx.maxY += queueBox.H() + cr.queueBoxBorder*2
	ctx.TrackMaxX(cr.txtSmall.Dot.X)
}

func (cr *CharacterRenderer) renderSlots(ctx *RenderCtx, char *combat.BattleCharacter) {
	for i, move := range char.Details.Moves.AsSlice() {
		if move == nil {
			continue
		}
		s := move.Sprite(cr.game.Resources)
		s.Draw(ctx.target, pixel.IM.
			ScaledXY(pixel.ZV, pixel.V(cr.slotSize/s.Frame().W(), cr.slotSize/s.Frame().H())).
			Moved(pixel.V(float64(i)*(cr.slotSize+cr.slotMargin*2.0), ctx.maxY)).
			Moved(pixel.V(cr.slotSize/2, cr.slotSize/2)))
	}
	ctx.maxY += cr.slotSize + cr.slotMargin*2.0
}

func (cr *CharacterRenderer) flashSelected(a, b pixel.RGBA) pixel.RGBA {
	return style.Flash(a, b, cr.time, 3)
}

func (cr *CharacterRenderer) textShadowed(txt *text.Text, str string, color pixel.RGBA) {
	origD := txt.Dot
	txt.Color = style.ColorDark1
	txt.Dot = origD.Add(cr.textShadow)
	txt.WriteString(str)

	txt.Dot = origD
	txt.Color = color
	txt.WriteString(str)
}
