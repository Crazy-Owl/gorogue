package game

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Hero struct {
	X          uint32
	Y          uint32
	HP         int32
	Tile       *Tile
	CurrentMap *TiledMap
}

func (h *Hero) DrawAt(renderer *sdl.Renderer, x uint32, y uint32) {
	h.Tile.DrawAt(renderer, x, y)
}

func (h *Hero) Update(dt uint32) {
	return
}

func (h *Hero) MoveTo(x, y uint32) {
	srcCell := h.CurrentMap.Get(h.X, h.Y)
	destCell := h.CurrentMap.Get(x, y)
	if destCell != nil && (destCell.Is("passable") || destCell.Is("ground_passable")) {
		srcCell.RemoveObject(h)
		destCell.AddObject(h)
		h.X = x
		h.Y = y
	}
}

func (h *Hero) Move(dx, dy int32) {
	h.MoveTo(uint32(int32(h.X)+dx), uint32(int32(h.Y)+dy))
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
