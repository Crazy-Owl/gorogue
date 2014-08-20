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
	w       uint32
	h       uint32
}

type ViewPort struct {
	Map   *TiledMap
	TileW uint32
	TileH uint32
	X     uint32
	Y     uint32
	W     uint32
	H     uint32
}

type MapGenRoom struct {
	X uint32
	Y uint32
	W uint32
	H uint32
}

func CreateMap(w, h uint32) TiledMap {
	objects := make([]Entity, 16)
	cells := make([][]Cell, h)
	tm := TiledMap{cells, objects, w, h}
	manager := GetResourceManager(GetRenderer())

	empty_tile := manager.GetTileOrNil("gray_nothing")

	for c_row := range cells {
		cells[c_row] = make([]Cell, w)
		for c_col := range cells[c_row] {
			cell := Cell{empty_tile, make([]Entity, 1), make([]string, 16)}
			cell.properties = cell.tile.properties
			tm.SetAt(cell, uint32(c_col), uint32(c_row))
		}
	}
	// TODO: generation/loading
	return tm
}

func (tm *TiledMap) GenerateRoom(maxW, maxH uint32) MapGenRoom {
	roomX := uint32(rand.Intn(int(tm.w - maxW - 2)))
	roomY := uint32(rand.Intn(int(tm.h - maxH - 2)))
	roomW := uint32(rand.Intn(int(maxW)) + 2)
	roomH := uint32(rand.Intn(int(maxH)) + 2)
	return MapGenRoom{roomX + 1, roomY + 1, roomW, roomH}
}

func (tm *TiledMap) GenerateMap() {
	manager := GetResourceManager(GetRenderer())
	// tiles := manager.FilterTiles("terrain")
	passable_tiles := manager.FilterTiles("passable")
	wall_tiles := manager.FilterTiles("wall")

	room := tm.GenerateRoom(64, 64)
	selected_floor := passable_tiles[rand.Intn(len(passable_tiles))]
	selected_wall := wall_tiles[rand.Intn(len(wall_tiles))]

	// set up room floor
	for x := room.X; x < room.X+room.W+1; x++ {
		for y := room.Y; y < room.Y+room.H+1; y++ {
			cell := Cell{selected_floor, make([]Entity, 16), make([]string, 16)}
			cell.properties = cell.tile.properties
			tm.SetAt(cell, x, y)
		}
	}

	// set up walls
	for x := room.X - 1; x < room.X+room.W+2; x++ {
		cell := Cell{selected_wall, make([]Entity, 1), make([]string, 16)}
		cell.properties = cell.tile.properties
		tm.SetAt(cell, x, room.Y-1)
		cell = Cell{selected_wall, make([]Entity, 1), make([]string, 16)}
		cell.properties = cell.tile.properties
		tm.SetAt(cell, x, room.Y+room.H+1)
	}
	for y := room.Y - 1; y < room.Y+room.H+2; y++ {
		cell := Cell{selected_wall, make([]Entity, 1), make([]string, 16)}
		cell.properties = cell.tile.properties
		tm.SetAt(cell, room.X-1, y)
		cell = Cell{selected_wall, make([]Entity, 1), make([]string, 16)}
		cell.properties = cell.tile.properties
		tm.SetAt(cell, room.X+room.W+1, y)
	}
}

func (tm *TiledMap) SetAt(cell Cell, x, y uint32) {
	if !(y < tm.h && y >= 0 && x < tm.w && x >= 0) {
		return
	}
	tm.cells[y][x] = cell
}

func (tm *TiledMap) Get(x, y uint32) *Cell {
	if !(y < tm.h && y >= 0 && x < tm.w && x >= 0) {
		return nil
	}
	return &tm.cells[y][x]
}

func (tm *TiledMap) AddObject(o Entity, x, y uint32) {
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
	for x := uint32(0); x < vp.W; x++ {
		for y := uint32(0); y < vp.H; y++ {
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

func (vp *ViewPort) PointAt(x, y uint32) {
	if x > vp.Map.w-vp.W {
		x = vp.Map.w - vp.W
	}
	if y > vp.Map.h-vp.H {
		y = vp.Map.h - vp.H
	}
	vp.X = x
	vp.Y = y
}

func (vp *ViewPort) Move(dx, dy int32) {
	nx := uint32(int32(vp.X) + dx)
	ny := uint32(int32(vp.Y) + dy)
	vp.PointAt(nx, ny)
}

func (vp *ViewPort) CenterAt(x, y uint32) {
	vp.PointAt(x-uint32(vp.W/2), y-uint32(vp.H/2))
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
