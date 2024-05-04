package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Sprite struct {
	x           float32
	y           float32
	width       float32
	height      float32
	frame_count int
}

// Sprite atlas
// Load at runtime in setup
// Sprites
var sprite_atlas rl.Texture2D
var (
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
		height:      32,
		frame_count: 8,
	}
	attack1_anim = Sprite{
		x:           0,
		y:           80,
		width:       64,
		height:      48,
		frame_count: 6,
	}
	attack2_anim = Sprite{
		x:           0,
		y:           128,
		width:       64,
		height:      48,
		frame_count: 4,
	}
	attack3_anim = Sprite{
		x:           0,
		y:           80,
		width:       64,
		height:      48,
		frame_count: 6,
	}
)

// Assets
var (
	sprite_table = [][][]*rl.Image{
		{
			// Player sprites
			{
				rl.LoadImage("assets/character/idle_1.png"),
				rl.LoadImage("assets/character/idle_2.png"),
				rl.LoadImage("assets/character/idle_3.png"),
				rl.LoadImage("assets/character/idle_4.png"),
				rl.LoadImage("assets/character/idle_5.png"),
				rl.LoadImage("assets/character/idle_6.png"),
			},
			{
				rl.LoadImage("assets/character/run_0.png"),
				rl.LoadImage("assets/character/run_1.png"),
				rl.LoadImage("assets/character/run_2.png"),
				rl.LoadImage("assets/character/run_3.png"),
				rl.LoadImage("assets/character/run_4.png"),
				rl.LoadImage("assets/character/run_5.png"),
				rl.LoadImage("assets/character/run_6.png"),
				rl.LoadImage("assets/character/run_7.png"),
			},
			{
				rl.LoadImage("assets/character/attack_0.png"),
				rl.LoadImage("assets/character/attack_1.png"),
				rl.LoadImage("assets/character/attack_2.png"),
				rl.LoadImage("assets/character/attack_3.png"),
				rl.LoadImage("assets/character/attack_4.png"),
				rl.LoadImage("assets/character/attack_5.png"),
				rl.LoadImage("assets/character/attack_6.png"),
			},
			{
				rl.LoadImage("assets/character/attack2_0.png"),
				rl.LoadImage("assets/character/attack2_1.png"),
				rl.LoadImage("assets/character/attack2_2.png"),
				rl.LoadImage("assets/character/attack2_3.png"),
				rl.LoadImage("assets/character/attack2_4.png"),
			},
			{
				rl.LoadImage("assets/character/attack3_0.png"),
				rl.LoadImage("assets/character/attack3_1.png"),
				rl.LoadImage("assets/character/attack3_2.png"),
				rl.LoadImage("assets/character/attack3_3.png"),
				rl.LoadImage("assets/character/attack3_4.png"),
				rl.LoadImage("assets/character/attack3_5.png"),
				// rl.LoadImage("assets/character/attack3_6.png"),
			},
		},
		{
			// Red axe demon sprites
			{
				rl.LoadImage("assets/demon_axe_red/ready_0.png"),
				rl.LoadImage("assets/demon_axe_red/ready_1.png"),
				rl.LoadImage("assets/demon_axe_red/ready_2.png"),
				rl.LoadImage("assets/demon_axe_red/ready_3.png"),
				rl.LoadImage("assets/demon_axe_red/ready_4.png"),
				rl.LoadImage("assets/demon_axe_red/ready_5.png"),
			},
		},
	}
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

func get_anim_sprite(entity *Entity, speed int) *rl.Texture2D {
	entity.animation_phase = int(frameCount) / speed % len(sprite_table[entity.sprite_set][entity.state])
	return get_texture(sprite_table[entity.sprite_set][entity.state][entity.animation_phase])
}

func get_anim_transform(entity *Entity, speed int) rl.Rectangle {
	entity.animation_phase = int(frameCount) / speed % entity.current_sprite.frame_count
	entity_width_factor := 1
	if !entity.isFacingRight {
		entity_width_factor = -1
	}
	return rl.Rectangle{
		X:      entity.current_sprite.x + (entity.current_sprite.width * float32(Player.animation_phase)),
		Y:      entity.current_sprite.y,
		Width:  entity.current_sprite.width * float32(entity_width_factor),
		Height: entity.current_sprite.height,
	}
}
