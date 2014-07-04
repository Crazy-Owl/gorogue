package game

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Hero struct {
	X          int32
	Y          int32
	HP         int32
	Tile       *Tile
	CurrentMap *TiledMap
}

func (h *Hero) DrawAt(renderer *sdl.Renderer, x int32, y int32) {
	h.Tile.DrawAt(renderer, x, y)
}

func (h *Hero) Update(dt uint32) {
	return
}

func (h *Hero) MoveTo(x, y int32) {
	srcCell := h.CurrentMap.Get(h.X, h.Y)
	destCell := h.CurrentMap.Get(x, y)
	if destCell != nil {
		srcCell.RemoveObject(h)
		destCell.AddObject(h)
		h.X = x
		h.Y = y
	}
}

func (h *Hero) Move(dx, dy int32) {
	h.MoveTo(h.X+dx, h.Y+dy)
}

func (h *Hero) Interact(e Entity) {
	return
}

func (h *Hero) Texture() *sdl.Texture {
	return h.Tile.texture
}

func (h *Hero) Rect() *sdl.Rect {
	return h.Tile.rect
}
