package main

import rl "github.com/gen2brain/raylib-go/raylib"

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
