package main

import (
	// "fmt"
	"github.com/Crazy-Owl/gorogue/game"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"math/rand"
	// "log"
)

const (
	winW = 1200
	winH = 900
	mapW = 256
	mapH = 256
)

func main() {
	var event sdl.Event
	var ts, ts2, dt uint32
	var hX, hY uint32

	sdl.Init(sdl.INIT_EVERYTHING)
	img.Init(img.INIT_PNG)
	ttf.Init()

	window := sdl.CreateWindow("GoRogue", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, winW, winH, sdl.WINDOW_SHOWN)
	defer window.Destroy()

	renderer := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	defer renderer.Destroy()

	game.SetRenderer(renderer)

	manager := game.GetResourceManager(renderer)
	m := game.CreateMap(mapW, mapH)
	m.GenerateMap()

	vp := game.ViewPort{&m, 32, 32, 0, 0, winW / 32, winH / 32}

	heroCanBePlaced := false
	for !heroCanBePlaced {
		hX = rand.Uint32() % mapW
		hY = rand.Uint32() % mapH
		if m.Get(hX, hY).Is("passable") {
			heroCanBePlaced = true
		}
	}

	hero := game.Hero{hX, hY, 10, manager.GetTileOrNil("fire_fiend"), &m}

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
				gs.KeyPressed(ev)
			}
		}
		renderer.SetDrawColor(255, 0, 0, 0)
		running = gs.Update(dt)
		gs.Render(renderer)
		renderer.Present()
	}
}
