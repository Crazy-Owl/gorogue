package game

import (
	"github.com/veandco/go-sdl2/sdl"
)

type State interface {
	Update(int32)
	Render(*sdl.Renderer)
	KeyPressed(*sdl.KeyDownEvent)
}

type Entity interface {
	DrawAt(*sdl.Renderer, uint32, uint32)
	Update(uint32)
	Move(int32, int32)
	MoveTo(uint32, uint32)
	Interact(Entity)
	Texture() *sdl.Texture
	Rect() *sdl.Rect
}

type MenuState struct {
	items   map[string]func()
	printer func(string, int32, int32)
}

type GameState struct {
	objects    []Entity
	currentMap *TiledMap
	viewPort   *ViewPort
	hero       Entity
	running    bool
}

func MakeGameState(m *TiledMap, vp *ViewPort, h *Hero) GameState {
	return GameState{make([]Entity, 0), m, vp, h, true}
}

func (gs *GameState) Update(dt uint32) (running bool) {
	if !gs.running {
		running = false
		return
	}
	hero := gs.hero.(*Hero)
	gs.viewPort.CenterAt(hero.X, hero.Y)
	running = true
	return
}

func (gs *GameState) Render(renderer *sdl.Renderer) {
	gs.viewPort.Render(renderer)
}

func (gs *GameState) KeyPressed(ev *sdl.KeyDownEvent) {
	ks := ev.Keysym
	switch ks.Sym {
	case sdl.K_LEFT:
		gs.hero.Move(-1, 0)
	case sdl.K_RIGHT:
		gs.hero.Move(1, 0)
	case sdl.K_UP:
		gs.hero.Move(0, -1)
	case sdl.K_DOWN:
		gs.hero.Move(0, 1)
	case sdl.K_1:
		gs.hero.Move(-1, 1)
	case sdl.K_2:
		gs.hero.Move(0, 1)
	case sdl.K_3:
		gs.hero.Move(1, 1)
	case sdl.K_4:
		gs.hero.Move(-1, 0)
	case sdl.K_6:
		gs.hero.Move(1, 0)
	case sdl.K_7:
		gs.hero.Move(-1, -1)
	case sdl.K_8:
		gs.hero.Move(0, -1)
	case sdl.K_9:
		gs.hero.Move(1, -1)
	case sdl.K_ESCAPE:
		gs.running = false
	}
}
