package main

import (
	// "fmt"
	"github.com/Crazy-Owl/gorogue/game"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
	// "log"
)

const (
	winW = 1200
	winH = 900
)

func main() {
	var event sdl.Event
	var ts, ts2, dt uint32

	sdl.Init(sdl.INIT_EVERYTHING)
	img.Init(img.INIT_PNG)
	ttf.Init()

	window := sdl.CreateWindow("GoRogue", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winW, winH, sdl.WINDOW_SHOWN)
	defer window.Destroy()

	renderer := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer renderer.Destroy()

	game.SetRenderer(renderer)

	manager := game.GetResourceManager(renderer)
	m := game.CreateMap(256, 256)
	vp := game.ViewPort{&m, 32, 32, 0, 0, int32(winW / 32), int32(winH / 32)}

	hero := game.Hero{5, 5, 10, manager.GetTileOrNil("ice_fiend"), &m}
	m.AddObject(&hero, hero.X, hero.Y)

	gs := game.MakeGameState(&m, &vp, &hero)

	// main loop
	running := true
	ts = sdl.GetTicks()
	for running {
		renderer.SetDrawColor(0, 0, 0, 0)
		renderer.Clear()
		ts2 = sdl.GetTicks()
		dt = ts2 - ts
		ts = ts2
		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyDownEvent:
				ev := event.(*sdl.KeyDownEvent)
				ks := ev.Keysym
				switch ks.Sym {
				// case sdl.K_LEFT:
				// 	vp.Move(-2, 0)
				// case sdl.K_RIGHT:
				// 	vp.Move(2, 0)
				// case sdl.K_UP:
				// 	vp.Move(0, -2)
				// case sdl.K_DOWN:
				// 	vp.Move(0, 2)
				case sdl.K_ESCAPE:
					running = false
				default:
					gs.KeyPressed(ev)
				}
			}
		}
		renderer.SetDrawColor(255, 0, 0, 0)
		gs.Update(dt)
		gs.Render(renderer)
		renderer.Present()
	}
}
