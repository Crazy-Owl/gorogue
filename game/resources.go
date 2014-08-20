package game

import (
	"encoding/json"
	"errors"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_image"
	"github.com/veandco/go-sdl2/sdl_ttf"
	"io/ioutil"
	"log"
	"strings"
)

// Game structs

type ResourceManager struct {
	textures map[string]*sdl.Texture
	fonts    map[string]*ttf.Font
	tiles    map[string]*Tile
	renderer *sdl.Renderer
	prefix   string
}

type Tile struct {
	texture    *sdl.Texture
	rect       *sdl.Rect
	properties []string
}

// JSON-parsing structs

type TileMap struct {
	Filename  string
	TileSpecs []TileSpec
	W         int32
	H         int32
}

type TileSpec struct {
	X          int32
	Y          int32
	Name       string
	Properties []string
}

var manager *ResourceManager

func GetResourceManager(renderer *sdl.Renderer) *ResourceManager {
	if manager != nil {
		return manager
	}
	manager = &ResourceManager{make(map[string]*sdl.Texture), make(map[string]*ttf.Font), make(map[string]*Tile), renderer, "data/"}
	manager.Init()
	return manager
}

func (manager *ResourceManager) Init() {
	files, err := ioutil.ReadDir(manager.prefix)
	if err != nil {
		log.Fatal("Can't load data dir")
	}
	for _, fileInfo := range files {
		if !fileInfo.IsDir() {
			if strings.HasSuffix(fileInfo.Name(), ".png") {
				manager.LoadTexture(fileInfo.Name())
			}
			if strings.HasSuffix(fileInfo.Name(), ".ttf") {
				manager.LoadFont(fileInfo.Name())
			}
		}
	}
	for _, fileInfo := range files {
		if !fileInfo.IsDir() {
			if strings.HasSuffix(fileInfo.Name(), ".tiles") {
				manager.LoadTileMap(fileInfo.Name())
			}
		}
	}
}

func (manager *ResourceManager) LoadTexture(filepath string) {
	texture := img.LoadTexture(manager.renderer, manager.prefix+filepath)
	manager.textures[filepath] = texture
}

func (manager *ResourceManager) LoadFont(filepath string) {
	font, err := ttf.OpenFont(manager.prefix+filepath, 14)
	if err != nil {
		log.Println(err)
		log.Fatal("Can't load font")
	}
	manager.fonts[filepath] = font
}

func (manager *ResourceManager) LoadTileMap(filepath string) {
	var tm TileMap
	data, err := ioutil.ReadFile(manager.prefix + filepath)
	if err != nil {
		log.Fatal("Can't open tilemap file")
	}
	err = json.Unmarshal(data, &tm)
	if err != nil {
		log.Println(err)
		log.Fatal("Can't load tile map")
	}
	for _, tilespec := range tm.TileSpecs {
		tex, err := manager.GetTexture(tm.Filename)
		if err != nil {
			log.Fatal("Can't load texture for tile")
		}
		tile := Tile{tex, &sdl.Rect{tilespec.X, tilespec.Y, tm.W, tm.H}, tilespec.Properties}
		manager.tiles[tilespec.Name] = &tile
	}
}

func (manager *ResourceManager) GetTexture(key string) (tex *sdl.Texture, err error) {
	var ok bool
	tex, ok = manager.textures[key]
	if !ok {
		err = errors.New("Texture not found")
	}
	return
}

func (manager *ResourceManager) GetTile(key string) (tile *Tile, err error) {
	tile, ok := manager.tiles[key]
	if !ok {
		err = errors.New("Tile not found")
	}
	return
}

func (manager *ResourceManager) GetFont(key string) (font *ttf.Font, err error) {
	var ok bool
	font, ok = manager.fonts[key]
	if !ok {
		err = errors.New("Font not found")
	}
	return
}

func (manager *ResourceManager) Print(font *ttf.Font, renderer *sdl.Renderer, color sdl.Color, message string, x, y int32) {
	surface := font.RenderText_Solid(message, color)
	string_texture := renderer.CreateTextureFromSurface(surface)
	renderer.Copy(string_texture, &sdl.Rect{0, 0, surface.W, surface.H}, &sdl.Rect{x, y, surface.W, surface.H})
}

func (manager *ResourceManager) MakePrinter(font *ttf.Font, renderer *sdl.Renderer, color sdl.Color) func(string, int32, int32) {
	return func(message string, x, y int32) {
		manager.Print(font, renderer, color, message, x, y)
	}
}

func (manager *ResourceManager) GetTiles() []*Tile {
	var tiles = make([]*Tile, 0)
	for _, v := range manager.tiles {
		tiles = append(tiles, v)
	}
	return tiles
}

func (manager *ResourceManager) FilterTiles(prop string) []*Tile {
	var tiles = make([]*Tile, 0)
	for _, v := range manager.tiles {
		if v.Is(prop) {
			tiles = append(tiles, v)
		}
	}
	return tiles
}

func (manager *ResourceManager) GetTileOrNil(key string) (tile *Tile) {
	tile, _ = manager.GetTile(key)
	return
}

func (tile *Tile) DrawAt(renderer *sdl.Renderer, x, y uint32) {
	dest := sdl.Rect{int32(x), int32(y), tile.rect.W, tile.rect.H}
	renderer.Copy(tile.texture, tile.rect, &dest)
}

func (tile *Tile) Is(prop string) (is bool) {
	is = false
	for _, x := range tile.properties {
		if x == prop {
			return true
		}
	}
	return
}
