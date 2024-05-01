package main

import (
	"fmt"
	"image/color"
	"runtime"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type App struct {
	frameCount int32
}

// Global variables
var (
	windowSize = rl.Vector2{
		X: 1280,
		Y: 720,
	}
	Player = Entity{
		transform: rl.Rectangle{
			X:      100,
			Y:      windowSize.Y - (44 * 3),
			Width:  64 * 3,
			Height: 44 * 3,
		},
		speed:         20,
		isRunnging:    false,
		isFacingRight: true,
		isAttacking:   false,
		attackType:    0,
		scale:         3,
		sprite:        nil,
	}
	deltaTime           float32 = 0
	animationPhase      int     = 0
	betweenAttacksTimer float32 = TIME_FOR_ATTACK_2
)

// Constants
const (
	TIME_FOR_ATTACK_2 float32 = 0.1
	ATTACK_NONE       int32   = 0
	ATTACK_LIGHT      int32   = 1
	ATTACK_LIGHT_BACK int32   = 2
	ATTACK_HEAVY      int32   = 3
	VELOCITY
)

func main() {
	app := App{}
	rl.InitWindow(int32(windowSize.X), int32(windowSize.Y), "Demon Slayer")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	// Main loop for window / game
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Update deltaTime
		deltaTime = 1.0 / float32(rl.GetFPS())

		// Call respective methods for udpating and drawing
		app.Update()
		app.Draw()

		rl.EndDrawing()
	}
}

func (a *App) Update() {
	// Update app frame count
	a.frameCount++

	// Check player input
	if Player.isAttacking {
		if Player.attackType == ATTACK_LIGHT && animationPhase < len(playerSprites[ATTACK_ANIM])-1 {
			a.play_attack_animation()
		} else if Player.attackType == ATTACK_LIGHT_BACK && animationPhase < len(playerSprites[ATTACK2_ANIM])-1 {
			a.play_attack_animation()
		} else if Player.attackType == ATTACK_HEAVY && animationPhase < len(playerSprites[ATTACK3_ANIM])-1 {
			a.play_attack_animation()
		} else if betweenAttacksTimer > 0 {
			if rl.IsMouseButtonDown(rl.MouseButtonLeft) && Player.attackType != ATTACK_LIGHT_BACK {
				Player.attackType = ATTACK_LIGHT_BACK
				a.play_attack_animation()
			}
			betweenAttacksTimer -= 1 * deltaTime
		} else {
			Player.isAttacking = false
			a.reset_for_animation()
			Player.attackType = ATTACK_NONE
			betweenAttacksTimer = TIME_FOR_ATTACK_2
		}
	} else if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		Player.attackType = ATTACK_LIGHT
		a.frameCount = 0
		Player.isAttacking = true
	} else if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
		Player.attackType = ATTACK_HEAVY
		a.reset_for_animation()
		Player.isAttacking = true
	} else if rl.IsKeyDown(rl.KeyA) {
		a.run_left()
	} else if rl.IsKeyDown(rl.KeyD) {
		a.run_right()
	} else {
		Player.isRunnging = false
		Player.isAttacking = false
	}

	if !Player.isRunnging && !Player.isAttacking {
		a.idle_player()
	}

	// Debug print output
	m := runtime.MemStats(*new(runtime.MemStats))
	runtime.ReadMemStats(&m)
	fmt.Print("\033[H\033[2J")
	fmt.Printf("Memory usage: %d MB\n", m.Alloc/1024/1024)
	println("Player.sprite: ", Player.sprite)
	fmt.Printf("Player position: x {%f} y{%f} \n", Player.transform.X, Player.transform.Y)
	println("Player is using attack: ", Player.attackType)
	println("Player is attacking: ", Player.isAttacking)
	println("Player is running: ", Player.isRunnging)
	println("Animation phase: ", animationPhase)
	fmt.Printf("Between attacks timer: %f \n", betweenAttacksTimer)
}

func (a *App) Draw() {
	// Draw fps for debugging
	rl.DrawFPS(10, 10)

	// Draw the player sprite
	if Player.isFacingRight {
		rl.DrawTextureRec(*Player.sprite, rl.Rectangle{X: 0, Y: 0, Width: Player.transform.Width, Height: Player.transform.Height}, rl.Vector2{X: Player.transform.X, Y: Player.transform.Y}, color.RGBA{255, 255, 255, 255})
	} else {
		rl.DrawTextureRec(*Player.sprite, rl.Rectangle{X: 0, Y: 0, Width: -Player.transform.Width, Height: Player.transform.Height}, rl.Vector2{X: Player.transform.X, Y: Player.transform.Y}, color.RGBA{255, 255, 255, 255})
	}

}

// Player specific behaviours
func (a *App) idle_player() {
	Player.sprite = a.get_sprite(IDLE_ANIM, 8)
}
func (a *App) run_left() {
	Player.transform.X -= Player.speed * 25 * deltaTime
	Player.isRunnging = true
	Player.isFacingRight = false
	Player.sprite = a.get_sprite(RUN_ANIM, 100/int(Player.speed))
}
func (a *App) run_right() {
	Player.transform.X += Player.speed * 25 * deltaTime
	Player.isRunnging = true
	Player.isFacingRight = true
	Player.sprite = a.get_sprite(RUN_ANIM, 100/int(Player.speed))
}
func (a *App) play_attack_animation() {
	switch Player.attackType {
	case ATTACK_LIGHT:
		Player.sprite = a.get_sprite(ATTACK_ANIM, 5)
	case ATTACK_LIGHT_BACK:
		Player.sprite = a.get_sprite(ATTACK2_ANIM, 4)
	case ATTACK_HEAVY:
		Player.sprite = a.get_sprite(ATTACK3_ANIM, 4)
	}
}

// Utility functions
func get_texture(sprite *rl.Image) *rl.Texture2D {
	if Player.sprite != nil {
		rl.UnloadTexture(*Player.sprite)
	}
	result := rl.LoadTextureFromImage(sprite)
	result.Width = Player.transform.ToInt32().Width
	result.Height = Player.transform.ToInt32().Height
	return &result
}
func (a *App) reset_for_animation() {
	animationPhase = 0
	a.frameCount = 0
}
