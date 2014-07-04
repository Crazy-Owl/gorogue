package game

import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
)

var renderer *sdl.Renderer

func InitRendering() {
	// create window and rendering mechanisms here
}

func GetRenderer() *sdl.Renderer {
	if renderer == nil {
		log.Fatal("Game was not properly initialized, can't get current renderer")
	}
	return renderer
}

func SetRenderer(r *sdl.Renderer) {
	renderer = r
}
