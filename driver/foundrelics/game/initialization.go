package game

import (
	"found-relics/driver/foundrelics/assets"
	"found-relics/driver/foundrelics/state"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

var img *ebiten.Image

func init() {
	var err error
	img, _, err = ebitenutil.NewImageFromFile("res/sprites/temp_bg.png")
	if err != nil {
		log.Fatal(err)
	}
}

func Start() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Found Relics")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	state := &InitState{}
	g := &FoundRelics{
		currentState: state,
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

type InitState struct {
	started bool
	done    bool
}

func (s *InitState) Update(game state.Game, dt float64) {
	if !s.started {
		go func() {
			assets.Initialize()
			s.done = true
		}()
		s.started = true
	}
	if s.done {
		game.EnterSelector()
	}
}

func (s *InitState) Draw(game state.Game, screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Loading...")
}
