package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Sprite struct {
	x           float32
	y           float32
	width       float32
	height      float32
	frame_count int64
}

// Sprite atlas
// Load at runtime in setup
// Sprites
var sprite_atlas rl.Texture2D
var (
	// :s_player
	idle_anim = Sprite{
		x:           0,
		y:           0,
		width:       64,
		height:      48,
		frame_count: 6,
	}
	run_anim = Sprite{
		x:           0,
		y:           48,
		width:       64,
		height:      48,
		frame_count: 8,
	}
	attack1_anim = Sprite{
		x:           0,
		y:           96,
		width:       64,
		height:      48,
		frame_count: 5,
	}
	attack2_anim = Sprite{
		x:           0,
		y:           144,
		width:       64,
		height:      48,
		frame_count: 4,
	}
	attack3_anim = Sprite{
		x:           0,
		y:           96,
		width:       64,
		height:      48,
		frame_count: 6,
	}
	// :s_demon
	demon_idle = Sprite{
		x:           0,
		y:           192,
		width:       48,
		height:      48,
		frame_count: 6,
	}
	// :s_shop
	shop = Sprite{
		x:           880,
		y:           896,
		width:       120,
		height:      104,
		frame_count: 1,
	}
	// :s_ui
	panel = Sprite{
		x:           0,
		y:           880,
		width:       240,
		height:      120,
		frame_count: 1,
	}
)

// Assets
var (
	backgroundSprites = []*rl.Image{
		rl.LoadImage("assets/background/back.png"),
		rl.LoadImage("assets/background/back-back.png"),
		rl.LoadImage("assets/background/clouds.png"),
		rl.LoadImage("assets/background/back-mid.png"),
		rl.LoadImage("assets/background/mid.png"),
		rl.LoadImage("assets/background/mid-mid.png"),
		rl.LoadImage("assets/background/mid-fore.png"),
		rl.LoadImage("assets/background/fore.png"),
		rl.LoadImage("assets/background/fore-fore.png"),
	}
)

func get_anim_transform(entity *Entity, speed int) rl.Rectangle {
	entity.animation_phase = entity.frame_count / int64(speed) % entity.current_sprite.frame_count
	entity_width_factor := 1
	if !entity.isFacingRight {
		entity_width_factor = -1
	}
	return rl.Rectangle{
		X:      entity.current_sprite.x + (entity.current_sprite.width * float32(entity.animation_phase)),
		Y:      entity.current_sprite.y,
		Width:  entity.current_sprite.width * float32(entity_width_factor),
		Height: entity.current_sprite.height,
	}
}
