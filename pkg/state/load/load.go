package load

import (
	"found-relics/assets"
	"found-relics/pkg/drawutil"
	"found-relics/pkg/state"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

const dotSpeed = 0.333
const maxDots = 3

type loadingStage int

const (
	firstTick loadingStage = iota
	loading
	complete
)

type State struct {
	toLoad state.LoadableState

	stage loadingStage

	text string
	dots int
	time float64
}

func LoadState(toLoad state.LoadableState) *State {
	s := &State{
		toLoad: toLoad,
	}
	s.reset()
	return s
}

func (s *State) Update(game state.Game, dt float64) {
	if s.stage == firstTick {
		s.stage = loading
		go func() {
			s.toLoad.Load()
			s.stage = complete
		}()
	}
	if s.stage == complete {
		game.EnterState(s.toLoad)
	}
	s.time += dt
	if s.time > dotSpeed {
		s.time -= dotSpeed
		if s.dots >= maxDots {
			s.reset()
		} else {
			s.dot()
		}
	}
}

func (s *State) dot() {
	s.dots++
	s.text += "."
}

func (s *State) reset() {
	s.dots = 0
	s.text = "Loading"
}

func (s *State) Draw(game state.Game, screen *ebiten.Image) {
	drawutil.DrawString(s.text, assets.Fonts.TitleLarge, 100, 100, color.White, screen)
}
