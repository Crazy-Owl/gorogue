package game

import (
	"github.com/veandco/go-sdl2/sdl"
)

// TODO: generalized State interface should have access to objects etc,
// so we can use same commands signature in every state
var commands map[string]func(*GameState) = make(map[string]func(*GameState))

type Bindings map[sdl.Keycode]string

var GameBindings Bindings = make(Bindings)

func RegisterCommand(key string, fun func(*GameState)) {
	commands[key] = fun
}

func GetCommand(key string) func(*GameState) {
	return commands[key]
}

func (bindings Bindings) RegisterBinding(key sdl.Keycode, command string) {
	bindings[key] = command
}

func (bindings Bindings) GetBinding(key sdl.Keycode) string {
	return bindings[key]
}

// Game commands
func gsMove(x, y int32) func(*GameState) {
	return func(gs *GameState) { gs.hero.Move(x, y) }
}

func init() {
	RegisterCommand("move_nw", gsMove(-1, -1))
	RegisterCommand("move_n", gsMove(0, -1))
	RegisterCommand("move_ne", gsMove(1, -1))
	RegisterCommand("move_e", gsMove(1, 0))
	RegisterCommand("move_se", gsMove(1, 1))
	RegisterCommand("move_s", gsMove(0, 1))
	RegisterCommand("move_sw", gsMove(-1, 1))
	RegisterCommand("move_w", gsMove(-1, 0))
	RegisterCommand("stop_running", func(gs *GameState) { gs.running = false })

	GameBindings.RegisterBinding(sdl.K_7, "move_nw")
	GameBindings.RegisterBinding(sdl.K_8, "move_n")
	GameBindings.RegisterBinding(sdl.K_9, "move_ne")
	GameBindings.RegisterBinding(sdl.K_6, "move_e")
	GameBindings.RegisterBinding(sdl.K_3, "move_se")
	GameBindings.RegisterBinding(sdl.K_2, "move_s")
	GameBindings.RegisterBinding(sdl.K_1, "move_sw")
	GameBindings.RegisterBinding(sdl.K_4, "move_w")

	GameBindings.RegisterBinding(sdl.K_UP, "move_n")
	GameBindings.RegisterBinding(sdl.K_RIGHT, "move_e")
	GameBindings.RegisterBinding(sdl.K_DOWN, "move_s")
	GameBindings.RegisterBinding(sdl.K_LEFT, "move_w")

	GameBindings.RegisterBinding(sdl.K_ESCAPE, "stop_running")
}
