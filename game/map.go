package game

import (
	"github.com/veandco/go-sdl2/sdl"
	"math/rand"
)

type Cell struct {
	tile       *Tile
	objects    []Entity
	properties []string
}

type TiledMap struct {
	cells   [][]Cell
	objects []Entity
	w       int32
	h       int32
}

type ViewPort struct {
	Map   *TiledMap
	TileW int32
	TileH int32
	X     int32
	Y     int32
	W     int32
	H     int32
}

func CreateMap(w, h int32) TiledMap {
	objects := make([]Entity, 16)
	cells := make([][]Cell, h)
	tm := TiledMap{cells, objects, w, h}
	manager := GetResourceManager(GetRenderer())
	// tiles := manager.FilterTiles("terrain")
	passable_tiles := manager.FilterTiles("passable")
	wall_tiles := manager.FilterTiles("wall")

	for c_row := range cells {
		cells[c_row] = make([]Cell, w)
		var tile *Tile
		// random generation for now
		for c_col := range cells[c_row] {
			if rand.Intn(100) > 15 {
				tile = passable_tiles[rand.Intn(len(passable_tiles))]
			} else {
				tile = wall_tiles[rand.Intn(len(wall_tiles))]
			}
			cell := Cell{tile, make([]Entity, 16), make([]string, 16)}
			cell.properties = cell.tile.properties
			tm.SetAt(cell, int32(c_col), int32(c_row))
		}
	}
	// TODO: generation/loading
	return tm
}

func (tm *TiledMap) SetAt(cell Cell, x, y int32) {
	if !(y < tm.h && y >= 0 && x < tm.w && x >= 0) {
		return
	}
	tm.cells[y][x] = cell
}

func (tm *TiledMap) Get(x, y int32) *Cell {
	if !(y < tm.h && y >= 0 && x < tm.w && x >= 0) {
		return nil
	}
	return &tm.cells[y][x]
}

func (tm *TiledMap) AddObject(o Entity, x, y int32) {
	n := make([]Entity, 0)
	n = append(n, o)
	for _, x := range tm.objects {
		if x != o {
			n = append(n, x)
		}
	}
	tm.objects = n
	cell := tm.Get(x, y)
	if cell != nil {
		cell.AddObject(o)
	}
}

func (vp *ViewPort) Render(renderer *sdl.Renderer) {
	m := vp.Map
	for x := int32(0); x < vp.W; x++ {
		for y := int32(0); y < vp.H; y++ {
			cell := m.Get(vp.X+x, vp.Y+y)
			if cell != nil {
				if cell.tile != nil {
					cell.tile.DrawAt(renderer, x*vp.TileW, y*vp.TileH)
				}
				for _, o := range cell.objects {
					if o != nil {
						o.DrawAt(renderer, x*vp.TileW, y*vp.TileH)
					}
				}
			}
		}
	}
}

func (vp *ViewPort) PointAt(x, y int32) {
	if x < 0 {
		x = 0
	}
	if x > vp.Map.w-vp.W {
		x = vp.Map.w - vp.W
	}
	if y < 0 {
		y = 0
	}
	if y > vp.Map.h-vp.H {
		y = vp.Map.h - vp.H
	}
	vp.X = x
	vp.Y = y
}

func (vp *ViewPort) Move(dx, dy int32) {
	nx := vp.X + dx
	ny := vp.Y + dy
	vp.PointAt(nx, ny)
}

func (vp *ViewPort) CenterAt(x, y int32) {
	vp.PointAt(x-int32(vp.W/2), y-int32(vp.H/2))
}

func (c *Cell) Is(prop string) (is bool) {
	is = false
	for _, x := range c.properties {
		if x == prop {
			return true
		}
	}
	return
}

func (c *Cell) Unset(prop string) {
	props := make([]string, 0)
	for _, x := range c.properties {
		if x != prop {
			props = append(props, x)
		}
	}
	c.properties = props
}

func (c *Cell) Set(prop string) {
	props := make([]string, 0)
	props = append(props, prop)
	for _, x := range c.properties {
		if x != prop {
			props = append(props, x)
		}
	}
	c.properties = props
}

func (c *Cell) AddObject(e Entity) {
	ents := make([]Entity, 0)
	ents = append(ents, e)
	for _, x := range c.objects {
		if x != e {
			ents = append(ents, x)
		}
	}
	c.objects = ents
}

func (c *Cell) RemoveObject(e Entity) {
	ents := make([]Entity, 0)
	for _, x := range c.objects {
		if x != e {
			ents = append(ents, x)
		}
	}
	c.objects = ents
}
