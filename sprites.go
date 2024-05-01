package main

import rl "github.com/gen2brain/raylib-go/raylib"

// Assets
var (
	playerSprites = [][]*rl.Image{
		{
			rl.LoadImage("assets/warrior/idle/Warrior_Idle_0.png"),
			rl.LoadImage("assets/warrior/idle/Warrior_Idle_1.png"),
			rl.LoadImage("assets/warrior/idle/Warrior_Idle_2.png"),
			rl.LoadImage("assets/warrior/idle/Warrior_Idle_3.png"),
			rl.LoadImage("assets/warrior/idle/Warrior_Idle_4.png"),
			rl.LoadImage("assets/warrior/idle/Warrior_Idle_5.png"),
		},
		{
			rl.LoadImage("assets/warrior/Run/Warrior_Run_0.png"),
			rl.LoadImage("assets/warrior/Run/Warrior_Run_1.png"),
			rl.LoadImage("assets/warrior/Run/Warrior_Run_2.png"),
			rl.LoadImage("assets/warrior/Run/Warrior_Run_3.png"),
			rl.LoadImage("assets/warrior/Run/Warrior_Run_4.png"),
			rl.LoadImage("assets/warrior/Run/Warrior_Run_5.png"),
			rl.LoadImage("assets/warrior/Run/Warrior_Run_6.png"),
			rl.LoadImage("assets/warrior/Run/Warrior_Run_7.png"),
		},
		{
			rl.LoadImage("assets/warrior/Attack/Warrior_Attack_0.png"),
			rl.LoadImage("assets/warrior/Attack/Warrior_Attack_1.png"),
			rl.LoadImage("assets/warrior/Attack/Warrior_Attack_2.png"),
			rl.LoadImage("assets/warrior/Attack/Warrior_Attack_3.png"),
			rl.LoadImage("assets/warrior/Attack/Warrior_Attack_4.png"),
			rl.LoadImage("assets/warrior/Attack/Warrior_Attack_5.png"),
			rl.LoadImage("assets/warrior/Attack/Warrior_Attack_6.png"),
			rl.LoadImage("assets/warrior/Attack/Warrior_Attack_7.png"),
		},
		{
			rl.LoadImage("assets/warrior/Attack/Warrior_Attack2_0.png"),
			rl.LoadImage("assets/warrior/Attack/Warrior_Attack2_1.png"),
			rl.LoadImage("assets/warrior/Attack/Warrior_Attack2_2.png"),
			rl.LoadImage("assets/warrior/Attack/Warrior_Attack2_3.png"),
			rl.LoadImage("assets/warrior/Attack/Warrior_Attack2_4.png"),
			rl.LoadImage("assets/warrior/Attack/Warrior_Attack2_5.png"),
			rl.LoadImage("assets/warrior/Attack/Warrior_Attack2_6.png"),
		},
		{
			rl.LoadImage("assets/warrior/Dash Attack/Warrior_Dash-Attack_0.png"),
			rl.LoadImage("assets/warrior/Dash Attack/Warrior_Dash-Attack_1.png"),
			rl.LoadImage("assets/warrior/Dash Attack/Warrior_Dash-Attack_2.png"),
			rl.LoadImage("assets/warrior/Dash Attack/Warrior_Dash-Attack_3.png"),
			rl.LoadImage("assets/warrior/Dash Attack/Warrior_Dash-Attack_4.png"),
			rl.LoadImage("assets/warrior/Dash Attack/Warrior_Dash-Attack_5.png"),
			rl.LoadImage("assets/warrior/Dash Attack/Warrior_Dash-Attack_6.png"),
			rl.LoadImage("assets/warrior/Dash Attack/Warrior_Dash-Attack_7.png"),
			rl.LoadImage("assets/warrior/Dash Attack/Warrior_Dash-Attack_8.png"),
			rl.LoadImage("assets/warrior/Dash Attack/Warrior_Dash-Attack_9.png"),
		},
	}
)

const (
	IDLE_ANIM    = 0
	RUN_ANIM     = 1
	ATTACK_ANIM  = 2
	ATTACK2_ANIM = 3
	ATTACK3_ANIM = 4
)

func (a *App) get_sprite(anim int, speed int) *rl.Texture2D {
	animationPhase = int(a.frameCount) / speed % len(playerSprites[anim])
	return get_texture(playerSprites[anim][animationPhase])
}
