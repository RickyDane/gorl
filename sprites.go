package main

import rl "github.com/gen2brain/raylib-go/raylib"

// Assets
var (
	playerSprites = [][]*rl.Image{
		{
			rl.LoadImage("assets/character/idle_0.png"),
			rl.LoadImage("assets/character/idle_1.png"),
			rl.LoadImage("assets/character/idle_2.png"),
			rl.LoadImage("assets/character/idle_3.png"),
			rl.LoadImage("assets/character/idle_4.png"),
			rl.LoadImage("assets/character/idle_5.png"),
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
			rl.LoadImage("assets/character/attack2_5.png"),
			rl.LoadImage("assets/character/attack2_6.png"),
		},
		{
			rl.LoadImage("assets/character/attack3_0.png"),
			rl.LoadImage("assets/character/attack3_1.png"),
			rl.LoadImage("assets/character/attack3_2.png"),
			rl.LoadImage("assets/character/attack3_3.png"),
			rl.LoadImage("assets/character/attack3_4.png"),
			rl.LoadImage("assets/character/attack3_5.png"),
			rl.LoadImage("assets/character/attack3_6.png"),
			rl.LoadImage("assets/character/attack3_8.png"),
			rl.LoadImage("assets/character/attack3_9.png"),
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

const (
	IDLE_ANIM    = 0
	RUN_ANIM     = 1
	ATTACK_ANIM  = 2
	ATTACK2_ANIM = 3
	ATTACK3_ANIM = 4
)

func (a *App) get_anim_sprite(anim int, speed int) *rl.Texture2D {
	animationPhase = int(a.frameCount) / speed % len(playerSprites[anim])
	return get_texture(playerSprites[anim][animationPhase])
}
