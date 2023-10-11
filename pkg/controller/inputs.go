package controller

import (
	"found-relics/pkg/state"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type inputController struct {
	mappings map[state.InputKey][]ebiten.Key
}

func FromEbiten() state.Controller {
	return &inputController{
		mappings: map[state.InputKey][]ebiten.Key{
			state.InputUp:    {ebiten.KeyUp, ebiten.KeyW},
			state.InputDown:  {ebiten.KeyDown, ebiten.KeyS},
			state.InputLeft:  {ebiten.KeyLeft, ebiten.KeyA},
			state.InputRight: {ebiten.KeyRight, ebiten.KeyD},

			state.InputOpt1: {ebiten.Key1, ebiten.KeyKP1, ebiten.KeyEnter, ebiten.KeyKPEnter},
			state.InputOpt2: {ebiten.Key2, ebiten.KeyKP2},
			state.InputOpt3: {ebiten.Key3, ebiten.KeyKP3},
			state.InputOpt4: {ebiten.Key4, ebiten.KeyKP4},

			state.InputAlt:  {ebiten.KeySpace, ebiten.KeyShift, ebiten.KeyShiftRight, ebiten.KeyShiftLeft},
			state.InputBack: {ebiten.KeyBackspace},

			state.InputMenu:    {ebiten.KeyEscape},
			state.InputInspect: {ebiten.KeyTab},
		},
	}
}

func (i *inputController) JustPressed(ik state.InputKey) bool {
	eks, ok := i.mappings[ik]
	if !ok {
		return false
	}
	for _, ek := range eks {
		if inpututil.IsKeyJustPressed(ek) {
			return true
		}
	}
	return false
}

func (i *inputController) IsPressed(ik state.InputKey) bool {
	eks, ok := i.mappings[ik]
	if !ok {
		return false
	}
	for _, ek := range eks {
		if inpututil.KeyPressDuration(ek) > 0 {
			return true
		}
	}
	return false
}
