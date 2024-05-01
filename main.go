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
			X:      windowSize.X / 4,
			Y:      windowSize.Y - (44 * 5),
			Width:  64 * 4,
			Height: 44 * 4,
		},
		world_position: rl.Vector2{X: windowSize.X / 4, Y: windowSize.Y - (44 * 5)},
		speed:          20,
		isRunnging:     false,
		isFacingRight:  true,
		isAttacking:    false,
		attackType:     0,
		scale:          5,
		sprite:         nil,
	}
	deltaTime           float32 = 0
	animationPhase      int     = 0
	betweenAttacksTimer float32 = TIME_FOR_ATTACK_2
)

// Background layers
var (
	bg_layer_0          *rl.Texture2D
	bg_layer_1          *rl.Texture2D
	bg_layer_2          *rl.Texture2D
	bg_layer_3          *rl.Texture2D
	bg_layer_4          *rl.Texture2D
	bg_layer_5          *rl.Texture2D
	bg_layer_6          *rl.Texture2D
	bg_layer_7          *rl.Texture2D
	bg_layer_8          *rl.Texture2D
	bg_layer_9          *rl.Texture2D
	bg_layer_10         *rl.Texture2D
	scrolling_backback  = 0
	scrolling_back      = 0
	scrolling_clouds    = 0
	scrolling_backmid   = 0
	scrolling_mid       = 0
	scrolling_midmid    = 0
	scrolling_midmidmid = 0
	scrolling_midfore   = 0
	scrolling_fore      = 0
	scrolling_forefore  = 0
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

	app.Setup()

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
		if app.frameCount%500 == 0 {
			runtime.GC()
		}
	}
}

func (a *App) Setup() {
	// Set target FPS
	rl.SetTargetFPS(60)

	// Load sprites for background
	bg_layer_0 = get_texture(backgroundSprites[0])
	bg_layer_1 = get_texture(backgroundSprites[1])
	bg_layer_2 = get_texture(backgroundSprites[2])
	bg_layer_3 = get_texture(backgroundSprites[3])
	bg_layer_4 = get_texture(backgroundSprites[4])
	bg_layer_5 = get_texture(backgroundSprites[5])
	bg_layer_6 = get_texture(backgroundSprites[6])
	bg_layer_7 = get_texture(backgroundSprites[7])
	bg_layer_8 = get_texture(backgroundSprites[8])
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
				a.reset_for_animation()
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
	fmt.Print("\033[H\033[2J")
	PrintMemUsage()
	fmt.Printf("Player position: x {%f} y{%f} \n", Player.transform.X, Player.transform.Y)
	fmt.Printf("Player world position: x: {%f} y: {%f}\n", Player.world_position.X, Player.world_position.Y)
	println("Player is using attack: ", Player.attackType)
	println("Player is attacking: ", Player.isAttacking)
	println("Player is running: ", Player.isRunnging)
	println("Animation phase: ", animationPhase)
	fmt.Printf("Between attacks timer: %f \n", betweenAttacksTimer)
}

func (a *App) Draw() {
	draw_background()

	// Draw the player sprite
	if Player.isFacingRight {
		rl.DrawTextureRec(*Player.sprite, rl.Rectangle{X: 0, Y: 0, Width: Player.transform.Width, Height: Player.transform.Height}, rl.Vector2{X: Player.transform.X, Y: Player.transform.Y}, color.RGBA{255, 255, 255, 255})
	} else {
		rl.DrawTextureRec(*Player.sprite, rl.Rectangle{X: 0, Y: 0, Width: -Player.transform.Width, Height: Player.transform.Height}, rl.Vector2{X: Player.transform.X, Y: Player.transform.Y}, color.RGBA{255, 255, 255, 255})
	}

	scrolling_clouds += 1

	// Draw foreground over player
	rl.DrawTextureRec(*bg_layer_8, rl.Rectangle{X: float32(scrolling_forefore), Y: 0, Width: float32(bg_layer_6.Width), Height: float32(bg_layer_6.Height)}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})

	// Draw fps for debugging
	rl.DrawFPS(10, 10)
}

// Player specific behaviours
func (a *App) idle_player() {
	sprite := a.get_anim_sprite(IDLE_ANIM, 8)
	sprite.Width = Player.transform.ToInt32().Width
	sprite.Height = Player.transform.ToInt32().Height
	Player.sprite = sprite
}
func (a *App) run_left() {
	Player.world_position.X -= Player.speed
	Player.isRunnging = true
	Player.isFacingRight = false
	scroll_background(false)
	sprite := a.get_anim_sprite(RUN_ANIM, 100/int(Player.speed))
	sprite.Width = Player.transform.ToInt32().Width
	sprite.Height = Player.transform.ToInt32().Height
	Player.sprite = sprite
}
func (a *App) run_right() {
	Player.world_position.X += Player.speed
	Player.isRunnging = true
	Player.isFacingRight = true
	scroll_background(true)
	sprite := a.get_anim_sprite(RUN_ANIM, 100/int(Player.speed))
	sprite.Width = Player.transform.ToInt32().Width
	sprite.Height = Player.transform.ToInt32().Height
	Player.sprite = sprite
}
func (a *App) play_attack_animation() {
	sprite := rl.Texture2D{}
	switch Player.attackType {
	case ATTACK_LIGHT:
		sprite = *a.get_anim_sprite(ATTACK_ANIM, 4)
	case ATTACK_LIGHT_BACK:
		sprite = *a.get_anim_sprite(ATTACK2_ANIM, 2)
	case ATTACK_HEAVY:
		sprite = *a.get_anim_sprite(ATTACK3_ANIM, 4)
	}
	sprite.Width = Player.transform.ToInt32().Width
	sprite.Height = Player.transform.ToInt32().Height
	Player.sprite = &sprite
}

func draw_background() {
	if int32(scrolling_back) <= -bg_layer_0.Width*2 {
		scrolling_back = 0
	}
	if int32(scrolling_backmid) <= -bg_layer_1.Width*2 {
		scrolling_backmid = 0
	}
	if int32(scrolling_mid) <= -bg_layer_2.Width*2 {
		scrolling_mid = 0
	}
	if int32(scrolling_midfore) <= -bg_layer_3.Width*2 {
		scrolling_midfore = 0
	}
	if int32(scrolling_fore) <= -bg_layer_4.Width*2 {
		scrolling_fore = 0
	}
	if int32(scrolling_forefore) <= -bg_layer_5.Width*2 {
		scrolling_forefore = 0
	}

	bg_layer_0.Width = int32(windowSize.X)
	bg_layer_0.Height = int32(windowSize.Y)
	bg_layer_1.Width = int32(windowSize.X)
	bg_layer_1.Height = int32(windowSize.Y)
	bg_layer_2.Width = int32(windowSize.X)
	bg_layer_2.Height = int32(windowSize.Y)
	bg_layer_3.Width = int32(windowSize.X)
	bg_layer_3.Height = int32(windowSize.Y)
	bg_layer_4.Width = int32(windowSize.X)
	bg_layer_4.Height = int32(windowSize.Y)
	bg_layer_5.Width = int32(windowSize.X)
	bg_layer_5.Height = int32(windowSize.Y)
	bg_layer_6.Width = int32(windowSize.X)
	bg_layer_6.Height = int32(windowSize.Y)
	bg_layer_7.Width = int32(windowSize.X)
	bg_layer_7.Height = int32(windowSize.Y)
	bg_layer_8.Width = int32(windowSize.X)
	bg_layer_8.Height = int32(windowSize.Y)

	rl.DrawTextureRec(*bg_layer_0, rl.Rectangle{X: float32(scrolling_backback), Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
	rl.DrawTextureRec(*bg_layer_1, rl.Rectangle{X: float32(scrolling_back), Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
	rl.DrawTextureRec(*bg_layer_2, rl.Rectangle{X: float32(scrolling_clouds), Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
	rl.DrawTextureRec(*bg_layer_3, rl.Rectangle{X: float32(scrolling_backmid), Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
	rl.DrawTextureRec(*bg_layer_4, rl.Rectangle{X: float32(scrolling_mid), Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
	rl.DrawTextureRec(*bg_layer_5, rl.Rectangle{X: float32(scrolling_midmid), Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
	rl.DrawTextureRec(*bg_layer_6, rl.Rectangle{X: float32(scrolling_midmidmid), Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
	rl.DrawTextureRec(*bg_layer_7, rl.Rectangle{X: float32(scrolling_midfore), Y: 0, Width: windowSize.X, Height: windowSize.Y}, rl.Vector2{X: 0, Y: 0}, color.RGBA{255, 255, 255, 255})
}
func scroll_background(is_right bool) {
	scroll_factor := -1
	if is_right {
		scroll_factor = 1
	}
	scrolling_backback += 1 * scroll_factor
	scrolling_back += 2 * scroll_factor
	scrolling_backmid += 3 * scroll_factor
	scrolling_mid += 4 * scroll_factor
	scrolling_midmid += 5 * scroll_factor
	scrolling_midmidmid += 6 * scroll_factor
	scrolling_midfore += 7 * scroll_factor
	scrolling_fore += int(Player.speed*15*deltaTime) * scroll_factor
	scrolling_forefore += int(Player.speed*30*deltaTime) * scroll_factor
}

// Utility functions
func get_texture(sprite *rl.Image) *rl.Texture2D {
	result := rl.LoadTextureFromImage(sprite)
	return &result
}
func (a *App) reset_for_animation() {
	animationPhase = 0
	a.frameCount = 0
}
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %f MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %f MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %f MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
func bToMb(b uint64) float32 {
	return float32(b) / 1024.0 / 1024.0
}
