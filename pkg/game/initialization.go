package game

import (
	"fmt"
	"found-relics/assets"
	"found-relics/pkg/controller"
	"found-relics/pkg/global"
	"found-relics/pkg/state"
	"found-relics/pkg/testdata"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
	"math"
)

const scale = 4

func Start() {

	dsf := int(math.Max(1, ebiten.DeviceScaleFactor()))
	fmt.Printf("dsf: %d, scale: %d\n", dsf, scale)
	ebiten.SetWindowSize(global.Width/dsf*scale, global.Height/dsf*scale)
	ebiten.SetWindowTitle("Found Relics")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)
	state := &InitState{}
	g := &FoundRelics{
		currentState: state,
		controller:   controller.FromEbiten(),
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func (g *FoundRelics) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return global.Width, global.Height
	/* high dpi
	s := ebiten.DeviceScaleFactor()
	return int(float64(outsideWidth) * s), int(float64(outsideHeight) * s)
	*/
}

type InitState struct {
	started bool
	done    bool
}

func (s *InitState) Update(game state.Game, dt float64) {
	if !s.started {
		go func() {
			assets.Initialize()
			testdata.Initialize()
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
