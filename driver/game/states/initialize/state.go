package initialize

import (
	"found-relics/driver/game"
	"found-relics/driver/game/states/battle"
	"found-relics/pkg/assets"
	"found-relics/pkg/gfxutil"
	"found-relics/pkg/style"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
)

const animSpeed = 0.33
const maxDots = 5

func New() *State {
	return &State{}
}

type State struct {
	txt                *text.Text
	sinceLastAnimation float64
	dotCount           int
	done               bool
}

func (s *State) Init(game *game.Game, win *pixelgl.Window) {
	atlas, err := gfxutil.LoadTTF("res/fonts/alagard-16.ttf", 32)
	if err != nil {
		panic(err)
	}
	s.txt = text.New(pixel.ZV, atlas)
	s.txt.Color = style.ColorBright2
	s.resetAnim()
	go func() {
		game.Resources = assets.LoadAll()
		s.done = true
	}()
}

func (s *State) Tick(game *game.Game, win *pixelgl.Window, dt float64) {
	win.Clear(style.ColorDark1)
	s.sinceLastAnimation += dt
	if s.sinceLastAnimation > animSpeed {
		s.sinceLastAnimation -= animSpeed
		if s.dotCount >= maxDots {
			s.resetAnim()
		} else {
			s.addDot()
		}
	}
	s.txt.Draw(win, pixel.IM.Moved(pixel.V(50, 50)))
	if s.done {
		game.SwapState(battle.NewExampleBattles())
	}
}

func (s *State) resetAnim() {
	s.txt.Clear()
	s.txt.WriteString("Loading")
	s.dotCount = 0
}

func (s *State) addDot() {
	s.txt.WriteString(".")
	s.dotCount++
}
