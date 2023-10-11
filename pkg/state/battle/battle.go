package battle

import (
	"fmt"
	"found-relics/assets"
	"found-relics/pkg/drawutil"
	"found-relics/pkg/rpg/combat"
	"found-relics/pkg/state"
	"github.com/hajimehoshi/ebiten/v2"
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

	if game.Controller().JustPressed(state.InputMenu) {
		game.EnterSelector()
	}
}

func (s *State) Draw(game state.Game, screen *ebiten.Image) {
	bgImg := assets.Images.TempBattleBackground
	bgOpt := &ebiten.DrawImageOptions{}
	scale := math.Max(float64(screen.Bounds().Dx())/float64(bgImg.Bounds().Dx()), float64(screen.Bounds().Dy())/float64(bgImg.Bounds().Dy()))
	bgOpt.GeoM.Scale(scale, scale)
	screen.DrawImage(bgImg, bgOpt)

	s.hudRenderer.Draw(20, 20, game, s.battle, screen, s.playController.selected)

	s.debugSummary(screen)
}

func (s *State) debugSummary(screen *ebiten.Image) {
	text := drawutil.NewTextDrawer(assets.Fonts.TextSmall.Regular, screen.Bounds().Dx()-400, 20, color.White).
		Shadowed(1, 1, colornames.Black)
	render := func(char *combat.BattleCharacter, color color.Color) {
		text.Color(color)
		b := text.Draw(summarizeCharacter(char), screen)
		text.Move(0, b.Dy())
	}
	//for i, c := range s.battle.PlayerTeam {
	//	if s.playController.selected == i {
	//		render(c, colornames.Cyan)
	//	} else {
	//		render(c, colornames.White)
	//	}
	//}
	for _, c := range s.battle.OpponentTeam {
		render(c, colornames.Pink)
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
