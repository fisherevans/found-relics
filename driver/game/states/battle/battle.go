package battle

import (
	"combat/driver/game"
	"combat/pkg/gfxutil"
	"combat/pkg/rpg/combat"
	"combat/pkg/style"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	_ "image/png"
	"math"
	"time"
)

const statsPadding = 10.0

type State struct {
	Battle *combat.Battle

	playerCtrl   *PlayerControls
	opponentCtrl *OpponentControls

	characterRenderer       *CharacterRenderer
	characterRendererCanvas *pixelgl.Canvas

	debugTxt *text.Text

	bg *pixel.Sprite

	lastBattleTime combat.Time
}

func (s *State) Init(game *game.Game, win *pixelgl.Window) {
	bgPic, err := gfxutil.LoadPicture("res/sprites/temp_bg.png")
	s.bg = pixel.NewSprite(bgPic, bgPic.Bounds())

	atlas16, err := gfxutil.LoadTTF("res/fonts/alagard-16.ttf", 16)
	if err != nil {
		panic(err)
	}
	atlasLarge, err := gfxutil.LoadTTF("res/fonts/raidercrusader.ttf", 32)
	if err != nil {
		panic(err)
	}
	atlasMedium, err := gfxutil.LoadTTF("res/fonts/raidercrusader.ttf", 24)
	if err != nil {
		panic(err)
	}
	atlasSmall, err := gfxutil.LoadTTF("res/fonts/Mont-HeavyDEMO.ttf", 16)
	if err != nil {
		panic(err)
	}

	s.characterRenderer = NewCharacterRenderer(atlasSmall, atlasMedium, atlasLarge)
	s.characterRendererCanvas = pixelgl.NewCanvas(win.Bounds())

	s.debugTxt = text.New(pixel.V(0, 0), atlas16)

	s.playerCtrl = &PlayerControls{}
	s.opponentCtrl = &OpponentControls{}

	s.lastBattleTime = -1
}

func (s *State) Tick(game *game.Game, win *pixelgl.Window, dt float64) {
	currentBattleTime := realTimeToCombatTime(time.Now())
	var elapsed combat.Time = 0
	if s.lastBattleTime > 0 { // first tick
		//*
		elapsed = currentBattleTime - s.lastBattleTime
		//*/elapsed = combat.Time(float64(currentBattleTime-s.lastBattleTime) * 0.25) // slow it down for debugging
	}
	s.lastBattleTime = currentBattleTime

	s.playerCtrl.Update(win, s.Battle)
	s.opponentCtrl.Update(win, s.Battle)
	s.Battle.Progress(elapsed)
	s.characterRenderer.Update(dt)

	win.Clear(pixel.RGB(0.5, 0.5, 0.5))

	bgScale := math.Max(win.Bounds().W()/s.bg.Frame().W(), win.Bounds().H()/s.bg.Frame().H())
	s.bg.Draw(win, pixel.IM.Scaled(pixel.ZV, bgScale).Moved(win.Bounds().Center()))

	s.debugTxt.Clear()
	s.debugTxt.Color = colornames.Gray
	s.debugTxt.WriteString("Enemies:\n")
	s.debugTxt.Color = colornames.Red
	for _, c := range s.Battle.OpponentTeam {
		RenderCharacter(s.debugTxt, c)
	}
	s.debugTxt.Draw(win, pixel.IM.Moved(pixel.V(win.Bounds().W()-300, s.debugTxt.Bounds().H()+statsPadding)))

	s.characterRendererCanvas.Clear(style.Transparent)
	statsBounds := s.characterRenderer.Render(s.characterRendererCanvas, pixel.V(0, 0), s.Battle.PlayerTeam, s.playerCtrl.selected)
	s.characterRendererCanvas.Draw(win, pixel.IM.Moved(s.characterRendererCanvas.Bounds().Center().Add(pixel.V(statsPadding, win.Bounds().H()-statsBounds.Y-statsPadding))))

	if win.JustPressed(pixelgl.KeyEscape) {
		game.SwapState(NewExampleBattles())
	}
}

func realTimeToCombatTime(t time.Time) combat.Time {
	return combat.Time(t.UnixMilli())
}

func RenderCharacter(txt *text.Text, character *combat.CharacterInstance) {
	txt.WriteString(character.Details.Name + "\n")
	txt.WriteString(fmt.Sprintf("HP: %4d / %-4d\n\n", character.Life, character.Details.MaxLife))
	for _, m := range character.MoveQueue {
		txt.WriteString(fmt.Sprintf(" - %s (%4d/%4d)\n", m.Move.Name, m.ElapsedTime, m.Move.Duration))
	}
}
