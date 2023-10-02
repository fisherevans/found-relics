package battle

import (
	"combat/driver/game"
	"combat/pkg/gfxutil"
	"combat/pkg/rpg/combat"
	"combat/pkg/rpg/combat/moves"
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

type SelectState struct {
	Options []Option
	txt     *text.Text
}

type Option struct {
	Name    string
	Factory func() *combat.Battle
}

var keys = []pixelgl.Button{
	pixelgl.Key1,
	pixelgl.Key2,
	pixelgl.Key3,
	pixelgl.Key4,
	pixelgl.Key5,
	pixelgl.Key6,
	pixelgl.Key7,
	pixelgl.Key8,
	pixelgl.Key9,
}

func NewExampleBattles() *SelectState {
	return &SelectState{
		Options: []Option{
			{
				Name: "Hello, World!",
				Factory: func() *combat.Battle {
					return &combat.Battle{
						PlayerTeam: []*combat.CharacterInstance{
							{
								Details: combat.CharacterDetails{
									Name: "Hank the Tank",
									Moves: combat.AvailableMoves{
										Slot1: moves.HeavySwing,
										Slot2: moves.Flourish,
										Slot3: moves.SitThere,
									},
									MaxLife:           5000,
									MaxMoveQueueDepth: 15,
								},
								Life: 5000,
							}, //*
							{
								Details: combat.CharacterDetails{
									Name: "Reluctant Healer",
									Moves: combat.AvailableMoves{
										Slot1: moves.HealSome,
										Slot2: moves.BloodRitual,
										Slot3: moves.SitThere,
									},
									MaxLife:           3000,
									MaxMoveQueueDepth: 10,
								},
								Life: 3000,
							}, //*/
						},
						OpponentTeam: []*combat.CharacterInstance{
							{
								Details: combat.CharacterDetails{
									Name: "Bad Guy",
									Moves: combat.AvailableMoves{
										Slot1: moves.Eat,
										Slot2: moves.Poke,
										Slot3: moves.SitThere,
									},
									MaxLife: 2000,
								},
								Life: 2000,
							},
						},
					}
				},
			},
		},
	}
}

func (s *SelectState) Init(game *game.Game, win *pixelgl.Window) {
	atlas, err := gfxutil.LoadTTF("res/fonts/alagard-16.ttf", 32)
	if err != nil {
		panic(err)
	}
	s.txt = text.New(pixel.ZV, atlas)
}

func (s *SelectState) Tick(game *game.Game, win *pixelgl.Window, dt float64) {
	win.Clear(colornames.Black)
	s.txt.Clear()
	s.txt.WriteString("Select a battle:\n")
	for idx, opt := range s.Options {
		s.txt.WriteString(fmt.Sprintf(" %d) %s", idx+1, opt.Name))
		if win.JustPressed(keys[idx]) {
			game.SwapState(&State{Battle: opt.Factory()})
		}
	}
	s.txt.Draw(win, pixel.IM.Moved(pixel.V(20, win.Bounds().H()-s.txt.Bounds().H())))

	if win.JustPressed(pixelgl.KeyEscape) {
		win.SetClosed(true)
	}

	imd := imdraw.New(nil)

	gfxutil.RectShape{
		Bounds:         gfxutil.Box(pixel.V(100, 100), 100, 50),
		Fill:           true,
		FillColor:      pixel.RGB(1, 0, 0),
		StrokeWidth:    10,
		StrokeColor:    pixel.RGB(0, 0, 1),
		StrokeShape:    imdraw.SharpEndShape,
		StrokePosition: gfxutil.StrokeInner,
	}.Draw(imd)

	gfxutil.RectShape{
		Bounds:    gfxutil.Box(pixel.V(100, 100), 100, 50),
		Fill:      true,
		FillColor: pixel.Alpha(0.5),
	}.Draw(imd)

	gfxutil.RectShape{
		Bounds:         gfxutil.Box(pixel.V(300, 100), 100, 50),
		Fill:           true,
		FillColor:      pixel.RGB(1, 0, 0),
		StrokeWidth:    10,
		StrokeColor:    pixel.RGB(0, 0, 1),
		StrokeShape:    imdraw.SharpEndShape,
		StrokePosition: gfxutil.StrokeCenter,
	}.Draw(imd)

	gfxutil.RectShape{
		Bounds:    gfxutil.Box(pixel.V(300, 100), 100, 50),
		Fill:      true,
		FillColor: pixel.Alpha(0.5),
	}.Draw(imd)

	gfxutil.RectShape{
		Bounds:         gfxutil.Box(pixel.V(500, 100), 100, 50),
		Fill:           true,
		FillColor:      pixel.RGB(1, 0, 0),
		StrokeWidth:    10,
		StrokeColor:    pixel.RGB(0, 0, 1),
		StrokeShape:    imdraw.SharpEndShape,
		StrokePosition: gfxutil.StrokeOuter,
	}.Draw(imd)

	gfxutil.RectShape{
		Bounds:    gfxutil.Box(pixel.V(500, 100), 100, 50),
		Fill:      true,
		FillColor: pixel.Alpha(0.5),
	}.Draw(imd)

	gfxutil.RectShape{
		Bounds:         gfxutil.Box(pixel.V(700, 100), 20, 20),
		Fill:           true,
		FillColor:      pixel.RGB(1, 0, 0),
		StrokeWidth:    8,
		StrokeColor:    pixel.RGB(0, 0, 1),
		StrokeShape:    imdraw.SharpEndShape,
		StrokePosition: gfxutil.StrokeInner,
	}.Draw(imd)

	gfxutil.RectShape{
		Bounds:         gfxutil.Box(pixel.V(700, 200), 20, 20),
		Fill:           true,
		FillColor:      pixel.RGB(1, 0, 0),
		StrokeWidth:    30,
		StrokeColor:    pixel.RGB(0, 0, 1),
		StrokeShape:    imdraw.SharpEndShape,
		StrokePosition: gfxutil.StrokeInner,
	}.Draw(imd)

	imd.Draw(win)

	s.txt.Clear()

	gfxutil.Text{
		String: "Hi!",
		Color:  pixel.RGB(0, 1, 0),
		Bounds: gfxutil.Box(pixel.V(100, 100), 100, 50),
		VAlign: gfxutil.Top,
		HAlign: gfxutil.Left,
	}.Draw(s.txt, pixel.ZV)

	gfxutil.Text{
		String: "Hi!",
		Color:  pixel.RGB(0, 1, 0),
		Bounds: gfxutil.Box(pixel.V(300, 100), 100, 50),
		VAlign: gfxutil.Middle,
		HAlign: gfxutil.Center,
	}.Draw(s.txt, pixel.ZV)

	gfxutil.Text{
		String: "Hi!",
		Color:  pixel.RGB(0, 1, 0),
		Bounds: gfxutil.Box(pixel.V(500, 100), 100, 50),
		VAlign: gfxutil.Bottom,
		HAlign: gfxutil.Right,
	}.Draw(s.txt, pixel.ZV)

	s.txt.Draw(win, pixel.IM)
}
