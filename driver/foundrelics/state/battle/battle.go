package battle

import (
	"fmt"
	"found-relics/driver/foundrelics/assets"
	"found-relics/driver/foundrelics/drawutil"
	"found-relics/driver/foundrelics/state"
	"found-relics/pkg/rpg/combat"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	"time"
)

type State struct {
	battle          *combat.Battle
	combatTimeScale float64

	playController     PlayerController
	opponentController OpponentController

	hudRenderer *HudRenderer

	lastBattleTime combat.Time
}

func NewBattle(battle *combat.Battle) *State {
	return &State{
		battle:             battle,
		combatTimeScale:    1.0,
		playController:     PlayerController{},
		opponentController: OpponentController{},
		hudRenderer:        NewHudRenderer(),
	}
}

func (s *State) Load() {
	time.Sleep(time.Millisecond * 100)
}

func (s *State) Update(game state.Game, dt float64) {
	currentBattleTime := combat.Time(int(float64(time.Now().UnixMilli()) * s.combatTimeScale))
	var elapsed combat.Time = 0
	if s.lastBattleTime > 0 { // first tick
		elapsed = currentBattleTime - s.lastBattleTime
	}
	s.lastBattleTime = currentBattleTime

	s.battle.Progress(elapsed)
	s.playController.Update(game, s.battle, dt, elapsed)
	s.opponentController.Update(game, s.battle, dt, elapsed)
	s.hudRenderer.Update(game, s.battle, dt, elapsed)

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		game.EnterSelector()
	}
}

func (s *State) Draw(game state.Game, screen *ebiten.Image) {
	bgImg := assets.Images.TempBattleBackground
	bgOpt := &ebiten.DrawImageOptions{}
	scale := math.Max(float64(screen.Bounds().Dx())/float64(bgImg.Bounds().Dx()), float64(screen.Bounds().Dy())/float64(bgImg.Bounds().Dy()))
	bgOpt.GeoM.Scale(scale, scale)
	screen.DrawImage(bgImg, bgOpt)

	s.hudRenderer.Draw(20, 20, game, s.battle, screen)

	s.debugSummary(screen)
}

func (s *State) debugSummary(screen *ebiten.Image) {
	text := drawutil.
		NewTextDrawer(assets.Fonts.TextMedium.Regular, 100, 700, color.White).
		Shadowed(2, 2, colornames.Gray)
	for i, c := range s.battle.PlayerTeam {
		text.Color(colornames.White)
		if s.playController.selected == i {
			text.Color(colornames.Cyan)
		}
		text.Draw(summarizeCharacter(c), screen)
		text.Move(700, 120)
	}
	for _, c := range s.battle.OpponentTeam {
		ebitenutil.DrawRect(screen, float64(text.X), float64(text.Y), 100, 100, color.RGBA{100, 0, 100, 100})
		text.Color(colornames.Red)
		text.Bounded(100, 100, drawutil.AlignRight, drawutil.AlignBottom)
		b := text.Draw(summarizeCharacter(c), screen)
		text.Move(700, 120)
		ebitenutil.DrawRect(screen, float64(b.Min.X), float64(b.Min.Y), float64(b.Dx()), float64(b.Dy()), color.RGBA{0, 100, 100, 100})
	}
}

func summarizeCharacter(character *combat.BattleCharacter) string {
	summary := character.Details.Name + "\n"
	summary += fmt.Sprintf("HP: %4d / %-4d\n", character.Life, character.Details.MaxLife)
	summary += fmt.Sprintf("Queue: %2d / %-2d\n", character.MoveQueueTimeDepth.ToBeatRoundedUp(), character.Details.MaxMoveQueueDepth)
	for _, m := range character.MoveQueue {
		summary += fmt.Sprintf(" - %s (%4d/%4d)\n", m.Move.Name, m.ElapsedTime.ToBeatRoundedDown(), m.Move.Duration)
	}
	return summary
}
